#!/usr/bin/env nix-shell
#! nix-shell -i bash -p entr
# shellcheck shell=bash

# Watch source files and recompile when any change.
export PORT=8080
git ls-files | entr -ccr -s "go build && go test && go run ."
