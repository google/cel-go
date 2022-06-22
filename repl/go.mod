module github.com/google/cel-go/repl

go 1.17

require (
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20220418222510-f25a4f6275ed
	github.com/chzyer/readline v1.5.0
	github.com/google/cel-go v0.11.4
	google.golang.org/genproto v0.0.0-20220614165028-45ed7f3ff16e
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/google/cel-go => ../.
