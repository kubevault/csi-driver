module kubevault.dev/csi-driver

go 1.12

require (
	github.com/appscode/go v0.0.0-20191119085241-0887d8ec2ecc
	github.com/appscode/pat v0.0.0-20170521084856-48ff78925b79
	github.com/container-storage-interface/spec v1.1.0
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/vault/api v1.0.5-0.20191108163347-bdd38fca2cff
	github.com/juju/errors v0.0.0-20190930114154-d42613fe1ab9
	github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8 // indirect
	github.com/juju/testing v0.0.0-20191001232224-ce9dec17d28b // indirect
	github.com/kubernetes-csi/csi-lib-utils v0.6.1
	github.com/kubernetes-csi/csi-test v2.2.0+incompatible
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	google.golang.org/grpc v1.23.0
	k8s.io/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/apiserver v0.0.0-20191114103151-9ca1dc586682
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v0.4.0
	k8s.io/kubernetes v1.16.3
	kmodules.xyz/client-go v0.0.0-20191127054604-26981530831d
	kmodules.xyz/custom-resources v0.0.0-20191130062942-f41b54f62419
	kmodules.xyz/offshoot-api v0.0.0-20191130175124-5676c4869928 // indirect
	kubevault.dev/operator v0.0.0-20191130190210-a4a972e3ecdd
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.38.0
	git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
	github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v32.5.0+incompatible
	github.com/Azure/go-ansiterm => github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible
	github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.9.0
	github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.5.0
	github.com/Azure/go-autorest/autorest/azure/auth => github.com/Azure/go-autorest/autorest/azure/auth v0.2.0
	github.com/Azure/go-autorest/autorest/date => github.com/Azure/go-autorest/autorest/date v0.1.0
	github.com/Azure/go-autorest/autorest/mocks => github.com/Azure/go-autorest/autorest/mocks v0.2.0
	github.com/Azure/go-autorest/autorest/to => github.com/Azure/go-autorest/autorest/to v0.2.0
	github.com/Azure/go-autorest/autorest/validation => github.com/Azure/go-autorest/autorest/validation v0.1.0
	github.com/Azure/go-autorest/logger => github.com/Azure/go-autorest/logger v0.1.0
	github.com/Azure/go-autorest/tracing => github.com/Azure/go-autorest/tracing v0.5.0
	github.com/kubernetes-csi/csi-test => github.com/kubevault/csi-test v1.1.2-0.20191201160333-490f2b7805bb
	github.com/kubernetes-csi/external-snapshotter => github.com/kmodules/external-snapshotter v1.2.1-0.20191128100451-0265c5fa679a
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	google.golang.org/api => google.golang.org/api v0.6.1-0.20190607001116-5213b8090861
	k8s.io/api => k8s.io/api v0.0.0-20191114100352-16d7abae0d2a
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191114105449-027877536833
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20191119091232-0553326db082
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20191119111000-36ac3646ae82
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191114110141-0a35778df828
	k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191114112024-4bbba8331835
	k8s.io/component-base => k8s.io/component-base v0.0.0-20191114102325-35a9586014f7
	k8s.io/klog => k8s.io/klog v0.4.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191114103820-f023614fb9ea
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf
	k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191114113550-6123e1c827f7
	k8s.io/kubernetes => github.com/kmodules/kubernetes v1.17.0-alpha.0.0.20191127022853-9d027e3886fd
	k8s.io/metrics => k8s.io/metrics v0.0.0-20191114105837-a4a2842dc51b
	k8s.io/repo-infra => k8s.io/repo-infra v0.0.0-20181204233714-00fe14e3d1a3
	k8s.io/utils => k8s.io/utils v0.0.0-20190801114015-581e00157fb1
	sigs.k8s.io/kustomize => sigs.k8s.io/kustomize v2.0.3+incompatible
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v0.0.0-20190817042607-6149e4549fca
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.1.0
)
