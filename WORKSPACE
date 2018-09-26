workspace(name = "cel_go")

# Must come before go_rules_dependencies()
#git_repository(
#  name = "org_golang_google_genproto",
#  commit = "221a8d4f74948678f06caaa13c9d41d22e069ae8",
#  remote = "https://github.com/google/go-genproto.git",
#  build_file_content = """
#go_library(
#  name = "googleapis_api_expr_v1alpha1",
#  srcs = [
#  ],
#  importpath = "",
#  visibility = ["//visibility:public"],
#  deps = [
#  ],
#)
#"""
#)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "8b68d0630d63d95dacc0016c3bb4b76154fe34fca93efd65d1c366de3fcb4294",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.12.1/rules_go-0.12.1.tar.gz"],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "ddedc7aaeb61f2654d7d7d4fd7940052ea992ccdb031b8f9797ed143ac7e8d43",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.12.0/bazel-gazelle-0.12.0.tar.gz"],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "org_golang_google_genproto",
    build_file_proto_mode = "disable",
  #  commit = "0e822944c569bf5c9afd034adaa56208bd2906ac",
  commit = "221a8d4f74948678f06caaa13c9d41d22e069ae8",
    importpath = "google.golang.org/genproto",
)

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

#http_archive(
#    name = "io_bazel_rules_go",
#    url = "https://github.com/bazelbuild/rules_go/releases/download/0.12.0/rules_go-0.12.0.tar.gz",
#    sha256 = "c1f52b8789218bb1542ed362c4f7de7052abcf254d865d96fb7ba6d44bc15ee3",
#)
#load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
#go_rules_dependencies()
#go_register_toolchains()

#load("@io_bazel_rules_go//go:def.bzl", "go_repository")
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
