#!/bin/bash
echo $TEST_TMPDIR
cd ../../../../../../../../../../
cd ./execroot/cel_go/
echo $PWD
mkdir -p $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  cp -r "$folder" $TEST_UNDECLARED_OUTPUTS_DIR/testlogs
done
cd $TEST_UNDECLARED_OUTPUTS_DIR
#echo $PWD
#cd ../../../../../../
#echo $PWD
#cd .cache
#echo $PWD
#cd bazel
#echo $PWD
#cd _bazel_$USER
echo $PWD
tar -czf ./tests.tar.gz ./testlogs
#cp tests.tar.gz ~/.cache/bazel/_bazel_$USER
rm -r testlogs
