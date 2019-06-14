#!/bin/bash
cd ../../../../../../../../../../
cd ./execroot/cel_go/
#echo $TEST_TMPDIR
mkdir -p $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  cp -r "$folder" $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
done
echo $PWD
cd $TEST_UNDECLARED_OUTPUTS_DIR
tar -czf ./tests.tar.gz ./testlogs
#cp tests.tar.gz ~/go/src/github.com/google/cel-go/
#rm -r testlogs
