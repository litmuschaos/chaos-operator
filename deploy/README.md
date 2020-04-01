#### Sample Chaos Experiment 

This folder has the required manifests to create, and run a simple pod delete experiment.
For more details, about Litmus and other experiments, please ref: docs.litmuschaos.io


### Steps to follow to run this experiment:

    - Apply the following YAML files
        `kubectl apply -f 01-chaos-operator.yaml 02-annotated-nginx-deploy.yaml 03-pod-delete-experiment 04-rbac-manifest.yaml 04-chaos-engine.yaml`

    - To check the status of the this experiment, use this command
        `kubectl get chaosengine engine -n litmus -oyaml`
      And check the status of variou experiments executed.