workspace(name = "celgo")

git_repository(
    name = "io_bazel_rules_go",
    commit = "9cf23e2aab101f86e4f51d8c5e0f14c012c2161c",  # Oct 12, 2017 (Add `build_external` option to `go_repository`)
    remote = "https://github.com/bazelbuild/rules_go.git",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")

proto_register_toolchains()

load("@io_bazel_rules_go//go:def.bzl", "go_repository")

git_repository(
    name = "io_bazel_rules_docker",
    commit = "9dd92c73e7c8cf07ad5e0dca89a3c3c422a3ab7d",  # Sep 27, 2017 (v0.3.0)
    remote = "https://github.com/bazelbuild/rules_docker.git",
)

load(":adapter_author_deps.bzl", "mixer_adapter_repositories")

mixer_adapter_repositories()


go_repository(
    name = "com_github_antlr",
    commit = "d432f94d3b01e713acce655c04a98524e405b3df",
    importpath = "github.com/antlr/antlr4",
)

load("//:googleapis.bzl", "go_googleapis_repositories")
load("//:istio_api.bzl", "go_istio_api_repositories")

go_googleapis_repositories()

go_istio_api_repositories(False)
