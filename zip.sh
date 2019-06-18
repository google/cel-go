#!/bin/bash
touch tests.xml
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  echo "$folder"
  for folder2 in "$folder"/*;do
    echo "$folder2"
    for file in "$folder2";do
      echo "$file"
      if [ "$file" = "test.xml" ]
      then
          cat "$file" >> tests.xml
      fi
    done
  done
done
#ls -la
#tar -czf ./tests.tar.gz ./testlogs
