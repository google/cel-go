#!/bin/bash
exec bazel test ...
cd bazel-testlogs
zip -r . ./*
