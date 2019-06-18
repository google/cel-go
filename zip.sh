#!/bin/bash
touch tests.xml
echo "<?xml version="1.0" encoding="UTF-8"?>" >> tests.xml
echo "<testsuites>" >> tests.xml
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  for folder2 in "$folder"/*;do
    for file in "$folder2"/*;do
      echo "${file##*/}"
      if [ "${file##*/}" = "test.xml" ]
      then
          sed -e'2d' "$file" >> tests.xml
          cat tests.xml
      fi
    done
  done
done
echo "</testsuites>" >> tests.xml
#ls -la
#tar -czf ./tests.tar.gz ./testlogs
