#!/usr/bin/env bash

kubectl apply -f hack/deploy/releases/csi-vault-v0.1.3.yaml

kubectl apply -f hack/deploy/example/kv/csi-storageclass.yaml

kubectl apply -f hack/deploy/example/kv/csi-pvc.yaml

kubectl apply -f hack/deploy/example/kv/csi-app.yaml


