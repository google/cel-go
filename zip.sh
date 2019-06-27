#!/bin/bash
mkdir -p artifacts
touch artifacts/junit_01.xml
cat ./bazel-out/k8-fastbuild/testlogs/conformance/ct/test.xml > artifacts/junit_01.xml

quote='"'
comma=","
startdate=$(date +%s)
timestamp='"timestamp": "'
time_string="$timestamp$startdate$quote$comma"

result='"result": "'
test_string=$(tail -n 4 artifacts/junit_01.xml | head -n 1 | cut -c1-4)

echo $test_string

if [ $test_string = "PASS" ]
then
  test="SUCCESS"
else
  test="FAILURE"
fi
result_string="$result$test$quote"

echo "{" > finished.json
echo "$time_string" >> finished.json
echo "$result_string" >> finished.json
echo "}" >> finished.json

touch build-log.txt

mkdir -p $1
mv started.json ./$1/started.json
mv finished.json ./$1/finished.json
mv build-log.txt ./$1/build-log.txt
mv artifacts ./$1/artifacts
