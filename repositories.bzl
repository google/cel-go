# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""repositories for loading non bzlmod dependencies"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def cel_policy_dependency():
    cel_policy_tag = "569292f1c4eaa41894c1e37ee94eb146e284bcfa"
    cel_policy_sha = "5a68318d906f6ce18492ad6f82b5f8bb083fd9d694cf567d399216c11da03157"
    http_archive(
        name = "cel_policy",
        sha256 = cel_policy_sha,
        strip_prefix = "cel-policy-%s" % cel_policy_tag,
        url = "https://github.com/cel-expr/cel-policy/archive/%s.tar.gz" % cel_policy_tag,
    )

def _non_module_dependencies_impl(_ctx):
    cel_policy_dependency()

non_module_dependencies = module_extension(
    implementation = _non_module_dependencies_impl,
)
