#!/bin/bash
#cd ../../../../../../../../../../
#cd ./execroot/cel_go/
mkdir -p testlogs
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  cp -r "$folder" testlogs
done
ls -la
tar -czf ./tests.tar.gz ./testlogs
#rm -r testlogs
