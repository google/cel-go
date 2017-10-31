load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "declaration.go",
        "function.go",
        "ident.go",
        "reference.go",
        "semantics.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//ast:go_default_library",
        "//common:go_default_library",
        "//semantics/types:go_default_library",
    ],
)
