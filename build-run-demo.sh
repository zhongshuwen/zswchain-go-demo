#!/bin/sh
goreleaser build --single-target --rm-dist --snapshot -o ./dist/zswchain-go-demo

source ./.env.sh
./dist/zswchain-go-demo full
