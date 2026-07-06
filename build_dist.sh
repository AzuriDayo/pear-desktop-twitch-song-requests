#!/bin/bash
set -ex

# Wails v2 uses native WebViews and CGO, so cross-compilation is not supported:
# this script builds the desktop app for the HOST platform only. CI builds every
# platform on its own native runner (see .github/workflows/).

if [ "$GIT_SHORT_TAG" = '' ]; then
    GIT_SHORT_TAG="v0.0.0+$(git rev-parse --short HEAD)"
fi

cd "$(git rev-parse --show-toplevel)"

# Build the frontend and copy it into the Go embed directory.
vp run build

# Build the desktop app for the host platform (output in cmd/main/build/bin).
cd cmd/main
wails build -s -skipbindings -m -trimpath -ldflags "-X main.version=${GIT_SHORT_TAG}"
