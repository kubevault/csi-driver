#!/usr/bin/env bash

pushd $GOPATH/src/kubevault.dev/csi-driver/hack/gendocs
go run main.go
popd
