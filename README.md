[![Go Report Card](https://goreportcard.com/badge/github.com/kubevault/csi-driver)](https://goreportcard.com/report/github.com/kubevault/csi-driver)
[![Build Status](https://travis-ci.org/kubevault/csi-driver.svg?branch=master)](https://travis-ci.org/kubevault/csi-driver)
[![codecov](https://codecov.io/gh/kubevault/csi-driver/branch/master/graph/badge.svg)](https://codecov.io/gh/kubevault/csi-driver)
[![Docker Pulls](https://img.shields.io/docker/pulls/kubevault/csi-driver.svg)](https://hub.docker.com/r/kubevault/csi-driver/)
[![Slack](http://slack.kubernetes.io/badge.svg)](http://slack.kubernetes.io/#pharmer)
[![Twitter](https://img.shields.io/twitter/follow/appscodehq.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=AppsCodeHQ)


# csi-driver


#### Issue tracking

https://github.com/kubernetes/kubernetes/issues/66362


A Container Storage Interface ([CSI](https://github.com/container-storage-interface/spec)) Driver for `Vault`, which will act as a source of secrets of kubernetes cluster.
The CSI plugin allows you to use `Vault` with your preferred Container Orchestrator.


## Installing to Kubernetes

**Requirements:**

* Kubernetes v1.10 minimum
* `--allow-privileged` flag must be set to true for both the API server and the kubelet
* (if you use Docker) the Docker daemon of the cluster nodes must allow shared mounts
* Pre-installed `vault`. To install vault on kubernetes, follow [this](docs/vault-install.md)
* pass `--feature-gates=CSIDriverRegistry=true,CSINodeInfo=true` to kubelet and kube-apiserver



### 1. Create a secret with your Vault root token

Replace the placeholder string with your own token and save it as `secret.yaml`

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: vault
  namespace: kube-system
stringData:
  token: "___REPLACE_ME___"
```

and create the secret using `kubectl` :

```bash
$ kubectl create -f ./secret.yaml
secret "vault" created
```

You should now see the `vault` secret in the `kube-system` namespace along with other secrets

```sh
$ kubectl -n kube-system get secrets
NAME                  TYPE                                  DATA      AGE
default-token-jskxx   kubernetes.io/service-account-token   3         18h
vault                 Opaque                                1         18h
```


#### 2. Deploy the CSI plugin and sidecars:

Before you continue, be sure to checkout to a [tagged release](https://github.com/kubevault/csi-driver/releases). For
example, to use the version `v0.0.1` you can execute the following command:

```sh
kubectl apply -f https://raw.githubusercontent.com/kubevault/csi-driver/master/hack/deploy/releases/csi-vault-v0.0.1.yaml
```


#### 3. Deploy storage class of your choice

create a policy on `vault` using following capabilities:
```hcl
# capability to create a token against the "nginx" role
path "auth/token/create/nginx" {
  capabilities = ["create", "read", "update", "delete", "list"]
}

path "auth/token/roles/nginx" {
  capabilities = ["read"]
}

# capability to list roles
path "auth/token/roles" {
  capabilities = ["read", "list"]
}

# capability of get secret
path "kv/*" {
  capabilities = ["read"]
}
```

If you have a KV secrets on your vault and you also have certain policy to access that secrets, you have to create a `storage-class.yaml` file and put the following data

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: vault-kv-storage
  annotations:
    storageclass.kubernetes.io/is-default-class: "false"
provisioner: com.vault.csi.vaultdbs
parameters:
  fsType: tmpfs
  policy: nginx #policy name which exists on vault
  secretEngine: KV # vault engine name
  secretName: my-secret # secret name on vault which you want get access

```


then create the storage class using `kubectl`.


#### 4. Test and verify

Create secret on vault with following command:

```bash
$ vault kv put kv/my-secret my-value=s3cr3t
```


Create a PersistentVolumeClaim. This makes sure a volume is created and provisioned on your behalf:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: csi-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: vault-kv-storage
  volumeMode: DirectoryOrCreate

```

After that create a Pod that refers to this volume. When the Pod is created, the volume will be attached, formatted and mounted to the specified Container

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: my-vault-app
spec:
  containers:
    - name: my-frontend
      image: busybox
      volumeMounts:
      - mountPath: "/testdata"
        name: my-vault-volume
        readOnly: true
      command: [ "sleep", "1000000" ]
  volumes:
    - name: my-vault-volume
      persistentVolumeClaim:
        claimName: csi-pvc
```

Check if the pod is running successfully:

```sh
kubectl describe pods/my-csi-app
```


Check inside the app container:

```sh
$ kubectl exec -ti my-csi-app /bin/sh
/ # ls /testdata/
my-value
/ # cat /testdata/my-value
s3cr3t
```
