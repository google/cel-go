module github.com/google/cel-go/codelab

go 1.18

require (
	github.com/golang/glog v1.0.0
	github.com/google/cel-go v0.12.5
	google.golang.org/genproto v0.0.0-20230106154932-a12b697841d9
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
)

replace github.com/google/cel-go => ../.
