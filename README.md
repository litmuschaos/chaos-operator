# Litmus chaos-operator for injecting chaos experiments on Kubernetes

[![Slack Channel](https://img.shields.io/badge/Slack-Join-purple)](https://slack.litmuschaos.io)
![GitHub Workflow](https://github.com/litmuschaos/chaos-operator/actions/workflows/push.yml/badge.svg?branch=master)
[![Docker Pulls](https://img.shields.io/docker/pulls/litmuschaos/chaos-operator.svg)](https://hub.docker.com/r/litmuschaos/chaos-operator)
[![GitHub issues](https://img.shields.io/github/issues/litmuschaos/chaos-operator)](https://github.com/litmuschaos/chaos-operator/issues)
[![Twitter Follow](https://img.shields.io/twitter/follow/litmuschaos?style=social)](https://twitter.com/LitmusChaos)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/2597079b1b5240d3866a6deb4112a2f2)](https://www.codacy.com/manual/litmuschaos/chaos-operator?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=litmuschaos/chaos-operator&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/litmuschaos/chaos-operator)](https://goreportcard.com/report/github.com/litmuschaos/chaos-operator)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5290/badge)](https://bestpractices.coreinfrastructure.org/projects/5290)
[![BCH compliance](https://bettercodehub.com/edge/badge/litmuschaos/chaos-operator?branch=master)](https://bettercodehub.com/)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Flitmuschaos%2Fchaos-operator.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Flitmuschaos%2Fchaos-operator?ref=badge_shield)
[![codecov](https://codecov.io/gh/litmuschaos/chaos-operator/branch/master/graph/badge.svg)](https://codecov.io/gh/litmuschaos/chaos-operator)
[![YouTube Channel](https://img.shields.io/badge/YouTube-Subscribe-red)](https://www.youtube.com/channel/UCa57PMqmz_j0wnteRa9nCaw)
<br><br>
  
Litmus chaos operator is used by Kubernetes application developers and SREs to inject chaos into the applications 
and Kubernetes infrastructure in a managed fashion. Its objective is to make the process of validation and 
hardening of application workloads on Kubernetes easy by automating the execution of chaos experiments. A sample chaos 
injection workflow could be as simple as:

-  Install the Litmus infrastructure components (RBAC, CRDs), the Operator & Experiment custom resource bundles via the operator manifest
-  Annotate the application under test (AUT), enabling it for chaos
-  Create a ChaosEngine custom resource tied to the AUT, which describes the experiment to be executed 

Benefits provided by the Chaos Operator include: 

-  Standardised chaos experiment spec 
-  Categorized chaos bundles for stateless/stateful/vendor-specific
-  Test-Run resiliency 
-  Ability to chaos run as a background service based on annotations

## What is a chaos operator and how is it built?

The Chaos Operator is a Kubernetes Operator, which are nothing but custom-controllers with direct access to Kubernetes API
that can manage the lifecycle of certain resources or applications, while always trying to ensure the resource is in the "desired
state". The logic that ensures this is commonly called "reconcile" function.

The Chaos Operator is built using the popular [Operator-SDK](https://github.com/operator-framework/operator-sdk/) framework, 
which provides bootstrap support for new operator projects, allowing teams to focus on business/operational logic. 

The Litmus Chaos Operator helps reconcile the state of the ChaosEngine, a custom resource that holds the chaos intent 
specified by a developer/devops engineer against a particular stateless/stateful Kubernetes deployment. The operator performs
specific actions upon CRUD of the ChaosEngine, its primary resource. The operator also defines a secondary resource (the engine 
runner pod), which is created & managed by it in order to implement the reconcile functions. 

## What is a chaos engine?

The ChaosEngine is the core schema that defines the chaos workflow for a given application. Currently, it defines the following:

-  Application info (namespace, labels, kind) of primary (AUT) and auxiliary (dependent) applications 
-  ServiceAccount used for execution of the experiment
-  Flag to turn on/off chaos annotation checks on applications
-  Chaos Experiment to be executed on the application
-  Attributes of the experiments (overrides defaults specified in the experiment CRs)
-  Flag to retain/cleanup chaos resources after experiment execution

The ChaosEngine is the referenced as the owner of the secondary (reconcile) resource with Kubernetes deletePropagation 
ensuring these also are removed upon deletion of the ChaosEngine CR.

Here is a sample ChaosEngineSpec for reference: <https://v1-docs.litmuschaos.io/docs/getstarted/#prepare-chaosengine>

## What is a litmus chaos chart and how can I use it?

Litmus Chaos Charts are used to install "Chaos Experiment Bundles" & are categorized based on the nature
of the experiments (general Kubernetes chaos, vendor/provider specific chaos - such as, OpenEBS or 
application-specific chaos, say NuoDB). They consist of custom resources that hold low-level chaos(test) 
parameters which are queried by the operator in order to execute the experiments. The spec.definition._fields_
and their corresponding _values_ are used to construct the eventual execution artifact that runs the chaos 
experiment (typically, the litmusbook, which is a K8s job resource). It also defines the permissions necessary 
to execute the experiment.  

Here is a sample ChaosEngineSpec for reference:

```yaml
apiVersion: litmuschaos.io/v1alpha1
description:
  message: |
    Deletes a pod belonging to a deployment/statefulset/daemonset
kind: ChaosExperiment
metadata:
  name: pod-delete
  labels:
    name: pod-delete
    app.kubernetes.io/part-of: litmus
    app.kubernetes.io/component: chaosexperiment
    app.kubernetes.io/version: latest
spec:
  definition:
    scope: Namespaced
    permissions:
      - apiGroups:
          - ""
          - "apps"
          - "batch"
          - "litmuschaos.io"
        resources:
          - "deployments"
          - "jobs"
          - "pods"
          - "configmaps"
          - "chaosengines"
          - "chaosexperiments"
          - "chaosresults"
        verbs:
          - "create"
          - "list"
          - "get"
          - "patch"
          - "update"
          - "delete"
    image: "litmuschaos/go-runner:latest"
    imagePullPolicy: Always
    args:
    - -c
    - ./experiments -name pod-delete
    command:
    - /bin/bash
    env:

    - name: TOTAL_CHAOS_DURATION
      value: '15'

    # Period to wait before/after injection of chaos in sec
    - name: RAMP_TIME
      value: ''

    - name: FORCE
      value: 'true'

    - name: CHAOS_INTERVAL
      value: '5'

    ## percentage of total pods to target
    - name: PODS_AFFECTED_PERC
      value: ''

    - name: LIB
      value: 'litmus'    

    - name: TARGET_PODS
      value: ''

    ## it defines the sequence of chaos execution for multiple target pods
    ## supported values: serial, parallel
    - name: SEQUENCE
      value: 'parallel'
    labels:
      name: pod-delete
      app.kubernetes.io/part-of: litmus
      app.kubernetes.io/component: experiment-job
      app.kubernetes.io/version: latest
```

## How to get started?

Refer the LitmusChaos documentation [litmus docs](https://docs.litmuschaos.io)

## How do I contribute?

You can contribute by raising issues, improving the documentation, contributing to the core framework and tooling, etc.

Head over to the [Contribution guide](CONTRIBUTING.md)

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Flitmuschaos%2Fchaos-operator.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Flitmuschaos%2Fchaos-operator?ref=badge_large)
