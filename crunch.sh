#!/bin/sh
# Build and crunch the binary

CGO_ENABLED=0 go install -ldflags="-s -w" github.com/thechriswalker/system-stats
upx --best --ultra-brute ${GOBIN}/system-stats
