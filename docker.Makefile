# Copyright 2014-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You
# may not use this file except in compliance with the License. A copy of
# the License is located at
#
#       http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is
# distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF
# ANY KIND, either express or implied. See the License for the specific
# language governing permissions and limitations under the License.
#

# build binary

# unit-test
unit-test:
	go test -v -cover -race -timeout 300s ./pkg/awsutils/...
	go test -v -cover -race -timeout 10s ./plugins/routed-eni/...
	go test -v -cover -race -timeout 10s ./plugins/routed-eni/driver
	go test -v -cover -race -timeout 10s ./pkg/k8sapi/...
	go test -v -cover -race -timeout 10s ./pkg/networkutils/...
	go test -v -cover -race -timeout 90s ./ipamd/...
