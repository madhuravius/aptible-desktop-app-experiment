#!/usr/bin/env bash

echo "Rebuilding CLI and placing in correct location"

pushd cli
#cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ../public/wasm_exec.js
#GOOS=js GOARCH=wasm go build -o ../public/cli.wasm main.go
go build -o ../public/cli main.go
popd cli