#!/bin/bash
mkdir -p artifacts
touch artifacts/junit_01.xml

input="./bazel-out/k8-fastbuild/testlogs/conformance/ct/test.xml"
while IFS= read -r line
do
  testline=$(echo $line | cut -c1-3)
  testcase='<testcase name="'
  testend='">'
  close_testcase="</testcase>"

  if [ $testline = "---" ]
  then
    status=$(echo $line | cut -c4-8)
    name=$(echo $line | tail -c +11 | head -c -9)
    echo "$testcase$name$testend"
    echo $status
    if [ $status = "FAIL" ]
    then
      read line1
      echo $line1
    fi
    echo $close_testcase
  fi
done < "$input"

#cat ./bazel-out/k8-fastbuild/testlogs/conformance/ct/test.xml > artifacts/junit_01.xml
comma=","
quote='"'
startdate=$(date +%s)
timestamp='"timestamp": '
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
