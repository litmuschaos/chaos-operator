/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ChaosExperimentLister helps list ChaosExperiments.
// All objects returned here must be treated as read-only.
type ChaosExperimentLister interface {
	// List lists all ChaosExperiments in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ChaosExperiment, err error)
	// ChaosExperiments returns an object that can list and get ChaosExperiments.
	ChaosExperiments(namespace string) ChaosExperimentNamespaceLister
	ChaosExperimentListerExpansion
}

// chaosExperimentLister implements the ChaosExperimentLister interface.
type chaosExperimentLister struct {
	indexer cache.Indexer
}

// NewChaosExperimentLister returns a new ChaosExperimentLister.
func NewChaosExperimentLister(indexer cache.Indexer) ChaosExperimentLister {
	return &chaosExperimentLister{indexer: indexer}
}

// List lists all ChaosExperiments in the indexer.
func (s *chaosExperimentLister) List(selector labels.Selector) (ret []*v1alpha1.ChaosExperiment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ChaosExperiment))
	})
	return ret, err
}

// ChaosExperiments returns an object that can list and get ChaosExperiments.
func (s *chaosExperimentLister) ChaosExperiments(namespace string) ChaosExperimentNamespaceLister {
	return chaosExperimentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ChaosExperimentNamespaceLister helps list and get ChaosExperiments.
// All objects returned here must be treated as read-only.
type ChaosExperimentNamespaceLister interface {
	// List lists all ChaosExperiments in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ChaosExperiment, err error)
	// Get retrieves the ChaosExperiment from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ChaosExperiment, error)
	ChaosExperimentNamespaceListerExpansion
}

// chaosExperimentNamespaceLister implements the ChaosExperimentNamespaceLister
// interface.
type chaosExperimentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ChaosExperiments in the indexer for a given namespace.
func (s chaosExperimentNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ChaosExperiment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ChaosExperiment))
	})
	return ret, err
}

// Get retrieves the ChaosExperiment from the indexer for a given namespace and name.
func (s chaosExperimentNamespaceLister) Get(name string) (*v1alpha1.ChaosExperiment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("chaosexperiment"), name)
	}
	return obj.(*v1alpha1.ChaosExperiment), nil
}
