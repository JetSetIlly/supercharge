#!/bin/bash

GOOS="linux" GOARCH="amd64" go build -ldflags "-s -w" -o supercharge_linux_amd64 .
GOOS="windows" GOARCH="amd64" go build -ldflags "-s -w" -o supercharge_windows_amd64 .
