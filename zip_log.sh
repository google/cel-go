#!/bin/bash
cd ~
echo $PWD
cd .cache
echo $PWD
cd bazel
echo $PWD
ls
cd _bazel_root
echo $PWD
ls
cd eab0d61a99b6696edb3d2aff87b585e8
echo $PWD
ls
cd cel_go
echo $PWD
ls
cd bazel-output
echo $PWD
cd k8-fastbuild
echo $PWD
cd testlogs
echo $PWD
tar -czf ./tests.tar.gz .
