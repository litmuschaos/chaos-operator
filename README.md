# Litmus chaos-operator for injecting chaos experiments on Kubernetes

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5d3a1caf80454c55bfa4fa4f6b1b9a9f)](https://www.codacy.com/app/chandan.kumar/chaos-operator?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=litmuschaos/chaos-operator&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/litmuschaos/chaos-operator)](https://goreportcard.com/report/github.com/litmuschaos/chaos-operator)
[![BCH compliance](https://bettercodehub.com/edge/badge/litmuschaos/chaos-operator?branch=master)](https://bettercodehub.com/)
  
Litmus chaos operator is used by Kubernetes application developers and SREs to inject chaos into the applications 
and Kubernetes infrastructure in a managed fashion. Its objective is to make the process of validation and 
hardening of application workloads on Kubernetes easy by automating the execution of chaos experiments. A sample chaos 
injection workflow could be as simple as:

- Install the Litmus infrastructure components (RBAC, CRDs), the Operator & Experiment custom resource bundles via helm charts
- Annotate the application under test (AUT), enabling it for chaos
- Create a ChaosEngine custom resource tied to the AUT, which describes the experiment list to be executed 

Benefits provided by the Chaos Operator include: 

- Scheduled batch Run of Chaos
- Standardised chaos experiment spec 
- Categorized chaos bundles for stateless/stateful/vendor-specific
- Test-Run resiliency 
- Ability to chaos run as a background service based on annotations

## What is a chaos operator and how is it built?

The Chaos Operator is a Kubernetes Operator, which are nothing but custom-controllers with direct access to Kubernetes API
that can manage the lifecycle of certain resources or applications, while always trying to ensure the resource is in the "desired
state". The logic that ensures this is commonly called "reconcile" function.

The Chaos Operator is built using the popular [Operator-SDK](https://github.com/operator-framework/operator-sdk/) framework, 
which provides bootstrap support for new operator projects, allowing teams to focus on business/operational logic. 

The Litmus Chaos Operator helps reconcile the state of the ChaosEngine, a custom resource that holds the chaos intent 
specified by a developer/devops engineer against a particular stateless/stateful Kubernetes deployment. The operator performs
specific actions upon CRUD of the ChaosEngine, its primary resource. The operator also defines secondary resources (the engine 
runner pod and engine monitor service), which are created & managed by it in order to implement the reconcile functions. 

## What is a chaos engine?

The ChaosEngine is the core schema that defines the chaos workflow for a given application. Currently, it defines the following:

- Application Data (namespace, labels, kind)
- List of Chaos Experiments to be executed
- Attributes of the experiments, such as, rank/priority 
- Execution Schedule for the batch run of the experiments

The ChaosEngine is the referenced as the owner of the secondary (reconcile) resource with Kubernetes deletePropagation 
ensuring these also are removed upon deletion of the ChaosEngine CR.

Here is a sample ChaosEngineSpec for reference: 

  ```yaml
  apiVersion: litmuschaos.io/v1alpha1
  kind: ChaosEngine
  metadata:
    name: engine-nginx
  spec:
    appinfo: 
      appns: default
      applabel: "app=nginx"
    experiments:
      - name: pod-delete 
        spec:
          rank: 
      - name: container-kill
        spec:
          rank:  
  ```

## What is a litmus chaos chart and how can I use it?

Litmus Chaos Charts are used to install "Chaos Experiment Bundles" & are categorized based on the nature
of the experiments (general Kubernetes chaos, vendor/provider specific chaos - such as, OpenEBS or 
application-specific chaos, say NuoDB). They consist of custom resources that hold low-level chaos(test) 
parameters which are queried by the operator in order to execute the experiments. The spec.definition._fields_
and their corresponding _values_ are used to construct the eventual execution artifact that runs the chaos 
experiment (typically, the litmusbook, which is a K8s job resource). 

Here is a sample ChaosEngineSpec for reference:

```yaml
apiVersion: litmuschaos.io/v1alpha1
description:
  message: |
    Deletes a pod belonging to a deployment/statefulset/daemonset
kind: ChaosExperiment
metadata:
  labels:
    helm.sh/chart: k8sChaos-0.1.0
    litmuschaos.io/instance: dealing-butterfly
    litmuschaos.io/name: k8sChaos
  name: pod-delete
spec:
  definition:
    image: openebs/ansible-runner:ci
    litmusbook: /experiments/chaos/kubernetes/pod_delete/run_litmus_test.yml
    labels:
      name: pod-delete
    args:
    - -c
    - ansible-playbook ./experiments/chaos/kubernetes/pod_delete/test.yml -i /etc/ansible/hosts
      -vv; exit 0
    command:
    - /bin/bash
    env:
    - name: ANSIBLE_STDOUT_CALLBACK
      value: null
    - name: TOTAL_CHAOS_DURATION
      value: 15
    - name: CHAOS_INTERVAL
      value: 5
    - name: LIB
      value: ""
```

## What are the steps to get started?

- Install Litmus infrastructure (RBAC, CRD, Operator) components 

  ```
  helm repo add litmuschaos https://litmuschaos.github.io/chaos-charts
  helm repo update
  helm install litmuschaos/litmus --namespace=litmus
  ```

- Download the desired Chaos Experiment bundles, say, general Kubernetes chaos

  ```
  helm install litmuschaos/k8sChaos
  ```

- Annotate your application to enable chaos. For ex:

  ```
  kubectl annotate deploy/nginx-deployment litmuschaos.io/chaos="true"
  ```

- Create a ChaosEngine CR with application information & chaos experiment list with their respective attributes

  ```
  # engine-nginx.yaml is a chaosengine manifest file
  kubectl apply -f engine-nginx.yaml
  ``` 
- Refer the ChaosEngine Status (or alternately, the corresponding ChaosResult resource) to know the status 
  of each experiment. The `spec.verdict` is set to _Running_ when the experiment is in progress, eventually
  changing to _pass_ or _fail_.

  ```
  kubectl describe chaosresult engine-nginx-pod-delete

  Name:         engine-nginx-pod-delete
  Namespace:    default
  Labels:       <none>
  Annotations:  kubectl.kubernetes.io/last-applied-configuration:
                {"apiVersion":"litmuschaos.io/v1alpha1","kind":"ChaosResult","metadata":{"annotations":{},"name":"engine-nginx-pod-delete","namespace":"de...
  API Version:  litmuschaos.io/v1alpha1
  Kind:         ChaosResult
  Metadata:
    Creation Timestamp:  2019-05-22T12:10:19Z
    Generation:          9
    Resource Version:    8898730
    Self Link:           /apis/litmuschaos.io/v1alpha1/namespaces/default/chaosresults/engine-nginx-pod-delete
    UID:                 911ada69-7c8a-11e9-b37f-42010a80019f
  Spec:
    Experimentstatus:
      Phase:    <nil>
      Verdict:  pass
  Events:       <none>
  ```

## Where are the docs?

They are available at [litmus docs](https://docs.litmuschaos.io)

## How do I contribute?

The Chaos Operator is in _alpha_ stage and needs all the help you can provide! Please contribute by raising issues, 
improving the documentation, contributing to the core framework and tooling, etc.
