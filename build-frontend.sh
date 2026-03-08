#!/bin/bash
set -e
cd "$(git rev-parse --show-toplevel)"
bun run --workspaces build
rm -rf cmd/main/build
cp -r control-panel/build cmd/main