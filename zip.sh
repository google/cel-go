#!/bin/bash
#zip -r ./tests.zip ~/go/src/github.com/google/cel-go/bazel-testlogs/*
mkdir -p testlogs
cp -r ~/go/src/github.com/google/cel-go/bazel-testlogs/* testlogs
tar -czf ./tests.tar.gz testlogs
