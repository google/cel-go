#!/bin/bash
cd ../../../../../../../../../../
cd ./execroot/cel_go/
mkdir -p $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  cp -r "$folder" $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
done
cd $TEST_UNDECLARED_OUTPUTS_DIR
ls -la
#tar -czf ./tests.tar.gz ./testlogs
#rm -r testlogs
