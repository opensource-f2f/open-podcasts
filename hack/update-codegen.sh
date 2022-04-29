#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# 注意:
# 1. kubebuilder2.3.2版本生成的api目录结构code-generator无法直接使用(将api由api/${VERSION}移动至api/${GROUP}/${VERSION}即可)

# corresponding to go mod init <module>
MODULE=github.com/opensource-f2f/open-podcasts
# api package
APIS_PKG=api
# generated output package
OUTPUT_PKG=generated
# group-version such as foo:v1alpha1
GROUP=osf2f.my.domain
VERSION=v1alpha1
GROUP_VERSION=${GROUP}:${VERSION}

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}

rm -rf ${OUTPUT_PKG}/{clientset,informers,listers}

# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
#bash "${CODEGEN_PKG}"/generate-groups.sh "client,informer,lister" \
bash ./hack/generate-groups.sh all \
  ${MODULE}/${OUTPUT_PKG} ${MODULE}/${APIS_PKG} \
  ${GROUP_VERSION} \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt \
  --output-base "${SCRIPT_ROOT}"
#  --output-base "${SCRIPT_ROOT}/../../.."

rm -rf generated && mv ${MODULE}/generated . && rm -rf github.com
