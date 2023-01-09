module github.com/google/cel-go/server

go 1.18

require (
	github.com/google/cel-go v0.12.5
	github.com/google/cel-spec v0.7.1
	google.golang.org/genproto v0.0.0-20230106154932-a12b697841d9
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/bazelbuild/rules_go v0.35.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/net v0.0.0-20221014081412-f15817d10f9b // indirect
	golang.org/x/sys v0.0.0-20220909162455-aba9fc2a8ff2 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/grpc v1.51.0 // indirect
)

replace github.com/google/cel-go => ./..
