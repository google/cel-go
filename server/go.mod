module github.com/google/cel-go/server

go 1.18

require (
	github.com/google/cel-go v0.12.5
	github.com/google/cel-spec v0.7.1
	google.golang.org/genproto v0.0.0-20220909194730-69f6226f97e5
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/bazelbuild/rules_go v0.35.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/net v0.0.0-20220909164309-bea034e7d591 // indirect
	golang.org/x/sys v0.0.0-20220909162455-aba9fc2a8ff2 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/grpc v1.49.0 // indirect
)

replace github.com/google/cel-go => ./..
