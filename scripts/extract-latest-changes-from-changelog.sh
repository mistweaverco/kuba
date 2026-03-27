#!/usr/bin/env bash

# Fetch latest version from git tags and strip the leading 'v' if present
VERSION=${VERSION:-$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')}
CHANGELOG_SOURCE_FILE=${SOURCE_FILE:-"CHANGELOG.md"}

if [[ -f "$CHANGELOG_SOURCE_FILE" ]]; then
  if [[ -z "$VERSION" ]]; then
    echo "Error: could not determine VERSION from git tags." >&2
    exit 1
  fi

  # NOTE:
  # We start printing at the matching heading and stop when the next version heading starts.
  extract_for_version() {
    local v="$1"
    awk -v version="$v" '
    BEGIN { in_section = 0; printed = 0 }
    # Match headings like:
    #   ## 1.10.0 (YYYY-MM-DD)
    #   ## <small>1.8.3 (YYYY-MM-DD)</small>
    #   ## [1.8.3](...) (YYYY-MM-DD)
    /^##[[:space:]]/ && match($0, /[0-9]+\.[0-9]+\.[0-9]+/) {
      current = ""
      current = substr($0, RSTART, RLENGTH)
      if (in_section) exit
      if (current == version) {
        in_section = 1
      }
    }
    in_section { print; printed = 1 }
    END {
      if (!printed) {
        exit 2
      }
    }
  ' "$CHANGELOG_SOURCE_FILE"
  }

  extract_for_version "$VERSION"
  rc=$?
  if [[ $rc -ne 0 ]]; then
    if [[ $rc -ne 2 ]]; then
      exit $rc
    fi

    # NOTE:
    # If the latest git tag isn't in CHANGELOG.md yet
    # (common in local development), fall back
    # to the newest version present in the changelog.
    latest_in_changelog=$(
      awk '
        /^##[[:space:]]/ && match($0, /[0-9]+\.[0-9]+\.[0-9]+/) {
          print substr($0, RSTART, RLENGTH)
          exit
        }
      ' "$CHANGELOG_SOURCE_FILE"
    )

    if [[ -z "$latest_in_changelog" ]]; then
      echo "Error: could not parse any versions from $CHANGELOG_SOURCE_FILE" >&2
      exit 1
    fi

    extract_for_version "$latest_in_changelog"
  fi
  else
    echo "Error: $CHANGELOG_SOURCE_FILE not found."
    exit 1
fi

