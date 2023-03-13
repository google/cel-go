module github.com/google/cel-go/repl

go 1.18

require (
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230305170008-8188dc5388df
	github.com/chzyer/readline v1.5.1
	github.com/google/cel-go v0.13.0
	google.golang.org/genproto v0.0.0-20230306155012-7f2fa6fef1f4
	google.golang.org/protobuf v1.29.0
)

require (
	github.com/stoewer/go-strcase v1.2.1 // indirect
	golang.org/x/exp v0.0.0-20230310171629-522b1b587ee0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
)

replace github.com/google/cel-go => ../.
