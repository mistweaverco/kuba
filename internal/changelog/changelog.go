package changelog

import _ "embed"

// Markdown contains the baked-in changelog content.
//
// It is refreshed from ./CHANGELOG.md by scripts/sync-embedded-changelog.sh
// (invoked from scripts/build.sh after `vp run changelog`).
//
//go:embed changelog.md
var Markdown string
