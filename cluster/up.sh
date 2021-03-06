#!/bin/bash
#
# Copyright 2018-2019 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

source ./cluster/common.sh
source ./cluster/kubevirtci.sh
kubevirtci::install

$(kubevirtci::path)/cluster-up/up.sh

if [[ "$KUBEVIRT_PROVIDER" =~ (ocp|okd)- ]]; then
    echo 'Remove components we do not need to save some resources'
    ./cluster/kubectl.sh delete ns openshift-monitoring --wait=false
    ./cluster/kubectl.sh delete ns openshift-marketplace --wait=false
    ./cluster/kubectl.sh delete ns openshift-cluster-samples-operator --wait=false
fi

ensure_golang

install_cnao

install_cdi
install_kubevirt
install_imageio
