#!/bin/sh

ANTLR_JAR="/home/jdtatum/bin/antlr-4.10.1-complete.jar"
antlr4="java -Xmx500M -cp \"$ANTLR_JAR:$CLASSPATH\" org.antlr.v4.Tool"

$antlr4 \
    -Dlanguage=Go \
    -lib ../../parser/gen/ \
    -visitor \
    -package parser \
    ./Types.g4
