module github.com/litmuschaos/chaos-operator

go 1.20

require (
	github.com/go-logr/logr v1.2.3
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/jpillora/go-ogle-analytics v0.0.0-20161213085824-14b04e0594ef
	github.com/litmuschaos/elves v0.0.0-20230607095010-c7119636b529
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	k8s.io/api v0.26.2
	k8s.io/apimachinery v0.26.2
	k8s.io/client-go v0.26.2
	sigs.k8s.io/controller-runtime v0.14.5
)

require (
	github.com/google/martian v2.1.0+incompatible
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.24.2
	github.com/stretchr/testify v1.8.2
	k8s.io/klog v1.0.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/zapr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo/v2 v2.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/term v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.29.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.26.2 // indirect
	k8s.io/component-base v0.26.2 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/utils v0.0.0-20230711102312-30195339c3c7 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

// Pinned to kubernetes-1.21.2
//replace (
//	github.com/go-logr/logr => github.com/go-logr/logr v0.4.0
//	k8s.io/api => k8s.io/api v0.21.2
//	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.2
//	k8s.io/apimachinery => k8s.io/apimachinery v0.21.2
//	k8s.io/apiserver => k8s.io/apiserver v0.21.2
//	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.2
//	k8s.io/client-go => k8s.io/client-go v0.21.2
//	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.2
//	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.2
//	k8s.io/code-generator => k8s.io/code-generator v0.21.2
//	k8s.io/component-base => k8s.io/component-base v0.21.2
//	k8s.io/cri-api => k8s.io/cri-api v0.21.2
//	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.2
//	k8s.io/klog/v2 => k8s.io/klog/v2 v2.9.0
//	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.2
//	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.2
//	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.2
//	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.2
//	k8s.io/kubectl => k8s.io/kubectl v0.21.2
//	k8s.io/kubelet => k8s.io/kubelet v0.21.2
//	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.2
//	k8s.io/metrics => k8s.io/metrics v0.21.2
//	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.2
//)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309 // Required by Helm

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible

replace golang.org/x/net => golang.org/x/net v0.7.0
