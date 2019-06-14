#!/bin/bash
#cd ../../../../../../../../../../
#cd ./execroot/cel_go/
echo $PWD
mkdir -p $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  cp -r "$folder" $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
done
echo $PWD
cd ../../../../../../../../../../
echo $PWD
tar -czf ./tests.tar.gz $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
#cp tests.tar.gz ~/.cache/bazel/_bazel_$USER
#rm -r testlogs
