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

package bdd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/martian/log"
	"github.com/litmuschaos/litmus-go/pkg/utils/retry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	scheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	"github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	chaosClient "github.com/litmuschaos/chaos-operator/pkg/client/clientset/versioned/typed/litmuschaos/v1alpha1"
	restclient "k8s.io/client-go/rest"
)

var (
	kubeconfig string
	config     *restclient.Config
	client     *kubernetes.Clientset
	clientSet  *chaosClient.LitmuschaosV1alpha1Client
)

func TestChaos(t *testing.T) {

	RegisterFailHandler(Fail)
	RunSpecs(t, "BDD test")
}

var _ = BeforeSuite(func() {

	var err error
	kubeconfig = os.Getenv("HOME") + "/.kube/config"
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	Expect(err).To(BeNil(), "failed to get config")

	client, err = kubernetes.NewForConfig(config)
	Expect(err).To(BeNil(), "failed to get client")

	clientSet, err = chaosClient.NewForConfig(config)
	Expect(err).To(BeNil(), "failed to get clientSet")

	err = v1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).To(BeNil())

	//Creating crds
	By("creating all crds")
	err = exec.Command("kubectl", "apply", "-f", "../../deploy/chaos_crds.yaml").Run()
	Expect(err).To(BeNil())
	klog.Infoln("CRDS created")

	//Creating rbacs
	err = exec.Command("kubectl", "apply", "-f", "../../deploy/rbac.yaml").Run()
	Expect(err).To(BeNil())
	klog.Infoln("RBAC created")

	//Creating Chaos-Operator
	By("creating operator")
	err = exec.Command("kubectl", "apply", "-f", "../../deploy/operator.yaml").Run()
	Expect(err).To(BeNil())
	klog.Infoln("chaos-operator created successfully")

	//Creating pod delete service account
	By("creating pod delete sa")
	err = exec.Command("kubectl", "apply", "-f", "../manifest/pod_delete_rbac.yaml").Run()
	Expect(err).To(BeNil())
	klog.Infoln("pod-delete-sa created")

	err = retry.
		Times(uint(180 / 2)).
		Wait(time.Duration(2) * time.Second).
		Try(func(attempt uint) error {
			podSpec, err := client.CoreV1().Pods("litmus").List(context.Background(), metav1.ListOptions{LabelSelector: "name=chaos-operator"})
			if err != nil || len(podSpec.Items) == 0 {
				return fmt.Errorf("Unable to list chaos-operator, err: %v", err)
			}
			for _, v := range podSpec.Items {
				if v.Status.Phase != "Running" {
					return fmt.Errorf("chaos-operator is not in running state, phase: %v", v.Status.Phase)
				}
			}
			return nil
		})

	Expect(err).To(BeNil(), "the chaos-operator is not in running state")
	klog.Infoln("Chaos-Operator is in running state")
})

