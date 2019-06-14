#!/bin/bash
#zip -r ./tests.zip ~/go/src/github.com/google/cel-go/bazel-testlogs/*
mkdir -p testlogs
for folder in ./bazel-testlogs/*;do
  cp -r "$folder" ./testlogs
done
tar -czf ./tests.tar.gz testlogs
#rm -r testlogs
