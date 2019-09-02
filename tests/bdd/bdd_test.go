package bdd

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	v1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis"
	clientV1alpha1 "github.com/litmuschaos/chaos-operator/pkg/clientset/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

func TestChaosOperator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ChaosOperator BDD")
}

var _ = BeforeSuite(func() {

	//creating crd of chaos engine
	By("Creating chaos engine")
	crdcmd := exec.Command("kubectl", "create", "-f", "../../deploy/crds/chaosengine_crd.yaml").Run()
	Expect(crdcmd).To(BeNil())

	// creating namespace litmus
	nscmd := exec.Command("kubectl", "create", "ns", "litmus").Run()
	Expect(nscmd).To(BeNil())

	// creating chaos
	chaosenginecmd := exec.Command("kubectl", "create", "-f", "../../deploy/crds/chaosengine.yaml").Run()
	Expect(chaosenginecmd).To(BeNil())

})

var _ = Describe("Chaos operator Suites", func() {

	Context("Chaos Engine Liviness check", func() {
		It("should fetch chaosengine (engine-nginx)", func() {

			By("Checking kubeconfig")

			kubeconfig, err := GetConfigPath()
			Expect(err).To(BeNil())

			config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				fmt.Println("KubeConfig Path is wrong", err)
				os.Exit(1)
			}
			Expect(err).To(BeNil())

			v1alpha1.AddToScheme(scheme.Scheme)

			clientSet, err := clientV1alpha1.NewForConfig(config)
			Expect(err).To(BeNil())

			By("Checking ChaosEngine")
			engine, err := clientSet.ChaosEngines("litmus").List(metav1.ListOptions{})
			Expect(err).To(BeNil())
			Expect(string(engine.Items[0].Name)).To(Equal("engine-nginx"))

		})
	})

})

var _ = AfterSuite(func() {
	// command for delete the custom resource definition
	deletecrdcmd := exec.Command("kubectl", "delete", "crd", "--all").Run()
	Expect(deletecrdcmd).To(BeNil())

	// command for delete the namespace
	deletenscmd := exec.Command("kubectl", "delete", "ns", "litmus").Run()
	Expect(deletenscmd).To(BeNil())

})
