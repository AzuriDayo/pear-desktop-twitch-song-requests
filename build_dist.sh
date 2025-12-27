#!/bin/bash
set -ex
appname=pear-desktop-twitch-song-requests
archfname=(amd64 arm64)
osfname=(linux windows darwin)

go mod download

for os in "${osfname[@]}"; do
    for arch in "${archfname[@]}"; do
        execname=${appname}_${os}_${arch}
        if [ "$os" = 'windows' ] && [ "$arch" = 'arm64' ]; then
            execname=${appname}_${os}_on_arm.exe
        elif [ "$os" = 'windows' ] && [ "$arch" = 'amd64' ]; then
            execname=${appname}_${os}_x64.exe
        elif [ "$os" = 'darwin' ] && [ "$arch" = 'arm64' ]; then
            execname=${appname}_macos_apple_silicon.exe
        elif [ "$os" = 'darwin' ] && [ "$arch" = 'amd64' ]; then
            execname=${appname}_macos_intel.exe
        fi
        CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags="-s -w -X main.version=${GIT_SHORT_TAG}" -trimpath -o "$execname" ./cmd/main &
    done
done

wait
