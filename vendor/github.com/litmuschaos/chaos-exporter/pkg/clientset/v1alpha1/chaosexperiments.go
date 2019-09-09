package v1alpha1

import (
	"github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ChaosExperimentInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ChaosExperimentList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.ChaosExperiment, error)
	Create(*v1alpha1.ChaosExperiment) (*v1alpha1.ChaosExperiment, error)
	// ...
}

type chaosExperimentClient struct {
	restClient rest.Interface
	ns         string
}

func (c *chaosExperimentClient) List(opts metav1.ListOptions) (*v1alpha1.ChaosExperimentList, error) {
	result := v1alpha1.ChaosExperimentList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("chaosexperiments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *chaosExperimentClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.ChaosExperiment, error) {
	result := v1alpha1.ChaosExperiment{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("chaosexpriments").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *chaosExperimentClient) Create(chaosexperiment *v1alpha1.ChaosExperiment) (*v1alpha1.ChaosExperiment, error) {
	result := v1alpha1.ChaosExperiment{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("chaosexperiments").
		Body(chaosexperiment).
		Do().
		Into(&result)

	return &result, err
}
