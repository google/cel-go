#!/bin/bash
mkdir -p artifacts
touch artifacts/junit_01.xml

input="./bazel-out/k8-fastbuild/testlogs/conformance/ct/test.xml"

head -2 $input > artifacts/junit_01.xml
echo "<testsuite>" >> artifacts/junit_01.xml
while IFS= read -r line
do
  testline=$(echo $line | cut -c1-3)
  extendline=$(echo $line | cut -c1-4)
  testcase='<testcase name="'
  testend='">'
  close_testcase="</testcase>"

  if [ $testline = "---" ] && [ $extendline != '----' ]
  then
    status=$(echo $line | cut -c4-8)
    name=$(echo $line | tail -c +11 | head -c -9)
    echo "$testcase$name$testend" >> artifacts/junit_01.xml
    echo $status >> artifacts/junit_01.xml
    if [ $status = "FAIL" ]
    then
      read line1
      echo $line1 >> artifacts/junit_01.xml
    fi
    echo $close_testcase >> artifacts/junit_01.xml
  fi
done < "$input"
tail -2 $input >> artifacts/junit_01.xml

comma=","
quote='"'
startdate=$(date +%s)
timestamp='"timestamp": '
time_string="$timestamp$startdate$comma"

result='"result": "'
test_string=$(tail -n 4 artifacts/junit_01.xml | head -n 1 | cut -c1-4)

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

filedir=$(cat _DATE)

mkdir -p $filedir
mv started.json ./$filedir/started.json
mv finished.json ./$filedir/finished.json
mv build-log.txt ./$filedir/build-log.txt
mv artifacts ./$filedir/artifacts
