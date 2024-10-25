#!/bin/bash

if [ -z $1 ]; then
  echo "Version is required"
  exit
fi

go run ./cmd/release ./cmd/release --ldflags="-X main.Version=${1}" --extra-arches -c tar.gz