//BDD Tests to check secondary resources
var _ = Describe("BDD on chaos-operator", func() {

	// BDD TEST CASE 1
	Context("Check for the custom resources", func() {

		It("Should check for creation of runner pod", func() {

			//creating nginx deployment
			deployment := &appv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "nginx",
					Namespace: "litmus",
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
							ServiceAccountName: "litmus",
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

			_, err := client.AppsV1().Deployments("litmus").Create(context.Background(), deployment, metav1.CreateOptions{})
			Expect(err).To(
				BeNil(),
				"while creating nginx deployment in namespace litmus",
			)
			klog.Infoln("nginx deployment created")

			//creating chaos-experiment for pod-delete
			By("Creating ChaosExperiments")
			ChaosExperiment := &v1alpha1.ChaosExperiment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod-delete",
					Namespace: "litmus",
					Labels: map[string]string{
						"litmuschaos.io/name": "kubernetes",
					},
				},
				Spec: v1alpha1.ChaosExperimentSpec{
					Definition: v1alpha1.ExperimentDef{
						Image: "litmuschaos/go-runner:ci",

						Scope: "Namespaced",

						Permissions: []rbacV1.PolicyRule{},

						Args:    []string{"-c", "./experiments -name pod-delete"},
						Command: []string{"/bin/bash"},

						ENVList: []v1.EnvVar{
							{
								Name:  "TOTAL_CHAOS_DURATION",
								Value: "30",
							},
							{
								Name:  "CHAOS_INTERVAL",
								Value: "5",
							},
							{
								Name:  "LIB",
								Value: "litmus",
							},
							{
								Name:  "FORCE",
								Value: "true",
							},
						},

						Labels: map[string]string{
							"name":                      "pod-delete",
							"app.kubernetes.io/part-of": "litmus",
						},
					},
				},
			}

			_, err = clientSet.ChaosExperiments("litmus").Create(context.Background(), ChaosExperiment)
			Expect(err).To(BeNil())
			klog.Infoln("ChaosExperiment created successfully...")

			//Creating chaosEngine
			By("Creating ChaosEngine")
			chaosEngine := &v1alpha1.ChaosEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "engine-nginx",
					Namespace: "litmus",
				},
				Spec: v1alpha1.ChaosEngineSpec{
					Appinfo: v1alpha1.ApplicationParams{
						Appns:    "litmus",
						Applabel: "app=nginx",
						AppKind:  "deployment",
					},
					ChaosServiceAccount: "pod-delete-sa",
					Components: v1alpha1.ComponentParams{
						Runner: v1alpha1.RunnerInfo{
							Image: "litmuschaos/chaos-runner:ci",
							Type:  "go",
						},
					},
					JobCleanUpPolicy: "retain",
					EngineState:      "active",
					Experiments: []v1alpha1.ExperimentList{
						{
							Name: "pod-delete",
						},
					},
				},
			}

			_, err = clientSet.ChaosEngines("litmus").Create(context.Background(), chaosEngine)
			Expect(err).To(BeNil())
			klog.Infoln("Chaosengine created successfully...")

			err = retry.
				Times(uint(180 / 2)).
				Wait(time.Duration(2) * time.Second).
				Try(func(attempt uint) error {
					pod, err := client.CoreV1().Pods("litmus").Get(context.Background(), "engine-nginx-runner", metav1.GetOptions{})
					if err != nil {
						return fmt.Errorf("unable to get chaos-runner pod, err: %v", err)
					}
					if pod.Status.Phase != v1.PodRunning && pod.Status.Phase != v1.PodSucceeded {
						return fmt.Errorf("chaos runner is not in running state, phase: %v", pod.Status.Phase)
					}
					return nil
				})

			if err != nil {
				log.Errorf("The chaos-runner is not in running state, err: %v", err)
			}
			klog.Infoln("runner pod created")

			// Check for EngineStatus
			engine, err := clientSet.ChaosEngines("litmus").Get(context.Background(), "engine-nginx", metav1.GetOptions{})
			Expect(err).To(BeNil())

			isInit := engine.Status.EngineStatus == v1alpha1.EngineStatusInitialized
			Expect(isInit).To(BeTrue())
		})
	})

	Context("Setting the EngineState of ChaosEngine as Stop", func() {

		It("Should delete chaos-resources", func() {

			engine, err := clientSet.ChaosEngines("litmus").Get(context.Background(), "engine-nginx", metav1.GetOptions{})
			Expect(err).To(BeNil())

			// setting the EngineState of chaosEngine to stop
			engine.Spec.EngineState = v1alpha1.EngineStateStop

			_, err = clientSet.ChaosEngines("litmus").Update(context.Background(), engine)
			Expect(err).To(BeNil())
			klog.Infoln("Chaosengine updated successfully...")

			err = retry.
				Times(uint(180 / 2)).
				Wait(time.Duration(2) * time.Second).
				Try(func(attempt uint) error {
					engine, err := clientSet.ChaosEngines("litmus").Get(context.Background(), "engine-nginx", metav1.GetOptions{})
					if err != nil {
						return fmt.Errorf("unable to get chaosengine, err: %v", err)
					}
					if engine.Spec.EngineState != v1alpha1.EngineStateStop {
						return fmt.Errorf("chaos engine is not in stopped state, state: %v", engine.Spec.EngineState)
					}
					return nil
				})
			Expect(err).To(BeNil())
		})

	})

	Context("Checking Default ChaosResources", func() {

		It("Should delete chaos-runner pod", func() {

			err := retry.
				Times(uint(180 / 2)).
				Wait(time.Duration(2) * time.Second).
				Try(func(attempt uint) error {
					_, err := client.CoreV1().Pods("litmus").Get(context.Background(), "engine-nginx-runner", metav1.GetOptions{})
					isNotFound := errors.IsNotFound(err)
					if isNotFound {
						return nil
					}
					return fmt.Errorf("chaos-runner is not deleted yet, err: %v", err)
				})
			Expect(err).To(BeNil())
			klog.Infoln("chaos-runner pod deletion verified")
		})

		It("Should change the engineStatus ", func() {

			err := retry.
				Times(uint(180 / 2)).
				Wait(time.Duration(2) * time.Second).
				Try(func(attempt uint) error {
					//Fetching engineStatus
					engine, err := clientSet.ChaosEngines("litmus").Get(context.Background(), "engine-nginx", metav1.GetOptions{})
					Expect(err).To(BeNil())
					if engine.Status.EngineStatus != v1alpha1.EngineStatusStopped {
						fmt.Printf("engine is not in stopped state")
					}
					return nil
				})
			Expect(err).To(BeNil())
		})
	})

	Context("Deletion of ChaosEngine", func() {
		It("Should delete chaos engine", func() {

			err := clientSet.ChaosEngines("litmus").Delete(context.Background(), "engine-nginx", &metav1.DeleteOptions{})
			Expect(err).To(BeNil())
			err = retry.
				Times(uint(180 / 2)).
				Wait(time.Duration(2) * time.Second).
				Try(func(attempt uint) error {
					_, err := clientSet.ChaosEngines("litmus").Get(context.Background(), "engine-nginx", metav1.GetOptions{})
					if err != nil && !k8serrors.IsNotFound(err) {
						return fmt.Errorf("unable to get chaosengine, err: %v", err)
					}
					return nil
				})
			Expect(err).To(BeNil())
			klog.Infoln("chaos engine deleted successfully")
		})

	})

	Context("Creation of ChaosEngine with invalid experiment", func() {
		It("Should create invalid chaos engine", func() {

			//Creating chaosEngine
			By("Creating ChaosEngine")
			chaosEngine := &v1alpha1.ChaosEngine{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "engine-nginx-1",
					Namespace: "litmus",
				},
				Spec: v1alpha1.ChaosEngineSpec{
					Appinfo: v1alpha1.ApplicationParams{
						Appns:    "litmus",
						Applabel: "app=nginx",
						AppKind:  "deployment",
					},
					ChaosServiceAccount: "pod-delete-sa",
					Components: v1alpha1.ComponentParams{
						Runner: v1alpha1.RunnerInfo{
							Image: "litmuschaos/chaos-runner:ci",
							Type:  "go",
						},
					},
					JobCleanUpPolicy: "delete",
					EngineState:      "active",
					Experiments: []v1alpha1.ExperimentList{
						{
							Name: "pod-delete-1",
						},
					},
				},
			}

			_, err := clientSet.ChaosEngines("litmus").Create(context.Background(), chaosEngine)
			Expect(err).To(BeNil())
		})
	})

	Context("Check for Chaos Resources for invalid engine", func() {
		It("Should delete chaos-runner pod", func() {

			err := retry.
				Times(uint(180 / 2)).
				Wait(time.Duration(2) * time.Second).
				Try(func(attempt uint) error {
					//Fetching engine-nginx-runner pod
					_, err := client.CoreV1().Pods("litmus").Get(context.Background(), "engine-nginx-1-runner", metav1.GetOptions{})
					isNotFound := errors.IsNotFound(err)
					if isNotFound {
						return nil
					}
					return fmt.Errorf("chaos-runner is not deleted yet, err: %v", err)
				})
			Expect(err).To(BeNil())
		})

		It("Should change EngineStatus ", func() {

			err := retry.
				Times(uint(180 / 2)).
				Wait(time.Duration(2) * time.Second).
				Try(func(attempt uint) error {
					//Fetching engineStatus
					engine, err := clientSet.ChaosEngines("litmus").Get(context.Background(), "engine-nginx-1", metav1.GetOptions{})
					if err != nil {
						return err
					}
					if engine.Status.EngineStatus != v1alpha1.EngineStatusCompleted {
						return fmt.Errorf("engine is not in completed state")
					}
					return nil
				})
			Expect(err).To(BeNil())
		})
	})

	Context("Validate via Chaos-Operator Logs", func() {
		It("Should Generate Operator logs", func() {

			pods, err := client.CoreV1().Pods("litmus").List(context.Background(), metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%v=%v", "name", "chaos-operator"),
			})
			Expect(err).To(BeNil())

			if len(pods.Items) > 1 {
				klog.Infof("Multiple Chaos-Operator Pods found")
				return
			}
			if len(pods.Items) < 1 {
				klog.Infof("Unable to find Chaos-Operator Pod")
				return
			}

			podName := pods.Items[0].Name
			Expect(podName).To(
				Not(BeEmpty()),
				"Unable to get the operator pod name",
			)

			klog.Infof("Got Pod Name: %v\n", podName)

			podLogOpts := v1.PodLogOptions{}

			req := client.CoreV1().Pods("litmus").GetLogs(podName, &podLogOpts)

			podLogs, err := req.Stream(context.Background())
			Expect(err).To(BeNil())

			defer podLogs.Close()

			buf := new(bytes.Buffer)
			_, err = io.Copy(buf, podLogs)
			Expect(err).To(BeNil())

			str := buf.String()

			klog.Infof("Chaos Operator Logs:\n%v\n", str)

		})
	})

})

