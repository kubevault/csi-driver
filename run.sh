#!/usr/bin/env bash

kubectl delete -f hack/deploy/example/kv/csi-app.yaml

kubectl delete -f hack/deploy/example/kv/csi-pvc.yaml

kubectl delete -f hack/deploy/example/kv/csi-storageclass.yaml

kubectl delete -f hack/deploy/releases/csi-vault-v0.1.3.yaml

