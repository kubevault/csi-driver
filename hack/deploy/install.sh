#!/bin/bash
set -eou pipefail

crds=(
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

function cleanup() {
  rm -rf $ONESSL ca.crt ca.key server.crt server.key
}

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
trap cleanup EXIT

# ref: https://github.com/appscodelabs/libbuild/blob/master/common/lib.sh#L55
inside_git_repo() {
  git rev-parse --is-inside-work-tree >/dev/null 2>&1
  inside_git=$?
  if [ "$inside_git" -ne 0 ]; then
    echo "Not inside a git repository"
    exit 1
  fi
}

detect_tag() {
  inside_git_repo

  # http://stackoverflow.com/a/1404862/3476121
  git_tag=$(git describe --exact-match --abbrev=0 2>/dev/null || echo '')

  commit_hash=$(git rev-parse --verify HEAD)
  git_branch=$(git rev-parse --abbrev-ref HEAD)
  commit_timestamp=$(git show -s --format=%ct)

  if [ "$git_tag" != '' ]; then
    TAG=$git_tag
    TAG_STRATEGY='git_tag'
  elif [ "$git_branch" != 'master' ] && [ "$git_branch" != 'HEAD' ] && [[ "$git_branch" != release-* ]]; then
    TAG=$git_branch
    TAG_STRATEGY='git_branch'
  else
    hash_ver=$(git describe --tags --always --dirty)
    TAG="${hash_ver}"
    TAG_STRATEGY='commit_hash'
  fi

  export TAG
  export TAG_STRATEGY
  export git_tag
  export git_branch
  export commit_hash
  export commit_timestamp
}

onessl_found() {
  # https://stackoverflow.com/a/677212/244009
  if [ -x "$(command -v onessl)" ]; then
    onessl wait-until-has -h >/dev/null 2>&1 || {
      # old version of onessl found
      echo "Found outdated onessl"
      return 1
    }
    export ONESSL=onessl
    return 0
  fi
  return 1
}

onessl_found || {
  echo "Downloading onessl ..."
  # ref: https://stackoverflow.com/a/27776822/244009
  case "$(uname -s)" in
    Darwin)
      curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.9.0/onessl-darwin-amd64
      chmod +x onessl
      export ONESSL=./onessl
      ;;

    Linux)
      curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.9.0/onessl-linux-amd64
      chmod +x onessl
      export ONESSL=./onessl
      ;;

    CYGWIN* | MINGW32* | MSYS*)
      curl -fsSL -o onessl.exe https://github.com/kubepack/onessl/releases/download/0.9.0/onessl-windows-amd64.exe
      chmod +x onessl.exe
      export ONESSL=./onessl.exe
      ;;
    *)
      echo 'other OS'
      ;;
  esac
}


# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export VAULT_CSI_DRIVER_NAMESPACE=kube-system
export CSI_VAULT_DOCKER_REGISTRY=kubevault
export CSI_VAULT_DOCKER_REPOSITORY=csi-vault
export CSI_VAULT_IMAGE_TAG=0.1.0
export CSI_VAULT_IMAGE_PULL_SECRET=
export CSI_VAULT_IMAGE_PULL_POLICY=IfNotPresent
export CSI_ATTACHER_DOCKER_REGISTRY=quay.io/k8scsi
export CSI_ATTACHER_DOCKER_REPOSITORY=csi-attacher
export CSI_ATTACHER_IMAGE_TAG=v0.2.0
export CSI_ATTACHER_IMAGE_PULL_SECRET=
export CSI_ATTACHER_IMAGE_PULL_POLICY=IfNotPresent
export CSI_PROVISIONER_DOCKER_REGISTRY=quay.io/k8scsi
export CSI_PROVISIONER_DOCKER_REPOSITORY=csi-provisioner
export CSI_PROVISIONER_IMAGE_TAG=v0.2.1
export CSI_PROVISIONER_IMAGE_PULL_SECRET=
export CSI_PROVISIONER_IMAGE_PULL_POLICY=IfNotPresent
export CSI_REGISTRAR_DOCKER_REGISTRY=quay.io/k8scsi
export CSI_REGISTRAR_DOCKER_REPOSITORY=driver-registrar
export CSI_REGISTRAR_IMAGE_TAG=v0.3.0
export CSI_REGISTRAR_IMAGE_PULL_SECRET=
export CSI_REGISTRAR_IMAGE_PULL_POLICY=IfNotPresent
export VAULT_CSI_DRIVER_NAME=com.kubevault.csi.secrets
export VAULT_CSI_DRIVER_UNINSTALL=0
export VAULT_CSI_DRIVER_PURGE=0
export CSI_REQUIRED_ATTACHMENT=false
export REQUIRED_APPBINDING_INSTALL=true

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
export SCRIPT_LOCATION="curl -fsSL https://raw.githubusercontent.com/kubevault/csi-driver/0.1.0/"
if [ "$APPSCODE_ENV" = "dev" ]; then
  detect_tag
  export SCRIPT_LOCATION="cat "
  export CSI_VAULT_IMAGE_TAG=$TAG
  export CSI_VAULT_IMAGE_PULL_POLICY=Always
fi

