#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

vendor/k8s.io/code-generator/generate-groups.sh \
  "defaulter,client,lister,informer" \
  "github.com/goodrain/servicemesh-operator/pkg/generated" \
  "github.com/goodrain/servicemesh-operator/pkg/api" \
  "rainbond:v1alpha1" \
  --go-header-file "./hack/boilerplate.go.txt" \
  $@
