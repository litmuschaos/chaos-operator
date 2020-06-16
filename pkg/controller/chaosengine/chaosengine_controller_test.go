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

package chaosengine

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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
	r          *ReconcileChaosEngine
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

func TestNewRunnerPodForCR(t *testing.T) {
	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          true,
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: false,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Positive-3": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						AnnotationCheck:     "false",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Positive-4": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						AnnotationCheck:     "true",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
		"Test Negative-2 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
					},
				},
				AppUUID:        "",
				AppExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
		"Test Negative-3 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{},
			},
			isErr: true,
		},
		"Test Negative-4 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "",
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{},
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := newRunnerPodForCR(&mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}
func TestInitializeApplicationInfo(t *testing.T) {
	tests := map[string]struct {
		instance *litmuschaosv1alpha1.ChaosEngine
		isErr    bool
	}{
		"Test Positive": {
			instance: &litmuschaosv1alpha1.ChaosEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-monitor",
					Namespace: "test",
				},
				Spec: litmuschaosv1alpha1.ChaosEngineSpec{
					Appinfo: litmuschaosv1alpha1.ApplicationParams{
						Applabel: "key=value",
					},
				},
			},
			isErr: false,
		},
		"Test Negative": {
			instance: nil,
			isErr:    true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			appInfo := &chaosTypes.ApplicationInfo{
				Namespace: "namespace",
				Label:     map[string]string{"fake_id": "aa"},
				ExperimentList: []litmuschaosv1alpha1.ExperimentList{
					{
						Name: "fake_name",
					},
				},
				ServiceAccountName: "fake-service-account-name",
			}
			_, err := initializeApplicationInfo(mock.instance, appInfo)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				fmt.Println(err)
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}
func TestGetChaosRunnerENV(t *testing.T) {
	fakeEngineName := "Fake Engine"
	fakeNameSpace := "Fake NameSpace"
	fakeServiceAcc := "Fake Service Account"
	fakeAppLabel := "Fake Label"
	fakeAExList := []string{"fake string"}
	fakeAuxilaryAppInfo := "ns1:name=percona,ns2:run=nginx"
	fakeClientUUID := "12345678-9012-3456-7890-123456789012"

	tests := map[string]struct {
		instance       *litmuschaosv1alpha1.ChaosEngine
		aExList        []string
		expectedResult []corev1.EnvVar
	}{
		"Test Positive": {
			instance: &litmuschaosv1alpha1.ChaosEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fakeEngineName,
					Namespace: fakeNameSpace,
				},
				Spec: litmuschaosv1alpha1.ChaosEngineSpec{
					ChaosServiceAccount: fakeServiceAcc,
					Appinfo: litmuschaosv1alpha1.ApplicationParams{
						Applabel: fakeAppLabel,
						Appns:    fakeNameSpace,
					},
					AuxiliaryAppInfo: fakeAuxilaryAppInfo,
				},
			},
			aExList: fakeAExList,
			expectedResult: []corev1.EnvVar{
				{
					Name:  "CHAOSENGINE",
					Value: fakeEngineName,
				},
				{
					Name:  "APP_LABEL",
					Value: fakeAppLabel,
				},
				{
					Name:  "APP_NAMESPACE",
					Value: fakeNameSpace,
				},
				{
					Name:  "EXPERIMENT_LIST",
					Value: fmt.Sprint(strings.Join(fakeAExList, ",")),
				},
				{
					Name:  "CHAOS_SVC_ACC",
					Value: fakeServiceAcc,
				},
				{
					Name:  "AUXILIARY_APPINFO",
					Value: fakeAuxilaryAppInfo,
				},
				{
					Name:  "CLIENT_UUID",
					Value: fakeClientUUID,
				},
				{
					Name:  "CHAOS_NAMESPACE",
					Value: fakeNameSpace,
				},
			},
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			actualResult := getChaosRunnerENV(mock.instance, mock.aExList, fakeClientUUID)
			println(actualResult)
			if len(actualResult) != 8 {
				t.Fatalf("Test %q failed: expected array length to be 8", name)
			}
			for index, result := range actualResult {
				if result.Value != mock.expectedResult[index].Value {
					t.Fatalf("Test %q failed: actual result %q, received result %q", name, result, mock.expectedResult[index])
				}
			}
		})
	}
}

func TestGetApplicationDetail(t *testing.T) {
	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Positive": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-monitor",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "key=value",
						},
					},
				},
			},
			isErr: false,
		},
		"Test Negative": {
			engine: chaosTypes.EngineInfo{
				Instance: nil,
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			err := getApplicationDetail(&mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				fmt.Println(err)
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}

func TestGetAnnotationCheck(t *testing.T) {
	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          true,
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "run=nginx",
						},
						AnnotationCheck: "true",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: false,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "run=nginx",
						},
						AnnotationCheck: "false",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
					},
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
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						AnnotationCheck:     "fakeCheck",
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "run=nginx",
						},
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
							},
						},
					},
				},

				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			err := getAnnotationCheck(&mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				fmt.Println(err)
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}

