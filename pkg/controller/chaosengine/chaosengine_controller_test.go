package chaosengine

import (
	"fmt"
	"testing"

	litmuschaosv1alpha1 "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Testadd(t *testing.T)
{
	tests := map[string]struct 
	{
		engine engineInfo
		isErr  bool
	}
	{
		"Test Positive": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-monitor",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
					},
				},
			},
			isErr: false,
		},
		"Test Negative": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring: true,
					},
				},
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		name, mock := name, mock
		t.Run(name, func(t *testing.T) {

			_, err := Add(mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}
// End Add() unit test

func TestNewRunnerPodForCR(t *testing.T) {
	tests := map[string]struct {
		engine engineInfo
		isErr  bool
	}{
		"Test Positive-1": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          true,
					},
				},
				appUUID:        "fake_id",
				appExperiments: []string{"exp-1"},
			},
			isErr: false,
		},
		"Test Positive-2": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
					},
				},
				appUUID:        "fake_id",
				appExperiments: []string{"exp-1"},
			},

			isErr: false,
		},
		"Test Negative-1": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{},
				},
				appUUID:        "fake_id",
				appExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
		"Test Negative-2 ": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
					},
				},
				appUUID:        "",
				appExperiments: []string{"exp-1"},
			},
			isErr: true,
		},
		"Test Negative-3 ": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-runner",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
					},
				},
				appUUID:        "fake_id",
				appExperiments: []string{},
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		name, mock := name, mock
		t.Run(name, func(t *testing.T) {
			_, err := newRunnerPodForCR(mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}
func TestNewMonitorServiceForCR(t *testing.T) {
	tests := map[string]struct {
		engine engineInfo
		isErr  bool
	}{
		"Test Positive": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-monitor",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
					},
				},
			},
			isErr: false,
		},
		"Test Negative": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring: true,
					},
				},
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		name, mock := name, mock
		t.Run(name, func(t *testing.T) {

			_, err := newMonitorServiceForCR(mock.engine)
			if mock.isErr && err == nil {
				t.Fatalf("Test %q failed: expected error not to be nil", name)
			}
			if !mock.isErr && err != nil {
				t.Fatalf("Test %q failed: expected error to be nil", name)
			}
		})
	}
}
func TestNewMonitorPodForCR(t *testing.T) {
	tests := map[string]struct {
		engine engineInfo
		isErr  bool
	}{
		"Test Positive": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-monitor",
						Namespace: "test",
					},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						ChaosServiceAccount: "fake-serviceAccount",
						Monitoring:          false,
					},
				},
			},
			isErr: false,
		},
		"Test Negative": {
			engine: engineInfo{
				instance: &litmuschaosv1alpha1.ChaosEngine{
					ObjectMeta: metav1.ObjectMeta{},
					Spec: litmuschaosv1alpha1.ChaosEngineSpec{
						Monitoring: true,
					},
				},
			},
			isErr: true,
		},
	}
	for name, mock := range tests {
		name, mock := name, mock
		t.Run(name, func(t *testing.T) {

			_, err := newMonitorPodForCR(mock.engine)
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
		name, mock := name, mock
		t.Run(name, func(t *testing.T) {
			appInfo := &applicationInfo{
				namespace: "namespace",
				label:     map[string]string{"fake_id": "aa"},
				experimentList: []litmuschaosv1alpha1.ExperimentList{
					{
						Name: "fake_name",
					},
				},
				serviceAccountName: "fake-serviceaccountname",
			}
			_, err := appInfo.initializeApplicationInfo(mock.instance)
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
