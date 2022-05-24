
grammar Types;

import CEL;

start : t=type EOF;

type :
    id=typeId params=typeParamList? ;

typeId :
    leadingDot='.'? id=(IDENTIFIER|NUL) ('.' qualifiers+=IDENTIFIER )* ;

typeParamList:
    '(' types+=type (',' types+=type)* ')' ;