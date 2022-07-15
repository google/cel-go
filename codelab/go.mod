module github.com/google/cel-go/codelab

go 1.17

require (
	github.com/golang/glog v1.0.0
	github.com/google/cel-go v0.12.4
	google.golang.org/genproto v0.0.0-20220714152414-ccd2914cffd4
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20220626175859-9abda183db8e // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/google/cel-go => ../.
