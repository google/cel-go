module github.com/google/cel-go/repl/appengine

go 1.18

require github.com/google/cel-go/repl v0.0.0-20230406155237-b081aea03865

require (
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/google/cel-go v0.18.1 // indirect
	github.com/google/cel-spec/tests v0.14.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	golang.org/x/exp v0.0.0-20230515195305-f3d0a9c9a5cc // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace github.com/google/cel-go => ../../.

replace github.com/google/cel-go/repl => ../.

replace github.com/google/cel-spec/tests v0.14.0 => github.com/l46kok/cel-spec/tests v0.0.0-20240319231845-95c1c0162ada
