#!/bin/sh

ANTLR_JAR="antlr-4.11.1-complete.jar"
export CLASSPATH=".:$ANTLR_JAR:$CLASSPATH"
antlr4="java -Xmx500M -cp \"$ANTLR_JAR:$CLASSPATH\" org.antlr.v4.Tool"

$antlr4 \
    -Dlanguage=Go \
    -lib ../../parser/gen \
    -package parser \
    -visitor \
    ./Commands.g4
