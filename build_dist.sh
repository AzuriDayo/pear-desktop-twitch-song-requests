#!/bin/bash

GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o main.exe cmd/main/main.go