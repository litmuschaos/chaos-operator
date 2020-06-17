#### Sample Chaos Experiment 

This folder has the required manifests to create, and run a simple pod delete chaos experiment.
For more details, about Litmus and other experiments, please ref: <https://docs.litmuschaos.io>

### Steps to follow to run this experiment

Apply the Kubernetes manifests in the described order to trigger the experiment. 

-   Deploy the Litmus chaos operator YAML. This will install the chaos CRDs in the cluster & also chaos operator deployment in `litmus` namespace

    ```yaml
    kubectl apply -f 01-chaos-operator.yaml`
    ```

-   Create a sample nginx deployment which is annotated for chaos

    ```yaml
    kubectl apply -f 02-annotated-nginx-deploy.yaml
    ```

-   Install the pod-delete chaosexperiment custom resource. 

    ```yaml
    kubectl apply -f 03-pod-delete-experiment.yaml 
    ```

-   Create a serviceaccount with just enough permissions to execute the experiment

    ```yaml
    kubectl apply -f 04-rbac-manifest.yaml
    ```
        
-   Create the chaosengine custom resource which ties the nginx app instance with the pod-delete experiment specification.

    ```yaml
    kubectl apply -f 05-chaos-engine.yaml
    ```
        
-   To check the status of the this experiment, refer to the chaosengine status
        
    ```yaml
    kubectl describe chaosengine engine -n litmus
    ```
    ```yaml
    ...
    Status:
      Engine Status:  completed
      Experiments:
        Last Update Time:  2020-03-31T15:05:39Z
        Name:              pod-delete-ceqkir
        Status:            Execution Successful
        Verdict:           Pass
    ```