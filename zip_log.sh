#!/bin/bash
echo workspace
ls -la
echo bazel-out
cd bazel-out
ls -la
echo k8
cd k8-fastbuild
ls -la
echo tests
cd testlogs
ls -la
echo zip
cd zip_test
ls -la
cd test.outputs
ls -la

