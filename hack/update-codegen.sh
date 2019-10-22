#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

${GOPATH}/src/github.com/litmuschaos/chaos-operator/vendor/k8s.io/code-generator/generate-groups.sh client,lister,informer \
  github.com/litmuschaos/chaos-operator/pkg/client github.com/litmuschaos/chaos-operator/pkg/apis \
  litmuschaos:v1alpha1 

