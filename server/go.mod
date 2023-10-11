module github.com/google/cel-go/server

go 1.18

require (
	github.com/google/cel-go v0.13.0
	github.com/google/cel-spec v0.12.0
	google.golang.org/genproto/googleapis/api v0.0.0-20230803162519-f966b187b2e5
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230305170008-8188dc5388df // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/stoewer/go-strcase v1.2.0 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20230726155614-23370e0ffb3e // indirect
	google.golang.org/grpc v1.57.0 // indirect
)

replace github.com/google/cel-go => ./..
