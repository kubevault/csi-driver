module github.com/kubevault/csi-driver

go 1.12

require (
	github.com/appscode/go v0.0.0-20190424183524-60025f1135c9
	github.com/appscode/pat v0.0.0-20170521084856-48ff78925b79
	github.com/container-storage-interface/spec v1.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.1
	github.com/gophercloud/gophercloud v0.0.0-20190509032623-7892efa714f1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/hashicorp/vault v1.0.1
	github.com/juju/errors v0.0.0-20190207033735-e65537c515d7
	github.com/juju/loggo v0.0.0-20190212223446-d976af380377 // indirect
	github.com/juju/testing v0.0.0-20190429233213-dfc56b8c09fc // indirect
	github.com/kubedb/apimachinery v0.0.0-20190508221312-5ba915343400
	github.com/kubernetes-csi/csi-lib-utils v0.0.0-20190415202911-789e4ed466cf // indirect
	github.com/kubernetes-csi/csi-test v2.0.0+incompatible
	github.com/kubernetes-csi/livenessprobe v1.0.1
	github.com/kubevault/operator v0.0.0-20190509030635-7f32eefb5188
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	google.golang.org/grpc v1.20.1
	k8s.io/api v0.0.0-20190503110853-61630f889b3c
	k8s.io/apiextensions-apiserver v0.0.0-20190509064156-0d7f274f68cb // indirect
	k8s.io/apimachinery v0.0.0-20190509063443-7d8f8feb49c5
	k8s.io/apiserver v0.0.0-20190509063909-3b296809833b
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/component-base v0.0.0-20190509023737-8de8845fb642 // indirect
	k8s.io/kubernetes v1.14.1
	kmodules.xyz/client-go v0.0.0-20190508091620-0d215c04352f
	kmodules.xyz/custom-resources v0.0.0-20190508103408-464e8324c3ec
)

replace (
	github.com/graymeta/stow => github.com/appscode/stow v0.0.0-20190506085026-ca5baa008ea3
	github.com/kubernetes-csi/csi-test => github.com/kubevault/csi-test v2.0.0+incompatible
	gopkg.in/robfig/cron.v2 => github.com/appscode/cron v0.0.0-20170717094345-ca60c6d796d4
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
