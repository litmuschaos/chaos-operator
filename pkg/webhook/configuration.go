/*
Copyright 2019 The LitmusChaos Authors.

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

package webhook

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	version "github.com/litmuschaos/chaos-operator/pkg/version"
)

const (
	validatorServiceName = "admission-controller-svc"
	validatorWebhook     = "litmuschaos-validation-webhook-cfg"
	validatorSecret      = "admission-controller-secret"
	webhookHandlerName   = "admission-controller.litmuschaos.io"
	validationPath       = "/validate"
	validationPort       = 8443
	webhookLabel         = "litmuschaos.io/component-name" + "=" + "admission-controller"
	webhooksvcLabel      = "litmuschaos.io/component-name" + "=" + "admission-controller-svc"
	// AdmissionNameEnvVar is the constant for env variable ADMISSION_WEBHOOK_NAME
	// which is the name of the current admission webhook
	AdmissionNameEnvVar = "ADMISSION_WEBHOOK_NAME"
	appCrt              = "app.crt"
	appKey              = "app.pem"
	rootCrt             = "ca.crt"
	litmuschaosVersion  = "litmuschaos.io/version"
)

type transformSvcFunc func(*corev1.Service)
type transformSecretFunc func(*corev1.Secret)
type transformConfigFunc func(*v1beta1.ValidatingWebhookConfiguration)

var (
	// TimeoutSeconds specifies the timeout for this webhook. After the timeout passes,
	// the webhook call will be ignored or the API call will fail based on the
	// failure policy.
	// The timeout value must be between 1 and 30 seconds.
	five = int32(5)
	// Ignore means that an error calling the webhook is ignored.
	Ignore = v1beta1.Ignore
	// transformation function lists to upgrade webhook resources
	transformSecret = []transformSecretFunc{}
	transformSvc    = []transformSvcFunc{}
	transformConfig = []transformConfigFunc{}
)

// createWebhookService creates our webhook Service resource if it does not
// exist.
func createWebhookService(
	ownerReference metav1.OwnerReference,
	serviceName string,
	namespace string,
	kubeClient *kubernetes.Clientset,
) error {

	_, err := kubeClient.CoreV1().Services(namespace).
		Get(serviceName, metav1.GetOptions{})

	if err == nil {
		return nil
	}

	// error other than 'not found', return err
	if !k8serror.IsNotFound(err) {
		return errors.Wrapf(
			err,
			"failed to get webhook service {%v}",
			serviceName,
		)
	}

	// create service resource that refers to admission server pod
	serviceLabels := map[string]string{"app": "admission-controller"}
	svcObj := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      serviceName,
			Labels: map[string]string{
				"app":                           "admission-controller",
				"litmuschaos.io/component-name": "admission-controller-svc",
				string(litmuschaosVersion):      version.Current(),
			},
			OwnerReferences: []metav1.OwnerReference{ownerReference},
		},
		Spec: corev1.ServiceSpec{
			Selector: serviceLabels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       443,
					TargetPort: intstr.FromInt(validationPort),
				},
			},
		},
	}
	_, err = kubeClient.CoreV1().Services(namespace).
		Create(svcObj)
	return err
}

// createAdmissionService creates our ValidatingWebhookConfiguration resource
// if it does not exist.
func createAdmissionService(
	ownerReference metav1.OwnerReference,
	validatorWebhook string,
	namespace string,
	serviceName string,
	signingCert []byte,
	kubeClient *kubernetes.Clientset,
) error {

	_, err := GetValidatorWebhook(validatorWebhook, kubeClient)
	// validator object already present, no need to do anything
	if err == nil {
		return nil
	}

	// error other than 'not found', return err
	if !k8serror.IsNotFound(err) {
		return errors.Wrapf(
			err,
			"failed to get webhook validator {%v}",
			validatorWebhook,
		)
	}

	webhookHandler := v1beta1.ValidatingWebhook{
		Name: webhookHandlerName,
		Rules: []v1beta1.RuleWithOperations{{
			Operations: []v1beta1.OperationType{
				v1beta1.Create,
			},
			Rule: v1beta1.Rule{
				APIGroups:   []string{"litmuschaos.io"},
				APIVersions: []string{"*"},
				Resources:   []string{"chaosengines"},
			},
		},
		},
		ClientConfig: v1beta1.WebhookClientConfig{
			Service: &v1beta1.ServiceReference{
				Namespace: namespace,
				Name:      serviceName,
				Path:      StrPtr(validationPath),
			},
			CABundle: signingCert,
		},
		TimeoutSeconds: &five,
		FailurePolicy:  &Ignore,
	}

	validator := &v1beta1.ValidatingWebhookConfiguration{
		TypeMeta: metav1.TypeMeta{
			Kind:       "validatingWebhookConfiguration",
			APIVersion: "admissionregistration.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: validatorWebhook,
			Labels: map[string]string{
				"app":                           "admission-controller",
				"litmuschaos.io/component-name": "admission-controller",
				string(litmuschaosVersion):      version.Current(),
			},
			OwnerReferences: []metav1.OwnerReference{ownerReference},
		},
		Webhooks: []v1beta1.ValidatingWebhook{webhookHandler},
	}

	_, err = kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Create(validator)

	return err
}

// createCertsSecret creates a self-signed certificate and stores it as a
// secret resource in Kubernetes.
func createCertsSecret(
	ownerReference metav1.OwnerReference,
	secretName string,
	serviceName string,
	namespace string,
	kubeClient *kubernetes.Clientset,
) (*corev1.Secret, error) {
	// Create a signing certificate
	caKeyPair, err := NewCA(fmt.Sprintf("%s-ca", serviceName))
	if err != nil {
		return nil, fmt.Errorf("failed to create root-ca: %v", err)
	}

	// Create app certs signed through the certificate created above
	apiServerKeyPair, err := NewServerKeyPair(
		caKeyPair,
		strings.Join([]string{serviceName, namespace, "svc"}, "."),
		serviceName,
		namespace,
		"cluster.local",
		[]string{},
		[]string{},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create server key pair: %v", err)
	}

	// create an opaque secret resource with certificate(s) created above
	secretObj := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":                           "admission-controller",
				"litmuschaos.io/component-name": "admission-controller",
				string(litmuschaosVersion):      version.Current(),
			},
			OwnerReferences: []metav1.OwnerReference{ownerReference},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			appCrt:  EncodeCertPEM(apiServerKeyPair.Cert),
			appKey:  EncodePrivateKeyPEM(apiServerKeyPair.Key),
			rootCrt: EncodeCertPEM(caKeyPair.Cert),
		},
	}
	return kubeClient.CoreV1().Secrets(namespace).Create(secretObj)
}

// GetValidatorWebhook fetches the webhook validator resource
func GetValidatorWebhook(
	validator string,
	kubeClient *kubernetes.Clientset,
) (*v1beta1.ValidatingWebhookConfiguration, error) {

	return kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Get(validator, metav1.GetOptions{})
}

// StrPtr convert a string to a pointer
func StrPtr(s string) *string {
	return &s
}

// InitValidationServer creates secret, service and admission validation k8s
// resources. All these resources are created in the same namespace where
// litmus components is running.
func InitValidationServer(ownerReference metav1.OwnerReference, kubeClient *kubernetes.Clientset) error {

	// Fetch our namespace
	litmusNamespace, err := getLitmusNamespace()
	if err != nil {
		return err
	}

	err = preUpgrade(litmusNamespace, kubeClient)
	if err != nil {
		return err
	}

	// Check to see if webhook secret is already present
	certSecret, err := GetSecret(litmusNamespace, validatorSecret, kubeClient)
	if err != nil {
		if k8serror.IsNotFound(err) {
			// Secret not found, create certs and the secret object
			certSecret, err = createCertsSecret(
				ownerReference,
				validatorSecret,
				validatorServiceName,
				litmusNamespace,
				kubeClient,
			)
			if err != nil {
				return fmt.Errorf(
					"failed to create secret(%s) resource %v",
					validatorSecret,
					err,
				)
			}
		} else {
			// Unable to read secret object
			return fmt.Errorf(
				"unable to read secret object %s: %v",
				validatorSecret,
				err,
			)
		}
	}

	signingCertBytes, ok := certSecret.Data[rootCrt]
	if !ok {
		return fmt.Errorf(
			"%s value not found in %s secret",
			rootCrt,
			validatorSecret,
		)
	}

	serviceErr := createWebhookService(
		ownerReference,
		validatorServiceName,
		litmusNamespace,
		kubeClient,
	)
	if serviceErr != nil {
		return fmt.Errorf(
			"failed to create Service{%s}: %v",
			validatorServiceName,
			serviceErr,
		)
	}

	validatorErr := createAdmissionService(
		ownerReference,
		validatorWebhook,
		litmusNamespace,
		validatorServiceName,
		signingCertBytes,
		kubeClient,
	)
	if validatorErr != nil {
		return fmt.Errorf(
			"failed to create validator{%s}: %v",
			validatorWebhook,
			validatorErr,
		)
	}

	return nil
}

// GetSecret fetches the secret resource in the given namespace.
func GetSecret(
	namespace string,
	secretName string,
	kubeClient *kubernetes.Clientset,
) (*corev1.Secret, error) {

	return kubeClient.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
}

// getLitmusNamespace gets the namespace ADMISSION_NAMESPACE env value which is
// set by the downward API where admission server has been deployed
func getLitmusNamespace() (string, error) {
	ns, found := os.LookupEnv("LITMUS_NAMESPACE")
	if !found {
		return "", fmt.Errorf("%s must be set", "LITMUS_NAMESPACE")
	}
	return ns, nil
}

// GetAdmissionName return the admission server name
func GetAdmissionName() (string, error) {
	admissionName, found := os.LookupEnv(AdmissionNameEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", AdmissionNameEnvVar)
	}
	if len(admissionName) == 0 {
		return "", fmt.Errorf("%s must not be empty", AdmissionNameEnvVar)
	}
	return admissionName, nil
}

// GetAdmissionReference is a utility function to fetch a reference
// to the admission webhook deployment object
func GetAdmissionReference(kubeClient *kubernetes.Clientset) (*metav1.OwnerReference, error) {

	// Fetch our namespace
	litmusNamespace, err := getLitmusNamespace()
	if err != nil {
		return nil, err
	}

	// Fetch our admission server deployment object
	admdeployList, err := kubeClient.AppsV1().Deployments(litmusNamespace).List(metav1.ListOptions{LabelSelector: webhookLabel})
	if err != nil {
		return nil, fmt.Errorf("failed to list admission deployment: %s", err.Error())
	}
	for _, admdeploy := range admdeployList.Items {
		if len(admdeploy.Name) != 0 {
			return metav1.NewControllerRef(admdeploy.GetObjectMeta(), schema.GroupVersionKind{
				Group:   appsv1.SchemeGroupVersion.Group,
				Version: appsv1.SchemeGroupVersion.Version,
				Kind:    "Deployment",
			}), nil

		}
	}
	return nil, fmt.Errorf("failed to create deployment ownerReference")
}

// preUpgrade checks for the required older webhook configs,older
// then 1.3.0 if exists delete them.
func preUpgrade(litmusNamespace string, kubeClient *kubernetes.Clientset) error {
	secretlist, err := kubeClient.CoreV1().Secrets(litmusNamespace).List(metav1.ListOptions{LabelSelector: webhookLabel})
	if err != nil {
		return fmt.Errorf("failed to list old secret: %s", err.Error())
	}

	for _, scrt := range secretlist.Items {
		if scrt.Labels[string(litmuschaosVersion)] != version.Current() {
			if scrt.Labels[string(litmuschaosVersion)] == "" {
				err = kubeClient.CoreV1().Secrets(litmusNamespace).Delete(scrt.Name, &metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("failed to delete old secret %s: %s", scrt.Name, err.Error())
				}
			} else {
				newScrt := scrt
				for _, t := range transformSecret {
					t(&newScrt)
				}
				newScrt.Labels[string(litmuschaosVersion)] = version.Current()
				_, err := kubeClient.CoreV1().Secrets(litmusNamespace).Update(&newScrt)
				if err != nil {
					return fmt.Errorf("failed to update old secret %s: %s", scrt.Name, err.Error())
				}
			}
		}
	}
	svcList, err := kubeClient.CoreV1().Services(litmusNamespace).List(metav1.ListOptions{LabelSelector: webhooksvcLabel})
	if err != nil {
		return fmt.Errorf("failed to list old service: %s", err.Error())
	}
	for _, service := range svcList.Items {
		if service.Labels[string(litmuschaosVersion)] != version.Current() {
			if service.Labels[string(litmuschaosVersion)] == "" {
				err = kubeClient.CoreV1().Services(litmusNamespace).Delete(service.Name, &metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("failed to delete old service %s: %s", service.Name, err.Error())
				}
			} else {
				newSvc := service
				for _, t := range transformSvc {
					t(&newSvc)
				}
				newSvc.Labels[string(litmuschaosVersion)] = version.Current()
				_, err = kubeClient.CoreV1().Services(litmusNamespace).Update(&newSvc)
				if err != nil {
					return fmt.Errorf("failed to update old service %s: %s", service.Name, err.Error())
				}
			}
		}
	}
	webhookConfigList, err := kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().List(metav1.ListOptions{LabelSelector: webhookLabel})
	if err != nil {
		return fmt.Errorf("failed to list older webhook config: %s", err.Error())
	}

	for _, config := range webhookConfigList.Items {
		if config.Labels[string(litmuschaosVersion)] != version.Current() {
			if config.Labels[string(litmuschaosVersion)] == "" {
				err = kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Delete(config.Name, &metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("failed to delete older webhook config %s: %s", config.Name, err.Error())
				}
			} else {
				newConfig := config
				for _, t := range transformConfig {
					t(&newConfig)
				}
				newConfig.Labels[string(litmuschaosVersion)] = version.Current()
				_, err = kubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Update(&newConfig)
				if err != nil {
					return fmt.Errorf("failed to update older webhook config %s: %s", config.Name, err.Error())
				}
			}
		}
	}

	return nil
}
