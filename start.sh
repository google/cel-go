#!/bin/bash
comma='",'
startdate=$(date +%s)
timestamp='"timestamp": "'
time_string="$timestamp$startdate$comma"

pull='"pull": "'
pull_string="$pull$1$comma"

echo "{" > started.json
echo "$time_string" >> started.json
echo "$pull_string" >> started.json
echo "}" >> started.json
