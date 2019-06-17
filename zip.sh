#!/bin/bash
mkdir -p testlogs
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  cp -r "$folder" testlogs
done
ls -la
tar -czf ./tests.tar.gz ./testlogs
