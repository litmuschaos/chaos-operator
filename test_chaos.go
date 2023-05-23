package unit_tests

import (
    "testing"
    "chaos-operator/chaos_operator"
    "chaos-operator/chaos_experiment"
)

func TestChaosOperator(t *testing.T) {
    // Create a new chaos-operator instance.
    operator := chaos_operator.NewChaosOperator()

    // Set the operator's configuration.
    operator.Config.KubeconfigPath = "/path/to/kubeconfig"
    operator.Config.Namespace = "default"

    // Start the operator.
    err := operator.Start()
    if err != nil {
        t.Errorf("Error starting operator: %v", err)
    }

    // Create a new chaos experiment.
    experiment := chaos_experiment.NewChaosExperiment()
    experiment.Name = "test-experiment"
    experiment.ChaosType = "pod-delete"
    experiment.TargetPods = []string{"nginx"}

    // Apply the chaos experiment.
    err = operator.ApplyChaosExperiment(experiment)
    if err != nil {
        t.Errorf("Error applying chaos experiment: %v", err)
    }

    // Wait for the chaos experiment to complete.
    err = operator.WaitForChaosExperimentCompletion(experiment)
    if err != nil {
        t.Errorf("Error waiting for chaos experiment to complete: %v", err)
    }

    // Check the chaos experiment's status.
    if experiment.Status != chaos_experiment.StatusCompleted {
        t.Errorf("Expected chaos experiment status to be Completed, got %s", experiment.Status)
    }

    // Stop the operator.
    operator.Stop()
}