#!/bin/bash
set -e
cd "$(git rev-parse --show-toplevel)"
vp run build
rm -rf cmd/main/build
cp -r control-panel/build cmd/main