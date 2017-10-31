#!/bin/sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Generate AntLR artifacts.
java -Xmx500M -cp ${DIR}/antlr-4.7-complete.jar org.antlr.v4.Tool  \
    -Dlanguage=Go \
    -package gen \
    -o ${DIR}/../parser/gen \
    -visitor ${DIR}/../parser/gen/CEL.g4

