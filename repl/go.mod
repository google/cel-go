module github.com/google/cel-go/repl

go 1.18

require (
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230321174746-8dcc6526cfb1
	github.com/chzyer/readline v1.5.1
	github.com/google/cel-go v0.18.1
	google.golang.org/genproto/googleapis/api v0.0.0-20230803162519-f966b187b2e5
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/exp v0.0.0-20230321023759-10a507213a29 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
)

replace github.com/google/cel-go => ../.
