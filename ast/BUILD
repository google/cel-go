load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "base.go",
        "call.go",
        "comprehension.go",
        "constants.go",
        "debug.go",
        "error.go",
        "expression.go",
        "ident.go",
        "list.go",
        "message.go",
        "select.go",
        "struct.go",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//common:go_default_library",
    ],
)
