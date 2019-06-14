#!/bin/bash
#zip -r ./tests.zip ~/go/src/github.com/google/cel-go/bazel-testlogs/*
mkdir -p testlogs
for folder in ~/go/src/github.com/google/cel-go/bazel-testlogs/*;do
  cp -r "$folder" ./testlogs
done
tar -czf ./tests.tar.gz testlogs/

