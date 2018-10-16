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

* Kubernetes v1.12 minimum
* `--allow-privileged` flag must be set to true for both the API server and the kubelet
* (if you use Docker) the Docker daemon of the cluster nodes must allow shared mounts
* Pre-installed `vault`. To install vault on kubernetes, follow [this](docs/vault-install.md)
* Pass `--feature-gates=CSIDriverRegistry=true,CSINodeInfo=true` to kubelet and kube-apiserver


#### 1. Install CSI driver on cluster

To install `csidriver` and `csinodeinfo` crds, apply this [file](hack/deploy/csi-crd.yaml) by running

```sh
kubectl apply -f https://raw.githubusercontent.com/kubevault/csi-driver/master/hack/deploy/csi-crd.yaml
```


#### 2. Create a secret with your Vault root token and address

Replace the placeholder string with your own token and save it as `secret.yaml`

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: vault
  namespace: kube-system
stringData:
  token: "___REPLACE_ME___"
  url: "http://REPLACE_ME__"
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


#### 3. Deploy the CSI plugin and sidecars:

Before you continue, be sure to checkout to a [tagged release](https://github.com/kubevault/csi-driver/releases). For
example, to use the version `v0.1.1` you can execute the following command:

```sh
kubectl apply -f https://raw.githubusercontent.com/kubevault/csi-driver/master/hack/deploy/releases/csi-vault-v0.1.1.yaml
```

#### 4. Create policy and role for service account

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

# capability to get aws credentials
path "aws/*" {
  capabilities = ["read"]
}

```
run

```bash
$ vault policy write test-policy policy.hcl
```
then create a file `serviceaccount.yaml` with following contents

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: role-tokenreview-binding
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:auth-delegator
subjects:
- kind: ServiceAccount
  name: postgres-vault
  namespace: default
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: postgres-vault
```

After that run `kubectl apply -f serviceaccount.yaml` to create service account.

To enable Kubernetes auth backend by extracting the token reviewer JWT, Kubernetes CA certificate and Kubernetes host,

```bash
$ export VAULT_SA_NAME=$(kubectl get sa postgres-vault -o jsonpath="{.secrets[*]['name']}")

$ export SA_JWT_TOKEN=$(kubectl get secret $VAULT_SA_NAME -o jsonpath="{.data.token}" | base64 --decode; echo)

$ export SA_CA_CRT=$(kubectl get secret $VAULT_SA_NAME -o jsonpath="{.data['ca\.crt']}" | base64 --decode; echo)

$ export K8S_HOST=$(kubectl exec consul-consul-0 -- sh -c 'echo $KUBERNETES_SERVICE_HOST')
$ export K8s_PORT=6443
```

Next we can enable the kubernetes authentication backend and create vault role that is attached to service account

```bash
$ vault auth enable kubernetes
$ vault write auth/kubernetes/config \
    token_reviewer_jwt="$SA_JWT_TOKEN" \
    kubernetes_host="https://$K8S_HOST:$k8s_PORT" \
    kubernetes_ca_cert="$SA_CA_CRT"

$ vault write auth/kubernetes/role/testrole \
      bound_service_account_names=postgres-vault \
      bound_service_account_namespaces=default \
      policies=test-policy \
      ttl=24h
```

#### 4. Deploy storage class of your choice


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
  authRole: testrole #vault role for authentication
  secretEngine: KV # vault engine name
  secretName: my-secret # secret name on vault which you want get access

```


then create the storage class using `kubectl`.


#### 5. Test and verify

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
      storage: 1Gi
  storageClassName: vault-kv-storage
  volumeMode: DirectoryOrCreate

```

After that create a Pod that refers to this volume. When the Pod is created, the volume will be attached, formatted and mounted to the specified Container

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mypod
spec:
  containers:
  - name: mypod
    image: redis
    volumeMounts:
    - name: my-vault-volume
      mountPath: "/etc/foo"
      readOnly: true
  serviceAccountName: postgres-vault
  volumes:
    - name: my-vault-volume
      persistentVolumeClaim:
        claimName: csi-pvc
```

Check if the pod is running successfully:

```sh
kubectl describe pods/my-pod
```


Check inside the app container:

```sh
$ kubectl exec -ti mypod /bin/sh
/ # ls /etc/foo
my-value
/ # cat /etc/foo/my-value
s3cr3t
```


* To setup AWS secret engine on vault click [here](docs/engines/aws.md)
* To setup PKI secret engine on vault click [here](docs/engines/pki.md)