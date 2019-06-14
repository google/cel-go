#!/bin/bash
mkdir -p testlogs
echo $PWD
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  cp -r "$folder" ./testlogs
done
tar -czf ./tests.tar.gz testlogs
#rm -r testlogs
