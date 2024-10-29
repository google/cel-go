module github.com/google/cel-go/repl

go 1.21.1

require (
	cel.dev/expr v0.17.0
	github.com/antlr4-go/antlr/v4 v4.13.0
	github.com/chzyer/readline v1.5.1
	github.com/google/cel-go v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/api v0.0.0-20240826202546-f6391c0de4c7
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240826202546-f6391c0de4c7
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace github.com/google/cel-go => ../.

replace cel.dev/expr => ../../cel-spec
