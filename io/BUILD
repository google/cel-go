load("@io_bazel_rules_go//go:def.bzl", "go_library")

package(
    default_visibility = ["//parser:__subpackages__"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "serializer.go",
    ],
    deps = [
        "//ast:go_default_library",
        "//common:go_default_library",
        "//proto:go_default_library",
        "@com_github_antlr//runtime/Go/antlr:go_default_library",
    ],
)
