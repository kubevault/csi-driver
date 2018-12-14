
# CSI Vault

[CSI Driver for Vault by AppsCode](https://github.com/kubevault/csi-driver)

## TL;DR;

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm install appscode/csi-vault --name csi-vault --namespace kube-system
```

## Introduction

This chart bootstraps a [Vault CSI Driver](https://github.com/kubevault/csi-driver) on a [Kubernetes](http://kubernetes.io)  cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes v1.12+
- `--allow-privileged` flag must be set to true for both the API server and the kubelet
- (If you use Docker) The Docker daemon of the cluster nodes must allow shared mounts
- Pre-installed HasiCorp Vault server.
- Pass `--feature-gates=CSIDriverRegistry=true,CSINodeInfo=true` to kubelet and kube-apiserver


## Installing the Chart

To install the chart with the release name `csi-vault`

```console
$ helm install appscode/csi-vault --name csi-vault
```

This command deploys CSI Driver for Vault on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `csi-vault`:

```console
$ helm delete csi-vault
```

The command removes all the Kubernetes components associated with the chart and deletes the release.


## Configuration

The following table lists the configurable parameters of the Stash chart and their default values.


| Parameter                 | Description                                                        | Default                                    |
| --------------------------| ------------------------------------------------------------------ | -------------------------------------------|
| `replicaCount`            | Number of Vault operator replicas to create (only 1 is supported)  | `1`                                        |
| `csi.registry`            | Docker registry used to pull Vault CSI driver image                | `kubevault`                                |
| `csi.repository`          | Vault CSI driver container image                                   | `csi-vault`                                |
| `csi.tag`                 | Vault CSI driver container image tag                               | `0.1.0`                                    |
| `csi.pullPolicy`          | Vault CSI driver container image pull policy                       | `IfNotPresent`                             |
| `attacher.registry`       | Docker registry used to pull CSI attacher image                    | `quay.io/k8scsi`                           |
| `attacher.repository`     | CSI attacher container image                                       | `csi-attacher`                             |
| `attacher.tag`            | CSI attacher container image tag                                   | `v0.2.0`                                   |
| `attacher.pullPolicy`     | CSI attacher container image pull policy                           | `IfNotPresent`                             |
| `provisioner.registry`    | Docker registry used to pull CSI provisioner image                 | `quay.io/k8scsi`                           |
| `provisioner.repository`  | CSI provisioner container image                                    | `csi-provisioner`                          |
| `provisioner.tag`         | CSI provisioner container image tag                                | `v0.2.1`                                   |
| `provisioner.pullPolicy`  | CSI provisioner container image pull policy                        | `IfNotPresent`                             |
| `registrar.registry`      | Docker registry used to pull CSI driver registrar image            | `quay.io/k8scsi`                           |
| `registrar.repository`    | CSI driver registrar container image                               | `driver-registrar`                         |
| `registrar.tag`           | CSI driver registrar container image tag                           | `v0.3.0`                                   |
| `registrar.pullPolicy`    | CSI driver registrar container image pull policy                   | `IfNotPresent`                             |
| `logLevel`                | Log level for container                                            | `5`                                        |
| `driverName`              | Vault CSI driver name                                              | `com.kubevault.csi.secrets`                   |
| `pluginAddress`           | Vault CSI driver endpoint address                                  | `/var/lib/csi/sockets/pluginproxy/csi.sock`|
| `pluginDir`               | Vault CSI driver plugin directory                                  | `/var/lib/csi/sockets/pluginproxy/`        |
| `attachRequired`          | Indicates CSI volume driver requires an attach operation           | `false`                                    |
| `appbinding.create`       | If true, AppBinding CRD will be created                            | `true`                                     |
| `affinity`                | Affinity rules for pod assignment                                  | `{}`                                       |
| `nodeSelector`            | Node labels for pod assignment                                     | `{}`                                       |
| `tolerations`             | Tolerations used pod assignment                                    | `{}`                                       |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install --name csi-vault --set csi.tag=v0.1.0 appscode/csi-vault

```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example:

```console
$ helm install --name csi-vault --values values.yaml appscode/csi-vault
```
