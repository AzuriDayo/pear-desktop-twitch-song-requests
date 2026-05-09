#!/bin/bash
set -e
cd "$(git rev-parse --show-toplevel)"
pnpm run -r build
rm -rf cmd/main/build
cp -r control-panel/build cmd/main