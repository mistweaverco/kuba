#!/usr/bin/env bash

# Fetch latest version from git tags and strip the leading 'v' if present
VERSION=${VERSION:-$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')}
SOURCE_FILE="CHANGELOG.md"
TARGET_FILE="internal/changelog/changelog.md"
TEMP_CONFIG=$(mktemp --suffix=".json")
echo "{\"version\": \"$VERSION\", \"date\": \"$(date +%Y-%m-%d)\"}" > "$TEMP_CONFIG"

echo "Generating changelog for version: ${VERSION}"
echo "Using PKG_VERSION: ${PKG_VERSION}"

if [[ -n "$CI" ]]; then
  # HACK:
  # In CI environments,
  # we might have a detached HEAD state
  # This will cause conventional-changelog to
  # only generate a changelog for previous version
  # and see the current version unreleased
  # To fix this, we remove the existing tag for
  # the current version (if it exists)
  # and restore it after generating the changelog
  git tag -d "v$VERSION" 2>/dev/null
fi

./node_modules/.bin/conventional-changelog \
  -i "$SOURCE_FILE" \
  -s \
  -r 0 \
  -u \
  -k "$TEMP_CONFIG" \
  -c "$TEMP_CONFIG"

rm "$TEMP_CONFIG"

if [[ -n "$CI" ]]; then
  # HACK:
  # Restore the tag for the current version after generating the changelog
  if git rev-parse "v$VERSION" >/dev/null 2>&1; then
    echo "Tag v$VERSION already exists, skipping tag creation"
  else
    git tag "v$VERSION" 2>/dev/null
  fi
fi

if [[ ! -f "${SOURCE_FILE}" ]]; then
  echo "ERROR: ${SOURCE_FILE} not found"
  exit 1
fi

mkdir -p "$(dirname "${TARGET_FILE}")"
cp "${SOURCE_FILE}" "${TARGET_FILE}"
echo "Synced embedded changelog to ${TARGET_FILE}"
