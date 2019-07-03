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
  end='">'
  close_testcase="</testcase>"

  failure='<failure message="'
  close_failure="</failure>"

  if [ $testline = "---" ] && [ $extendline != '----' ]
  then
    status=$(echo $line | cut -c4-8)
    name=$(echo $line | tail -c +11 | head -c -9)
    echo "$testcase$name$end" >> artifacts/junit_01.xml
    if [ $status = "FAIL" ]
    then
      read line1
      message=$(echo "$line1" | sed 's/</ /g; s/>/ /g; s/"//g')
      echo "$failure$message$end$close_failure" >> artifacts/junit_01.xml
    else
      echo $status >> artifacts/junit_01.xml
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
test_string=$(tail -n 5 $input)

if [[ $test_string = *"FAIL"* ]]
then
  test="FAILURE"
else
  test="SUCCESS"
fi
result_string="$result$test$quote"

filedir=$(cat _DATE)
mkdir -p $filedir

touch $filedir/build-log.txt

echo "{" > $filedir/finished.json
echo "$time_string" >> $filedir/finished.json
echo "$result_string" >> $filedir/finished.json
echo "}" >> $filedir/finished.json

mv started.json $filedir/started.json
mv artifacts $filedir/artifacts
