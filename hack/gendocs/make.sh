#!/usr/bin/env bash

pushd $GOPATH/src/github.com/kubevault/csi-driver/hack/gendocs
go run main.go
popd
