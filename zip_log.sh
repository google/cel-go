#!/bin/bash
cd ~
echo $PWD
cd .cache
cd bazel
tar -czf ./tests.tar.gz ./_bazel_root/eab0d61a99b6696edb3d2aff87b585e8/execroot/cel_go/bazel-output/k8-fastbuild/testlogs