show_help() {
  echo "install.sh -install vault csi driver"
  echo " "
  echo "install.sh [options]"
  echo " "
  echo "options:"
  echo "-h, --help                                show brief help"
  echo "-n, --namespace=NAMESPACE                 specify namespace (default: kube-system)"
  echo "    --csi-vault-docker-registry           docker registry used to pull csi-vault image (default: kubevault)"
  echo "    --csi-vault-image-pull-secret         name of secret used to pull csi-vault images"
  echo "    --csi-vault-image-tag                 docker image version of csi vault"
  echo "    --csi-attacher-docker-registry        docker registry used to pull csi attacher image (default: quay.io/k8scsi)"
  echo "    --csi-attacher-image-pull-secret      name of secret used to pull csi attacher image"
  echo "    --csi-attacher-image-tag              docker image version of csi attacher"
  echo "    --csi-provisioner-docker-registry     docker registry used to pull csi provisioner image (default: quay.io/k8scsi)"
  echo "    --csi-provisioner-image-pull-secret   name of secret used to pull csi provisioner image"
  echo "    --csi-provisioner-image-tag           docker image version of csi provisioner"
  echo "    --csi-registrar-docker-registry       docker registry used to pull csi registrar image (default: quay.io/k8scsi)"
  echo "    --csi-registrar-image-pull-secret     name of secret used to pull csi registrar image"
  echo "    --csi-registrar-image-tag             docker image version of csi registrar"
  echo "    --csi-driver-name                     name of csi driver to install (default: com.kubevault.csi.secrets)"
  echo "    --csi-required-attachment             indicates csi volume driver requires an attach operation (default: false)"
  echo "    --install-appbinding                  indicates appbinding crd need to be installed (default: true)"
  echo "    --uninstall                           uninstall vault csi driver"
  echo "    --purge                               purges csi driver crd objects and crds"
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
        export VAULT_CSI_DRIVER_NAMESPACE=$1
      else
        echo "no namespace specified"
        exit 1
      fi
      shift
      ;;
    --namespace*)
      export VAULT_CSI_DRIVER_NAMESPACE=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-vault-docker-registry*)
      export CSI_VAULT_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-vault-image-pull-secret*)
      secret=$(echo $1 | sed -e 's/^[^=]*=//g')
      export CSI_VAULT_IMAGE_PULL_SECRET="name: '$secret'"
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
      secret=$(echo $1 | sed -e 's/^[^=]*=//g')
      export CSI_ATTACHER_IMAGE_PULL_SECRET="name: '$secret'"
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
      secret=$(echo $1 | sed -e 's/^[^=]*=//g')
      export CSI_PROVISIONER_IMAGE_PULL_SECRET="name: '$secret'"
      shift
      ;;
    --csi-provisioner-image-tag*)
      export CSI_PROVISIONER_IMAGE_TAG=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-registrar-docker-registry*)
      export CSI_REGISTRAR_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-registrar-image-pull-secret*)
      secret=$(echo $1 | sed -e 's/^[^=]*=//g')
      export CSI_REGISTRAR_IMAGE_PULL_SECRET="name: '$secret'"
      shift
      ;;
    --csi-registrar-image-tag*)
      export CSI_REGISTRAR_IMAGE_TAG=$(echo $1 | sed -e 's/^[^=]*=//g')
      shift
      ;;
    --csi-driver-name*)
      export VAULT_CSI_DRIVER_NAME=$(echo $1 | sed -e 's/^[^=]*=//g')
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
      export VAULT_CSI_DRIVER_UNINSTALL=1
      shift
      ;;
    --purge)
      export VAULT_CSI_DRIVER_PURGE=1
      shift
      ;;
    *)
      echo "Error: unknown flag:" $1
      show_help
      exit 1
      ;;
  esac
done

if [ "$VAULT_CSI_DRIVER_UNINSTALL" -eq 1 ]; then
   kubectl delete csidrivers.csi.storage.k8s.io ${VAULT_CSI_DRIVER_NAME}

   ${SCRIPT_LOCATION}hack/deploy/csi-attacher.yaml | $ONESSL envsubst  | kubectl delete -f -
   ${SCRIPT_LOCATION}hack/deploy/csi-provisioner.yaml | $ONESSL envsubst  | kubectl delete -f -
   ${SCRIPT_LOCATION}hack/deploy/csi-driver-registrar.yaml | $ONESSL envsubst  | kubectl delete -f -

  if [ "$VAULT_CSI_DRIVER_PURGE" -eq 1 ]; then
    for crd in "${crds[@]}"; do
      kubectl delete crd ${crd} --ignore-not-found=true
    done
    if [ "$REQUIRED_APPBINDING_INSTALL" = true ]; then
        kubectl delete -f https://raw.githubusercontent.com/kmodules/custom-resources/master/api/crds/appbinding.yaml
    fi
  fi

  echo
  echo "Successfully uninstalled Vault  CSI driver!"
  exit 0
fi

${SCRIPT_LOCATION}hack/deploy/driver-crd.yaml | $ONESSL envsubst  | kubectl apply -f -

if [ "$REQUIRED_APPBINDING_INSTALL" = true ]; then
  kubectl apply -f https://raw.githubusercontent.com/kmodules/custom-resources/master/api/crds/appbinding.yaml
fi

${SCRIPT_LOCATION}hack/deploy/csi-attacher.yaml | $ONESSL envsubst  | kubectl apply -f -
${SCRIPT_LOCATION}hack/deploy/csi-provisioner.yaml | $ONESSL envsubst  | kubectl apply -f -
${SCRIPT_LOCATION}hack/deploy/csi-driver-registrar.yaml | $ONESSL envsubst  | kubectl apply -f -

echo
echo "Successfully installed Vault CSI driver in $VAULT_CSI_DRIVER_NAMESPACE namespace!"