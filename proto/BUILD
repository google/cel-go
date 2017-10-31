package(default_visibility = ["//visibility:public"])

load("@org_pubref_rules_protobuf//gogo:rules.bzl", "gogoslick_proto_library")

load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_library")

go_prefix("github.com/googleapis/googleapis")

gogoslick_proto_library(
    name = "go_default_library",
    importmap = {
        "gogoproto/gogo.proto": "github.com/gogo/protobuf/gogoproto",
        "google/protobuf/any.proto": "github.com/gogo/protobuf/types",
        "google/protobuf/empty.proto": "github.com/gogo/protobuf/types",
        "google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
        "google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
        "google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
    },
    imports = [
        "external/com_github_google_protobuf/src",
    ],
    inputs = [
        "@com_github_google_protobuf//:well_known_protos",
        "@com_github_gogo_protobuf//gogoproto:go_default_library_protos",
    ],
    protos = [
         "checked.proto",
         "syntax.proto",
    ],
    verbose = 0,
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//sortkeys:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
    ],
)
