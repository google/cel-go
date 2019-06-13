#!/bin/bash
#zip -r ./tests.zip ~/go/src/github.com/google/cel-go/bazel-testlogs/*
tar -czf ./tests.tar.gz ./bazel-testlogs/*
