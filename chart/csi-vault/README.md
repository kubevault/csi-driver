
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

- Kubernetes v1.13+
- `--allow-privileged` flag must be set to true for both the API server and the kubelet
- (If you use Docker) The Docker daemon of the cluster nodes must allow shared mounts
- Pre-installed HashiCorp Vault server.
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


| Parameter                               | Description                                                        | Default                                    |
| --------------------------------------- | ------------------------------------------------------------------ | -------------------------------------------|
| `replicaCount`                          | Number of Vault operator replicas to create (only 1 is supported)  | `1`                                        |
| `attacher.name`                         | Name of the attacher component                                     | `attacher`                                 |
| `attacher.registry`                     | Docker registry used to pull CSI attacher image                    | `quay.io/k8scsi`                           |
| `attacher.repository`                   | CSI attacher container image                                       | `csi-attacher`                             |
| `attacher.tag`                          | CSI attacher container image tag                                   | `v1.0.1`                                   |
| `attacher.pullPolicy`                   | CSI attacher container image pull policy                           | `IfNotPresent`                             |
| `plugin.name`                           | Name of the plugin component                                       | `plugin`                                   |
| `plugin.registry`                       | Docker registry used to pull Vault CSI driver image                | `kubevault`                                |
| `plugin.repository`                     | Vault CSI driver container image                                   | `csi-vault`                                |
| `plugin.tag`                            | Vault CSI driver container image tag                               | `0.2.0`                                    |
| `plugin.pullPolicy`                     | Vault CSI driver container image pull policy                       | `IfNotPresent`                             |
| `provisioner.name`                      | Name of the provisioner component                                  | `provisioner`                              |
| `provisioner.registry`                  | Docker registry used to pull CSI provisioner image                 | `quay.io/k8scsi`                           |
| `provisioner.repository`                | CSI provisioner container image                                    | `csi-provisioner`                          |
| `provisioner.tag`                       | CSI provisioner container image tag                                | `v1.0.1`                                   |
| `provisioner.pullPolicy`                | CSI provisioner container image pull policy                        | `IfNotPresent`                             |
| `clusterRegistrar.registry`             | Docker registry used to pull CSI driver cluster registrar image    | `quay.io/k8scsi`                           |
| `clusterRregistrar.repository`          | CSI driver cluster registrar container image                       | `csi-cluster-driver-registrar`             |
| `clusterRregistrar.tag`                 | CSI driver cluster registrar container image tag                   | `v1.0.1`                                   |
| `clusterRregistrar.pullPolicy`          | CSI driver cluster registrar container image pull policy           | `IfNotPresent`                             |
| `nodeRegistrar.registry`                | Docker registry used to pull CSI driver node registrar image       | `quay.io/k8scsi`                           |
| `nodeRregistrar.repository`             | CSI driver node registrar container image                          | `csi-node-driver-registrar`                |
| `nodeRregistrar.tag`                    | CSI driver node registrar container image tag                      | `v1.0.1`                                   |
| `nodeRregistrar.pullPolicy`             | CSI driver node registrar container image pull policy              | `IfNotPresent`                             |
| `driverName`                            | Vault CSI driver name                                              | `com.kubevault.csi.secrets`                |
| `pluginAddress`                         | Vault CSI driver endpoint address                                  | `/var/lib/csi/sockets/pluginproxy/csi.sock`|
| `pluginDir`                             | Vault CSI driver plugin directory                                  | `/var/lib/csi/sockets/pluginproxy/`        |
| `attachRequired`                        | Indicates CSI volume driver requires an attach operation           | `false`                                    |
| `appbinding.create`                     | If true, AppBinding CRD will be created                            | `true`                                     |
| `imagePullSecrets`                      | Specify image pull secrets                                         | `nil` (does not add image pull secrets to deployed pods) |
| `criticalAddon`                         | If true, installs Vault CSI driver as critical addon               | `false`                                    |
| `logLevel`                              | Log level for CSI driver                                           | `3`                                        |
| `affinity`                              | Affinity rules for pod assignment                                  | `{}`                                       |
| `nodeSelector`                          | Node labels for pod assignment                                     | `{}`                                       |
| `tolerations`                           | Tolerations used pod assignment                                    | `{}`                                       |
| `apiserver.useKubeapiserverFqdnForAks`  | If true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 | `true`             |
| `apiserver.healthcheck.enabled`         | Enable readiness and liveliness probes                             | `true`                                     |
| `enableAnalytics`                       | Send usage events to Google Analytics                              | `true`                                     |
| `monitoring.agent`                      | Specify which monitoring agent to use for monitoring Vault. It accepts either `prometheus.io/builtin` or `prometheus.io/coreos-operator`.                                  | `none`                                                    |
| `monitoring.node`                       | Specify whether to monitor Vault CSI driver node plugin.              | `false`                                    |
| `monitoring.controller`                 | Specify whether to monitor Vault CSI driver controllerplugin.                | `false`                                    |
| `monitoring.prometheus.namespace`       | Specify the namespace where Prometheus server is running or will be deployed.                                                                                              | Release namespace                                         |
| `monitoring.serviceMonitor.labels`      | Specify the labels for ServiceMonitor. Prometheus crd will select ServiceMonitor using these labels. Only usable when monitoring agent is `prometheus.io/coreos-operator`. | `app: <generated app name>` and `release: <release name>` |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install --name csi-vault --set plugin.tag=v0.2.0 appscode/csi-vault

```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example:

```console
$ helm install --name csi-vault --values values.yaml appscode/csi-vault
```
