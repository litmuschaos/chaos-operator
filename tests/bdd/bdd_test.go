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
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	//Creating rbacs
	err = exec.Command("kubectl", "apply", "-f", "../../deploy/rbac.yaml").Run()
	Expect(err).To(BeNil())

	//Creating Chaos-Operator
	By("creating operator")
	err = exec.Command("kubectl", "apply", "-f", "../../deploy/operator.yaml").Run()
	Expect(err).To(BeNil())
	klog.Infoln("chaos-operator created successfully")

	//Creating pod delete service account
	By("creating pod delete sa")
	err = exec.Command("kubectl", "apply", "-f", "../manifest/pod_delete_rbac.yaml").Run()
	Expect(err).To(BeNil())

	//Wait for the creation of chaos-operator
	time.Sleep(50 * time.Second)

	//Check for the status of the chaos-operator
	operator, _ := client.CoreV1().Pods("litmus").List(metav1.ListOptions{LabelSelector: "name=chaos-operator"})
	for _, v := range operator.Items {

		Expect(string(v.Status.Phase)).To(Equal("Running"))
		break
	}
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

			_, err := client.AppsV1().Deployments("litmus").Create(deployment)
			if err != nil {
				klog.Infoln("Deployment is not created and error is ", err)
			}

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

						ENVList: []v1alpha1.ENVPair{
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
							"name": "pod-delete",
							"app.kubernetes.io/part-of": "litmus",
						},
					},
				},
			}

			_, err = clientSet.ChaosExperiments("litmus").Create(ChaosExperiment)
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
					Monitoring:       true,
					EngineState:      "active",
					Experiments: []v1alpha1.ExperimentList{
						{
							Name: "pod-delete",
						},
					},
				},
			}

			_, err = clientSet.ChaosEngines("litmus").Create(chaosEngine)
			Expect(err).To(BeNil())

			klog.Infoln("Chaosengine created successfully...")

			//Wait till the creation of runner pod resource
			time.Sleep(5 * time.Second)

			var runnerStatus v1.PodPhase

			//Wait for 90s for runner to start running, before failing the test
			for i := 0; i < 90; i++ {
				runner, err := client.CoreV1().Pods("litmus").Get("engine-nginx-runner", metav1.GetOptions{})
				runnerStatus = runner.Status.Phase
				Expect(err).To(BeNil())
				klog.Infof("Runner state is: %s\n", string(runnerStatus))
				if string(runnerStatus) != "Running" {
					time.Sleep(1 * time.Second)
				} else {
					break
				}
			}


			Expect(string(runnerStatus)).To(Or(Equal("Running"), Equal("Succeeded")))

			// Check for EngineStatus
			engine, err := clientSet.ChaosEngines("litmus").Get("engine-nginx", metav1.GetOptions{})
			Expect(err).To(BeNil())

			isInit := engine.Status.EngineStatus == v1alpha1.EngineStatusInitialized
			Expect(isInit).To(BeTrue())
		})
	})

	Context("Setting the EngineState of ChaosEngine as Stop", func() {

		It("Should delete chaos-resources", func() {

			engine, err := clientSet.ChaosEngines("litmus").Get("engine-nginx", metav1.GetOptions{})
			Expect(err).To(BeNil())

			// setting the EngineState of chaosEngine to stop
			engine.Spec.EngineState = v1alpha1.EngineStateStop

			_, err = clientSet.ChaosEngines("litmus").Update(engine)
			Expect(err).To(BeNil())

			klog.Infoln("Chaosengine updated successfully...")

			//Wait till the creation of runner pod
			time.Sleep(50 * time.Second)

		})

	})

	Context("Checking Default ChaosResources", func() {

		It("Should delete chaos-runner pod", func() {

			//Fetching engine-nginx-runner pod
			_, err := client.CoreV1().Pods("litmus").Get("engine-nginx-runner", metav1.GetOptions{})
			klog.Infof("%v\n", err)
			isNotFound := errors.IsNotFound(err)
			Expect(isNotFound).To(BeTrue())
			klog.Infoln("chaos-runner pod deletion verified")

		})

		It("Should change the engineStatus ", func() {

			//Fetching engineStatus
			engine, err := clientSet.ChaosEngines("litmus").Get("engine-nginx", metav1.GetOptions{})
			Expect(err).To(BeNil())

			isStopped := engine.Status.EngineStatus == v1alpha1.EngineStatusStopped
			Expect(isStopped).To(BeTrue())
		})
	})

	Context("Deletion of ChaosEngine", func() {

		It("Should delete chaos engine", func() {

			err := clientSet.ChaosEngines("litmus").Delete("engine-nginx", &metav1.DeleteOptions{})
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
					JobCleanUpPolicy: "delete",
					Monitoring:       true,
					EngineState:      "active",
					Experiments: []v1alpha1.ExperimentList{
						{
							Name: "pod-delete-1",
						},
					},
				},
			}

			_, err := clientSet.ChaosEngines("litmus").Create(chaosEngine)
			Expect(err).To(BeNil())

			time.Sleep(50 * time.Second)

		})
	})

	Context("Check for Chaos Resources for invalid engine", func() {

		It("Should delete chaos-runner pod", func() {

			//Fetching engine-nginx-runner pod
			_, err := client.CoreV1().Pods("litmus").Get("engine-nginx-runner", metav1.GetOptions{})
			klog.Infof("%v\n", err)
			isNotFound := errors.IsNotFound(err)
			Expect(isNotFound).To(BeTrue())
			klog.Infoln("chaos-runner pod deletion verified")

		})

		It("Should change EngineStatus ", func() {

			//Fetching engineStatus
			engine, err := clientSet.ChaosEngines("litmus").Get("engine-nginx", metav1.GetOptions{})
			Expect(err).To(BeNil())

			isComplete := engine.Status.EngineStatus == v1alpha1.EngineStatusCompleted
			Expect(isComplete).To(BeTrue())

		})
	})

	Context("Validate via Chaos-Operator Logs", func() {

		It("Should Generate Operator logs", func() {
			pods, err := client.CoreV1().Pods("litmus").List(metav1.ListOptions{
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

			podLogs, err := req.Stream()
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
