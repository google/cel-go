// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

grammar CEL;

// Note: grammar here should be rather final, but lexis is a moving target.

// Grammar Rules
// =============

start
    : e=expr EOF
    ;

expr
    : e=conditionalOr (op='?' e1=conditionalOr ':' e2=expr)?
    ;

conditionalOr
    : e=conditionalAnd (ops+='||' e1+=conditionalAnd)*
    ;

conditionalAnd
    : e=relation (ops+='&&' e1+=relation)*
    ;

relation
    : calc
    | relation op=('<'|'<='|'>='|'>'|'=='|'!='|'in') relation
    ;

calc
    : unary
    | calc op=('*'|'/'|'%') calc
    | calc op=('+'|'-') calc
    ;

unary
    : statement                                                     # StatementExpr
    | (ops+='!')+ statement                                         # LogicalNot
    | (ops+='-')+ statement                                         # Negate
    ;

statement
    : primary                                                       # PrimaryExpr
    | statement op='.' id=IDENTIFIER (open='(' args=exprList? ')')? # SelectOrCall
    | statement op='[' index=expr ']'                               # Index
    | statement op='{' entries=fieldInitializerList? '}'            # CreateMessage
    ;

primary
    : leadingDot='.'? id=IDENTIFIER (op='(' args=exprList? ')')?    # IdentOrGlobalCall
    | '(' e=expr ')'                                                # Nested
    | op='[' elems=exprList? ','? ']'                               # CreateList
    | op='{' entries=mapInitializerList? '}'                        # CreateStruct
    | literal                                                       # ConstantLiteral
    ;

exprList
    : e+=expr (',' e+=expr)*
    ;

fieldInitializerList
    : fields+=IDENTIFIER cols+=':' values+=expr (',' fields+=IDENTIFIER cols+=':' values+=expr)*
    ;

mapInitializerList
    : keys+=expr cols+=':' values+=expr (',' keys+=expr cols+=':' values+=expr)*
    ;

literal
    : tok=NUM_INT   # Int
    | tok=NUM_UINT  # Uint
    | tok=NUM_FLOAT # Double
    | tok=STRING    # String
    | tok=BYTES     # Bytes
    | tok='true'    # BoolTrue
    | tok='false'   # BoolFalse
    | tok='null'    # Null
    ;

// Lexer Rules
// ===========

// TODO(b/63180832): align implementation with spec at go/api-expr/langdef.md.
//   At this point, it would be premature to do this as the spec is still a moving
//   target. The current rules were c&p from some other code in google3, and are
//   only approximative.

EQUALS : '==';
NOT_EQUALS : '!=';
LESS : '<';
LESS_EQUALS : '<=';
GREATER_EQUALS : '>=';
GREATER : '>';
LOGICAL_AND : '&&';
LOGICAL_OR : '||';

LBRACKET : '[';
RPRACKET : ']';
LBRACE : '{';
RBRACE : '}';
LPAREN : '(';
RPAREN : ')';
DOT : '.';
COMMA : ',';
MINUS : '-';
EXCLAM : '!';
QUESTIONMARK : '?';
COLON : ':';
PLUS : '+';
STAR : '*';
SLASH : '/';
PERCENT : '%';
TRUE : 'true';
FALSE : 'false';
NULL : 'null';

fragment EXPONENT : ('e' | 'E') ( '+' | '-' )? DIGIT+ ;
fragment DIGIT  : '0'..'9' ;
fragment LETTER : 'A'..'Z' | 'a'..'z' ;
fragment HEXDIGIT : ('0'..'9'|'a'..'f'|'A'..'F') ;
fragment OCTDIGIT : '0'..'7' ;

// TODO: Ensure escape sequences mirror those supported in Google SQL.
fragment ESC_SEQ
    : '\\' ~('0'..'9'|'a'..'f'|'A'..'F' | '\r' | '\n' | '\t' | 'u')
    | UNI_SEQ
    | OCT_SEQ
    ;

fragment SPECIAL_ESC_CHAR
    : ('b'|'t'|'n'|'f'|'r'|'\"'|'\''|'\\')
    ;

fragment OCT_SEQ
    : '\\' ('0'..'3') ('0'..'7') ('0'..'7')
    | '\\' ('0'..'7') ('0'..'7')
    | '\\' ('0'..'7')
    ;

fragment UNI_SEQ
    : '\\' 'u' HEXDIGIT HEXDIGIT HEXDIGIT HEXDIGIT
    ;

WHITESPACE : ( '\t' | ' ' | '\r' | '\n'| '\u000C' )+ -> channel(HIDDEN) ;
COMMENT : '//' (~'\n')* -> channel(HIDDEN) ;

NUM_FLOAT
  : ( DIGIT+ ('.' DIGIT+) EXPONENT?
    | DIGIT+ EXPONENT
    | '.' DIGIT+ EXPONENT?
    )
  ;

NUM_INT
  : ( DIGIT+ | '0x' HEXDIGIT+ );

NUM_UINT
   : DIGIT+ ( 'u' | 'U' )
   | '0x' HEXDIGIT+ ( 'u' | 'U' )
   ;

STRING
  : '"' (ESC_SEQ | ~('"'|'\n'|'\r'))* '"'
  | '\'' (ESC_SEQ | ~('\''|'\n'|'\r'))* '\''
  | '"""' (ESC_SEQ | ~('\\'))*? '"""'
  | '\'\'\'' (ESC_SEQ | ~('\\'))*? '\'\'\''
  | RAW '"' ~('"'|'\n'|'\r')* '"'
  | RAW '\'' ~('\''|'\n'|'\r')* '\''
  | RAW '"""' .*? '"""'
  | RAW '\'\'\'' .*? '\'\'\''
  ;

fragment RAW : 'r' | 'R';

BYTES : ('b' | 'B') STRING;

IDENTIFIER : (LETTER | '_') ( LETTER | DIGIT | '_')*;