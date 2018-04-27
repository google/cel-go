workspace(name = "cel_go")

http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.10.1/rules_go-0.10.1.tar.gz",
    sha256 = "4b14d8dd31c6dbaf3ff871adcd03f28c3274e42abc855cb8fb4d01233c0154dc",
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
  commit = "3769a0b59441e6a2ebe747154dd1a4c85ce65ae0",
  remote = "https://github.com/google/cel-spec.git",
)
