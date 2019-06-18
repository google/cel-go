#!/bin/bash
touch tests.xml
for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
  echo $PWD
  for folder2 in "$folder"/*;do
    echo $PWD
    for file in "$folder2";do
      if ["$file" = "test.xml"]
      then
          cat "$file" >> tests.xml
      fi
    done
  done
done
#ls -la
#tar -czf ./tests.tar.gz ./testlogs
