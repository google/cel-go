// Package main declares the executable entry point for the CEL server.
package main

import (
	"github.com/google/cel-go/server"
	"github.com/google/cel-spec/tools/celrpc"
)

func main() {
	celrpc.RunServer(&server.ConformanceServer{})
}
