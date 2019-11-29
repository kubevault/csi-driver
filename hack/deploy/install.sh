#!/bin/bash

# Copyright The KubeVault Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -eou pipefail

driver_crds=(
  csidrivers.csi.storage.k8s.io
  csinodeinfos.csi.storage.k8s.io
)

echo "checking kubeconfig context"
kubectl config current-context || {
  echo "Set a context (kubectl use-context <context>) out of the following:"
  echo
  kubectl config get-contexts
  exit 1
}
echo ""

OS=""
ARCH=""
DOWNLOAD_URL=""
DOWNLOAD_DIR=""
TEMP_DIRS=()
ONESSL=""
ONESSL_VERSION=v0.13.1

# http://redsymbol.net/articles/bash-exit-traps/
function cleanup() {
  rm -rf ca.crt ca.key server.crt server.key
  # remove temporary directories
  for dir in "${TEMP_DIRS[@]}"; do
    rm -rf "${dir}"
  done
}

# detect operating system
# ref: https://raw.githubusercontent.com/helm/helm/master/scripts/get
function detectOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    cygwin* | mingw* | msys*) OS='windows';;
  esac
}

# detect machine architecture
function detectArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

detectOS
detectArch

# download file pointed by DOWNLOAD_URL variable
# store download file to the directory pointed by DOWNLOAD_DIR variable
# you have to sent the output file name as argument. i.e. downloadFile myfile.tar.gz
function downloadFile() {
  if curl --output /dev/null --silent --head --fail "$DOWNLOAD_URL"; then
    curl -fsSL ${DOWNLOAD_URL} -o $DOWNLOAD_DIR/$1
  else
    echo "File does not exist"
    exit 1
  fi
}

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
trap cleanup EXIT

onessl_found() {
  # https://stackoverflow.com/a/677212/244009
  if [ -x "$(command -v onessl)" ]; then
    onessl version --check=">=${ONESSL_VERSION}" >/dev/null 2>&1 || {
      # old version of onessl found
      echo "Found outdated onessl"
      return 1
    }
    export ONESSL=onessl
    return 0
  fi
  return 1
}

# download onessl if it does not exist
onessl_found || {
  echo "Downloading onessl ..."

  ARTIFACT="https://github.com/kubepack/onessl/releases/download/${ONESSL_VERSION}"
  ONESSL_BIN=onessl-${OS}-${ARCH}
  case "$OS" in
    cygwin* | mingw* | msys*)
      ONESSL_BIN=${ONESSL_BIN}.exe
    ;;
  esac

  DOWNLOAD_URL=${ARTIFACT}/${ONESSL_BIN}
  DOWNLOAD_DIR="$(mktemp -dt onessl-XXXXXX)"
  TEMP_DIRS+=($DOWNLOAD_DIR) # store DOWNLOAD_DIR to cleanup later

  downloadFile $ONESSL_BIN # downloaded file name will be saved as the value of ONESSL_BIN variable

  export ONESSL=${DOWNLOAD_DIR}/${ONESSL_BIN}
  chmod +x $ONESSL
}

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export CSI_VAULT_NAMESPACE=kube-system
export CSI_VAULT_DOCKER_REGISTRY=${CSI_VAULT_DOCKER_REGISTRY:-kubevault}
export CSI_VAULT_DOCKER_REPOSITORY=csi-vault
export CSI_VAULT_IMAGE_TAG=${CSI_VAULT_IMAGE_TAG:-0.2.0}
export CSI_VAULT_IMAGE_PULL_SECRET_NAME=
export CSI_VAULT_IMAGE_PULL_POLICY=IfNotPresent
export CSI_ATTACHER_DOCKER_REGISTRY=${CSI_ATTACHER_DOCKER_REGISTRY:-quay.io/k8scsi}
export CSI_ATTACHER_DOCKER_REPOSITORY=csi-attacher
export CSI_ATTACHER_IMAGE_TAG=v1.2.0
export CSI_ATTACHER_IMAGE_PULL_SECRET_NAME=
export CSI_ATTACHER_IMAGE_PULL_POLICY=IfNotPresent
export CSI_PROVISIONER_DOCKER_REGISTRY=quay.io/k8scsi
export CSI_PROVISIONER_DOCKER_REPOSITORY=csi-provisioner
export CSI_PROVISIONER_IMAGE_TAG=v1.3.0
export CSI_PROVISIONER_IMAGE_PULL_SECRET_NAME=
export CSI_PROVISIONER_IMAGE_PULL_POLICY=IfNotPresent
export CSI_CLUSTER_REGISTRAR_DOCKER_REGISTRY=quay.io/k8scsi
export CSI_CLUSTER_REGISTRAR_DOCKER_REPOSITORY=csi-cluster-driver-registrar
export CSI_CLUSTER_REGISTRAR_IMAGE_TAG=v1.0.1
export CSI_CLUSTER_REGISTRAR_IMAGE_PULL_SECRET_NAME=
export CSI_CLUSTER_REGISTRAR_IMAGE_PULL_POLICY=IfNotPresent
export CSI_NODE_REGISTRAR_DOCKER_REGISTRY=quay.io/k8scsi
export CSI_NODE_REGISTRAR_DOCKER_REPOSITORY=csi-node-driver-registrar
export CSI_NODE_REGISTRAR_IMAGE_TAG=v1.1.0
export CSI_NODE_REGISTRAR_IMAGE_PULL_SECRET_NAME=
export CSI_NODE_REGISTRAR_IMAGE_PULL_POLICY=IfNotPresent
export CSI_VAULT_DRIVER_NAME=secrets.csi.kubevault.com
export CSI_VAULT_UNINSTALL=0
export CSI_VAULT_PURGE=0
export CSI_REQUIRED_ATTACHMENT=false
export REQUIRED_APPBINDING_INSTALL=true
export CSI_VAULT_PRIORITY_CLASS=system-cluster-critical

