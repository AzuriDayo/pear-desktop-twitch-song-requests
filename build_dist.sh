#!/bin/bash
set -ex
appname=pear-desktop-twitch-song-requests
archfname=(amd64 arm64)
osfname=(linux windows darwin)

go mod download

for os in "${osfname[@]}"; do
    for arch in "${archfname[@]}"; do
        execname=${appname}_${os}_${arch}
        if [ "$os" = 'windows' ]; then execname=$execname.exe ; fi
        CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags="-s -w -X main.version=${GIT_SHORT_TAG}" -trimpath -o "$execname" ./cmd/main &
    done
done

wait
