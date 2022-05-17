module github.com/google/cel-go/repl

go 1.15

require (
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20220418222510-f25a4f6275ed
	github.com/chzyer/readline v1.5.0
	github.com/google/cel-go v0.11.4
	google.golang.org/genproto v0.0.0-20220502173005-c8bf987b8c21
	google.golang.org/protobuf v1.28.0
)

replace github.com/google/cel-go => ../.