export CSI_VAULT_USE_KUBEAPISERVER_FQDN_FOR_AKS=false
export CSI_VAULT_ENABLE_ANALYTICS=false

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
export SCRIPT_LOCATION="curl -fsSL https://raw.githubusercontent.com/kubevault/csi-driver/0.2.0/"
if [ "$APPSCODE_ENV" = "dev" ]; then
  export SCRIPT_LOCATION="cat "
  export CSI_VAULT_IMAGE_PULL_POLICY=IfNotPresent
fi

KUBE_APISERVER_VERSION=$(kubectl version -o=json | $ONESSL jsonpath '{.serverVersion.gitVersion}')
$ONESSL semver --check='>=1.13.0' $KUBE_APISERVER_VERSION || {
  echo "This release of Vault CSI driver does not support Kubernetes version ${KUBE_APISERVER_VERSION}."
  echo
  exit 1
}
echo ""

MONITORING_AGENT_NONE="none"
MONITORING_AGENT_BUILTIN="prometheus.io/builtin"
MONITORING_AGENT_COREOS_OPERATOR="prometheus.io/coreos-operator"

export MONITORING_AGENT=${MONITORING_AGENT:-$MONITORING_AGENT_NONE}
export MONITOR_CONTROLLER_PLUGIN=${MONITOR_CONTROLLER_PLUGIN:-false}
export MONITOR_NODE_PLUGIN=${MONITOR_NODE_PLUGIN:-false}
export SERVICE_MONITOR_LABEL_KEY="app"
export SERVICE_MONITOR_LABEL_VALUE="csi-vault"

