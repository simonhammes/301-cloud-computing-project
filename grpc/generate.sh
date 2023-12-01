#!/usr/bin/env bash

set -exuo pipefail

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    search/search.proto