//Deleting all unused resources
var _ = AfterSuite(func() {

	//Deleting Pod Delete sa
	By("Deleting pod delete sa")
	err := exec.Command("kubectl", "delete", "-f", "../manifest/pod_delete_rbac.yaml").Run()
	Expect(err).To(BeNil())

	//Deleting ChaosExperiments
	By("Deleting ChaosExperiments")
	err = exec.Command("kubectl", "delete", "chaosexperiments", "--all", "-n", "litmus").Run()
	Expect(err).To(BeNil())

	//Deleting ChaosEngines
	By("Deleting ChaosEngines")
	err = exec.Command("kubectl", "delete", "chaosengine", "--all", "-n", "litmus").Run()
	Expect(err).To(BeNil())

	//Deleting Chaos-Operator
	By("Deleting operator")
	err = exec.Command("kubectl", "delete", "-f", "../../deploy/operator.yaml").Run()
	Expect(err).To(BeNil())

	//Deleting rbacs
	By("Deleting RBAC's")
	err = exec.Command("kubectl", "delete", "-f", "../../deploy/rbac.yaml").Run()
	Expect(err).To(BeNil())

	//Deleting CRD's
	By("Deleting chaosengine crd")
	err = exec.Command("kubectl", "delete", "-f", "../../deploy/chaos_crds.yaml").Run()
	Expect(err).To(BeNil())

})