show_help() {
  echo "install.sh -install vault csi driver"
  echo " "
  echo "install.sh [options]"
  echo " "
  echo "options:"
  echo "-h, --help                                     show brief help"
  echo "-n, --namespace=NAMESPACE                      specify namespace (default: kube-system)"
  echo "    --csi-vault-docker-registry                docker registry used to pull csi-vault image (default: kubevault)"
  echo "    --csi-vault-image-pull-secret              name of secret used to pull csi-vault images"
  echo "    --csi-vault-image-tag                      docker image version of csi vault"
  echo "    --csi-attacher-docker-registry             docker registry used to pull csi attacher image (default: quay.io/k8scsi)"
  echo "    --csi-attacher-image-pull-secret           name of secret used to pull csi attacher image"
  echo "    --csi-attacher-image-tag                   docker image version of csi attacher"
  echo "    --csi-provisioner-docker-registry          docker registry used to pull csi provisioner image (default: quay.io/k8scsi)"
  echo "    --csi-provisioner-image-pull-secret        name of secret used to pull csi provisioner image"
  echo "    --csi-provisioner-image-tag                docker image version of csi provisioner"
  echo "    --csi-cluster-registrar-docker-registry    docker registry used to pull csi registrar image (default: quay.io/k8scsi)"
  echo "    --csi-cluster-registrar-image-pull-secret  name of secret used to pull csi registrar image"
  echo "    --csi-cluster-registrar-image-tag          docker image version of csi registrar"
  echo "    --csi-node-registrar-docker-registry       docker registry used to pull csi registrar image (default: quay.io/k8scsi)"
  echo "    --csi-node-registrar-image-pull-secret     name of secret used to pull csi registrar image"
  echo "    --csi-node-registrar-image-tag             docker image version of csi registrar"
  echo "    --csi-driver-name                          name of csi driver to install (default: secrets.csi.kubevault.com)"
  echo "    --csi-required-attachment                  indicates csi volume driver requires an attach operation (default: false)"
  echo "    --install-appbinding                       indicates appbinding crd need to be installed (default: true)"
  echo "    --monitoring-agent                         specify which monitoring agent to use (default: none)"
  echo "    --monitor-attacher                         specify whether to monitor Vault CSI driver attacher (default: false)"
  echo "    --monitor-plugin                           specify whether to monitor Vault CSI driver plugin (default: false)"
  echo "    --monitor-provisioner                      specify whether to monitor Vault CSI driver provisioner (default: false)"
  echo "    --prometheus-namespace                     specify the namespace where Prometheus server is running or will be deployed (default: same namespace as csi-vault)"
  echo "    --servicemonitor-label                     specify the label for ServiceMonitor crd. Prometheus crd will use this label to select the ServiceMonitor. (default: 'app: csi-vault')"
  echo "    --uninstall                                uninstall vault csi driver"
  echo "    --purge                                    purges csi driver crd objects and crds"
}

