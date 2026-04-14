// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

grammar Commands;

import CEL;

// parser rules:
startCommand: command EOF;

command: help |
         let |
         declare |
         delete |
         simple |
         compile |
         parse |
         exprCmd |
         empty;

help: '%help' | '%?';

let: '%let' ( (var=varDecl '=') | (fn=fnDecl '->') ) e=expr;

declare: '%declare' (var=varDecl | fn=fnDecl);

varDecl: id=qualId (':' t=type)?;

fnDecl: id=qualId '(' (params+=param (',' params+=param)*)? ')' ':' rType=type;

param: pid=IDENTIFIER ':' t=type;

delete: '%delete' (var=varDecl | fn=fnDecl);

simple: cmd=COMMAND (args+=FLAG | args+=STRING)*;

empty: ;

compile: '%compile' e=expr;

parse: '%parse' e=expr;

exprCmd: ('%eval' (flags+=FLAG (flags+=FLAG)* '--')?)? e=expr;

qualId: leadingDot='.'? rid=IDENTIFIER ('.' qualifiers+=IDENTIFIER)*;

// type sublanguage
startType : t=type EOF;

type :
    (typeId params=typeParamList?) | paramId ;

typeId :
    leadingDot='.'? id=(IDENTIFIER|NUL) ('.' qualifiers+=IDENTIFIER )* ;

paramId :
  PARAM_SPECIFIER;

typeParamList:
    '<' ( types+=type (',' types+=type)* )? '>' ;

// lexer rules:
COMMAND: '%' IDENTIFIER;
FLAG: '-' ('-' IDENTIFIER)+;
ARROW: '->';
EQUAL_ASSIGN: '=';
PARAM_SPECIFIER: '~' LETTER;