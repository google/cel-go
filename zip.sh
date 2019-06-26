#!/bin/bash
touch tests.xml
echo "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" > tests.xml
echo "<testsuites>" >> tests.xml
for file in $(find ./bazel-out/k8-fastbuild/testlogs -name '*.xml');do
  cat $file
  sed -e'1d;2d;$d' $file >> tests.xml
done
echo "</testsuites>" >> tests.xml
date +%s >> build-log.txt
$COMMIT_SHA >> build-log.txt
