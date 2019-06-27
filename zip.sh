#!/bin/bash
#ls
#touch tests.xml
#echo "<?xml version=\"1.0\" encoding=\"UTF-8\"?>" > tests.xml
#echo "<testsuites>" >> tests.xml
#for file in $(find ./bazel-out/k8-fastbuild/testlogs -name '*.xml');do
#  cat $file
mkdir -p artifacts
touch artifacts/junit_01.xml
cat ./bazel-out/k8-fastbuild/testlogs/conformance/ct/test.xml > artifacts/junit_01.xml
#done
#echo "</testsuites>" >> tests.xml
comma='",'
startdate=$(date +%s)
timestamp='"timestamp": "'
time_string="$timestamp$startdate$comma"

result='"result": "'
test_string=$(tail -n 4 artifacts/junit_01.xml | head -n 1 | cut -c1-4)

echo $test_string

if [ $test_string = "PASS" ]
then
  test="SUCCESS"
else
  test="FAILURE"
fi
result_string="$result$test$comma"

echo "{" > finished.json
echo "$time_string" >> finished.json
echo "$result_string" >> finished.json
echo "}" >> finished.json

touch build-log.txt

mkdir $1
mv started.json ./$1/started.json
mv finished.json ./$1/finished.json
mv artifacts ./$1/artifacts
