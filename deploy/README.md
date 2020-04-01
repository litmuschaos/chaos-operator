#### Sample Chaos Experiment 

This folder has the required manifests to create, and run a simple pod delete experiment.
For more details, about Litmus and other experiments, please ref: docs.litmuschaos.io


### Steps to follow to run this experiment:

    - Apply the following YAML files

        - Apply Chaos Operator: `kubectl apply -f 01-chaos-operator.yaml`

        - Apply Nginx Deployment: `kubectl apply -f 02-annotated-nginx-deploy.yaml`

        - Apply Pod Delete Chaos Experiment: `kubectl apply -f 03-pod-delete-experiment.yaml`

        - Apply RBAC Manifest for Chaos Engine: `kubectl apply -f 05-chaos-engine.yaml`
        
        - At last, apply the ChaosEngine: `kubectl apply -f 04-rbac-manifest.yaml`
        
        - To check the status of the this experiment, use this command
        
          `kubectl get chaosengine engine -n litmus -oyaml`

            And check the status of variou experiments executed.