func TestValidateAnnontatedApplication(t *testing.T) {
	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Positive-1": {
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
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: false,
		},
		"Test Positive-2": {
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
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Positive-3": {
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
			},
			isErr: false,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := clientSet.ChaosEngines(mock.engine.Instance.Namespace).Create(mock.engine.Instance)
			if err != nil {
				fmt.Printf("engine nai bna, err: %v", err)
			}
			err = r.validateAnnontatedApplication(&mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				fmt.Println(err)
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}

func TestUpdateEngineForComplete(t *testing.T) {
	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "engine-complete-p1",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx",
							AppKind:  "deployment",
						},
						EngineState:     litmuschaosv1alpha1.EngineStateActive,
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
					Status: litmuschaosv1alpha1.ChaosEngineStatus{
						EngineStatus: litmuschaosv1alpha1.EngineStatusCompleted,
					},
				},
			},
			isErr: false,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "engine-complete-p2",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx",
							AppKind:  "deployment",
						},
						EngineState:     litmuschaosv1alpha1.EngineStateActive,
						AnnotationCheck: "false",
						Experiments: []litmuschaosv1alpha1.ExperimentList{
							{
								Name: "exp-1",
							},
						},
					},
					Status: litmuschaosv1alpha1.ChaosEngineStatus{
						EngineStatus: litmuschaosv1alpha1.EngineStatusCompleted,
					},
				},
			},
			isErr: false,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := clientSet.ChaosEngines(mock.engine.Instance.Namespace).Create(mock.engine.Instance)
			if err != nil {
				fmt.Printf("engine nai bna, err: %v", err)
			}
			err = r.updateEngineForComplete(&mock.engine, true)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				fmt.Println(err)
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}

func TestNewGoRunnerPodForCR(t *testing.T) {
	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          true,
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
								Command: []string{
									"cmd1",
									"cmd2",
								},
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: false,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image:           "fake-runner-image",
								ImagePullPolicy: "Always",
								Args: []string{
									"args1",
									"args2",
								},
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Positive-3": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						AnnotationCheck:     "false",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image:           "fake-runner-image",
								ImagePullPolicy: "IfNotPresent",
								Command: []string{
									"cmd1",
									"cmd2",
								},
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Positive-4": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						AnnotationCheck:     "true",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image:           "fake-runner-image",
								ImagePullPolicy: "Never",
								Args: []string{
									"args1",
									"args2",
								},
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
		"Test Negative-2 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
					},
				},
				AppUUID:        "",
				AppExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
		"Test Negative-3 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{},
			},
			isErr: true,
		},
		"Test Negative-4 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "",
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{},
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := newGoRunnerPodForCR(&mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}

func TestNewAnsibleRunnerPodForCR(t *testing.T) {
	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          true,
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "fake-runner-image",
								Command: []string{
									"cmd1",
									"cmd2",
								},
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: false,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image:           "fake-runner-image",
								ImagePullPolicy: "Always",
								Args: []string{
									"args1",
									"args2",
								},
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Positive-3": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						AnnotationCheck:     "false",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image:           "fake-runner-image",
								ImagePullPolicy: "IfNotPresent",
								Command: []string{
									"cmd1",
									"cmd2",
								},
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Positive-4": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
						AnnotationCheck:     "true",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image:           "fake-runner-image",
								ImagePullPolicy: "Never",
								Args: []string{
									"args1",
									"args2",
								},
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Negative-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
		"Test Negative-2 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
					},
				},
				AppUUID:        "",
				AppExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
		"Test Negative-3 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{},
			},
			isErr: true,
		},
		"Test Negative-4 ": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Components: litmuschaosv1alpha1.ComponentParams{
							Runner: litmuschaosv1alpha1.RunnerInfo{
								Image: "",
							},
						},
					},
				},
				AppUUID:        "fake_id",
				AppExperiments: []string{},
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := newAnsibleRunnerPodForCR(&mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}

func TestInitEngine(t *testing.T) {
	tests := map[string]struct {
		engine chaosTypes.EngineInfo
		isErr  bool
	}{
		"Test Positive-1": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "engine-complete-p1",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx",
							AppKind:  "deployment",
						},
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
					Status: litmuschaosv1alpha1.ChaosEngineStatus{
						EngineStatus: litmuschaosv1alpha1.EngineStatusCompleted,
					},
				},
			},
			isErr: false,
		},
		"Test Positive-2": {
			engine: chaosTypes.EngineInfo{
				Instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "engine-complete-p2",
						Namespace: "default",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Appinfo: litmuschaosv1alpha1.ApplicationParams{
							Applabel: "app=nginx",
							AppKind:  "deployment",
						},
						EngineState:     "active",
						AnnotationCheck: "false",
						Experiments: []litmuschaosv1alpha1.ExperimentList{
							{
								Name: "exp-1",
							},
						},
					},
					Status: litmuschaosv1alpha1.ChaosEngineStatus{
						EngineStatus: litmuschaosv1alpha1.EngineStatusStopped,
					},
				},
			},
			isErr: false,
		},
	}
	for name, mock := range tests {
		t.Run(name, func(t *testing.T) {

			err := r.initEngine(&mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				fmt.Println(err)
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}
