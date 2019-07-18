module kubevault.dev/csi-driver

go 1.12

require (
	github.com/appscode/go v0.0.0-20190621064509-6b292c9166e3
	github.com/appscode/pat v0.0.0-20170521084856-48ff78925b79
	github.com/container-storage-interface/spec v1.1.0
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/emicklei/go-restful v2.9.6+incompatible // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.1
	github.com/gophercloud/gophercloud v0.0.0-20190509032623-7892efa714f1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-hclog v0.9.2 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hashicorp/vault v1.0.3
	github.com/juju/errors v0.0.0-20190207033735-e65537c515d7
	github.com/juju/loggo v0.0.0-20190212223446-d976af380377 // indirect
	github.com/juju/testing v0.0.0-20190429233213-dfc56b8c09fc // indirect
	github.com/kubernetes-csi/csi-lib-utils v0.0.0-20190415202911-789e4ed466cf // indirect
	github.com/kubernetes-csi/csi-test v2.0.0+incompatible
	github.com/kubernetes-csi/livenessprobe v1.0.1
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.4
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	go.etcd.io/bbolt v1.3.3 // indirect
	google.golang.org/grpc v1.21.1
	k8s.io/api v0.0.0-20190503110853-61630f889b3c
	k8s.io/apimachinery v0.0.0-20190509063443-7d8f8feb49c5
	k8s.io/apiserver v0.0.0-20190516230822-f89599b3f645
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/component-base v0.0.0-20190509023737-8de8845fb642 // indirect
	k8s.io/klog v0.3.2 // indirect
	k8s.io/kubernetes v1.14.3
	kmodules.xyz/client-go v0.0.0-20190715080709-7162a6c90b04
	kmodules.xyz/custom-resources v0.0.0-20190508103408-464e8324c3ec
	kubedb.dev/apimachinery v0.0.0-20190718012505-6973f90fefd0
	kubevault.dev/operator v0.0.0-20190718164535-ca947abcfc3f
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest/autorest v0.5.0
	github.com/kubernetes-csi/csi-test => github.com/kubevault/csi-test v2.0.1-0.20190520112505-d88b429c7476+incompatible
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190508082252-8397d761d4b5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190221042446-c2654d5206da
)
