#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
SRC=$GOPATH/src
BIN=$GOPATH/bin
ROOT=$GOPATH
REPO_ROOT=$GOPATH/src/github.com/kubevault/csi-driver

source "$REPO_ROOT/hack/libbuild/common/pharmer_image.sh"

APPSCODE_ENV=${APPSCODE_ENV:-dev}
IMG=csi-vault

DIST=$GOPATH/src/github.com/kubevault/csi-driver/dist
mkdir -p $DIST
if [ -f "$DIST/.tag" ]; then
  export $(cat $DIST/.tag | xargs)
fi

clean() {
  pushd $GOPATH/src/github.com/kubevault/csi-driver/hack/docker
  rm csi-vault Dockerfile
  popd
}

build_binary() {
  pushd $GOPATH/src/github.com/kubevault/csi-driver
  ./hack/builddeps.sh
  ./hack/make.py build
  detect_tag $DIST/.tag
  popd
}

build_docker() {
  pushd $GOPATH/src/github.com/kubevault/csi-driver/hack/docker
  cp $DIST/csi-vault/csi-vault-alpine-amd64 csi-vault
  chmod 755 csi-vault

  cat >Dockerfile <<EOL
FROM alpine

RUN set -x \
  && apk add --update --no-cache ca-certificates tzdata e2fsprogs findmnt


ADD csi-vault /usr/bin/csi-vault

ENTRYPOINT ["csi-vault"]

EOL
  local cmd="docker build -t pharmer/$IMG:$TAG ."
  echo $cmd
  $cmd

  rm csi-vault Dockerfile
  popd
}

build() {
  build_binary
  build_docker
}

docker_push() {
  if [ "$APPSCODE_ENV" = "prod" ]; then
    echo "Nothing to do in prod env. Are you trying to 'release' binaries to prod?"
    exit 0
  fi
  if [ "$TAG_STRATEGY" = "git_tag" ]; then
    echo "Are you trying to 'release' binaries to prod?"
    exit 1
  fi
  hub_canary
}

docker_release() {
  if [ "$APPSCODE_ENV" != "prod" ]; then
    echo "'release' only works in PROD env."
    exit 1
  fi
  if [ "$TAG_STRATEGY" != "git_tag" ]; then
    echo "'apply_tag' to release binaries and/or docker images."
    exit 1
  fi
  hub_up
}

source_repo $@
