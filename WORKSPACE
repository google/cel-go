workspace(name = "cel_go")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "099a9fb96a376ccbbb7d291ed4ecbdfd42f6bc822ab77ae6f1b5cb9e914e94fa",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.35.0/rules_go-v0.35.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "62ca106be173579c0a167deb23358fdfe71ffa1e4cfdddf5582af26520f1c66f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
    ],
)

# Latest version as of 10/27/2020
http_archive(
    name = "rules_proto",
    sha256 = "e017528fd1c91c5a33f15493e3a398181a9e821a804eb7ff5acdd1d2d6c2b18d",
    strip_prefix = "rules_proto-4.0.0-3.20.0",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/refs/tags/4.0.0-3.20.0.tar.gz",
    ],
)

# Latest version of googleapis as of 7/27/2020
http_archive(
    name = "com_google_googleapis",
    sha256 = "1f742f6cafe616fe73302db010e0b7ee6579cb1ce06010427b7d0995cbd80ce4",
    strip_prefix = "googleapis-6a813acf535e4746fa4a135ce23547bb6425c26d",
    urls = [
        "https://github.com/googleapis/googleapis/archive/6a813acf535e4746fa4a135ce23547bb6425c26d.tar.gz",
    ],
)

# protobuf
http_archive(
    name = "com_google_protobuf",
    sha256 = "8242327e5df8c80ba49e4165250b8f79a76bd11765facefaaecfca7747dc8da2",
    strip_prefix = "protobuf-3.21.5",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.21.5.zip"],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("@com_google_googleapis//:repository_rules.bzl", "switched_rules_by_language")
load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
    go = True,
)

# Do *not* call *_dependencies(), etc, yet.  See comment at the end.

go_repository(
    name = "org_golang_google_protobuf",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/protobuf",
    tag = "v1.28.1",
)

# Generated Google APIs protos for Golang 11/16/2020
go_repository(
    name = "org_golang_google_genproto",
    build_file_proto_mode = "disable_global",
    commit = "69f6226f97e558dbaa68715071622af0d86b3a17",
    importpath = "google.golang.org/genproto",
)

# gRPC deps for v1.49.0 (including x/text and x/net)
go_repository(
    name = "org_golang_google_grpc",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/grpc",
    tag = "v1.49.0",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:oWX7TPOiFAMXLq8o0ikBYfCJVlRHBcsciT5bXOrH628=",
    version = "v0.0.0-20190311183353-d8887717615a",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:olpwvP2KacW1ZWvsR7uQhoyTYvKAupfQrRGBFM352Gk=",
    version = "v0.3.7",
)

# ANTLR v4.11.1
go_repository(
    name = "com_github_antlr_antlr4_runtime_go_antlr",
    importpath = "github.com/antlr/antlr4/runtime/Go/antlr/v4",
    sum = "h1:et5J11AOyUn9qwkIAF9kcxTxjTO8Z9oSmlOqH7MVSPo=",
    version = "v4.0.0-20220911224424-aa1f1f12a846",
)

# CEL Spec deps
go_repository(
    name = "com_google_cel_spec",
    commit = "1ce60e0b12de519b6792621eb687fe31926927eb",
    importpath = "github.com/TristonianJones/cel-spec",
)

# strcase deps
go_repository(
    name = "com_github_stoewer_go_strcase",
    importpath = "github.com/stoewer/go-strcase",
    sum = "h1:Z2iHWqGXH00XYgqDmNgQbIBxf3wrNq0F3feEy0ainaU=",
    version = "v1.2.0",
)

# Readline for repl
go_repository(
    name = "com_github_chzyer_readline",
    importpath = "github.com/chzyer/readline",
    commit = "62c6fe6193755f722b8b8788aa7357be55a50ff1"  # v1.4
)

# golang.org/x/exp deps
go_repository(
    name = "org_golang_x_exp",
    importpath = "golang.org/x/exp",
    sum = "h1:+WEEuIdZHnUeJJmEUjyYC2gfUMj69yZXw17EnHg/otA=",
    version = "v0.0.0-20220722155223-a9213eeb770e"
)

# Run the dependencies at the end.  These will silently try to import some
# of the above repositories but at different versions, so ours must come first.
go_rules_dependencies()

go_register_toolchains(version = "1.19.1")

gazelle_dependencies()

rules_proto_dependencies()

rules_proto_toolchains()

protobuf_deps()
