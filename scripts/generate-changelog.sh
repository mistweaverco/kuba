#!/usr/bin/env bash

# Fetch latest version from git tags and strip the leading 'v' if present
VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
PKG_VERSION="{\"version\": \"${VERSION}\"}"
SOURCE_FILE="CHANGELOG.md"
TARGET_FILE="internal/changelog/changelog.md"

echo "Generating changelog for version: ${VERSION}"
echo "Using PKG_VERSION: ${PKG_VERSION}"

./node_modules/.bin/conventional-changelog -k <(echo "${PKG_VERSION}") -i ${SOURCE_FILE} -s -r 0

if [[ ! -f "${SOURCE_FILE}" ]]; then
  echo "ERROR: ${SOURCE_FILE} not found"
  exit 1
fi

mkdir -p "$(dirname "${TARGET_FILE}")"
cp "${SOURCE_FILE}" "${TARGET_FILE}"
echo "Synced embedded changelog to ${TARGET_FILE}"

echo "Contents of ${SOURCE_FILE}:"
cat "${SOURCE_FILE}"
echo

echo "Contents of ${TARGET_FILE}:"
cat "${TARGET_FILE}"
echo

echo "Changelog generation complete."
