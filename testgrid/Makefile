# Copyright 2018 The Knative Authors
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

TESTGRID_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

update-config:

# If recieving GOOGLE_APPLICATION_CREDENTIALS error, 
# run "export GOOGLE_APPLICATION_CREDENTIALS="path/to/google/application/credentials/json/file".
# This needs to be run during every session.

ifndef GOOGLE_APPLICATION_CREDENTIALS
	$(error GOOGLE_APPLICATION_CREDENTIALS not set)

endif
	bazel run @k8s//testgrid/cmd/configurator -- \
		--yaml=$(TESTGRID_DIR)/config.yaml \
		--output=gs://travis-content/config \
		--oneshot
