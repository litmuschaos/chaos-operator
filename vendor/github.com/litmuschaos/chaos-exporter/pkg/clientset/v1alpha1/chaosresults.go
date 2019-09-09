package v1alpha1

import (
	"github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ChaosResultInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ChaosResultList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.ChaosResult, error)
	Create(*v1alpha1.ChaosResult) (*v1alpha1.ChaosResult, error)
	// Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type chaosResultClient struct {
	restClient rest.Interface
	ns         string
}

func (c *chaosResultClient) List(opts metav1.ListOptions) (*v1alpha1.ChaosResultList, error) {
	result := v1alpha1.ChaosResultList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("chaosresults").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *chaosResultClient) Get(name string, opts metav1.GetOptions) (*v1alpha1.ChaosResult, error) {
	result := v1alpha1.ChaosResult{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("chaosresults").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *chaosResultClient) Create(chaosresult *v1alpha1.ChaosResult) (*v1alpha1.ChaosResult, error) {
	result := v1alpha1.ChaosResult{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("chaosresults").
		Body(chaosresult).
		Do().
		Into(&result)

	return &result, err
}

/*
func (c *projectClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
*/
