#!/usr/bin/env bash

# Fetch latest version from git tags and strip the leading 'v' if present
VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
PKG_VERSION="{\"version\": \"${VERSION}\"}"
echo "Generating changelog for version: ${VERSION}"
echo "Using PKG_VERSION: ${PKG_VERSION}"

./node_modules/.bin/conventional-changelog -k <(echo "$PKG_VERSION") -i CHANGELOG.md -s -r 0
