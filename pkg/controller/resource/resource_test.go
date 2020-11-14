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
	"testing"

	litmusFakeClientset "github.com/litmuschaos/chaos-operator/pkg/client/clientset/versioned/fake"
	appv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"

	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/controller/types"
)

var (
	//fake deploymentconfig
	gvfakedc = schema.GroupVersionResource{
		Group:    "apps.openshift.io",
		Version:  "v1",
		Resource: "deploymentconfigs",
	}
	//fake rollout
	gvfakero = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "rollouts",
	}

)

func print32(p int32) *int32 {
	return &p
}

func TestCheckChaosAnnotationDeployment(t *testing.T) {

	tests := map[string]struct {
		engine     chaosTypes.EngineInfo
		isErr      bool
		deployment []appv1.Deployment
		check      bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d1",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring:      false,
						AnnotationCheck: "true",
						EngineState:     "active",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx1",
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
					Label: map[string]string{
						"app": "nginx1",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			deployment: []appv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx1",
						},
						Annotations: map[string]string{
							"litmuschaos.io/chaos": "true",
						},
					},
					Spec: appv1.DeploymentSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx1",
							},
						},
					},
				},
			},

			isErr: false,
			check: true,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d2",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						EngineState:         "active",
						AnnotationCheck:     "false",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx2",
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
					Label: map[string]string{
						"app": "nginx2",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			deployment: []appv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx2",
						},
						Annotations: map[string]string{
							"litmuschaos.io/chaos": "true",
						},
					},
					Spec: appv1.DeploymentSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx2",
							},
						},
					},
				},
			},

			isErr: false,
			check: true,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d3",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx3",
							AppKind:  "deployment",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Label: map[string]string{
						"app": "nginx3",
					},
				},
			},
			isErr: true,
			check: false,
		},
		"Test Negative-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d4",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx4",
							AppKind:  "deployment",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Label: map[string]string{
						"app": "nginx4",
					},
				},
			},
			deployment: []appv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx4",
						},
					},
					Spec: appv1.DeploymentSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx4",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx1",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx4",
						},
					},
					Spec: appv1.DeploymentSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx4",
							},
						},
					},
				},
			},

			isErr: true,
			check: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {

			f := newFixture(t)
			f.SetFakeClient()
			if mock.check == true {
				for _, deploy := range mock.deployment {
					_, err := f.k8sClient.AppsV1().Deployments(deploy.Namespace).Create(&deploy)
					if err != nil {
						fmt.Printf("deployment not created, err: %v", err)
					}
				}
			}
			_, err := f.litmusClient.LitmuschaosV1alpha1().ChaosEngines(mock.engine.Instance.Namespace).Create(mock.engine.Instance)
			if err != nil {
				fmt.Printf("engine not created, err: %v", err)
			}

			engine, err := CheckChaosAnnotation(&mock.engine, f.k8sClient, f.dynamicClient)
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

func TestCheckChaosAnnotationStatefulSet(t *testing.T) {

	tests := map[string]struct {
		engine      chaosTypes.EngineInfo
		isErr       bool
		statefulSet []appv1.StatefulSet
		check       bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d1",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring:      false,
						AnnotationCheck: "true",
						EngineState:     "active",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx1",
							AppKind:  "statefulset",
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
					Kind: "statefulset",
					Label: map[string]string{
						"app": "nginx1",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			statefulSet: []appv1.StatefulSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx1",
						},
						Annotations: map[string]string{
							"litmuschaos.io/chaos": "true",
						},
					},
					Spec: appv1.StatefulSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx1",
							},
						},
					},
				},
			},

			isErr: false,
			check: true,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d2",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						EngineState:         "active",
						AnnotationCheck:     "false",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx2",
							AppKind:  "statefulset",
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
					Kind: "statefulset",
					Label: map[string]string{
						"app": "nginx2",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			statefulSet: []appv1.StatefulSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx2",
						},
						Annotations: map[string]string{
							"litmuschaos.io/chaos": "true",
						},
					},
					Spec: appv1.StatefulSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx2",
							},
						},
					},
				},
			},

			isErr: false,
			check: true,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d3",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx3",
							AppKind:  "statefulset",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Kind: "statefulset",
					Label: map[string]string{
						"app": "nginx3",
					},
				},
			},
			isErr: true,
			check: false,
		},
		"Test Negative-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d4",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx4",
							AppKind:  "statefulset",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Kind: "statefulset",
					Label: map[string]string{
						"app": "nginx4",
					},
				},
			},
			statefulSet: []appv1.StatefulSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx4",
						},
					},
					Spec: appv1.StatefulSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx4",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx1",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx4",
						},
					},
					Spec: appv1.StatefulSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx4",
							},
						},
					},
				},
			},

			isErr: true,
			check: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {

			f := newFixture(t)
			f.SetFakeClient()
			if mock.check == true {
				for _, sts := range mock.statefulSet {
					_, err := f.k8sClient.AppsV1().StatefulSets(sts.Namespace).Create(&sts)
					if err != nil {
						fmt.Printf("statefulset not created, err: %v", err)
					}
				}
			}
			_, err := f.litmusClient.LitmuschaosV1alpha1().ChaosEngines(mock.engine.Instance.Namespace).Create(mock.engine.Instance)
			if err != nil {
				fmt.Printf("engine not created, err: %v", err)
			}

			engine, err := CheckChaosAnnotation(&mock.engine, f.k8sClient, f.dynamicClient)
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

func TestCheckChaosAnnotationDaemonset(t *testing.T) {

	tests := map[string]struct {
		engine    chaosTypes.EngineInfo
		isErr     bool
		daemonset []appv1.DaemonSet
		check     bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d1",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring:      false,
						AnnotationCheck: "true",
						EngineState:     "active",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx1",
							AppKind:  "daemonset",
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
					Kind: "daemonset",
					Label: map[string]string{
						"app": "nginx1",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			daemonset: []appv1.DaemonSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx1",
						},
						Annotations: map[string]string{
							"litmuschaos.io/chaos": "true",
						},
					},
					Spec: appv1.DaemonSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx1",
							},
						},
					},
				},
			},

			isErr: false,
			check: true,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d2",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						EngineState:         "active",
						AnnotationCheck:     "false",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx2",
							AppKind:  "daemonset",
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
					Kind: "daemonset",
					Label: map[string]string{
						"app": "nginx2",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			daemonset: []appv1.DaemonSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx2",
						},
						Annotations: map[string]string{
							"litmuschaos.io/chaos": "true",
						},
					},
					Spec: appv1.DaemonSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx2",
							},
						},
					},
				},
			},

			isErr: false,
			check: true,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d3",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx3",
							AppKind:  "daemonset",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Kind: "daemonset",
					Label: map[string]string{
						"app": "nginx3",
					},
				},
			},
			isErr: true,
			check: false,
		},
		"Test Negative-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d4",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx4",
							AppKind:  "daemonset",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Kind: "daemonset",
					Label: map[string]string{
						"app": "nginx4",
					},
				},
			},
			daemonset: []appv1.DaemonSet{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx4",
						},
					},
					Spec: appv1.DaemonSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx4",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx1",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx4",
						},
					},
					Spec: appv1.DaemonSetSpec{
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": "nginx4",
							},
						},
					},
				},
			},

			isErr: true,
			check: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {

			f := newFixture(t)
			f.SetFakeClient()
			if mock.check == true {
				for _, ds := range mock.daemonset {
					_, err := f.k8sClient.AppsV1().DaemonSets(ds.Namespace).Create(&ds)
					if err != nil {
						fmt.Printf("daemonset not created, err: %v", err)
					}
				}
			}
			_, err := f.litmusClient.LitmuschaosV1alpha1().ChaosEngines(mock.engine.Instance.Namespace).Create(mock.engine.Instance)
			if err != nil {
				fmt.Printf("engine not created, err: %v", err)
			}

			engine, err := CheckChaosAnnotation(&mock.engine, f.k8sClient, f.dynamicClient)
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

func TestCheckChaosAnnotationDeploymentConfigs(t *testing.T) {

	tests := map[string]struct {
		engine           chaosTypes.EngineInfo
		isErr            bool
		deploymentconfig []unstructured.Unstructured
		check            bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d1",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring:      false,
						AnnotationCheck: "true",
						EngineState:     "active",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx1",
							AppKind:  "deploymentconfig",
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
					Kind: "deploymentconfig",
					Label: map[string]string{
						"app": "nginx1",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			deploymentconfig: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "DeploymentConfig",
						"metadata": map[string]interface{}{
							"name":      "nginx",
							"namespace": "default",
							"labels": map[string]interface{}{
								"app": "nginx1",
							},
							"annotations": map[string]interface{}{
								"litmuschaos.io/chaos": "true",
							},
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"app": "nginx1",
							},
						},
					},
				},
			},
			isErr: false,
			check: true,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d2",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						EngineState:         "active",
						AnnotationCheck:     "false",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx2",
							AppKind:  "deploymentconfig",
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
					Kind: "deploymentconfig",
					Label: map[string]string{
						"app": "nginx2",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			deploymentconfig: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "DeploymentConfig",
						"metadata": map[string]interface{}{
							"name":      "nginx",
							"namespace": "default",
							"labels": map[string]interface{}{
								"app": "nginx2",
							},
							"annotations": map[string]interface{}{
								"litmuschaos.io/chaos": "true",
							},
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"app": "nginx2",
							},
						},
					},
				},
			},

			isErr: false,
			check: true,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d3",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx3",
							AppKind:  "deploymentconfig",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Kind: "deploymentconfig",
					Label: map[string]string{
						"app": "nginx3",
					},
				},
			},
			isErr: true,
			check: false,
		},
		"Test Negative-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d4",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx4",
							AppKind:  "deploymentconfig",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Kind: "deploymentconfig",
					Label: map[string]string{
						"app": "nginx4",
					},
				},
			},
			deploymentconfig: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "DeploymentConfig",
						"metadata": map[string]interface{}{
							"name":      "nginx",
							"namespace": "default",
							"labels": map[string]interface{}{
								"app": "nginx4",
							},
							"annotations": map[string]interface{}{
								"litmuschaos.io/chaos": "true",
							},
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"app": "nginx4",
							},
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "DeploymentConfig",
						"metadata": map[string]interface{}{
							"name":      "nginx1",
							"namespace": "default",
							"labels": map[string]interface{}{
								"app": "nginx4",
							},
							"annotations": map[string]interface{}{
								"litmuschaos.io/chaos": "true",
							},
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"app": "nginx4",
							},
						},
					},
				},
			},

			isErr: true,
			check: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {

			f := newFixture(t)
			f.SetFakeClient()
			dynamicClient := f.dynamicClient.Resource(gvfakedc)
			dynamic := dynamicClient.Namespace("default")

			if mock.check == true {
				for _, dc := range mock.deploymentconfig {
					_, err := dynamic.Create(&dc, metav1.CreateOptions{})
					if err != nil {
						fmt.Printf("deploymentconfig not created, err: %v", err)
					}

				}
			}
			_, err := f.litmusClient.LitmuschaosV1alpha1().ChaosEngines(mock.engine.Instance.Namespace).Create(mock.engine.Instance)
			if err != nil {
				fmt.Printf("engine not created, err: %v", err)
			}

			engine, err := CheckChaosAnnotation(&mock.engine, f.k8sClient, f.dynamicClient)
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


func TestCheckChaosAnnotationRollouts(t *testing.T) {

	tests := map[string]struct {
		engine           chaosTypes.EngineInfo
		isErr            bool
		rollout []unstructured.Unstructured
		check            bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d1",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring:      false,
						AnnotationCheck: "true",
						EngineState:     "active",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx1",
							AppKind:  "rollout",
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
					Kind: "rollout",
					Label: map[string]string{
						"app": "nginx1",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			rollout: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1alpha1",
						"kind":       "Rollout",
						"metadata": map[string]interface{}{
							"name":      "nginx",
							"namespace": "default",
							"labels": map[string]interface{}{
								"app": "nginx1",
							},
							"annotations": map[string]interface{}{
								"litmuschaos.io/chaos": "true",
							},
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"matchLabels": map[string]interface{}{
									"app": "nginx1",
								},
							},
						},
					},
				},
			},
			isErr: false,
			check: true,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d2",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						EngineState:         "active",
						AnnotationCheck:     "false",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx2",
							AppKind:  "rollout",
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
					Kind: "rollout",
					Label: map[string]string{
						"app": "nginx2",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			rollout: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1alpha1",
						"kind":       "Rollout",
						"metadata": map[string]interface{}{
							"name":      "nginx",
							"namespace": "default",
							"labels": map[string]interface{}{
								"app": "nginx2",
							},
							"annotations": map[string]interface{}{
								"litmuschaos.io/chaos": "true",
							},
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"matchLabels": map[string]interface{}{
									"app": "nginx2",
								},
							},
						},
					},
				},
			},

			isErr: false,
			check: true,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d3",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx3",
							AppKind:  "rollout",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Kind: "rollout",
					Label: map[string]string{
						"app": "nginx3",
					},
				},
			},
			isErr: true,
			check: false,
		},
		"Test Negative-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "check-chaos-annotation-d4",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx4",
							AppKind:  "rollout",
						},
						EngineState:     "active",
						AnnotationCheck: "true",
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
					Kind: "rollout",
					Label: map[string]string{
						"app": "nginx4",
					},
				},
			},
			rollout: []unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"apiVersion": "v1alpha1",
						"kind":       "Rollout",
						"metadata": map[string]interface{}{
							"name":      "nginx",
							"namespace": "default",
							"labels": map[string]interface{}{
								"app": "nginx4",
							},
							"annotations": map[string]interface{}{
								"litmuschaos.io/chaos": "true",
							},
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"matchLabels": map[string]interface{}{
									"app": "nginx4",
								},
							},
						},
					},
				},
				{
					Object: map[string]interface{}{
						"apiVersion": "v1alpha1",
						"kind":       "Rollout",
						"metadata": map[string]interface{}{
							"name":      "nginx1",
							"namespace": "default",
							"labels": map[string]interface{}{
								"app": "nginx4",
							},
							"annotations": map[string]interface{}{
								"litmuschaos.io/chaos": "true",
							},
						},
						"spec": map[string]interface{}{
							"selector": map[string]interface{}{
								"matchLabels": map[string]interface{}{
									"app": "nginx4",
								},
							},
						},
					},
				},
			},

			isErr: true,
			check: true,
		},
	}

	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {

			f := newFixture(t)
			f.SetFakeClient()
			dynamicClient := f.dynamicClient.Resource(gvfakero)
			dynamic := dynamicClient.Namespace("default")

			if mock.check == true {
				for _, ro := range mock.rollout {
					_, err := dynamic.Create(&ro, metav1.CreateOptions{})
					if err != nil {
						fmt.Printf("rollout not created, err: %v", err)
					}

				}
			}
			_, err := f.litmusClient.LitmuschaosV1alpha1().ChaosEngines(mock.engine.Instance.Namespace).Create(mock.engine.Instance)
			if err != nil {
				fmt.Printf("engine not created, err: %v", err)
			}

			engine, err := CheckChaosAnnotation(&mock.engine, f.k8sClient, f.dynamicClient)
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


type fixture struct {
	t *testing.T
	// k8sClient is the fake client set for k8s native objects.
	k8sClient *fake.Clientset
	// litmusClient is the fake client set for litmus cr objects.
	litmusClient *litmusFakeClientset.Clientset

	dynamicClient *dynamicfake.FakeDynamicClient

	k8sObjects     []runtime.Object
	litmusObjects  []runtime.Object
	dynamicObjects []runtime.Object
}

// SetFakeClient initilizes the fake required clientsets
func (f *fixture) SetFakeClient() {
	// Load kubernetes client set by preloading with k8s objects.
	f.k8sClient = fake.NewSimpleClientset(f.k8sObjects...)

	// Load litmus client set by preloading with litmus objects.
	f.litmusClient = litmusFakeClientset.NewSimpleClientset(f.litmusObjects...)

	// Load dynamic client set by preloading with litmus objects.
	f.dynamicClient = dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), f.dynamicObjects...)
}

func newFixture(t *testing.T) *fixture {
	f := &fixture{}
	f.t = t
	f.k8sObjects = []runtime.Object{}
	f.litmusObjects = []runtime.Object{}
	f.dynamicObjects = []runtime.Object{}
	return f
}
