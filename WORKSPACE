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
    sha256 = "ecba0f04f96b4960a5b250c8e8eeec42281035970aa8852dda73098274d14a1d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.29.0/bazel-gazelle-v0.29.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.29.0/bazel-gazelle-v0.29.0.tar.gz",
    ],
)

# rules_proto v3.20.0
http_archive(
    name = "rules_proto",
    sha256 = "e017528fd1c91c5a33f15493e3a398181a9e821a804eb7ff5acdd1d2d6c2b18d",
    strip_prefix = "rules_proto-4.0.0-3.20.0",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/refs/tags/4.0.0-3.20.0.tar.gz",
    ],
)

# googleapis as of 05/26/2023
http_archive(
    name = "com_google_googleapis",
    strip_prefix = "googleapis-07c27163ac591955d736f3057b1619ece66f5b99",
    sha256 = "bd8e735d881fb829751ecb1a77038dda4a8d274c45490cb9fcf004583ee10571",
    urls = [
        "https://github.com/googleapis/googleapis/archive/07c27163ac591955d736f3057b1619ece66f5b99.tar.gz",
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

# Generated Google APIs protos for Golang 05/25/2023
go_repository(
    name = "org_golang_google_genproto_googleapis_api",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/genproto/googleapis/api",
    sum = "h1:m8v1xLLLzMe1m5P+gCTF8nJB9epwZQUBERm20Oy1poQ=",
    version = "v0.0.0-20230525234035-dd9d682886f9",
)

# Generated Google APIs protos for Golang 05/25/2023
go_repository(
    name = "org_golang_google_genproto_googleapis_rpc",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/genproto/googleapis/rpc",
    sum = "h1:0nDDozoAU19Qb2HwhXadU8OcsiO/09cnTqhUtq2MEOM=",
    version = "v0.0.0-20230525234030-28d5490b6b19",
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

# ANTLR v4.12.0
go_repository(
    name = "com_github_antlr_antlr4_runtime_go_antlr_v4",
    importpath = "github.com/antlr/antlr4/runtime/Go/antlr/v4",
    sum = "h1:7RFfzj4SSt6nnvCPbCqijJi1nWCd+TqAT3bYCStRC18=",
    version = "v4.0.0-20230305170008-8188dc5388df",
)

# CEL Spec deps v0.9.0
go_repository(
    name = "com_google_cel_spec",
    importpath = "github.com/google/cel-spec",
    commit = "51af45e2b75a8aa2b3108b00f0e91cd172cfbea1",
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
