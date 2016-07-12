#!/usr/bin/env bash

BASH_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
ROOT_DIR=$BASH_DIR/..

# Get dependencies
echo "Getting dependencies..."
go get github.com/cyrillrx/go-rabbit/...
go test github.com/cyrillrx/go-rabbit/...

# Build
echo "Building..."
cd $GOPATH/bin
go build github.com/cyrillrx/go-rabbit
