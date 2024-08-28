module github.com/google/cel-go/repl/appengine

go 1.21

toolchain go1.23.0

require github.com/google/cel-go/repl v0.0.0-20230406155237-b081aea03865

require (
	cel.dev/expr v0.16.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/google/cel-go v0.0.0-00010101000000-000000000000 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240826202546-f6391c0de4c7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240826202546-f6391c0de4c7 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/google/cel-go => ../../.

replace github.com/google/cel-go/repl => ../.
