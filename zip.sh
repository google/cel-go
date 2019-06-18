#!/bin/bash
touch tests.xml
echo "<?xml version="1.0" encoding="UTF-8"?>" >> tests.xml
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  for folder2 in "$folder"/*;do
    for file in "$folder2"/*;do
      echo "${file##*/}"
      if [ "${file##*/}" = "test.xml" ]
      then
          sed -e'1d' "$file" >> tests.xml
      fi
    done
  done
done
#ls -la
#tar -czf ./tests.tar.gz ./testlogs
