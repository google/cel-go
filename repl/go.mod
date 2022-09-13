module github.com/google/cel-go/repl

go 1.18

require (
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10
	github.com/chzyer/readline v1.5.1
	github.com/google/cel-go v0.12.5
	google.golang.org/genproto v0.0.0-20220909194730-69f6226f97e5
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/sys v0.0.0-20220624220833-87e55d714810 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/google/cel-go => ../.
