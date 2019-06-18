#!/bin/bash
touch tests.xml
echo "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" > tests.xml
echo "<testsuites>" >> tests.xml
for file in $(find ./bazel-out -name '*.xml');do
  sed -e'1d;2d;$d' $file >> tests.xml
done

#for folder in ./bazel-out/k8-fastbuild/testlogs/*;do
#  for folder2 in "$folder"/*;do
#    for file in "$folder2"/*;do
#      if [ "${file##*/}" = "test.xml" ]
#      then
#          sed -e'1d;2d;$d' "$file" >> tests.xml
#      fi
#    done
#  done
#done
echo "</testsuites>" >> tests.xml
cat tests.xml
#tar -czf ./tests.tar.gz ./testlogs
