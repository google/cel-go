module github.com/google/cel-go/conformance

go 1.21.1

require (
	cel.dev/expr v0.18.0
	github.com/bazelbuild/rules_go v0.49.0
	github.com/google/cel-go v0.21.0
	github.com/google/go-cmp v0.6.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240826202546-f6391c0de4c7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240826202546-f6391c0de4c7 // indirect
)

replace github.com/google/cel-go => ./..
