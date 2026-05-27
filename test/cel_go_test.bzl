# Copyright 2025 Google LLC
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

"""
Starlark rule for triggering unit tests on CEL policies and expressions for go
runtime target.
"""

load("@bazel_skylib//lib:paths.bzl", "paths")
load("@io_bazel_rules_go//go:def.bzl", "go_test")

def _is_label(s):
    return s.startswith("//") or s.startswith(":") or s.startswith("@")

def cel_go_test(
        name,
        test_suite = "",
        config = "",
        test_src = "",
        cel_expr = "",
        is_raw_expr = False,
        base_config = "",
        test_data_path = "",
        file_descriptor_set = "",
        filegroup = "",
        enable_coverage = False,
        deps = [],
        data = [],
        **kwargs):
    """trigger tests for a CEL checked expression

    This macro wraps an invocation of the go_test rule and is used to trigger tests for a CEL
    checked expression.

    Args:
      name: str name for the generated artifact
      test_suite: str label of a file containing a cel.expr.conformance.test.TestSuite message in
          Textproto format. Alternatively, the file can also contain a test.Suite object in YAML format.
      config: str label of a file containing a google.api.expr.conformance.Environment message in
          Textproto format. Alternatively, the file can also contain the environment configuration represented
          as an env.Config object in YAML format.
      base_config: str label of a file containing a google.api.expr.conformance.Environment message
          in Textproto format. Alternatively, the file can also contain the environment configuration represented
          as an env.Config object in YAML format.
      test_src: source file containing the invocaton of the test runner.
      cel_expr: str label of a file containing either a CEL policy or a CEL expression or a checked
          expression or a raw CEL expression string.
      is_raw_expr: boolean indicating if the cel_expr is a raw CEL expression string.
      test_data_path: path of the directory containing the test files. This is needed only
          if the test files are not located in the same directory as the BUILD file.
      file_descriptor_set: str label or filename pointing to a file_descriptor_set message. Note:
          this must be in binary format with either a .binarypb or .pb or.fds extension. If you need
          to support a textformat file_descriptor_set, embed it in the environment file. (default None)
      filegroup: str label of a filegroup containing the test suite, config, and cel_expr.
      enable_coverage: boolean indicating if coverage should be enabled for the test.
      deps: list of dependencies for the go_test rule
      data: list of data dependencies for the go_test rule
      **kwargs: additional arguments to pass to the go_test rule
    """

    # Avoid mutating the original data list passed into the macro
    resolved_data = list(data)
    resolved_deps = list(deps)

    # Normalize paths
    pkg_name = native.package_name()
    test_data_dir = test_data_path.lstrip("/") if test_data_path else pkg_name

    # Add filegroup if provided
    if filegroup:
        resolved_data.append(filegroup)

    args = [
        "--enable_coverage=%s" % enable_coverage,
    ]

    def _process_file_arg(file_val, flag_name):
        """Helper to append args and resolve data targets for file inputs."""
        if not file_val:
            return

        if _is_label(file_val):
            args.append("--{}=$(rlocationpath {})".format(flag_name, file_val))
            resolved_data.append(file_val)
        else:
            args.append("--{}={}/{}".format(flag_name, test_data_dir, file_val))
            if not filegroup:
                target = file_val if test_data_dir == pkg_name else "//{}:{}".format(test_data_dir, file_val)
                resolved_data.append(target)

    _process_file_arg(test_suite, "test_suite_path")
    _process_file_arg(config, "config_path")
    _process_file_arg(base_config, "base_config_path")

    # Process cel_expr (has specialized fallback logic)
    _, cel_expr_format = paths.split_extension(cel_expr)
    is_valid_cel_ext = cel_expr_format in [".cel", ".celpolicy", ".yaml"]

    if _is_label(cel_expr):
        args.append("--cel_expr=$(rlocationpath {})".format(cel_expr))
        resolved_data.append(cel_expr)
    elif is_raw_expr:
        args.append("--cel_expr=%s" % cel_expr)
    elif is_valid_cel_ext:
        args.append("--cel_expr={}/{}".format(test_data_dir, cel_expr))
        if not filegroup:
            target = cel_expr if test_data_dir == pkg_name else "//{}:{}".format(test_data_dir, cel_expr)
            resolved_data.append(target)
    else:
        args.append("--cel_expr=$(rlocationpath {})".format(cel_expr))
        resolved_data.append(cel_expr)

    if file_descriptor_set:
        if _is_label(file_descriptor_set):
            resolved_data.append(file_descriptor_set)
            args.append("--file_descriptor_set=$(rlocationpath {})".format(file_descriptor_set))
        else:
            target = file_descriptor_set if test_data_dir == pkg_name else "//{}:{}".format(test_data_dir, file_descriptor_set)
            resolved_data.append(target)
            args.append("--file_descriptor_set=$(rlocationpath {})".format(target))

    go_test(
        name = name,
        srcs = [test_src],
        args = args,
        data = resolved_data,
        deps = resolved_deps + [
            "//cel:go_default_library",
            "//common/types:go_default_library",
            "//common/types/ref:go_default_library",
            "//interpreter:go_default_library",
            "//test:go_default_library",
            "//tools/celtest:go_default_library",
            "//tools/compiler:go_default_library",
            "@com_github_google_go_cmp//cmp:go_default_library",
            "@dev_cel_expr//:expr",
            "@dev_cel_expr//conformance/test:go_default_library",
            "@in_yaml_go_yaml_v3//:go_default_library",
            "@io_bazel_rules_go//go/runfiles",
            "@org_golang_google_genproto_googleapis_api//expr/v1alpha1:go_default_library",
            "@org_golang_google_protobuf//encoding/prototext:go_default_library",
            "@org_golang_google_protobuf//proto:go_default_library",
            "@org_golang_google_protobuf//reflect/protodesc:go_default_library",
            "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
            "@org_golang_google_protobuf//reflect/protoregistry:go_default_library",
            "@org_golang_google_protobuf//testing/protocmp:go_default_library",
            "@org_golang_google_protobuf//types/descriptorpb:go_default_library",
            "@org_golang_google_protobuf//types/dynamicpb:go_default_library",
        ],
        **kwargs
    )

