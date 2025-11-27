module github.com/litmuschaos/chaos-operator

go 1.22

require (
	github.com/go-logr/logr v1.4.2
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/jpillora/go-ogle-analytics v0.0.0-20161213085824-14b04e0594ef
	github.com/litmuschaos/elves v0.0.0-20230607095010-c7119636b529
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5 // indirect
	go.opentelemetry.io/otel v1.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.27.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.27.0
	go.opentelemetry.io/otel/sdk v1.27.0
	golang.org/x/oauth2 v0.20.0 // indirect
	k8s.io/api v0.26.15
	k8s.io/apimachinery v0.26.15
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.14.6
)

require (
	github.com/AdaLogics/go-fuzz-headers v0.0.0-20230811130428-ced1acdcaa24
	github.com/google/martian v2.1.0+incompatible
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.24.2
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/otel/trace v1.27.0
	k8s.io/klog v1.0.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/evanphx/json-patch/v5 v5.6.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-logr/zapr v1.2.3 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.20.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	go.opentelemetry.io/otel/metric v1.27.0 // indirect
	go.opentelemetry.io/proto/otlp v1.2.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/term v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240520151616-dc85e6b867a5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240515191416-fc5f0ca64291 // indirect
	google.golang.org/grpc v1.64.0 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.26.1 // indirect
	k8s.io/component-base v0.26.15 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

// Pinned to kubernetes-1.26
replace (
	github.com/go-logr/logr => github.com/go-logr/logr v1.4.2
	k8s.io/api => k8s.io/api v0.26.15
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.26.15
	k8s.io/apimachinery => k8s.io/apimachinery v0.26.15
	k8s.io/client-go => k8s.io/client-go v0.26.15
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.26.15
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.26.15
	k8s.io/component-base => k8s.io/component-base v0.26.15
	k8s.io/cri-api => k8s.io/cri-api v0.26.15
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.26.15
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.80.1
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.26.15
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.26.15
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.26.15
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.26.15
	k8s.io/kubelet => k8s.io/kubelet v0.26.15
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.26.15
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.26.15
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.14.6
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309 // Required by Helm

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
