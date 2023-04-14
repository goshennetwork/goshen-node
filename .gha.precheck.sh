#!/usr/bin/env bash

set -ex

VERSION=$(git describe --always --tags --long)
echo $VERSION

if [ $RUNNER_OS == 'Linux' ]; then
  echo "linux sys"
  env
  export GOPATH="/home/runner/go"
  export PATH=$PATH:$GOPATH/bin
  #go test -v ./...
  export GOPRIVATE=github.com/goshennetwork

  go mod tidy

  bash ./.gha.gofmt.sh

  make geth

  go run build/ci.go test -dlgo

  #quit when meet first fail test
  #for s in $(go list ./...); do if ! go test -failfast -v -p 1 $s; then exit 1; fi; done
  fi