while test $# -gt 0; do
  case "$1" in
    -h | --help)
      show_help
      exit 0
      ;;
    -n)
      shift
      if test $# -gt 0; then
        export CSI_VAULT_NAMESPACE=$1
      else
        echo "no namespace specified"
        exit 1
      fi
      shift
      ;;
    --namespace*)
      export CSI_VAULT_NAMESPACE=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-vault-docker-registry*)
      export CSI_VAULT_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-vault-image-pull-secret*)
      export CSI_VAULT_IMAGE_PULL_SECRET_NAME=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-vault-image-tag*)
      export CSI_VAULT_IMAGE_TAG=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-attacher-docker-registry*)
      export CSI_ATTACHER_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-attacher-image-pull-secret*)
      export CSI_ATTACHER_IMAGE_PULL_SECRET_NAME=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-attacher-image-tag*)
      export CSI_ATTACHER_IMAGE_TAG=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-provisioner-docker-registry*)
      export CSI_PROVISIONER_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-provisioner-image-pull-secret*)
      export CSI_PROVISIONER_IMAGE_PULL_SECRET_NAME=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-provisioner-image-tag*)
      export CSI_PROVISIONER_IMAGE_TAG=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-cluster-registrar-docker-registry*)
      export CSI_CLUSTER_REGISTRAR_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-cluster-registrar-image-pull-secret*)
      export CSI_CLUSTER_REGISTRAR_IMAGE_PULL_SECRET_NAME=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-cluster-registrar-image-tag*)
      export CSI_CLUSTER_REGISTRAR_IMAGE_TAG=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-node-registrar-docker-registry*)
      export CSI_NODE_REGISTRAR_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-node-registrar-image-pull-secret*)
      export CSI_NODE_REGISTRAR_IMAGE_PULL_SECRET_NAME=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-node-registrar-image-tag*)
      export CSI_NODE_REGISTRAR_IMAGE_TAG=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-driver-name*)
      export CSI_VAULT_DRIVER_NAME=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-required-attachment*)
      export CSI_REQUIRED_ATTACHMENT=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --install-appbinding*)
      export REQUIRED_APPBINDING_INSTALL=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --uninstall)
      export CSI_VAULT_UNINSTALL=1
      shift
      ;;
    --purge)
      export CSI_VAULT_PURGE=1
      shift
      ;;
    --monitoring-agent*)
       val=$(echo $1 | sed -e 's/^[^=]*=//g')
       if [ "$val" != "$MONITORING_AGENT_BUILTIN" ] && [ "$val" != "$MONITORING_AGENT_COREOS_OPERATOR" ]; then
         echo 'Invalid monitoring agent. Use "builtin" or "coreos-operator"'
         exit 1
       else
         export MONITORING_AGENT="$val"
       fi
       shift
       ;;
     --monitor-controller-plugin*)
       val=$(echo $1 | sed -e 's/^[^=]*=//g')
       if [ "$val" = "true" ]; then
         export MONITOR_CONTROLLER_PLUGIN="$val"
       fi
       shift
       ;;
     --monitor-node-plugin*)
       val=$(echo $1 | sed -e 's/^[^=]*=//g')
       if [ "$val" = "true" ]; then
         export MONITOR_NODE_PLUGIN="$val"
       fi
       shift
       ;;
     --prometheus-namespace*)
       export PROMETHEUS_NAMESPACE=$(echo $1 | sed -e 's/^[^=]*=//g')
       shift
       ;;
     --servicemonitor-label*)
       label=$(echo $1 | sed -e 's/^[^=]*=//g')
       # split label into key value pair
       IFS='='
       pair=($label)
       unset IFS
       # check if the label is valid
       if [ ! ${#pair[@]} = 2 ]; then
         echo "Invalid ServiceMonitor label format. Use '--servicemonitor-label=key=value'"
         exit 1
       fi
       export SERVICE_MONITOR_LABEL_KEY="${pair[0]}"
       export SERVICE_MONITOR_LABEL_VALUE="${pair[1]}"
       shift
       ;;
    *)
      echo "Error: unknown flag:" $1
      show_help
      exit 1
      ;;
  esac
done

export CSI_VAULT_IMAGE_PULL_SECRET=
if [ -n "$CSI_VAULT_IMAGE_PULL_SECRET_NAME" ]; then
  export CSI_VAULT_IMAGE_PULL_SECRET="name: '$CSI_VAULT_IMAGE_PULL_SECRET_NAME'"
fi

export CSI_ATTACHER_IMAGE_PULL_SECRET=
if [ -n "$CSI_ATTACHER_IMAGE_PULL_SECRET_NAME" ]; then
  export CSI_ATTACHER_IMAGE_PULL_SECRET="name: '$CSI_ATTACHER_IMAGE_PULL_SECRET_NAME'"
fi

export CSI_PROVISIONER_IMAGE_PULL_SECRET=
if [ -n "$CSI_PROVISIONER_IMAGE_PULL_SECRET_NAME" ]; then
  export CSI_PROVISIONER_IMAGE_PULL_SECRET="name: '$CSI_PROVISIONER_IMAGE_PULL_SECRET_NAME'"
fi

export CSI_CLUSTER_REGISTRAR_IMAGE_PULL_SECRET=
if [ -n "$CSI_CLUSTER_REGISTRAR_IMAGE_PULL_SECRET_NAME" ]; then
  export CSI_CLUSTER_REGISTRAR_IMAGE_PULL_SECRET="name: '$CSI_CLUSTER_REGISTRAR_IMAGE_PULL_SECRET_NAME'"
fi

export CSI_NODE_REGISTRAR_IMAGE_PULL_SECRET=
if [ -n "$CSI_NODE_REGISTRAR_IMAGE_PULL_SECRET_NAME" ]; then
  export CSI_NODE_REGISTRAR_IMAGE_PULL_SECRET="name: '$CSI_NODE_REGISTRAR_IMAGE_PULL_SECRET_NAME'"
fi

export PROMETHEUS_NAMESPACE=${PROMETHEUS_NAMESPACE:-$CSI_VAULT_NAMESPACE}

if [ "$CSI_VAULT_NAMESPACE" != "kube-system" ]; then
    export CSI_VAULT_PRIORITY_CLASS=""
fi

if [ "$CSI_VAULT_UNINSTALL" -eq 1 ]; then
  # delete monitoring resources. ignore error as they might not exist
  kubectl delete servicemonitor csi-vault-controller-servicemonitor --namespace $PROMETHEUS_NAMESPACE || true
  kubectl delete servicemonitor csi-vault-node-servicemonitor --namespace $PROMETHEUS_NAMESPACE || true
  kubectl delete secret csi-vault-apiserver-cert --namespace $PROMETHEUS_NAMESPACE || true

   if [ "$CSI_VAULT_PURGE" -eq 1 ]; then
      kubectl delete csidrivers.csi.storage.k8s.io ${CSI_VAULT_DRIVER_NAME} --ignore-not-found=true

      for crd in "${crds[@]}"; do
        kubectl delete crd ${crd} --ignore-not-found=true
      done
   fi

   ${SCRIPT_LOCATION}hack/deploy/controller-plugin.yaml | $ONESSL envsubst  | kubectl delete -f -
   ${SCRIPT_LOCATION}hack/deploy/node-plugin.yaml | $ONESSL envsubst  | kubectl delete -f -
   ${SCRIPT_LOCATION}hack/deploy/csi-driver.yaml | $ONESSL envsubst  | kubectl delete -f -

  echo
  echo "Successfully uninstalled Vault  CSI driver!"
  exit 0
fi

env | sort | grep CSI*
echo ""

if [ "$REQUIRED_APPBINDING_INSTALL" = true ]; then
  ${SCRIPT_LOCATION}hack/deploy/appbinding-crd.yaml | $ONESSL envsubst  | kubectl apply -f -

  echo "waiting until AppBinding crd is ready"
  crd=appbindings.appcatalog.appscode.com
  $ONESSL wait-until-ready crd "${crd}" || {
    echo "$crd crd failed to be ready"
    exit 1
  }
fi

# create necessary TLS certificates:
# - a local CA key and cert
# - a webhook server key and cert signed by the local CA
$ONESSL create ca-cert
$ONESSL create server-cert server --domains=csi-vault-controller.$CSI_VAULT_NAMESPACE,csi-vault-controller.$CSI_VAULT_NAMESPACE.svc,csi-vault-node.$CSI_VAULT_NAMESPACE,csi-vault-node.$CSI_VAULT_NAMESPACE.svc
export SERVICE_SERVING_CERT_CA=$(cat ca.crt | $ONESSL base64)
export TLS_SERVING_CERT=$(cat server.crt | $ONESSL base64)
export TLS_SERVING_KEY=$(cat server.key | $ONESSL base64)

${SCRIPT_LOCATION}hack/deploy/appcatalog-user-roles.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
${SCRIPT_LOCATION}hack/deploy/apiserver-cert.yaml | $ONESSL envsubst | kubectl apply -f -
${SCRIPT_LOCATION}hack/deploy/controller-plugin.yaml | $ONESSL envsubst  | kubectl apply -f -
${SCRIPT_LOCATION}hack/deploy/node-plugin.yaml | $ONESSL envsubst  | kubectl apply -f -
${SCRIPT_LOCATION}hack/deploy/csi-driver.yaml | $ONESSL envsubst  | kubectl apply -f -

# configure prometheus monitoring
if [ "$MONITORING_AGENT" != "$MONITORING_AGENT_NONE" ]; then
  # if operator monitoring is enabled and prometheus-namespace is provided,
  # create csi-vault-apiserver-cert there. this will be mounted on prometheus pod.
  if [ "$PROMETHEUS_NAMESPACE" != "$CSI_VAULT_NAMESPACE" ]; then
    ${SCRIPT_LOCATION}hack/deploy/monitor/apiserver-cert.yaml | $ONESSL envsubst | kubectl apply -f -
  fi

  case "$MONITORING_AGENT" in
    "$MONITORING_AGENT_BUILTIN")
      if [ "$MONITOR_CONTROLLER_PLUGIN" = "true" ]; then
        kubectl annotate service csi-vault-controller -n "$CSI_VAULT_NAMESPACE" --overwrite \
          prometheus.io/scrape="true" \
          prometheus.io/path="/metrics" \
          prometheus.io/port="8443" \
          prometheus.io/scheme="https"
      fi
      if [ "$MONITOR_NODE_PLUGIN" = "true" ]; then
        kubectl annotate service csi-vault-node -n "$CSI_VAULT_NAMESPACE" --overwrite \
          prometheus.io/scrape="true" \
          prometheus.io/path="/metrics" \
          prometheus.io/port="8443" \
          prometheus.io/scheme="https"
      fi
      ;;
    "$MONITORING_AGENT_COREOS_OPERATOR")
      if [ "$MONITOR_CONTROLLER_PLUGIN" = "true" ]; then
        ${SCRIPT_LOCATION}hack/deploy/monitor/servicemonitor-controller.yaml | $ONESSL envsubst | kubectl apply -f -
      fi
      if [ "$MONITOR_NODE_PLUGIN" = "true" ]; then
        ${SCRIPT_LOCATION}hack/deploy/monitor/servicemonitor-node.yaml | $ONESSL envsubst | kubectl apply -f -
      fi
      ;;
  esac
fi

echo
echo "Successfully installed Vault CSI driver in $CSI_VAULT_NAMESPACE namespace!"
