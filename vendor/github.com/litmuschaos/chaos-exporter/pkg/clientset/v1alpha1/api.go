package v1alpha1

import (
	"github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

//ExampleV1Alpha1Interface type defines chaosEngines & chaosResults
type ExampleV1Alpha1Interface interface {
	// ChaosEngines with namespace attribute
	ChaosEngines(namespace string) ChaosEngineInterface
	// ChaosResults with namespace attribute
	ChaosResults(namespace string) ChaosResultInterface
	//
	ChaosExperiments(namespace string) ChaosExperimentInterface
}

//ExampleV1Alpha1Client type defines the rest client for chaos resources
type ExampleV1Alpha1Client struct {
	restClient rest.Interface
}

//NewForConfig returns the kubeclient for the config provided
func NewForConfig(c *rest.Config) (*ExampleV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleV1Alpha1Client{restClient: client}, nil
}

func (c *ExampleV1Alpha1Client) ChaosEngines(namespace string) ChaosEngineInterface {
	return &chaosEngineClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
func (c *ExampleV1Alpha1Client) ChaosExperiments(namespace string) ChaosExperimentInterface {
	return &chaosExperimentClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *ExampleV1Alpha1Client) ChaosResults(namespace string) ChaosResultInterface {
	return &chaosResultClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
