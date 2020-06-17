/*
Copyright 2019 LitmusChaos Authors

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

package resource

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	"github.com/litmuschaos/chaos-operator/pkg/client/clientset/versioned/scheme"
	chaosClient "github.com/litmuschaos/chaos-operator/pkg/client/clientset/versioned/typed/litmuschaos/v1alpha1"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig string
	config     *restclient.Config
	clients    *kubernetes.Clientset
	clientSet  *chaosClient.LitmuschaosV1alpha1Client
)

func init() {
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)

	clients, _ = kubernetes.NewForConfig(config)

	clientSet, _ = chaosClient.NewForConfig(config)

	v1alpha1.AddToScheme(scheme.Scheme)

	// create chaosengine crds
	exec.Command("kubectl", "apply", "-f", "../../deploy/crds/chaosengine_crd.yaml").Run()

	// create sample nginx application
	deployment := &appv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
			Labels: map[string]string{
				"app": "nginx",
			},
			Annotations: map[string]string{
				"litmuschaos.io/chaos": "true",
			},
		},
		Spec: appv1.DeploymentSpec{
			Replicas: func(i int32) *int32 { return &i }(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
							Ports: []v1.ContainerPort{
								{

									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := clients.AppsV1().Deployments("default").Create(deployment)
	if err != nil {
		klog.Infoln("Deployment is not created and error is ", err)
	}
}


func TestValidateAnnontatedApplication(t *testing.T) {
	var labels map[string]string
	labels = make(map[string]string)
	labels["app"] = "nginx"

	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Postive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "validate-annotation-p1",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring:      true,
						AnnotationCheck: "true",
						EngineState:     "active",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx",
							AppKind:  "deployment",
						},
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
						Experiments: []litmuschaosv1alpha1.ExperimentList{
							{
								Name: "exp-1",
							},
						},
					},
				},
				AppInfo: &chaosTypes.ApplicationInfo{
					Kind: "deployment",
					Label: labels,
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: false,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "validate-annotation-p2",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						EngineState:         "active",
						AnnotationCheck:     "false",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx",
							AppKind:  "deployment",
						},
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
						Experiments: []litmuschaosv1alpha1.ExperimentList{
							{
								Name: "exp-1",
							},
						},
					},
				},
				AppInfo: &chaosTypes.ApplicationInfo{
					Kind: "deployment",
					Label: labels,
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Negative-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "validate-annotation-p3",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx",
							AppKind:  "deployment",
						},
						EngineState:     "active",
						AnnotationCheck: "false",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
						Experiments: []litmuschaosv1alpha1.ExperimentList{
							{
								Name: "exp-1",
							},
						},
					},
				},
				AppInfo: &chaosTypes.ApplicationInfo{
					Kind: "deployment",
					Label: labels,
				},
			},
			isErr: false,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := clientSet.ChaosEngines(mock.engine.Instance.Namespace).Create(mock.engine.Instance)
			if err != nil {
				fmt.Printf("engine not created, err: %v", err)
			}

			engine, err := CheckChaosAnnotation(&mock.engine)
			if mock.isErr && err == nil && engine != nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil && engine != nil {
				fmt.Println(err)
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}
