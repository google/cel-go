workspace(name = "cel_go")

http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.12.0/rules_go-0.12.0.tar.gz",
    sha256 = "c1f52b8789218bb1542ed362c4f7de7052abcf254d865d96fb7ba6d44bc15ee3",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@io_bazel_rules_go//go:def.bzl", "go_repository")
go_repository(
  name = "com_github_antlr",
  commit = "763a1242b7f5fca2c7a06f671ebe757580dacfb2",
  importpath = "github.com/antlr/antlr4",
)

git_repository(
  name = "com_google_cel_spec",
  commit = "3c25c4d4ffb504e2c24c1d84196c78ba3ac9e612",
  remote = "https://github.com/google/cel-spec.git",
)

new_http_archive(
    name = "com_google_googleapis",
    url = "https://github.com/googleapis/googleapis/archive/common-protos-1_3_1.zip",
    strip_prefix = "googleapis-common-protos-1_3_1/",
    build_file_content = """
load('@io_bazel_rules_go//proto:def.bzl', 'go_proto_library')

proto_library(
    name = 'rpc_status',
    srcs = ['google/rpc/status.proto'],
    deps = [
        '@com_google_protobuf//:any_proto',
        '@com_google_protobuf//:empty_proto'
    ],
    visibility = ['//visibility:public'],
)

go_proto_library(
    name = 'rpc_status_go_proto',
    # TODO: Switch to the correct import path when bazel rules fixed.
    #importpath = 'google.golang.org/genproto/googleapis/rpc/status',
    importpath = 'github.com/googleapis/googleapis/google/rpc',
    proto = ':rpc_status',
    visibility = ['//visibility:public'],
)
"""
)

git_repository(
  name = "org_pubref_rules_protobuf",
  remote = "https://github.com/pubref/rules_protobuf",
  tag = "v0.8.2",
)

load("@org_pubref_rules_protobuf//go:rules.bzl", "go_proto_repositories")
go_proto_repositories()

go_repository(
  name = "org_golang_google_grpc",
  importpath = "google.golang.org/grpc",
  tag = "v1.11.3",
  remote = "https://github.com/grpc/grpc-go.git",
)
