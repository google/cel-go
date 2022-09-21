// Code generated from ./Commands.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Commands
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type CommandsParser struct {
	*antlr.BaseParser
}

var commandsParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func commandsParserInit() {
	staticData := &commandsParserStaticData
	staticData.literalNames = []string{
		"", "'%let'", "'%declare'", "'%delete'", "'%eval'", "", "", "'->'",
		"'='", "'=='", "'!='", "'in'", "'<'", "'<='", "'>='", "'>'", "'&&'",
		"'||'", "'['", "']'", "'{'", "'}'", "'('", "')'", "'.'", "','", "'-'",
		"'!'", "'?'", "':'", "'+'", "'*'", "'/'", "'%'", "'true'", "'false'",
		"'null'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "", "COMMAND", "FLAG", "ARROW", "EQUAL_ASSIGN", "EQUALS",
		"NOT_EQUALS", "IN", "LESS", "LESS_EQUALS", "GREATER_EQUALS", "GREATER",
		"LOGICAL_AND", "LOGICAL_OR", "LBRACKET", "RPRACKET", "LBRACE", "RBRACE",
		"LPAREN", "RPAREN", "DOT", "COMMA", "MINUS", "EXCLAM", "QUESTIONMARK",
		"COLON", "PLUS", "STAR", "SLASH", "PERCENT", "CEL_TRUE", "CEL_FALSE",
		"NUL", "WHITESPACE", "COMMENT", "NUM_FLOAT", "NUM_INT", "NUM_UINT",
		"STRING", "BYTES", "IDENTIFIER",
	}
	staticData.ruleNames = []string{
		"startCommand", "command", "let", "declare", "varDecl", "fnDecl", "param",
		"delete", "simple", "empty", "exprCmd", "qualId", "startType", "type",
		"typeId", "typeParamList", "start", "expr", "conditionalOr", "conditionalAnd",
		"relation", "calc", "unary", "member", "primary", "exprList", "fieldInitializerList",
		"mapInitializerList", "literal",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 44, 353, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 3, 1, 68, 8, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 3, 2, 77, 8, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 3, 3, 84, 8, 3, 1, 4,
		1, 4, 1, 4, 3, 4, 89, 8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 5, 5, 96, 8,
		5, 10, 5, 12, 5, 99, 9, 5, 3, 5, 101, 8, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1,
		6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 3, 7, 114, 8, 7, 1, 8, 1, 8, 1,
		8, 5, 8, 119, 8, 8, 10, 8, 12, 8, 122, 9, 8, 1, 9, 1, 9, 1, 10, 3, 10,
		127, 8, 10, 1, 10, 1, 10, 1, 11, 3, 11, 132, 8, 11, 1, 11, 1, 11, 1, 11,
		5, 11, 137, 8, 11, 10, 11, 12, 11, 140, 9, 11, 1, 12, 1, 12, 1, 12, 1,
		13, 1, 13, 3, 13, 147, 8, 13, 1, 14, 3, 14, 150, 8, 14, 1, 14, 1, 14, 1,
		14, 5, 14, 155, 8, 14, 10, 14, 12, 14, 158, 9, 14, 1, 15, 1, 15, 1, 15,
		1, 15, 5, 15, 164, 8, 15, 10, 15, 12, 15, 167, 9, 15, 1, 15, 1, 15, 1,
		16, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 1, 17, 3, 17, 180,
		8, 17, 1, 18, 1, 18, 1, 18, 5, 18, 185, 8, 18, 10, 18, 12, 18, 188, 9,
		18, 1, 19, 1, 19, 1, 19, 5, 19, 193, 8, 19, 10, 19, 12, 19, 196, 9, 19,
		1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 5, 20, 204, 8, 20, 10, 20, 12,
		20, 207, 9, 20, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21, 1, 21,
		1, 21, 5, 21, 218, 8, 21, 10, 21, 12, 21, 221, 9, 21, 1, 22, 1, 22, 4,
		22, 225, 8, 22, 11, 22, 12, 22, 226, 1, 22, 1, 22, 4, 22, 231, 8, 22, 11,
		22, 12, 22, 232, 1, 22, 3, 22, 236, 8, 22, 1, 23, 1, 23, 1, 23, 1, 23,
		1, 23, 1, 23, 1, 23, 1, 23, 3, 23, 246, 8, 23, 1, 23, 3, 23, 249, 8, 23,
		1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 3, 23, 259, 8,
		23, 1, 23, 3, 23, 262, 8, 23, 1, 23, 5, 23, 265, 8, 23, 10, 23, 12, 23,
		268, 9, 23, 1, 24, 3, 24, 271, 8, 24, 1, 24, 1, 24, 1, 24, 3, 24, 276,
		8, 24, 1, 24, 3, 24, 279, 8, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1,
		24, 3, 24, 287, 8, 24, 1, 24, 3, 24, 290, 8, 24, 1, 24, 1, 24, 1, 24, 3,
		24, 295, 8, 24, 1, 24, 3, 24, 298, 8, 24, 1, 24, 1, 24, 3, 24, 302, 8,
		24, 1, 25, 1, 25, 1, 25, 5, 25, 307, 8, 25, 10, 25, 12, 25, 310, 9, 25,
		1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26, 5, 26, 319, 8, 26, 10,
		26, 12, 26, 322, 9, 26, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27,
		1, 27, 5, 27, 332, 8, 27, 10, 27, 12, 27, 335, 9, 27, 1, 28, 3, 28, 338,
		8, 28, 1, 28, 1, 28, 1, 28, 3, 28, 343, 8, 28, 1, 28, 1, 28, 1, 28, 1,
		28, 1, 28, 1, 28, 3, 28, 351, 8, 28, 1, 28, 0, 3, 40, 42, 46, 29, 0, 2,
		4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40,
		42, 44, 46, 48, 50, 52, 54, 56, 0, 4, 2, 0, 36, 36, 44, 44, 1, 0, 9, 15,
		1, 0, 31, 33, 2, 0, 26, 26, 30, 30, 383, 0, 58, 1, 0, 0, 0, 2, 67, 1, 0,
		0, 0, 4, 69, 1, 0, 0, 0, 6, 80, 1, 0, 0, 0, 8, 85, 1, 0, 0, 0, 10, 90,
		1, 0, 0, 0, 12, 106, 1, 0, 0, 0, 14, 110, 1, 0, 0, 0, 16, 115, 1, 0, 0,
		0, 18, 123, 1, 0, 0, 0, 20, 126, 1, 0, 0, 0, 22, 131, 1, 0, 0, 0, 24, 141,
		1, 0, 0, 0, 26, 144, 1, 0, 0, 0, 28, 149, 1, 0, 0, 0, 30, 159, 1, 0, 0,
		0, 32, 170, 1, 0, 0, 0, 34, 173, 1, 0, 0, 0, 36, 181, 1, 0, 0, 0, 38, 189,
		1, 0, 0, 0, 40, 197, 1, 0, 0, 0, 42, 208, 1, 0, 0, 0, 44, 235, 1, 0, 0,
		0, 46, 237, 1, 0, 0, 0, 48, 301, 1, 0, 0, 0, 50, 303, 1, 0, 0, 0, 52, 311,
		1, 0, 0, 0, 54, 323, 1, 0, 0, 0, 56, 350, 1, 0, 0, 0, 58, 59, 3, 2, 1,
		0, 59, 60, 5, 0, 0, 1, 60, 1, 1, 0, 0, 0, 61, 68, 3, 4, 2, 0, 62, 68, 3,
		6, 3, 0, 63, 68, 3, 14, 7, 0, 64, 68, 3, 16, 8, 0, 65, 68, 3, 20, 10, 0,
		66, 68, 3, 18, 9, 0, 67, 61, 1, 0, 0, 0, 67, 62, 1, 0, 0, 0, 67, 63, 1,
		0, 0, 0, 67, 64, 1, 0, 0, 0, 67, 65, 1, 0, 0, 0, 67, 66, 1, 0, 0, 0, 68,
		3, 1, 0, 0, 0, 69, 76, 5, 1, 0, 0, 70, 71, 3, 8, 4, 0, 71, 72, 5, 8, 0,
		0, 72, 77, 1, 0, 0, 0, 73, 74, 3, 10, 5, 0, 74, 75, 5, 7, 0, 0, 75, 77,
		1, 0, 0, 0, 76, 70, 1, 0, 0, 0, 76, 73, 1, 0, 0, 0, 77, 78, 1, 0, 0, 0,
		78, 79, 3, 34, 17, 0, 79, 5, 1, 0, 0, 0, 80, 83, 5, 2, 0, 0, 81, 84, 3,
		8, 4, 0, 82, 84, 3, 10, 5, 0, 83, 81, 1, 0, 0, 0, 83, 82, 1, 0, 0, 0, 84,
		7, 1, 0, 0, 0, 85, 88, 3, 22, 11, 0, 86, 87, 5, 29, 0, 0, 87, 89, 3, 26,
		13, 0, 88, 86, 1, 0, 0, 0, 88, 89, 1, 0, 0, 0, 89, 9, 1, 0, 0, 0, 90, 91,
		3, 22, 11, 0, 91, 100, 5, 22, 0, 0, 92, 97, 3, 12, 6, 0, 93, 94, 5, 25,
		0, 0, 94, 96, 3, 12, 6, 0, 95, 93, 1, 0, 0, 0, 96, 99, 1, 0, 0, 0, 97,
		95, 1, 0, 0, 0, 97, 98, 1, 0, 0, 0, 98, 101, 1, 0, 0, 0, 99, 97, 1, 0,
		0, 0, 100, 92, 1, 0, 0, 0, 100, 101, 1, 0, 0, 0, 101, 102, 1, 0, 0, 0,
		102, 103, 5, 23, 0, 0, 103, 104, 5, 29, 0, 0, 104, 105, 3, 26, 13, 0, 105,
		11, 1, 0, 0, 0, 106, 107, 5, 44, 0, 0, 107, 108, 5, 29, 0, 0, 108, 109,
		3, 26, 13, 0, 109, 13, 1, 0, 0, 0, 110, 113, 5, 3, 0, 0, 111, 114, 3, 8,
		4, 0, 112, 114, 3, 10, 5, 0, 113, 111, 1, 0, 0, 0, 113, 112, 1, 0, 0, 0,
		114, 15, 1, 0, 0, 0, 115, 120, 5, 5, 0, 0, 116, 119, 5, 6, 0, 0, 117, 119,
		5, 42, 0, 0, 118, 116, 1, 0, 0, 0, 118, 117, 1, 0, 0, 0, 119, 122, 1, 0,
		0, 0, 120, 118, 1, 0, 0, 0, 120, 121, 1, 0, 0, 0, 121, 17, 1, 0, 0, 0,
		122, 120, 1, 0, 0, 0, 123, 124, 1, 0, 0, 0, 124, 19, 1, 0, 0, 0, 125, 127,
		5, 4, 0, 0, 126, 125, 1, 0, 0, 0, 126, 127, 1, 0, 0, 0, 127, 128, 1, 0,
		0, 0, 128, 129, 3, 34, 17, 0, 129, 21, 1, 0, 0, 0, 130, 132, 5, 24, 0,
		0, 131, 130, 1, 0, 0, 0, 131, 132, 1, 0, 0, 0, 132, 133, 1, 0, 0, 0, 133,
		138, 5, 44, 0, 0, 134, 135, 5, 24, 0, 0, 135, 137, 5, 44, 0, 0, 136, 134,
		1, 0, 0, 0, 137, 140, 1, 0, 0, 0, 138, 136, 1, 0, 0, 0, 138, 139, 1, 0,
		0, 0, 139, 23, 1, 0, 0, 0, 140, 138, 1, 0, 0, 0, 141, 142, 3, 26, 13, 0,
		142, 143, 5, 0, 0, 1, 143, 25, 1, 0, 0, 0, 144, 146, 3, 28, 14, 0, 145,
		147, 3, 30, 15, 0, 146, 145, 1, 0, 0, 0, 146, 147, 1, 0, 0, 0, 147, 27,
		1, 0, 0, 0, 148, 150, 5, 24, 0, 0, 149, 148, 1, 0, 0, 0, 149, 150, 1, 0,
		0, 0, 150, 151, 1, 0, 0, 0, 151, 156, 7, 0, 0, 0, 152, 153, 5, 24, 0, 0,
		153, 155, 5, 44, 0, 0, 154, 152, 1, 0, 0, 0, 155, 158, 1, 0, 0, 0, 156,
		154, 1, 0, 0, 0, 156, 157, 1, 0, 0, 0, 157, 29, 1, 0, 0, 0, 158, 156, 1,
		0, 0, 0, 159, 160, 5, 22, 0, 0, 160, 165, 3, 26, 13, 0, 161, 162, 5, 25,
		0, 0, 162, 164, 3, 26, 13, 0, 163, 161, 1, 0, 0, 0, 164, 167, 1, 0, 0,
		0, 165, 163, 1, 0, 0, 0, 165, 166, 1, 0, 0, 0, 166, 168, 1, 0, 0, 0, 167,
		165, 1, 0, 0, 0, 168, 169, 5, 23, 0, 0, 169, 31, 1, 0, 0, 0, 170, 171,
		3, 34, 17, 0, 171, 172, 5, 0, 0, 1, 172, 33, 1, 0, 0, 0, 173, 179, 3, 36,
		18, 0, 174, 175, 5, 28, 0, 0, 175, 176, 3, 36, 18, 0, 176, 177, 5, 29,
		0, 0, 177, 178, 3, 34, 17, 0, 178, 180, 1, 0, 0, 0, 179, 174, 1, 0, 0,
		0, 179, 180, 1, 0, 0, 0, 180, 35, 1, 0, 0, 0, 181, 186, 3, 38, 19, 0, 182,
		183, 5, 17, 0, 0, 183, 185, 3, 38, 19, 0, 184, 182, 1, 0, 0, 0, 185, 188,
		1, 0, 0, 0, 186, 184, 1, 0, 0, 0, 186, 187, 1, 0, 0, 0, 187, 37, 1, 0,
		0, 0, 188, 186, 1, 0, 0, 0, 189, 194, 3, 40, 20, 0, 190, 191, 5, 16, 0,
		0, 191, 193, 3, 40, 20, 0, 192, 190, 1, 0, 0, 0, 193, 196, 1, 0, 0, 0,
		194, 192, 1, 0, 0, 0, 194, 195, 1, 0, 0, 0, 195, 39, 1, 0, 0, 0, 196, 194,
		1, 0, 0, 0, 197, 198, 6, 20, -1, 0, 198, 199, 3, 42, 21, 0, 199, 205, 1,
		0, 0, 0, 200, 201, 10, 1, 0, 0, 201, 202, 7, 1, 0, 0, 202, 204, 3, 40,
		20, 2, 203, 200, 1, 0, 0, 0, 204, 207, 1, 0, 0, 0, 205, 203, 1, 0, 0, 0,
		205, 206, 1, 0, 0, 0, 206, 41, 1, 0, 0, 0, 207, 205, 1, 0, 0, 0, 208, 209,
		6, 21, -1, 0, 209, 210, 3, 44, 22, 0, 210, 219, 1, 0, 0, 0, 211, 212, 10,
		2, 0, 0, 212, 213, 7, 2, 0, 0, 213, 218, 3, 42, 21, 3, 214, 215, 10, 1,
		0, 0, 215, 216, 7, 3, 0, 0, 216, 218, 3, 42, 21, 2, 217, 211, 1, 0, 0,
		0, 217, 214, 1, 0, 0, 0, 218, 221, 1, 0, 0, 0, 219, 217, 1, 0, 0, 0, 219,
		220, 1, 0, 0, 0, 220, 43, 1, 0, 0, 0, 221, 219, 1, 0, 0, 0, 222, 236, 3,
		46, 23, 0, 223, 225, 5, 27, 0, 0, 224, 223, 1, 0, 0, 0, 225, 226, 1, 0,
		0, 0, 226, 224, 1, 0, 0, 0, 226, 227, 1, 0, 0, 0, 227, 228, 1, 0, 0, 0,
		228, 236, 3, 46, 23, 0, 229, 231, 5, 26, 0, 0, 230, 229, 1, 0, 0, 0, 231,
		232, 1, 0, 0, 0, 232, 230, 1, 0, 0, 0, 232, 233, 1, 0, 0, 0, 233, 234,
		1, 0, 0, 0, 234, 236, 3, 46, 23, 0, 235, 222, 1, 0, 0, 0, 235, 224, 1,
		0, 0, 0, 235, 230, 1, 0, 0, 0, 236, 45, 1, 0, 0, 0, 237, 238, 6, 23, -1,
		0, 238, 239, 3, 48, 24, 0, 239, 266, 1, 0, 0, 0, 240, 241, 10, 3, 0, 0,
		241, 242, 5, 24, 0, 0, 242, 248, 5, 44, 0, 0, 243, 245, 5, 22, 0, 0, 244,
		246, 3, 50, 25, 0, 245, 244, 1, 0, 0, 0, 245, 246, 1, 0, 0, 0, 246, 247,
		1, 0, 0, 0, 247, 249, 5, 23, 0, 0, 248, 243, 1, 0, 0, 0, 248, 249, 1, 0,
		0, 0, 249, 265, 1, 0, 0, 0, 250, 251, 10, 2, 0, 0, 251, 252, 5, 18, 0,
		0, 252, 253, 3, 34, 17, 0, 253, 254, 5, 19, 0, 0, 254, 265, 1, 0, 0, 0,
		255, 256, 10, 1, 0, 0, 256, 258, 5, 20, 0, 0, 257, 259, 3, 52, 26, 0, 258,
		257, 1, 0, 0, 0, 258, 259, 1, 0, 0, 0, 259, 261, 1, 0, 0, 0, 260, 262,
		5, 25, 0, 0, 261, 260, 1, 0, 0, 0, 261, 262, 1, 0, 0, 0, 262, 263, 1, 0,
		0, 0, 263, 265, 5, 21, 0, 0, 264, 240, 1, 0, 0, 0, 264, 250, 1, 0, 0, 0,
		264, 255, 1, 0, 0, 0, 265, 268, 1, 0, 0, 0, 266, 264, 1, 0, 0, 0, 266,
		267, 1, 0, 0, 0, 267, 47, 1, 0, 0, 0, 268, 266, 1, 0, 0, 0, 269, 271, 5,
		24, 0, 0, 270, 269, 1, 0, 0, 0, 270, 271, 1, 0, 0, 0, 271, 272, 1, 0, 0,
		0, 272, 278, 5, 44, 0, 0, 273, 275, 5, 22, 0, 0, 274, 276, 3, 50, 25, 0,
		275, 274, 1, 0, 0, 0, 275, 276, 1, 0, 0, 0, 276, 277, 1, 0, 0, 0, 277,
		279, 5, 23, 0, 0, 278, 273, 1, 0, 0, 0, 278, 279, 1, 0, 0, 0, 279, 302,
		1, 0, 0, 0, 280, 281, 5, 22, 0, 0, 281, 282, 3, 34, 17, 0, 282, 283, 5,
		23, 0, 0, 283, 302, 1, 0, 0, 0, 284, 286, 5, 18, 0, 0, 285, 287, 3, 50,
		25, 0, 286, 285, 1, 0, 0, 0, 286, 287, 1, 0, 0, 0, 287, 289, 1, 0, 0, 0,
		288, 290, 5, 25, 0, 0, 289, 288, 1, 0, 0, 0, 289, 290, 1, 0, 0, 0, 290,
		291, 1, 0, 0, 0, 291, 302, 5, 19, 0, 0, 292, 294, 5, 20, 0, 0, 293, 295,
		3, 54, 27, 0, 294, 293, 1, 0, 0, 0, 294, 295, 1, 0, 0, 0, 295, 297, 1,
		0, 0, 0, 296, 298, 5, 25, 0, 0, 297, 296, 1, 0, 0, 0, 297, 298, 1, 0, 0,
		0, 298, 299, 1, 0, 0, 0, 299, 302, 5, 21, 0, 0, 300, 302, 3, 56, 28, 0,
		301, 270, 1, 0, 0, 0, 301, 280, 1, 0, 0, 0, 301, 284, 1, 0, 0, 0, 301,
		292, 1, 0, 0, 0, 301, 300, 1, 0, 0, 0, 302, 49, 1, 0, 0, 0, 303, 308, 3,
		34, 17, 0, 304, 305, 5, 25, 0, 0, 305, 307, 3, 34, 17, 0, 306, 304, 1,
		0, 0, 0, 307, 310, 1, 0, 0, 0, 308, 306, 1, 0, 0, 0, 308, 309, 1, 0, 0,
		0, 309, 51, 1, 0, 0, 0, 310, 308, 1, 0, 0, 0, 311, 312, 5, 44, 0, 0, 312,
		313, 5, 29, 0, 0, 313, 320, 3, 34, 17, 0, 314, 315, 5, 25, 0, 0, 315, 316,
		5, 44, 0, 0, 316, 317, 5, 29, 0, 0, 317, 319, 3, 34, 17, 0, 318, 314, 1,
		0, 0, 0, 319, 322, 1, 0, 0, 0, 320, 318, 1, 0, 0, 0, 320, 321, 1, 0, 0,
		0, 321, 53, 1, 0, 0, 0, 322, 320, 1, 0, 0, 0, 323, 324, 3, 34, 17, 0, 324,
		325, 5, 29, 0, 0, 325, 333, 3, 34, 17, 0, 326, 327, 5, 25, 0, 0, 327, 328,
		3, 34, 17, 0, 328, 329, 5, 29, 0, 0, 329, 330, 3, 34, 17, 0, 330, 332,
		1, 0, 0, 0, 331, 326, 1, 0, 0, 0, 332, 335, 1, 0, 0, 0, 333, 331, 1, 0,
		0, 0, 333, 334, 1, 0, 0, 0, 334, 55, 1, 0, 0, 0, 335, 333, 1, 0, 0, 0,
		336, 338, 5, 26, 0, 0, 337, 336, 1, 0, 0, 0, 337, 338, 1, 0, 0, 0, 338,
		339, 1, 0, 0, 0, 339, 351, 5, 40, 0, 0, 340, 351, 5, 41, 0, 0, 341, 343,
		5, 26, 0, 0, 342, 341, 1, 0, 0, 0, 342, 343, 1, 0, 0, 0, 343, 344, 1, 0,
		0, 0, 344, 351, 5, 39, 0, 0, 345, 351, 5, 42, 0, 0, 346, 351, 5, 43, 0,
		0, 347, 351, 5, 34, 0, 0, 348, 351, 5, 35, 0, 0, 349, 351, 5, 36, 0, 0,
		350, 337, 1, 0, 0, 0, 350, 340, 1, 0, 0, 0, 350, 342, 1, 0, 0, 0, 350,
		345, 1, 0, 0, 0, 350, 346, 1, 0, 0, 0, 350, 347, 1, 0, 0, 0, 350, 348,
		1, 0, 0, 0, 350, 349, 1, 0, 0, 0, 351, 57, 1, 0, 0, 0, 45, 67, 76, 83,
		88, 97, 100, 113, 118, 120, 126, 131, 138, 146, 149, 156, 165, 179, 186,
		194, 205, 217, 219, 226, 232, 235, 245, 248, 258, 261, 264, 266, 270, 275,
		278, 286, 289, 294, 297, 301, 308, 320, 333, 337, 342, 350,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// CommandsParserInit initializes any static state used to implement CommandsParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewCommandsParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func CommandsParserInit() {
	staticData := &commandsParserStaticData
	staticData.once.Do(commandsParserInit)
}

// NewCommandsParser produces a new parser instance for the optional input antlr.TokenStream.
func NewCommandsParser(input antlr.TokenStream) *CommandsParser {
	CommandsParserInit()
	this := new(CommandsParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &commandsParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "Commands.g4"

	return this
}

// CommandsParser tokens.
const (
	CommandsParserEOF            = antlr.TokenEOF
	CommandsParserT__0           = 1
	CommandsParserT__1           = 2
	CommandsParserT__2           = 3
	CommandsParserT__3           = 4
	CommandsParserCOMMAND        = 5
	CommandsParserFLAG           = 6
	CommandsParserARROW          = 7
	CommandsParserEQUAL_ASSIGN   = 8
	CommandsParserEQUALS         = 9
	CommandsParserNOT_EQUALS     = 10
	CommandsParserIN             = 11
	CommandsParserLESS           = 12
	CommandsParserLESS_EQUALS    = 13
	CommandsParserGREATER_EQUALS = 14
	CommandsParserGREATER        = 15
	CommandsParserLOGICAL_AND    = 16
	CommandsParserLOGICAL_OR     = 17
	CommandsParserLBRACKET       = 18
	CommandsParserRPRACKET       = 19
	CommandsParserLBRACE         = 20
	CommandsParserRBRACE         = 21
	CommandsParserLPAREN         = 22
	CommandsParserRPAREN         = 23
	CommandsParserDOT            = 24
	CommandsParserCOMMA          = 25
	CommandsParserMINUS          = 26
	CommandsParserEXCLAM         = 27
	CommandsParserQUESTIONMARK   = 28
	CommandsParserCOLON          = 29
	CommandsParserPLUS           = 30
	CommandsParserSTAR           = 31
	CommandsParserSLASH          = 32
	CommandsParserPERCENT        = 33
	CommandsParserCEL_TRUE       = 34
	CommandsParserCEL_FALSE      = 35
	CommandsParserNUL            = 36
	CommandsParserWHITESPACE     = 37
	CommandsParserCOMMENT        = 38
	CommandsParserNUM_FLOAT      = 39
	CommandsParserNUM_INT        = 40
	CommandsParserNUM_UINT       = 41
	CommandsParserSTRING         = 42
	CommandsParserBYTES          = 43
	CommandsParserIDENTIFIER     = 44
)

// CommandsParser rules.
const (
	CommandsParserRULE_startCommand         = 0
	CommandsParserRULE_command              = 1
	CommandsParserRULE_let                  = 2
	CommandsParserRULE_declare              = 3
	CommandsParserRULE_varDecl              = 4
	CommandsParserRULE_fnDecl               = 5
	CommandsParserRULE_param                = 6
	CommandsParserRULE_delete               = 7
	CommandsParserRULE_simple               = 8
	CommandsParserRULE_empty                = 9
	CommandsParserRULE_exprCmd              = 10
	CommandsParserRULE_qualId               = 11
	CommandsParserRULE_startType            = 12
	CommandsParserRULE_type                 = 13
	CommandsParserRULE_typeId               = 14
	CommandsParserRULE_typeParamList        = 15
	CommandsParserRULE_start                = 16
	CommandsParserRULE_expr                 = 17
	CommandsParserRULE_conditionalOr        = 18
	CommandsParserRULE_conditionalAnd       = 19
	CommandsParserRULE_relation             = 20
	CommandsParserRULE_calc                 = 21
	CommandsParserRULE_unary                = 22
	CommandsParserRULE_member               = 23
	CommandsParserRULE_primary              = 24
	CommandsParserRULE_exprList             = 25
	CommandsParserRULE_fieldInitializerList = 26
	CommandsParserRULE_mapInitializerList   = 27
	CommandsParserRULE_literal              = 28
)

// IStartCommandContext is an interface to support dynamic dispatch.
type IStartCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStartCommandContext differentiates from other interfaces.
	IsStartCommandContext()
}

type StartCommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartCommandContext() *StartCommandContext {
	var p = new(StartCommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_startCommand
	return p
}

func (*StartCommandContext) IsStartCommandContext() {}

func NewStartCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartCommandContext {
	var p = new(StartCommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_startCommand

	return p
}

func (s *StartCommandContext) GetParser() antlr.Parser { return s.parser }

func (s *StartCommandContext) Command() ICommandContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICommandContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICommandContext)
}

func (s *StartCommandContext) EOF() antlr.TerminalNode {
	return s.GetToken(CommandsParserEOF, 0)
}

func (s *StartCommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartCommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartCommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterStartCommand(s)
	}
}

func (s *StartCommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitStartCommand(s)
	}
}

func (s *StartCommandContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitStartCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) StartCommand() (localctx IStartCommandContext) {
	this := p
	_ = this

	localctx = NewStartCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, CommandsParserRULE_startCommand)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(58)
		p.Command()
	}
	{
		p.SetState(59)
		p.Match(CommandsParserEOF)
	}

	return localctx
}

// ICommandContext is an interface to support dynamic dispatch.
type ICommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCommandContext differentiates from other interfaces.
	IsCommandContext()
}

type CommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandContext() *CommandContext {
	var p = new(CommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_command
	return p
}

func (*CommandContext) IsCommandContext() {}

func NewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandContext {
	var p = new(CommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_command

	return p
}

func (s *CommandContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandContext) Let() ILetContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILetContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILetContext)
}

func (s *CommandContext) Declare() IDeclareContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeclareContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeclareContext)
}

func (s *CommandContext) Delete() IDeleteContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDeleteContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDeleteContext)
}

func (s *CommandContext) Simple() ISimpleContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISimpleContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISimpleContext)
}

func (s *CommandContext) ExprCmd() IExprCmdContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprCmdContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprCmdContext)
}

func (s *CommandContext) Empty() IEmptyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IEmptyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IEmptyContext)
}

func (s *CommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterCommand(s)
	}
}

func (s *CommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitCommand(s)
	}
}

func (s *CommandContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Command() (localctx ICommandContext) {
	this := p
	_ = this

	localctx = NewCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, CommandsParserRULE_command)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(67)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CommandsParserT__0:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(61)
			p.Let()
		}

	case CommandsParserT__1:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(62)
			p.Declare()
		}

	case CommandsParserT__2:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(63)
			p.Delete()
		}

	case CommandsParserCOMMAND:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(64)
			p.Simple()
		}

	case CommandsParserT__3, CommandsParserLBRACKET, CommandsParserLBRACE, CommandsParserLPAREN, CommandsParserDOT, CommandsParserMINUS, CommandsParserEXCLAM, CommandsParserCEL_TRUE, CommandsParserCEL_FALSE, CommandsParserNUL, CommandsParserNUM_FLOAT, CommandsParserNUM_INT, CommandsParserNUM_UINT, CommandsParserSTRING, CommandsParserBYTES, CommandsParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(65)
			p.ExprCmd()
		}

	case CommandsParserEOF:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(66)
			p.Empty()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ILetContext is an interface to support dynamic dispatch.
type ILetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVar returns the var rule contexts.
	GetVar() IVarDeclContext

	// GetFn returns the fn rule contexts.
	GetFn() IFnDeclContext

	// GetE returns the e rule contexts.
	GetE() IExprContext

	// SetVar sets the var rule contexts.
	SetVar(IVarDeclContext)

	// SetFn sets the fn rule contexts.
	SetFn(IFnDeclContext)

	// SetE sets the e rule contexts.
	SetE(IExprContext)

	// IsLetContext differentiates from other interfaces.
	IsLetContext()
}

type LetContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	var_   IVarDeclContext
	fn     IFnDeclContext
	e      IExprContext
}

func NewEmptyLetContext() *LetContext {
	var p = new(LetContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_let
	return p
}

func (*LetContext) IsLetContext() {}

func NewLetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LetContext {
	var p = new(LetContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_let

	return p
}

func (s *LetContext) GetParser() antlr.Parser { return s.parser }

func (s *LetContext) GetVar() IVarDeclContext { return s.var_ }

func (s *LetContext) GetFn() IFnDeclContext { return s.fn }

func (s *LetContext) GetE() IExprContext { return s.e }

func (s *LetContext) SetVar(v IVarDeclContext) { s.var_ = v }

func (s *LetContext) SetFn(v IFnDeclContext) { s.fn = v }

func (s *LetContext) SetE(v IExprContext) { s.e = v }

func (s *LetContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *LetContext) EQUAL_ASSIGN() antlr.TerminalNode {
	return s.GetToken(CommandsParserEQUAL_ASSIGN, 0)
}

func (s *LetContext) ARROW() antlr.TerminalNode {
	return s.GetToken(CommandsParserARROW, 0)
}

func (s *LetContext) VarDecl() IVarDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVarDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVarDeclContext)
}

func (s *LetContext) FnDecl() IFnDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFnDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFnDeclContext)
}

func (s *LetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterLet(s)
	}
}

func (s *LetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitLet(s)
	}
}

func (s *LetContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitLet(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Let() (localctx ILetContext) {
	this := p
	_ = this

	localctx = NewLetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, CommandsParserRULE_let)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		p.Match(CommandsParserT__0)
	}
	p.SetState(76)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(70)

			var _x = p.VarDecl()

			localctx.(*LetContext).var_ = _x
		}
		{
			p.SetState(71)
			p.Match(CommandsParserEQUAL_ASSIGN)
		}

	case 2:
		{
			p.SetState(73)

			var _x = p.FnDecl()

			localctx.(*LetContext).fn = _x
		}
		{
			p.SetState(74)
			p.Match(CommandsParserARROW)
		}

	}
	{
		p.SetState(78)

		var _x = p.Expr()

		localctx.(*LetContext).e = _x
	}

	return localctx
}

// IDeclareContext is an interface to support dynamic dispatch.
type IDeclareContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVar returns the var rule contexts.
	GetVar() IVarDeclContext

	// GetFn returns the fn rule contexts.
	GetFn() IFnDeclContext

	// SetVar sets the var rule contexts.
	SetVar(IVarDeclContext)

	// SetFn sets the fn rule contexts.
	SetFn(IFnDeclContext)

	// IsDeclareContext differentiates from other interfaces.
	IsDeclareContext()
}

type DeclareContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	var_   IVarDeclContext
	fn     IFnDeclContext
}

func NewEmptyDeclareContext() *DeclareContext {
	var p = new(DeclareContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_declare
	return p
}

func (*DeclareContext) IsDeclareContext() {}

func NewDeclareContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclareContext {
	var p = new(DeclareContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_declare

	return p
}

func (s *DeclareContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclareContext) GetVar() IVarDeclContext { return s.var_ }

func (s *DeclareContext) GetFn() IFnDeclContext { return s.fn }

func (s *DeclareContext) SetVar(v IVarDeclContext) { s.var_ = v }

func (s *DeclareContext) SetFn(v IFnDeclContext) { s.fn = v }

func (s *DeclareContext) VarDecl() IVarDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVarDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVarDeclContext)
}

func (s *DeclareContext) FnDecl() IFnDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFnDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFnDeclContext)
}

func (s *DeclareContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclareContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeclareContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterDeclare(s)
	}
}

func (s *DeclareContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitDeclare(s)
	}
}

func (s *DeclareContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitDeclare(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Declare() (localctx IDeclareContext) {
	this := p
	_ = this

	localctx = NewDeclareContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, CommandsParserRULE_declare)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(80)
		p.Match(CommandsParserT__1)
	}
	p.SetState(83)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(81)

			var _x = p.VarDecl()

			localctx.(*DeclareContext).var_ = _x
		}

	case 2:
		{
			p.SetState(82)

			var _x = p.FnDecl()

			localctx.(*DeclareContext).fn = _x
		}

	}

	return localctx
}

// IVarDeclContext is an interface to support dynamic dispatch.
type IVarDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetId returns the id rule contexts.
	GetId() IQualIdContext

	// GetT returns the t rule contexts.
	GetT() ITypeContext

	// SetId sets the id rule contexts.
	SetId(IQualIdContext)

	// SetT sets the t rule contexts.
	SetT(ITypeContext)

	// IsVarDeclContext differentiates from other interfaces.
	IsVarDeclContext()
}

type VarDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	id     IQualIdContext
	t      ITypeContext
}

func NewEmptyVarDeclContext() *VarDeclContext {
	var p = new(VarDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_varDecl
	return p
}

func (*VarDeclContext) IsVarDeclContext() {}

func NewVarDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarDeclContext {
	var p = new(VarDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_varDecl

	return p
}

func (s *VarDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *VarDeclContext) GetId() IQualIdContext { return s.id }

func (s *VarDeclContext) GetT() ITypeContext { return s.t }

func (s *VarDeclContext) SetId(v IQualIdContext) { s.id = v }

func (s *VarDeclContext) SetT(v ITypeContext) { s.t = v }

func (s *VarDeclContext) QualId() IQualIdContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualIdContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualIdContext)
}

func (s *VarDeclContext) COLON() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOLON, 0)
}

func (s *VarDeclContext) Type() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *VarDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VarDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VarDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterVarDecl(s)
	}
}

func (s *VarDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitVarDecl(s)
	}
}

func (s *VarDeclContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitVarDecl(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) VarDecl() (localctx IVarDeclContext) {
	this := p
	_ = this

	localctx = NewVarDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, CommandsParserRULE_varDecl)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)

		var _x = p.QualId()

		localctx.(*VarDeclContext).id = _x
	}
	p.SetState(88)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserCOLON {
		{
			p.SetState(86)
			p.Match(CommandsParserCOLON)
		}
		{
			p.SetState(87)

			var _x = p.Type()

			localctx.(*VarDeclContext).t = _x
		}

	}

	return localctx
}

// IFnDeclContext is an interface to support dynamic dispatch.
type IFnDeclContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetId returns the id rule contexts.
	GetId() IQualIdContext

	// Get_param returns the _param rule contexts.
	Get_param() IParamContext

	// GetRType returns the rType rule contexts.
	GetRType() ITypeContext

	// SetId sets the id rule contexts.
	SetId(IQualIdContext)

	// Set_param sets the _param rule contexts.
	Set_param(IParamContext)

	// SetRType sets the rType rule contexts.
	SetRType(ITypeContext)

	// GetParams returns the params rule context list.
	GetParams() []IParamContext

	// SetParams sets the params rule context list.
	SetParams([]IParamContext)

	// IsFnDeclContext differentiates from other interfaces.
	IsFnDeclContext()
}

type FnDeclContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	id     IQualIdContext
	_param IParamContext
	params []IParamContext
	rType  ITypeContext
}

func NewEmptyFnDeclContext() *FnDeclContext {
	var p = new(FnDeclContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_fnDecl
	return p
}

func (*FnDeclContext) IsFnDeclContext() {}

func NewFnDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FnDeclContext {
	var p = new(FnDeclContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_fnDecl

	return p
}

func (s *FnDeclContext) GetParser() antlr.Parser { return s.parser }

func (s *FnDeclContext) GetId() IQualIdContext { return s.id }

func (s *FnDeclContext) Get_param() IParamContext { return s._param }

func (s *FnDeclContext) GetRType() ITypeContext { return s.rType }

func (s *FnDeclContext) SetId(v IQualIdContext) { s.id = v }

func (s *FnDeclContext) Set_param(v IParamContext) { s._param = v }

func (s *FnDeclContext) SetRType(v ITypeContext) { s.rType = v }

func (s *FnDeclContext) GetParams() []IParamContext { return s.params }

func (s *FnDeclContext) SetParams(v []IParamContext) { s.params = v }

func (s *FnDeclContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserLPAREN, 0)
}

func (s *FnDeclContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserRPAREN, 0)
}

func (s *FnDeclContext) COLON() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOLON, 0)
}

func (s *FnDeclContext) QualId() IQualIdContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualIdContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualIdContext)
}

func (s *FnDeclContext) Type() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *FnDeclContext) AllParam() []IParamContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IParamContext); ok {
			len++
		}
	}

	tst := make([]IParamContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IParamContext); ok {
			tst[i] = t.(IParamContext)
			i++
		}
	}

	return tst
}

func (s *FnDeclContext) Param(i int) IParamContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *FnDeclContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOMMA)
}

func (s *FnDeclContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, i)
}

func (s *FnDeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FnDeclContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FnDeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterFnDecl(s)
	}
}

func (s *FnDeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitFnDecl(s)
	}
}

func (s *FnDeclContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitFnDecl(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) FnDecl() (localctx IFnDeclContext) {
	this := p
	_ = this

	localctx = NewFnDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, CommandsParserRULE_fnDecl)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(90)

		var _x = p.QualId()

		localctx.(*FnDeclContext).id = _x
	}
	{
		p.SetState(91)
		p.Match(CommandsParserLPAREN)
	}
	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserIDENTIFIER {
		{
			p.SetState(92)

			var _x = p.Param()

			localctx.(*FnDeclContext)._param = _x
		}
		localctx.(*FnDeclContext).params = append(localctx.(*FnDeclContext).params, localctx.(*FnDeclContext)._param)
		p.SetState(97)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == CommandsParserCOMMA {
			{
				p.SetState(93)
				p.Match(CommandsParserCOMMA)
			}
			{
				p.SetState(94)

				var _x = p.Param()

				localctx.(*FnDeclContext)._param = _x
			}
			localctx.(*FnDeclContext).params = append(localctx.(*FnDeclContext).params, localctx.(*FnDeclContext)._param)

			p.SetState(99)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(102)
		p.Match(CommandsParserRPAREN)
	}
	{
		p.SetState(103)
		p.Match(CommandsParserCOLON)
	}
	{
		p.SetState(104)

		var _x = p.Type()

		localctx.(*FnDeclContext).rType = _x
	}

	return localctx
}

// IParamContext is an interface to support dynamic dispatch.
type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetPid returns the pid token.
	GetPid() antlr.Token

	// SetPid sets the pid token.
	SetPid(antlr.Token)

	// GetT returns the t rule contexts.
	GetT() ITypeContext

	// SetT sets the t rule contexts.
	SetT(ITypeContext)

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}

type ParamContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	pid    antlr.Token
	t      ITypeContext
}

func NewEmptyParamContext() *ParamContext {
	var p = new(ParamContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_param
	return p
}

func (*ParamContext) IsParamContext() {}

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext {
	var p = new(ParamContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_param

	return p
}

func (s *ParamContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamContext) GetPid() antlr.Token { return s.pid }

func (s *ParamContext) SetPid(v antlr.Token) { s.pid = v }

func (s *ParamContext) GetT() ITypeContext { return s.t }

func (s *ParamContext) SetT(v ITypeContext) { s.t = v }

func (s *ParamContext) COLON() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOLON, 0)
}

func (s *ParamContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, 0)
}

func (s *ParamContext) Type() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *ParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterParam(s)
	}
}

func (s *ParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitParam(s)
	}
}

func (s *ParamContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitParam(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Param() (localctx IParamContext) {
	this := p
	_ = this

	localctx = NewParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, CommandsParserRULE_param)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(106)

		var _m = p.Match(CommandsParserIDENTIFIER)

		localctx.(*ParamContext).pid = _m
	}
	{
		p.SetState(107)
		p.Match(CommandsParserCOLON)
	}
	{
		p.SetState(108)

		var _x = p.Type()

		localctx.(*ParamContext).t = _x
	}

	return localctx
}

// IDeleteContext is an interface to support dynamic dispatch.
type IDeleteContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVar returns the var rule contexts.
	GetVar() IVarDeclContext

	// GetFn returns the fn rule contexts.
	GetFn() IFnDeclContext

	// SetVar sets the var rule contexts.
	SetVar(IVarDeclContext)

	// SetFn sets the fn rule contexts.
	SetFn(IFnDeclContext)

	// IsDeleteContext differentiates from other interfaces.
	IsDeleteContext()
}

type DeleteContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	var_   IVarDeclContext
	fn     IFnDeclContext
}

func NewEmptyDeleteContext() *DeleteContext {
	var p = new(DeleteContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_delete
	return p
}

func (*DeleteContext) IsDeleteContext() {}

func NewDeleteContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeleteContext {
	var p = new(DeleteContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_delete

	return p
}

func (s *DeleteContext) GetParser() antlr.Parser { return s.parser }

func (s *DeleteContext) GetVar() IVarDeclContext { return s.var_ }

func (s *DeleteContext) GetFn() IFnDeclContext { return s.fn }

func (s *DeleteContext) SetVar(v IVarDeclContext) { s.var_ = v }

func (s *DeleteContext) SetFn(v IFnDeclContext) { s.fn = v }

func (s *DeleteContext) VarDecl() IVarDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVarDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVarDeclContext)
}

func (s *DeleteContext) FnDecl() IFnDeclContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFnDeclContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFnDeclContext)
}

func (s *DeleteContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeleteContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeleteContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterDelete(s)
	}
}

func (s *DeleteContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitDelete(s)
	}
}

func (s *DeleteContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitDelete(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Delete() (localctx IDeleteContext) {
	this := p
	_ = this

	localctx = NewDeleteContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, CommandsParserRULE_delete)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(110)
		p.Match(CommandsParserT__2)
	}
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(111)

			var _x = p.VarDecl()

			localctx.(*DeleteContext).var_ = _x
		}

	case 2:
		{
			p.SetState(112)

			var _x = p.FnDecl()

			localctx.(*DeleteContext).fn = _x
		}

	}

	return localctx
}

// ISimpleContext is an interface to support dynamic dispatch.
type ISimpleContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetCmd returns the cmd token.
	GetCmd() antlr.Token

	// Get_FLAG returns the _FLAG token.
	Get_FLAG() antlr.Token

	// Get_STRING returns the _STRING token.
	Get_STRING() antlr.Token

	// SetCmd sets the cmd token.
	SetCmd(antlr.Token)

	// Set_FLAG sets the _FLAG token.
	Set_FLAG(antlr.Token)

	// Set_STRING sets the _STRING token.
	Set_STRING(antlr.Token)

	// GetArgs returns the args token list.
	GetArgs() []antlr.Token

	// SetArgs sets the args token list.
	SetArgs([]antlr.Token)

	// IsSimpleContext differentiates from other interfaces.
	IsSimpleContext()
}

type SimpleContext struct {
	*antlr.BaseParserRuleContext
	parser  antlr.Parser
	cmd     antlr.Token
	_FLAG   antlr.Token
	args    []antlr.Token
	_STRING antlr.Token
}

func NewEmptySimpleContext() *SimpleContext {
	var p = new(SimpleContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_simple
	return p
}

func (*SimpleContext) IsSimpleContext() {}

func NewSimpleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SimpleContext {
	var p = new(SimpleContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_simple

	return p
}

func (s *SimpleContext) GetParser() antlr.Parser { return s.parser }

func (s *SimpleContext) GetCmd() antlr.Token { return s.cmd }

func (s *SimpleContext) Get_FLAG() antlr.Token { return s._FLAG }

func (s *SimpleContext) Get_STRING() antlr.Token { return s._STRING }

func (s *SimpleContext) SetCmd(v antlr.Token) { s.cmd = v }

func (s *SimpleContext) Set_FLAG(v antlr.Token) { s._FLAG = v }

func (s *SimpleContext) Set_STRING(v antlr.Token) { s._STRING = v }

func (s *SimpleContext) GetArgs() []antlr.Token { return s.args }

func (s *SimpleContext) SetArgs(v []antlr.Token) { s.args = v }

func (s *SimpleContext) COMMAND() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMAND, 0)
}

func (s *SimpleContext) AllFLAG() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserFLAG)
}

func (s *SimpleContext) FLAG(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserFLAG, i)
}

func (s *SimpleContext) AllSTRING() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserSTRING)
}

func (s *SimpleContext) STRING(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserSTRING, i)
}

func (s *SimpleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SimpleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterSimple(s)
	}
}

func (s *SimpleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitSimple(s)
	}
}

func (s *SimpleContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitSimple(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Simple() (localctx ISimpleContext) {
	this := p
	_ = this

	localctx = NewSimpleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, CommandsParserRULE_simple)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(115)

		var _m = p.Match(CommandsParserCOMMAND)

		localctx.(*SimpleContext).cmd = _m
	}
	p.SetState(120)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserFLAG || _la == CommandsParserSTRING {
		p.SetState(118)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case CommandsParserFLAG:
			{
				p.SetState(116)

				var _m = p.Match(CommandsParserFLAG)

				localctx.(*SimpleContext)._FLAG = _m
			}
			localctx.(*SimpleContext).args = append(localctx.(*SimpleContext).args, localctx.(*SimpleContext)._FLAG)

		case CommandsParserSTRING:
			{
				p.SetState(117)

				var _m = p.Match(CommandsParserSTRING)

				localctx.(*SimpleContext)._STRING = _m
			}
			localctx.(*SimpleContext).args = append(localctx.(*SimpleContext).args, localctx.(*SimpleContext)._STRING)

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(122)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IEmptyContext is an interface to support dynamic dispatch.
type IEmptyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEmptyContext differentiates from other interfaces.
	IsEmptyContext()
}

type EmptyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEmptyContext() *EmptyContext {
	var p = new(EmptyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_empty
	return p
}

func (*EmptyContext) IsEmptyContext() {}

func NewEmptyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EmptyContext {
	var p = new(EmptyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_empty

	return p
}

func (s *EmptyContext) GetParser() antlr.Parser { return s.parser }
func (s *EmptyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EmptyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EmptyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterEmpty(s)
	}
}

func (s *EmptyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitEmpty(s)
	}
}

func (s *EmptyContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitEmpty(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Empty() (localctx IEmptyContext) {
	this := p
	_ = this

	localctx = NewEmptyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, CommandsParserRULE_empty)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)

	return localctx
}

// IExprCmdContext is an interface to support dynamic dispatch.
type IExprCmdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetE returns the e rule contexts.
	GetE() IExprContext

	// SetE sets the e rule contexts.
	SetE(IExprContext)

	// IsExprCmdContext differentiates from other interfaces.
	IsExprCmdContext()
}

type ExprCmdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IExprContext
}

func NewEmptyExprCmdContext() *ExprCmdContext {
	var p = new(ExprCmdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_exprCmd
	return p
}

func (*ExprCmdContext) IsExprCmdContext() {}

func NewExprCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprCmdContext {
	var p = new(ExprCmdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_exprCmd

	return p
}

func (s *ExprCmdContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprCmdContext) GetE() IExprContext { return s.e }

func (s *ExprCmdContext) SetE(v IExprContext) { s.e = v }

func (s *ExprCmdContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprCmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprCmdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprCmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterExprCmd(s)
	}
}

func (s *ExprCmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitExprCmd(s)
	}
}

func (s *ExprCmdContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitExprCmd(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ExprCmd() (localctx IExprCmdContext) {
	this := p
	_ = this

	localctx = NewExprCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, CommandsParserRULE_exprCmd)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserT__3 {
		{
			p.SetState(125)
			p.Match(CommandsParserT__3)
		}

	}
	{
		p.SetState(128)

		var _x = p.Expr()

		localctx.(*ExprCmdContext).e = _x
	}

	return localctx
}

// IQualIdContext is an interface to support dynamic dispatch.
type IQualIdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeadingDot returns the leadingDot token.
	GetLeadingDot() antlr.Token

	// GetRid returns the rid token.
	GetRid() antlr.Token

	// Get_IDENTIFIER returns the _IDENTIFIER token.
	Get_IDENTIFIER() antlr.Token

	// SetLeadingDot sets the leadingDot token.
	SetLeadingDot(antlr.Token)

	// SetRid sets the rid token.
	SetRid(antlr.Token)

	// Set_IDENTIFIER sets the _IDENTIFIER token.
	Set_IDENTIFIER(antlr.Token)

	// GetQualifiers returns the qualifiers token list.
	GetQualifiers() []antlr.Token

	// SetQualifiers sets the qualifiers token list.
	SetQualifiers([]antlr.Token)

	// IsQualIdContext differentiates from other interfaces.
	IsQualIdContext()
}

type QualIdContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	leadingDot  antlr.Token
	rid         antlr.Token
	_IDENTIFIER antlr.Token
	qualifiers  []antlr.Token
}

func NewEmptyQualIdContext() *QualIdContext {
	var p = new(QualIdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_qualId
	return p
}

func (*QualIdContext) IsQualIdContext() {}

func NewQualIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualIdContext {
	var p = new(QualIdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_qualId

	return p
}

func (s *QualIdContext) GetParser() antlr.Parser { return s.parser }

func (s *QualIdContext) GetLeadingDot() antlr.Token { return s.leadingDot }

func (s *QualIdContext) GetRid() antlr.Token { return s.rid }

func (s *QualIdContext) Get_IDENTIFIER() antlr.Token { return s._IDENTIFIER }

func (s *QualIdContext) SetLeadingDot(v antlr.Token) { s.leadingDot = v }

func (s *QualIdContext) SetRid(v antlr.Token) { s.rid = v }

func (s *QualIdContext) Set_IDENTIFIER(v antlr.Token) { s._IDENTIFIER = v }

func (s *QualIdContext) GetQualifiers() []antlr.Token { return s.qualifiers }

func (s *QualIdContext) SetQualifiers(v []antlr.Token) { s.qualifiers = v }

func (s *QualIdContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserIDENTIFIER)
}

func (s *QualIdContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, i)
}

func (s *QualIdContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserDOT)
}

func (s *QualIdContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserDOT, i)
}

func (s *QualIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QualIdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QualIdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterQualId(s)
	}
}

func (s *QualIdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitQualId(s)
	}
}

func (s *QualIdContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitQualId(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) QualId() (localctx IQualIdContext) {
	this := p
	_ = this

	localctx = NewQualIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, CommandsParserRULE_qualId)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserDOT {
		{
			p.SetState(130)

			var _m = p.Match(CommandsParserDOT)

			localctx.(*QualIdContext).leadingDot = _m
		}

	}
	{
		p.SetState(133)

		var _m = p.Match(CommandsParserIDENTIFIER)

		localctx.(*QualIdContext).rid = _m
	}
	p.SetState(138)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserDOT {
		{
			p.SetState(134)
			p.Match(CommandsParserDOT)
		}
		{
			p.SetState(135)

			var _m = p.Match(CommandsParserIDENTIFIER)

			localctx.(*QualIdContext)._IDENTIFIER = _m
		}
		localctx.(*QualIdContext).qualifiers = append(localctx.(*QualIdContext).qualifiers, localctx.(*QualIdContext)._IDENTIFIER)

		p.SetState(140)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IStartTypeContext is an interface to support dynamic dispatch.
type IStartTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetT returns the t rule contexts.
	GetT() ITypeContext

	// SetT sets the t rule contexts.
	SetT(ITypeContext)

	// IsStartTypeContext differentiates from other interfaces.
	IsStartTypeContext()
}

type StartTypeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	t      ITypeContext
}

func NewEmptyStartTypeContext() *StartTypeContext {
	var p = new(StartTypeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_startType
	return p
}

func (*StartTypeContext) IsStartTypeContext() {}

func NewStartTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartTypeContext {
	var p = new(StartTypeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_startType

	return p
}

func (s *StartTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *StartTypeContext) GetT() ITypeContext { return s.t }

func (s *StartTypeContext) SetT(v ITypeContext) { s.t = v }

func (s *StartTypeContext) EOF() antlr.TerminalNode {
	return s.GetToken(CommandsParserEOF, 0)
}

func (s *StartTypeContext) Type() ITypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *StartTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterStartType(s)
	}
}

func (s *StartTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitStartType(s)
	}
}

func (s *StartTypeContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitStartType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) StartType() (localctx IStartTypeContext) {
	this := p
	_ = this

	localctx = NewStartTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, CommandsParserRULE_startType)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(141)

		var _x = p.Type()

		localctx.(*StartTypeContext).t = _x
	}
	{
		p.SetState(142)
		p.Match(CommandsParserEOF)
	}

	return localctx
}

// ITypeContext is an interface to support dynamic dispatch.
type ITypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetId returns the id rule contexts.
	GetId() ITypeIdContext

	// GetParams returns the params rule contexts.
	GetParams() ITypeParamListContext

	// SetId sets the id rule contexts.
	SetId(ITypeIdContext)

	// SetParams sets the params rule contexts.
	SetParams(ITypeParamListContext)

	// IsTypeContext differentiates from other interfaces.
	IsTypeContext()
}

type TypeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	id     ITypeIdContext
	params ITypeParamListContext
}

func NewEmptyTypeContext() *TypeContext {
	var p = new(TypeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_type
	return p
}

func (*TypeContext) IsTypeContext() {}

func NewTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeContext {
	var p = new(TypeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_type

	return p
}

func (s *TypeContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeContext) GetId() ITypeIdContext { return s.id }

func (s *TypeContext) GetParams() ITypeParamListContext { return s.params }

func (s *TypeContext) SetId(v ITypeIdContext) { s.id = v }

func (s *TypeContext) SetParams(v ITypeParamListContext) { s.params = v }

func (s *TypeContext) TypeId() ITypeIdContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeIdContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeIdContext)
}

func (s *TypeContext) TypeParamList() ITypeParamListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeParamListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeParamListContext)
}

func (s *TypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterType(s)
	}
}

func (s *TypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitType(s)
	}
}

func (s *TypeContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Type() (localctx ITypeContext) {
	this := p
	_ = this

	localctx = NewTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, CommandsParserRULE_type)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(144)

		var _x = p.TypeId()

		localctx.(*TypeContext).id = _x
	}
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserLPAREN {
		{
			p.SetState(145)

			var _x = p.TypeParamList()

			localctx.(*TypeContext).params = _x
		}

	}

	return localctx
}

// ITypeIdContext is an interface to support dynamic dispatch.
type ITypeIdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetLeadingDot returns the leadingDot token.
	GetLeadingDot() antlr.Token

	// GetId returns the id token.
	GetId() antlr.Token

	// Get_IDENTIFIER returns the _IDENTIFIER token.
	Get_IDENTIFIER() antlr.Token

	// SetLeadingDot sets the leadingDot token.
	SetLeadingDot(antlr.Token)

	// SetId sets the id token.
	SetId(antlr.Token)

	// Set_IDENTIFIER sets the _IDENTIFIER token.
	Set_IDENTIFIER(antlr.Token)

	// GetQualifiers returns the qualifiers token list.
	GetQualifiers() []antlr.Token

	// SetQualifiers sets the qualifiers token list.
	SetQualifiers([]antlr.Token)

	// IsTypeIdContext differentiates from other interfaces.
	IsTypeIdContext()
}

type TypeIdContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	leadingDot  antlr.Token
	id          antlr.Token
	_IDENTIFIER antlr.Token
	qualifiers  []antlr.Token
}

func NewEmptyTypeIdContext() *TypeIdContext {
	var p = new(TypeIdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_typeId
	return p
}

func (*TypeIdContext) IsTypeIdContext() {}

func NewTypeIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeIdContext {
	var p = new(TypeIdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_typeId

	return p
}

func (s *TypeIdContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeIdContext) GetLeadingDot() antlr.Token { return s.leadingDot }

func (s *TypeIdContext) GetId() antlr.Token { return s.id }

func (s *TypeIdContext) Get_IDENTIFIER() antlr.Token { return s._IDENTIFIER }

func (s *TypeIdContext) SetLeadingDot(v antlr.Token) { s.leadingDot = v }

func (s *TypeIdContext) SetId(v antlr.Token) { s.id = v }

func (s *TypeIdContext) Set_IDENTIFIER(v antlr.Token) { s._IDENTIFIER = v }

func (s *TypeIdContext) GetQualifiers() []antlr.Token { return s.qualifiers }

func (s *TypeIdContext) SetQualifiers(v []antlr.Token) { s.qualifiers = v }

func (s *TypeIdContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserIDENTIFIER)
}

func (s *TypeIdContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, i)
}

func (s *TypeIdContext) NUL() antlr.TerminalNode {
	return s.GetToken(CommandsParserNUL, 0)
}

func (s *TypeIdContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserDOT)
}

func (s *TypeIdContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserDOT, i)
}

func (s *TypeIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeIdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeIdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterTypeId(s)
	}
}

func (s *TypeIdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitTypeId(s)
	}
}

func (s *TypeIdContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitTypeId(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) TypeId() (localctx ITypeIdContext) {
	this := p
	_ = this

	localctx = NewTypeIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, CommandsParserRULE_typeId)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(149)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserDOT {
		{
			p.SetState(148)

			var _m = p.Match(CommandsParserDOT)

			localctx.(*TypeIdContext).leadingDot = _m
		}

	}
	{
		p.SetState(151)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*TypeIdContext).id = _lt

		_la = p.GetTokenStream().LA(1)

		if !(_la == CommandsParserNUL || _la == CommandsParserIDENTIFIER) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*TypeIdContext).id = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserDOT {
		{
			p.SetState(152)
			p.Match(CommandsParserDOT)
		}
		{
			p.SetState(153)

			var _m = p.Match(CommandsParserIDENTIFIER)

			localctx.(*TypeIdContext)._IDENTIFIER = _m
		}
		localctx.(*TypeIdContext).qualifiers = append(localctx.(*TypeIdContext).qualifiers, localctx.(*TypeIdContext)._IDENTIFIER)

		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ITypeParamListContext is an interface to support dynamic dispatch.
type ITypeParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_type returns the _type rule contexts.
	Get_type() ITypeContext

	// Set_type sets the _type rule contexts.
	Set_type(ITypeContext)

	// GetTypes returns the types rule context list.
	GetTypes() []ITypeContext

	// SetTypes sets the types rule context list.
	SetTypes([]ITypeContext)

	// IsTypeParamListContext differentiates from other interfaces.
	IsTypeParamListContext()
}

type TypeParamListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	_type  ITypeContext
	types  []ITypeContext
}

func NewEmptyTypeParamListContext() *TypeParamListContext {
	var p = new(TypeParamListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_typeParamList
	return p
}

func (*TypeParamListContext) IsTypeParamListContext() {}

func NewTypeParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamListContext {
	var p = new(TypeParamListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_typeParamList

	return p
}

func (s *TypeParamListContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeParamListContext) Get_type() ITypeContext { return s._type }

func (s *TypeParamListContext) Set_type(v ITypeContext) { s._type = v }

func (s *TypeParamListContext) GetTypes() []ITypeContext { return s.types }

func (s *TypeParamListContext) SetTypes(v []ITypeContext) { s.types = v }

func (s *TypeParamListContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserLPAREN, 0)
}

func (s *TypeParamListContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserRPAREN, 0)
}

func (s *TypeParamListContext) AllType() []ITypeContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeContext); ok {
			len++
		}
	}

	tst := make([]ITypeContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeContext); ok {
			tst[i] = t.(ITypeContext)
			i++
		}
	}

	return tst
}

func (s *TypeParamListContext) Type(i int) ITypeContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeContext)
}

func (s *TypeParamListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOMMA)
}

func (s *TypeParamListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, i)
}

func (s *TypeParamListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeParamListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterTypeParamList(s)
	}
}

func (s *TypeParamListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitTypeParamList(s)
	}
}

func (s *TypeParamListContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitTypeParamList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) TypeParamList() (localctx ITypeParamListContext) {
	this := p
	_ = this

	localctx = NewTypeParamListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, CommandsParserRULE_typeParamList)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(159)
		p.Match(CommandsParserLPAREN)
	}
	{
		p.SetState(160)

		var _x = p.Type()

		localctx.(*TypeParamListContext)._type = _x
	}
	localctx.(*TypeParamListContext).types = append(localctx.(*TypeParamListContext).types, localctx.(*TypeParamListContext)._type)
	p.SetState(165)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserCOMMA {
		{
			p.SetState(161)
			p.Match(CommandsParserCOMMA)
		}
		{
			p.SetState(162)

			var _x = p.Type()

			localctx.(*TypeParamListContext)._type = _x
		}
		localctx.(*TypeParamListContext).types = append(localctx.(*TypeParamListContext).types, localctx.(*TypeParamListContext)._type)

		p.SetState(167)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(168)
		p.Match(CommandsParserRPAREN)
	}

	return localctx
}

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetE returns the e rule contexts.
	GetE() IExprContext

	// SetE sets the e rule contexts.
	SetE(IExprContext)

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IExprContext
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) GetE() IExprContext { return s.e }

func (s *StartContext) SetE(v IExprContext) { s.e = v }

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(CommandsParserEOF, 0)
}

func (s *StartContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitStart(s)
	}
}

func (s *StartContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitStart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Start() (localctx IStartContext) {
	this := p
	_ = this

	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, CommandsParserRULE_start)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(170)

		var _x = p.Expr()

		localctx.(*StartContext).e = _x
	}
	{
		p.SetState(171)
		p.Match(CommandsParserEOF)
	}

	return localctx
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOp returns the op token.
	GetOp() antlr.Token

	// SetOp sets the op token.
	SetOp(antlr.Token)

	// GetE returns the e rule contexts.
	GetE() IConditionalOrContext

	// GetE1 returns the e1 rule contexts.
	GetE1() IConditionalOrContext

	// GetE2 returns the e2 rule contexts.
	GetE2() IExprContext

	// SetE sets the e rule contexts.
	SetE(IConditionalOrContext)

	// SetE1 sets the e1 rule contexts.
	SetE1(IConditionalOrContext)

	// SetE2 sets the e2 rule contexts.
	SetE2(IExprContext)

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IConditionalOrContext
	op     antlr.Token
	e1     IConditionalOrContext
	e2     IExprContext
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) GetOp() antlr.Token { return s.op }

func (s *ExprContext) SetOp(v antlr.Token) { s.op = v }

func (s *ExprContext) GetE() IConditionalOrContext { return s.e }

func (s *ExprContext) GetE1() IConditionalOrContext { return s.e1 }

func (s *ExprContext) GetE2() IExprContext { return s.e2 }

func (s *ExprContext) SetE(v IConditionalOrContext) { s.e = v }

func (s *ExprContext) SetE1(v IConditionalOrContext) { s.e1 = v }

func (s *ExprContext) SetE2(v IExprContext) { s.e2 = v }

func (s *ExprContext) AllConditionalOr() []IConditionalOrContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConditionalOrContext); ok {
			len++
		}
	}

	tst := make([]IConditionalOrContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConditionalOrContext); ok {
			tst[i] = t.(IConditionalOrContext)
			i++
		}
	}

	return tst
}

func (s *ExprContext) ConditionalOr(i int) IConditionalOrContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionalOrContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionalOrContext)
}

func (s *ExprContext) COLON() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOLON, 0)
}

func (s *ExprContext) QUESTIONMARK() antlr.TerminalNode {
	return s.GetToken(CommandsParserQUESTIONMARK, 0)
}

func (s *ExprContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitExpr(s)
	}
}

func (s *ExprContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Expr() (localctx IExprContext) {
	this := p
	_ = this

	localctx = NewExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, CommandsParserRULE_expr)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(173)

		var _x = p.ConditionalOr()

		localctx.(*ExprContext).e = _x
	}
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserQUESTIONMARK {
		{
			p.SetState(174)

			var _m = p.Match(CommandsParserQUESTIONMARK)

			localctx.(*ExprContext).op = _m
		}
		{
			p.SetState(175)

			var _x = p.ConditionalOr()

			localctx.(*ExprContext).e1 = _x
		}
		{
			p.SetState(176)
			p.Match(CommandsParserCOLON)
		}
		{
			p.SetState(177)

			var _x = p.Expr()

			localctx.(*ExprContext).e2 = _x
		}

	}

	return localctx
}

// IConditionalOrContext is an interface to support dynamic dispatch.
type IConditionalOrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS17 returns the s17 token.
	GetS17() antlr.Token

	// SetS17 sets the s17 token.
	SetS17(antlr.Token)

	// GetOps returns the ops token list.
	GetOps() []antlr.Token

	// SetOps sets the ops token list.
	SetOps([]antlr.Token)

	// GetE returns the e rule contexts.
	GetE() IConditionalAndContext

	// Get_conditionalAnd returns the _conditionalAnd rule contexts.
	Get_conditionalAnd() IConditionalAndContext

	// SetE sets the e rule contexts.
	SetE(IConditionalAndContext)

	// Set_conditionalAnd sets the _conditionalAnd rule contexts.
	Set_conditionalAnd(IConditionalAndContext)

	// GetE1 returns the e1 rule context list.
	GetE1() []IConditionalAndContext

	// SetE1 sets the e1 rule context list.
	SetE1([]IConditionalAndContext)

	// IsConditionalOrContext differentiates from other interfaces.
	IsConditionalOrContext()
}

type ConditionalOrContext struct {
	*antlr.BaseParserRuleContext
	parser          antlr.Parser
	e               IConditionalAndContext
	s17             antlr.Token
	ops             []antlr.Token
	_conditionalAnd IConditionalAndContext
	e1              []IConditionalAndContext
}

func NewEmptyConditionalOrContext() *ConditionalOrContext {
	var p = new(ConditionalOrContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_conditionalOr
	return p
}

func (*ConditionalOrContext) IsConditionalOrContext() {}

func NewConditionalOrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalOrContext {
	var p = new(ConditionalOrContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_conditionalOr

	return p
}

func (s *ConditionalOrContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionalOrContext) GetS17() antlr.Token { return s.s17 }

func (s *ConditionalOrContext) SetS17(v antlr.Token) { s.s17 = v }

func (s *ConditionalOrContext) GetOps() []antlr.Token { return s.ops }

func (s *ConditionalOrContext) SetOps(v []antlr.Token) { s.ops = v }

func (s *ConditionalOrContext) GetE() IConditionalAndContext { return s.e }

func (s *ConditionalOrContext) Get_conditionalAnd() IConditionalAndContext { return s._conditionalAnd }

func (s *ConditionalOrContext) SetE(v IConditionalAndContext) { s.e = v }

func (s *ConditionalOrContext) Set_conditionalAnd(v IConditionalAndContext) { s._conditionalAnd = v }

func (s *ConditionalOrContext) GetE1() []IConditionalAndContext { return s.e1 }

func (s *ConditionalOrContext) SetE1(v []IConditionalAndContext) { s.e1 = v }

func (s *ConditionalOrContext) AllConditionalAnd() []IConditionalAndContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConditionalAndContext); ok {
			len++
		}
	}

	tst := make([]IConditionalAndContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConditionalAndContext); ok {
			tst[i] = t.(IConditionalAndContext)
			i++
		}
	}

	return tst
}

func (s *ConditionalOrContext) ConditionalAnd(i int) IConditionalAndContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConditionalAndContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConditionalAndContext)
}

func (s *ConditionalOrContext) AllLOGICAL_OR() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserLOGICAL_OR)
}

func (s *ConditionalOrContext) LOGICAL_OR(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserLOGICAL_OR, i)
}

func (s *ConditionalOrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionalOrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConditionalOrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterConditionalOr(s)
	}
}

func (s *ConditionalOrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitConditionalOr(s)
	}
}

func (s *ConditionalOrContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitConditionalOr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ConditionalOr() (localctx IConditionalOrContext) {
	this := p
	_ = this

	localctx = NewConditionalOrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, CommandsParserRULE_conditionalOr)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(181)

		var _x = p.ConditionalAnd()

		localctx.(*ConditionalOrContext).e = _x
	}
	p.SetState(186)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserLOGICAL_OR {
		{
			p.SetState(182)

			var _m = p.Match(CommandsParserLOGICAL_OR)

			localctx.(*ConditionalOrContext).s17 = _m
		}
		localctx.(*ConditionalOrContext).ops = append(localctx.(*ConditionalOrContext).ops, localctx.(*ConditionalOrContext).s17)
		{
			p.SetState(183)

			var _x = p.ConditionalAnd()

			localctx.(*ConditionalOrContext)._conditionalAnd = _x
		}
		localctx.(*ConditionalOrContext).e1 = append(localctx.(*ConditionalOrContext).e1, localctx.(*ConditionalOrContext)._conditionalAnd)

		p.SetState(188)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IConditionalAndContext is an interface to support dynamic dispatch.
type IConditionalAndContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS16 returns the s16 token.
	GetS16() antlr.Token

	// SetS16 sets the s16 token.
	SetS16(antlr.Token)

	// GetOps returns the ops token list.
	GetOps() []antlr.Token

	// SetOps sets the ops token list.
	SetOps([]antlr.Token)

	// GetE returns the e rule contexts.
	GetE() IRelationContext

	// Get_relation returns the _relation rule contexts.
	Get_relation() IRelationContext

	// SetE sets the e rule contexts.
	SetE(IRelationContext)

	// Set_relation sets the _relation rule contexts.
	Set_relation(IRelationContext)

	// GetE1 returns the e1 rule context list.
	GetE1() []IRelationContext

	// SetE1 sets the e1 rule context list.
	SetE1([]IRelationContext)

	// IsConditionalAndContext differentiates from other interfaces.
	IsConditionalAndContext()
}

type ConditionalAndContext struct {
	*antlr.BaseParserRuleContext
	parser    antlr.Parser
	e         IRelationContext
	s16       antlr.Token
	ops       []antlr.Token
	_relation IRelationContext
	e1        []IRelationContext
}

func NewEmptyConditionalAndContext() *ConditionalAndContext {
	var p = new(ConditionalAndContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_conditionalAnd
	return p
}

func (*ConditionalAndContext) IsConditionalAndContext() {}

func NewConditionalAndContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalAndContext {
	var p = new(ConditionalAndContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_conditionalAnd

	return p
}

func (s *ConditionalAndContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionalAndContext) GetS16() antlr.Token { return s.s16 }

func (s *ConditionalAndContext) SetS16(v antlr.Token) { s.s16 = v }

func (s *ConditionalAndContext) GetOps() []antlr.Token { return s.ops }

func (s *ConditionalAndContext) SetOps(v []antlr.Token) { s.ops = v }

func (s *ConditionalAndContext) GetE() IRelationContext { return s.e }

func (s *ConditionalAndContext) Get_relation() IRelationContext { return s._relation }

func (s *ConditionalAndContext) SetE(v IRelationContext) { s.e = v }

func (s *ConditionalAndContext) Set_relation(v IRelationContext) { s._relation = v }

func (s *ConditionalAndContext) GetE1() []IRelationContext { return s.e1 }

func (s *ConditionalAndContext) SetE1(v []IRelationContext) { s.e1 = v }

func (s *ConditionalAndContext) AllRelation() []IRelationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelationContext); ok {
			len++
		}
	}

	tst := make([]IRelationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelationContext); ok {
			tst[i] = t.(IRelationContext)
			i++
		}
	}

	return tst
}

func (s *ConditionalAndContext) Relation(i int) IRelationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *ConditionalAndContext) AllLOGICAL_AND() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserLOGICAL_AND)
}

func (s *ConditionalAndContext) LOGICAL_AND(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserLOGICAL_AND, i)
}

func (s *ConditionalAndContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionalAndContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConditionalAndContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterConditionalAnd(s)
	}
}

func (s *ConditionalAndContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitConditionalAnd(s)
	}
}

func (s *ConditionalAndContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitConditionalAnd(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ConditionalAnd() (localctx IConditionalAndContext) {
	this := p
	_ = this

	localctx = NewConditionalAndContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, CommandsParserRULE_conditionalAnd)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(189)

		var _x = p.relation(0)

		localctx.(*ConditionalAndContext).e = _x
	}
	p.SetState(194)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserLOGICAL_AND {
		{
			p.SetState(190)

			var _m = p.Match(CommandsParserLOGICAL_AND)

			localctx.(*ConditionalAndContext).s16 = _m
		}
		localctx.(*ConditionalAndContext).ops = append(localctx.(*ConditionalAndContext).ops, localctx.(*ConditionalAndContext).s16)
		{
			p.SetState(191)

			var _x = p.relation(0)

			localctx.(*ConditionalAndContext)._relation = _x
		}
		localctx.(*ConditionalAndContext).e1 = append(localctx.(*ConditionalAndContext).e1, localctx.(*ConditionalAndContext)._relation)

		p.SetState(196)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IRelationContext is an interface to support dynamic dispatch.
type IRelationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOp returns the op token.
	GetOp() antlr.Token

	// SetOp sets the op token.
	SetOp(antlr.Token)

	// IsRelationContext differentiates from other interfaces.
	IsRelationContext()
}

type RelationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	op     antlr.Token
}

func NewEmptyRelationContext() *RelationContext {
	var p = new(RelationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_relation
	return p
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_relation

	return p
}

func (s *RelationContext) GetParser() antlr.Parser { return s.parser }

func (s *RelationContext) GetOp() antlr.Token { return s.op }

func (s *RelationContext) SetOp(v antlr.Token) { s.op = v }

func (s *RelationContext) Calc() ICalcContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICalcContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICalcContext)
}

func (s *RelationContext) AllRelation() []IRelationContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRelationContext); ok {
			len++
		}
	}

	tst := make([]IRelationContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRelationContext); ok {
			tst[i] = t.(IRelationContext)
			i++
		}
	}

	return tst
}

func (s *RelationContext) Relation(i int) IRelationContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRelationContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRelationContext)
}

func (s *RelationContext) LESS() antlr.TerminalNode {
	return s.GetToken(CommandsParserLESS, 0)
}

func (s *RelationContext) LESS_EQUALS() antlr.TerminalNode {
	return s.GetToken(CommandsParserLESS_EQUALS, 0)
}

func (s *RelationContext) GREATER_EQUALS() antlr.TerminalNode {
	return s.GetToken(CommandsParserGREATER_EQUALS, 0)
}

func (s *RelationContext) GREATER() antlr.TerminalNode {
	return s.GetToken(CommandsParserGREATER, 0)
}

func (s *RelationContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(CommandsParserEQUALS, 0)
}

func (s *RelationContext) NOT_EQUALS() antlr.TerminalNode {
	return s.GetToken(CommandsParserNOT_EQUALS, 0)
}

func (s *RelationContext) IN() antlr.TerminalNode {
	return s.GetToken(CommandsParserIN, 0)
}

func (s *RelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RelationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterRelation(s)
	}
}

func (s *RelationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitRelation(s)
	}
}

func (s *RelationContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitRelation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Relation() (localctx IRelationContext) {
	return p.relation(0)
}

func (p *CommandsParser) relation(_p int) (localctx IRelationContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewRelationContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IRelationContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 40
	p.EnterRecursionRule(localctx, 40, CommandsParserRULE_relation, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(198)
		p.calc(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(205)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewRelationContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_relation)
			p.SetState(200)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(201)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*RelationContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CommandsParserEQUALS)|(1<<CommandsParserNOT_EQUALS)|(1<<CommandsParserIN)|(1<<CommandsParserLESS)|(1<<CommandsParserLESS_EQUALS)|(1<<CommandsParserGREATER_EQUALS)|(1<<CommandsParserGREATER))) != 0) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*RelationContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(202)
				p.relation(2)
			}

		}
		p.SetState(207)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext())
	}

	return localctx
}

// ICalcContext is an interface to support dynamic dispatch.
type ICalcContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOp returns the op token.
	GetOp() antlr.Token

	// SetOp sets the op token.
	SetOp(antlr.Token)

	// IsCalcContext differentiates from other interfaces.
	IsCalcContext()
}

type CalcContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	op     antlr.Token
}

func NewEmptyCalcContext() *CalcContext {
	var p = new(CalcContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_calc
	return p
}

func (*CalcContext) IsCalcContext() {}

func NewCalcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CalcContext {
	var p = new(CalcContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_calc

	return p
}

func (s *CalcContext) GetParser() antlr.Parser { return s.parser }

func (s *CalcContext) GetOp() antlr.Token { return s.op }

func (s *CalcContext) SetOp(v antlr.Token) { s.op = v }

func (s *CalcContext) Unary() IUnaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnaryContext)
}

func (s *CalcContext) AllCalc() []ICalcContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICalcContext); ok {
			len++
		}
	}

	tst := make([]ICalcContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICalcContext); ok {
			tst[i] = t.(ICalcContext)
			i++
		}
	}

	return tst
}

func (s *CalcContext) Calc(i int) ICalcContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICalcContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICalcContext)
}

func (s *CalcContext) STAR() antlr.TerminalNode {
	return s.GetToken(CommandsParserSTAR, 0)
}

func (s *CalcContext) SLASH() antlr.TerminalNode {
	return s.GetToken(CommandsParserSLASH, 0)
}

func (s *CalcContext) PERCENT() antlr.TerminalNode {
	return s.GetToken(CommandsParserPERCENT, 0)
}

func (s *CalcContext) PLUS() antlr.TerminalNode {
	return s.GetToken(CommandsParserPLUS, 0)
}

func (s *CalcContext) MINUS() antlr.TerminalNode {
	return s.GetToken(CommandsParserMINUS, 0)
}

func (s *CalcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CalcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CalcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterCalc(s)
	}
}

func (s *CalcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitCalc(s)
	}
}

func (s *CalcContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCalc(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Calc() (localctx ICalcContext) {
	return p.calc(0)
}

func (p *CommandsParser) calc(_p int) (localctx ICalcContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewCalcContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ICalcContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 42
	p.EnterRecursionRule(localctx, 42, CommandsParserRULE_calc, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(209)
		p.Unary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(219)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(217)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) {
			case 1:
				localctx = NewCalcContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_calc)
				p.SetState(211)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(212)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*CalcContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(((_la-31)&-(0x1f+1)) == 0 && ((1<<uint((_la-31)))&((1<<(CommandsParserSTAR-31))|(1<<(CommandsParserSLASH-31))|(1<<(CommandsParserPERCENT-31)))) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*CalcContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(213)
					p.calc(3)
				}

			case 2:
				localctx = NewCalcContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_calc)
				p.SetState(214)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(215)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*CalcContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == CommandsParserMINUS || _la == CommandsParserPLUS) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*CalcContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(216)
					p.calc(2)
				}

			}

		}
		p.SetState(221)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext())
	}

	return localctx
}

// IUnaryContext is an interface to support dynamic dispatch.
type IUnaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUnaryContext differentiates from other interfaces.
	IsUnaryContext()
}

type UnaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnaryContext() *UnaryContext {
	var p = new(UnaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_unary
	return p
}

func (*UnaryContext) IsUnaryContext() {}

func NewUnaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryContext {
	var p = new(UnaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_unary

	return p
}

func (s *UnaryContext) GetParser() antlr.Parser { return s.parser }

func (s *UnaryContext) CopyFrom(ctx *UnaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *UnaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LogicalNotContext struct {
	*UnaryContext
	s27 antlr.Token
	ops []antlr.Token
}

func NewLogicalNotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalNotContext {
	var p = new(LogicalNotContext)

	p.UnaryContext = NewEmptyUnaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*UnaryContext))

	return p
}

func (s *LogicalNotContext) GetS27() antlr.Token { return s.s27 }

func (s *LogicalNotContext) SetS27(v antlr.Token) { s.s27 = v }

func (s *LogicalNotContext) GetOps() []antlr.Token { return s.ops }

func (s *LogicalNotContext) SetOps(v []antlr.Token) { s.ops = v }

func (s *LogicalNotContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalNotContext) Member() IMemberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberContext)
}

func (s *LogicalNotContext) AllEXCLAM() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserEXCLAM)
}

func (s *LogicalNotContext) EXCLAM(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserEXCLAM, i)
}

func (s *LogicalNotContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterLogicalNot(s)
	}
}

func (s *LogicalNotContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitLogicalNot(s)
	}
}

func (s *LogicalNotContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitLogicalNot(s)

	default:
		return t.VisitChildren(s)
	}
}

type MemberExprContext struct {
	*UnaryContext
}

func NewMemberExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberExprContext {
	var p = new(MemberExprContext)

	p.UnaryContext = NewEmptyUnaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*UnaryContext))

	return p
}

func (s *MemberExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberExprContext) Member() IMemberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberContext)
}

func (s *MemberExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterMemberExpr(s)
	}
}

func (s *MemberExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitMemberExpr(s)
	}
}

func (s *MemberExprContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitMemberExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type NegateContext struct {
	*UnaryContext
	s26 antlr.Token
	ops []antlr.Token
}

func NewNegateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NegateContext {
	var p = new(NegateContext)

	p.UnaryContext = NewEmptyUnaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*UnaryContext))

	return p
}

func (s *NegateContext) GetS26() antlr.Token { return s.s26 }

func (s *NegateContext) SetS26(v antlr.Token) { s.s26 = v }

func (s *NegateContext) GetOps() []antlr.Token { return s.ops }

func (s *NegateContext) SetOps(v []antlr.Token) { s.ops = v }

func (s *NegateContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NegateContext) Member() IMemberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberContext)
}

func (s *NegateContext) AllMINUS() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserMINUS)
}

func (s *NegateContext) MINUS(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserMINUS, i)
}

func (s *NegateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterNegate(s)
	}
}

func (s *NegateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitNegate(s)
	}
}

func (s *NegateContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitNegate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Unary() (localctx IUnaryContext) {
	this := p
	_ = this

	localctx = NewUnaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, CommandsParserRULE_unary)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.SetState(235)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		localctx = NewMemberExprContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(222)
			p.member(0)
		}

	case 2:
		localctx = NewLogicalNotContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		p.SetState(224)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == CommandsParserEXCLAM {
			{
				p.SetState(223)

				var _m = p.Match(CommandsParserEXCLAM)

				localctx.(*LogicalNotContext).s27 = _m
			}
			localctx.(*LogicalNotContext).ops = append(localctx.(*LogicalNotContext).ops, localctx.(*LogicalNotContext).s27)

			p.SetState(226)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(228)
			p.member(0)
		}

	case 3:
		localctx = NewNegateContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		p.SetState(230)
		p.GetErrorHandler().Sync(p)
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(229)

					var _m = p.Match(CommandsParserMINUS)

					localctx.(*NegateContext).s26 = _m
				}
				localctx.(*NegateContext).ops = append(localctx.(*NegateContext).ops, localctx.(*NegateContext).s26)

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(232)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext())
		}
		{
			p.SetState(234)
			p.member(0)
		}

	}

	return localctx
}

// IMemberContext is an interface to support dynamic dispatch.
type IMemberContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMemberContext differentiates from other interfaces.
	IsMemberContext()
}

type MemberContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberContext() *MemberContext {
	var p = new(MemberContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_member
	return p
}

func (*MemberContext) IsMemberContext() {}

func NewMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberContext {
	var p = new(MemberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_member

	return p
}

func (s *MemberContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberContext) CopyFrom(ctx *MemberContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *MemberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SelectOrCallContext struct {
	*MemberContext
	op   antlr.Token
	id   antlr.Token
	open antlr.Token
	args IExprListContext
}

func NewSelectOrCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectOrCallContext {
	var p = new(SelectOrCallContext)

	p.MemberContext = NewEmptyMemberContext()
	p.parser = parser
	p.CopyFrom(ctx.(*MemberContext))

	return p
}

func (s *SelectOrCallContext) GetOp() antlr.Token { return s.op }

func (s *SelectOrCallContext) GetId() antlr.Token { return s.id }

func (s *SelectOrCallContext) GetOpen() antlr.Token { return s.open }

func (s *SelectOrCallContext) SetOp(v antlr.Token) { s.op = v }

func (s *SelectOrCallContext) SetId(v antlr.Token) { s.id = v }

func (s *SelectOrCallContext) SetOpen(v antlr.Token) { s.open = v }

func (s *SelectOrCallContext) GetArgs() IExprListContext { return s.args }

func (s *SelectOrCallContext) SetArgs(v IExprListContext) { s.args = v }

func (s *SelectOrCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectOrCallContext) Member() IMemberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberContext)
}

func (s *SelectOrCallContext) DOT() antlr.TerminalNode {
	return s.GetToken(CommandsParserDOT, 0)
}

func (s *SelectOrCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, 0)
}

func (s *SelectOrCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserRPAREN, 0)
}

func (s *SelectOrCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserLPAREN, 0)
}

func (s *SelectOrCallContext) ExprList() IExprListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprListContext)
}

func (s *SelectOrCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterSelectOrCall(s)
	}
}

func (s *SelectOrCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitSelectOrCall(s)
	}
}

func (s *SelectOrCallContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitSelectOrCall(s)

	default:
		return t.VisitChildren(s)
	}
}

type PrimaryExprContext struct {
	*MemberContext
}

func NewPrimaryExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PrimaryExprContext {
	var p = new(PrimaryExprContext)

	p.MemberContext = NewEmptyMemberContext()
	p.parser = parser
	p.CopyFrom(ctx.(*MemberContext))

	return p
}

func (s *PrimaryExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExprContext) Primary() IPrimaryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *PrimaryExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterPrimaryExpr(s)
	}
}

func (s *PrimaryExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitPrimaryExpr(s)
	}
}

func (s *PrimaryExprContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitPrimaryExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type IndexContext struct {
	*MemberContext
	op    antlr.Token
	index IExprContext
}

func NewIndexContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IndexContext {
	var p = new(IndexContext)

	p.MemberContext = NewEmptyMemberContext()
	p.parser = parser
	p.CopyFrom(ctx.(*MemberContext))

	return p
}

func (s *IndexContext) GetOp() antlr.Token { return s.op }

func (s *IndexContext) SetOp(v antlr.Token) { s.op = v }

func (s *IndexContext) GetIndex() IExprContext { return s.index }

func (s *IndexContext) SetIndex(v IExprContext) { s.index = v }

func (s *IndexContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IndexContext) Member() IMemberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberContext)
}

func (s *IndexContext) RPRACKET() antlr.TerminalNode {
	return s.GetToken(CommandsParserRPRACKET, 0)
}

func (s *IndexContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(CommandsParserLBRACKET, 0)
}

func (s *IndexContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *IndexContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterIndex(s)
	}
}

func (s *IndexContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitIndex(s)
	}
}

func (s *IndexContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitIndex(s)

	default:
		return t.VisitChildren(s)
	}
}

type CreateMessageContext struct {
	*MemberContext
	op      antlr.Token
	entries IFieldInitializerListContext
}

func NewCreateMessageContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateMessageContext {
	var p = new(CreateMessageContext)

	p.MemberContext = NewEmptyMemberContext()
	p.parser = parser
	p.CopyFrom(ctx.(*MemberContext))

	return p
}

func (s *CreateMessageContext) GetOp() antlr.Token { return s.op }

func (s *CreateMessageContext) SetOp(v antlr.Token) { s.op = v }

func (s *CreateMessageContext) GetEntries() IFieldInitializerListContext { return s.entries }

func (s *CreateMessageContext) SetEntries(v IFieldInitializerListContext) { s.entries = v }

func (s *CreateMessageContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateMessageContext) Member() IMemberContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMemberContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMemberContext)
}

func (s *CreateMessageContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(CommandsParserRBRACE, 0)
}

func (s *CreateMessageContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(CommandsParserLBRACE, 0)
}

func (s *CreateMessageContext) COMMA() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, 0)
}

func (s *CreateMessageContext) FieldInitializerList() IFieldInitializerListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFieldInitializerListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFieldInitializerListContext)
}

func (s *CreateMessageContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterCreateMessage(s)
	}
}

func (s *CreateMessageContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitCreateMessage(s)
	}
}

func (s *CreateMessageContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCreateMessage(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Member() (localctx IMemberContext) {
	return p.member(0)
}

func (p *CommandsParser) member(_p int) (localctx IMemberContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewMemberContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IMemberContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 46
	p.EnterRecursionRule(localctx, 46, CommandsParserRULE_member, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewPrimaryExprContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(238)
		p.Primary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(266)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(264)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext()) {
			case 1:
				localctx = NewSelectOrCallContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_member)
				p.SetState(240)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(241)

					var _m = p.Match(CommandsParserDOT)

					localctx.(*SelectOrCallContext).op = _m
				}
				{
					p.SetState(242)

					var _m = p.Match(CommandsParserIDENTIFIER)

					localctx.(*SelectOrCallContext).id = _m
				}
				p.SetState(248)
				p.GetErrorHandler().Sync(p)

				if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) == 1 {
					{
						p.SetState(243)

						var _m = p.Match(CommandsParserLPAREN)

						localctx.(*SelectOrCallContext).open = _m
					}
					p.SetState(245)
					p.GetErrorHandler().Sync(p)
					_la = p.GetTokenStream().LA(1)

					if ((_la-18)&-(0x1f+1)) == 0 && ((1<<uint((_la-18)))&((1<<(CommandsParserLBRACKET-18))|(1<<(CommandsParserLBRACE-18))|(1<<(CommandsParserLPAREN-18))|(1<<(CommandsParserDOT-18))|(1<<(CommandsParserMINUS-18))|(1<<(CommandsParserEXCLAM-18))|(1<<(CommandsParserCEL_TRUE-18))|(1<<(CommandsParserCEL_FALSE-18))|(1<<(CommandsParserNUL-18))|(1<<(CommandsParserNUM_FLOAT-18))|(1<<(CommandsParserNUM_INT-18))|(1<<(CommandsParserNUM_UINT-18))|(1<<(CommandsParserSTRING-18))|(1<<(CommandsParserBYTES-18))|(1<<(CommandsParserIDENTIFIER-18)))) != 0 {
						{
							p.SetState(244)

							var _x = p.ExprList()

							localctx.(*SelectOrCallContext).args = _x
						}

					}
					{
						p.SetState(247)
						p.Match(CommandsParserRPAREN)
					}

				}

			case 2:
				localctx = NewIndexContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_member)
				p.SetState(250)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(251)

					var _m = p.Match(CommandsParserLBRACKET)

					localctx.(*IndexContext).op = _m
				}
				{
					p.SetState(252)

					var _x = p.Expr()

					localctx.(*IndexContext).index = _x
				}
				{
					p.SetState(253)
					p.Match(CommandsParserRPRACKET)
				}

			case 3:
				localctx = NewCreateMessageContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_member)
				p.SetState(255)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(256)

					var _m = p.Match(CommandsParserLBRACE)

					localctx.(*CreateMessageContext).op = _m
				}
				p.SetState(258)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == CommandsParserIDENTIFIER {
					{
						p.SetState(257)

						var _x = p.FieldInitializerList()

						localctx.(*CreateMessageContext).entries = _x
					}

				}
				p.SetState(261)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == CommandsParserCOMMA {
					{
						p.SetState(260)
						p.Match(CommandsParserCOMMA)
					}

				}
				{
					p.SetState(263)
					p.Match(CommandsParserRBRACE)
				}

			}

		}
		p.SetState(268)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext())
	}

	return localctx
}

// IPrimaryContext is an interface to support dynamic dispatch.
type IPrimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPrimaryContext differentiates from other interfaces.
	IsPrimaryContext()
}

type PrimaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryContext() *PrimaryContext {
	var p = new(PrimaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_primary
	return p
}

func (*PrimaryContext) IsPrimaryContext() {}

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext {
	var p = new(PrimaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_primary

	return p
}

func (s *PrimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryContext) CopyFrom(ctx *PrimaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *PrimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type CreateListContext struct {
	*PrimaryContext
	op    antlr.Token
	elems IExprListContext
}

func NewCreateListContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateListContext {
	var p = new(CreateListContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *CreateListContext) GetOp() antlr.Token { return s.op }

func (s *CreateListContext) SetOp(v antlr.Token) { s.op = v }

func (s *CreateListContext) GetElems() IExprListContext { return s.elems }

func (s *CreateListContext) SetElems(v IExprListContext) { s.elems = v }

func (s *CreateListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateListContext) RPRACKET() antlr.TerminalNode {
	return s.GetToken(CommandsParserRPRACKET, 0)
}

func (s *CreateListContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(CommandsParserLBRACKET, 0)
}

func (s *CreateListContext) COMMA() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, 0)
}

func (s *CreateListContext) ExprList() IExprListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprListContext)
}

func (s *CreateListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterCreateList(s)
	}
}

func (s *CreateListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitCreateList(s)
	}
}

func (s *CreateListContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCreateList(s)

	default:
		return t.VisitChildren(s)
	}
}

type CreateStructContext struct {
	*PrimaryContext
	op      antlr.Token
	entries IMapInitializerListContext
}

func NewCreateStructContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateStructContext {
	var p = new(CreateStructContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *CreateStructContext) GetOp() antlr.Token { return s.op }

func (s *CreateStructContext) SetOp(v antlr.Token) { s.op = v }

func (s *CreateStructContext) GetEntries() IMapInitializerListContext { return s.entries }

func (s *CreateStructContext) SetEntries(v IMapInitializerListContext) { s.entries = v }

func (s *CreateStructContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateStructContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(CommandsParserRBRACE, 0)
}

func (s *CreateStructContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(CommandsParserLBRACE, 0)
}

func (s *CreateStructContext) COMMA() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, 0)
}

func (s *CreateStructContext) MapInitializerList() IMapInitializerListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMapInitializerListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMapInitializerListContext)
}

func (s *CreateStructContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterCreateStruct(s)
	}
}

func (s *CreateStructContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitCreateStruct(s)
	}
}

func (s *CreateStructContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCreateStruct(s)

	default:
		return t.VisitChildren(s)
	}
}

type ConstantLiteralContext struct {
	*PrimaryContext
}

func NewConstantLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ConstantLiteralContext {
	var p = new(ConstantLiteralContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *ConstantLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstantLiteralContext) Literal() ILiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *ConstantLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterConstantLiteral(s)
	}
}

func (s *ConstantLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitConstantLiteral(s)
	}
}

func (s *ConstantLiteralContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitConstantLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type NestedContext struct {
	*PrimaryContext
	e IExprContext
}

func NewNestedContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NestedContext {
	var p = new(NestedContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *NestedContext) GetE() IExprContext { return s.e }

func (s *NestedContext) SetE(v IExprContext) { s.e = v }

func (s *NestedContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NestedContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserLPAREN, 0)
}

func (s *NestedContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserRPAREN, 0)
}

func (s *NestedContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *NestedContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterNested(s)
	}
}

func (s *NestedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitNested(s)
	}
}

func (s *NestedContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitNested(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdentOrGlobalCallContext struct {
	*PrimaryContext
	leadingDot antlr.Token
	id         antlr.Token
	op         antlr.Token
	args       IExprListContext
}

func NewIdentOrGlobalCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentOrGlobalCallContext {
	var p = new(IdentOrGlobalCallContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *IdentOrGlobalCallContext) GetLeadingDot() antlr.Token { return s.leadingDot }

func (s *IdentOrGlobalCallContext) GetId() antlr.Token { return s.id }

func (s *IdentOrGlobalCallContext) GetOp() antlr.Token { return s.op }

func (s *IdentOrGlobalCallContext) SetLeadingDot(v antlr.Token) { s.leadingDot = v }

func (s *IdentOrGlobalCallContext) SetId(v antlr.Token) { s.id = v }

func (s *IdentOrGlobalCallContext) SetOp(v antlr.Token) { s.op = v }

func (s *IdentOrGlobalCallContext) GetArgs() IExprListContext { return s.args }

func (s *IdentOrGlobalCallContext) SetArgs(v IExprListContext) { s.args = v }

func (s *IdentOrGlobalCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentOrGlobalCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, 0)
}

func (s *IdentOrGlobalCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserRPAREN, 0)
}

func (s *IdentOrGlobalCallContext) DOT() antlr.TerminalNode {
	return s.GetToken(CommandsParserDOT, 0)
}

func (s *IdentOrGlobalCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserLPAREN, 0)
}

func (s *IdentOrGlobalCallContext) ExprList() IExprListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprListContext)
}

func (s *IdentOrGlobalCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterIdentOrGlobalCall(s)
	}
}

func (s *IdentOrGlobalCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitIdentOrGlobalCall(s)
	}
}

func (s *IdentOrGlobalCallContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitIdentOrGlobalCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Primary() (localctx IPrimaryContext) {
	this := p
	_ = this

	localctx = NewPrimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, CommandsParserRULE_primary)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(301)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CommandsParserDOT, CommandsParserIDENTIFIER:
		localctx = NewIdentOrGlobalCallContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(270)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserDOT {
			{
				p.SetState(269)

				var _m = p.Match(CommandsParserDOT)

				localctx.(*IdentOrGlobalCallContext).leadingDot = _m
			}

		}
		{
			p.SetState(272)

			var _m = p.Match(CommandsParserIDENTIFIER)

			localctx.(*IdentOrGlobalCallContext).id = _m
		}
		p.SetState(278)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(273)

				var _m = p.Match(CommandsParserLPAREN)

				localctx.(*IdentOrGlobalCallContext).op = _m
			}
			p.SetState(275)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if ((_la-18)&-(0x1f+1)) == 0 && ((1<<uint((_la-18)))&((1<<(CommandsParserLBRACKET-18))|(1<<(CommandsParserLBRACE-18))|(1<<(CommandsParserLPAREN-18))|(1<<(CommandsParserDOT-18))|(1<<(CommandsParserMINUS-18))|(1<<(CommandsParserEXCLAM-18))|(1<<(CommandsParserCEL_TRUE-18))|(1<<(CommandsParserCEL_FALSE-18))|(1<<(CommandsParserNUL-18))|(1<<(CommandsParserNUM_FLOAT-18))|(1<<(CommandsParserNUM_INT-18))|(1<<(CommandsParserNUM_UINT-18))|(1<<(CommandsParserSTRING-18))|(1<<(CommandsParserBYTES-18))|(1<<(CommandsParserIDENTIFIER-18)))) != 0 {
				{
					p.SetState(274)

					var _x = p.ExprList()

					localctx.(*IdentOrGlobalCallContext).args = _x
				}

			}
			{
				p.SetState(277)
				p.Match(CommandsParserRPAREN)
			}

		}

	case CommandsParserLPAREN:
		localctx = NewNestedContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(280)
			p.Match(CommandsParserLPAREN)
		}
		{
			p.SetState(281)

			var _x = p.Expr()

			localctx.(*NestedContext).e = _x
		}
		{
			p.SetState(282)
			p.Match(CommandsParserRPAREN)
		}

	case CommandsParserLBRACKET:
		localctx = NewCreateListContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(284)

			var _m = p.Match(CommandsParserLBRACKET)

			localctx.(*CreateListContext).op = _m
		}
		p.SetState(286)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if ((_la-18)&-(0x1f+1)) == 0 && ((1<<uint((_la-18)))&((1<<(CommandsParserLBRACKET-18))|(1<<(CommandsParserLBRACE-18))|(1<<(CommandsParserLPAREN-18))|(1<<(CommandsParserDOT-18))|(1<<(CommandsParserMINUS-18))|(1<<(CommandsParserEXCLAM-18))|(1<<(CommandsParserCEL_TRUE-18))|(1<<(CommandsParserCEL_FALSE-18))|(1<<(CommandsParserNUL-18))|(1<<(CommandsParserNUM_FLOAT-18))|(1<<(CommandsParserNUM_INT-18))|(1<<(CommandsParserNUM_UINT-18))|(1<<(CommandsParserSTRING-18))|(1<<(CommandsParserBYTES-18))|(1<<(CommandsParserIDENTIFIER-18)))) != 0 {
			{
				p.SetState(285)

				var _x = p.ExprList()

				localctx.(*CreateListContext).elems = _x
			}

		}
		p.SetState(289)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserCOMMA {
			{
				p.SetState(288)
				p.Match(CommandsParserCOMMA)
			}

		}
		{
			p.SetState(291)
			p.Match(CommandsParserRPRACKET)
		}

	case CommandsParserLBRACE:
		localctx = NewCreateStructContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(292)

			var _m = p.Match(CommandsParserLBRACE)

			localctx.(*CreateStructContext).op = _m
		}
		p.SetState(294)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if ((_la-18)&-(0x1f+1)) == 0 && ((1<<uint((_la-18)))&((1<<(CommandsParserLBRACKET-18))|(1<<(CommandsParserLBRACE-18))|(1<<(CommandsParserLPAREN-18))|(1<<(CommandsParserDOT-18))|(1<<(CommandsParserMINUS-18))|(1<<(CommandsParserEXCLAM-18))|(1<<(CommandsParserCEL_TRUE-18))|(1<<(CommandsParserCEL_FALSE-18))|(1<<(CommandsParserNUL-18))|(1<<(CommandsParserNUM_FLOAT-18))|(1<<(CommandsParserNUM_INT-18))|(1<<(CommandsParserNUM_UINT-18))|(1<<(CommandsParserSTRING-18))|(1<<(CommandsParserBYTES-18))|(1<<(CommandsParserIDENTIFIER-18)))) != 0 {
			{
				p.SetState(293)

				var _x = p.MapInitializerList()

				localctx.(*CreateStructContext).entries = _x
			}

		}
		p.SetState(297)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserCOMMA {
			{
				p.SetState(296)
				p.Match(CommandsParserCOMMA)
			}

		}
		{
			p.SetState(299)
			p.Match(CommandsParserRBRACE)
		}

	case CommandsParserMINUS, CommandsParserCEL_TRUE, CommandsParserCEL_FALSE, CommandsParserNUL, CommandsParserNUM_FLOAT, CommandsParserNUM_INT, CommandsParserNUM_UINT, CommandsParserSTRING, CommandsParserBYTES:
		localctx = NewConstantLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(300)
			p.Literal()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExprListContext is an interface to support dynamic dispatch.
type IExprListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_expr returns the _expr rule contexts.
	Get_expr() IExprContext

	// Set_expr sets the _expr rule contexts.
	Set_expr(IExprContext)

	// GetE returns the e rule context list.
	GetE() []IExprContext

	// SetE sets the e rule context list.
	SetE([]IExprContext)

	// IsExprListContext differentiates from other interfaces.
	IsExprListContext()
}

type ExprListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	_expr  IExprContext
	e      []IExprContext
}

func NewEmptyExprListContext() *ExprListContext {
	var p = new(ExprListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_exprList
	return p
}

func (*ExprListContext) IsExprListContext() {}

func NewExprListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprListContext {
	var p = new(ExprListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_exprList

	return p
}

func (s *ExprListContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprListContext) Get_expr() IExprContext { return s._expr }

func (s *ExprListContext) Set_expr(v IExprContext) { s._expr = v }

func (s *ExprListContext) GetE() []IExprContext { return s.e }

func (s *ExprListContext) SetE(v []IExprContext) { s.e = v }

func (s *ExprListContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *ExprListContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *ExprListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOMMA)
}

func (s *ExprListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, i)
}

func (s *ExprListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterExprList(s)
	}
}

func (s *ExprListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitExprList(s)
	}
}

func (s *ExprListContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitExprList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ExprList() (localctx IExprListContext) {
	this := p
	_ = this

	localctx = NewExprListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, CommandsParserRULE_exprList)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(303)

		var _x = p.Expr()

		localctx.(*ExprListContext)._expr = _x
	}
	localctx.(*ExprListContext).e = append(localctx.(*ExprListContext).e, localctx.(*ExprListContext)._expr)
	p.SetState(308)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(304)
				p.Match(CommandsParserCOMMA)
			}
			{
				p.SetState(305)

				var _x = p.Expr()

				localctx.(*ExprListContext)._expr = _x
			}
			localctx.(*ExprListContext).e = append(localctx.(*ExprListContext).e, localctx.(*ExprListContext)._expr)

		}
		p.SetState(310)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 39, p.GetParserRuleContext())
	}

	return localctx
}

// IFieldInitializerListContext is an interface to support dynamic dispatch.
type IFieldInitializerListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_IDENTIFIER returns the _IDENTIFIER token.
	Get_IDENTIFIER() antlr.Token

	// GetS29 returns the s29 token.
	GetS29() antlr.Token

	// Set_IDENTIFIER sets the _IDENTIFIER token.
	Set_IDENTIFIER(antlr.Token)

	// SetS29 sets the s29 token.
	SetS29(antlr.Token)

	// GetFields returns the fields token list.
	GetFields() []antlr.Token

	// GetCols returns the cols token list.
	GetCols() []antlr.Token

	// SetFields sets the fields token list.
	SetFields([]antlr.Token)

	// SetCols sets the cols token list.
	SetCols([]antlr.Token)

	// Get_expr returns the _expr rule contexts.
	Get_expr() IExprContext

	// Set_expr sets the _expr rule contexts.
	Set_expr(IExprContext)

	// GetValues returns the values rule context list.
	GetValues() []IExprContext

	// SetValues sets the values rule context list.
	SetValues([]IExprContext)

	// IsFieldInitializerListContext differentiates from other interfaces.
	IsFieldInitializerListContext()
}

type FieldInitializerListContext struct {
	*antlr.BaseParserRuleContext
	parser      antlr.Parser
	_IDENTIFIER antlr.Token
	fields      []antlr.Token
	s29         antlr.Token
	cols        []antlr.Token
	_expr       IExprContext
	values      []IExprContext
}

func NewEmptyFieldInitializerListContext() *FieldInitializerListContext {
	var p = new(FieldInitializerListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_fieldInitializerList
	return p
}

func (*FieldInitializerListContext) IsFieldInitializerListContext() {}

func NewFieldInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldInitializerListContext {
	var p = new(FieldInitializerListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_fieldInitializerList

	return p
}

func (s *FieldInitializerListContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldInitializerListContext) Get_IDENTIFIER() antlr.Token { return s._IDENTIFIER }

func (s *FieldInitializerListContext) GetS29() antlr.Token { return s.s29 }

func (s *FieldInitializerListContext) Set_IDENTIFIER(v antlr.Token) { s._IDENTIFIER = v }

func (s *FieldInitializerListContext) SetS29(v antlr.Token) { s.s29 = v }

func (s *FieldInitializerListContext) GetFields() []antlr.Token { return s.fields }

func (s *FieldInitializerListContext) GetCols() []antlr.Token { return s.cols }

func (s *FieldInitializerListContext) SetFields(v []antlr.Token) { s.fields = v }

func (s *FieldInitializerListContext) SetCols(v []antlr.Token) { s.cols = v }

func (s *FieldInitializerListContext) Get_expr() IExprContext { return s._expr }

func (s *FieldInitializerListContext) Set_expr(v IExprContext) { s._expr = v }

func (s *FieldInitializerListContext) GetValues() []IExprContext { return s.values }

func (s *FieldInitializerListContext) SetValues(v []IExprContext) { s.values = v }

func (s *FieldInitializerListContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserIDENTIFIER)
}

func (s *FieldInitializerListContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, i)
}

func (s *FieldInitializerListContext) AllCOLON() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOLON)
}

func (s *FieldInitializerListContext) COLON(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOLON, i)
}

func (s *FieldInitializerListContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *FieldInitializerListContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *FieldInitializerListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOMMA)
}

func (s *FieldInitializerListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, i)
}

func (s *FieldInitializerListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldInitializerListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldInitializerListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterFieldInitializerList(s)
	}
}

func (s *FieldInitializerListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitFieldInitializerList(s)
	}
}

func (s *FieldInitializerListContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitFieldInitializerList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) FieldInitializerList() (localctx IFieldInitializerListContext) {
	this := p
	_ = this

	localctx = NewFieldInitializerListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, CommandsParserRULE_fieldInitializerList)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(311)

		var _m = p.Match(CommandsParserIDENTIFIER)

		localctx.(*FieldInitializerListContext)._IDENTIFIER = _m
	}
	localctx.(*FieldInitializerListContext).fields = append(localctx.(*FieldInitializerListContext).fields, localctx.(*FieldInitializerListContext)._IDENTIFIER)
	{
		p.SetState(312)

		var _m = p.Match(CommandsParserCOLON)

		localctx.(*FieldInitializerListContext).s29 = _m
	}
	localctx.(*FieldInitializerListContext).cols = append(localctx.(*FieldInitializerListContext).cols, localctx.(*FieldInitializerListContext).s29)
	{
		p.SetState(313)

		var _x = p.Expr()

		localctx.(*FieldInitializerListContext)._expr = _x
	}
	localctx.(*FieldInitializerListContext).values = append(localctx.(*FieldInitializerListContext).values, localctx.(*FieldInitializerListContext)._expr)
	p.SetState(320)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(314)
				p.Match(CommandsParserCOMMA)
			}
			{
				p.SetState(315)

				var _m = p.Match(CommandsParserIDENTIFIER)

				localctx.(*FieldInitializerListContext)._IDENTIFIER = _m
			}
			localctx.(*FieldInitializerListContext).fields = append(localctx.(*FieldInitializerListContext).fields, localctx.(*FieldInitializerListContext)._IDENTIFIER)
			{
				p.SetState(316)

				var _m = p.Match(CommandsParserCOLON)

				localctx.(*FieldInitializerListContext).s29 = _m
			}
			localctx.(*FieldInitializerListContext).cols = append(localctx.(*FieldInitializerListContext).cols, localctx.(*FieldInitializerListContext).s29)
			{
				p.SetState(317)

				var _x = p.Expr()

				localctx.(*FieldInitializerListContext)._expr = _x
			}
			localctx.(*FieldInitializerListContext).values = append(localctx.(*FieldInitializerListContext).values, localctx.(*FieldInitializerListContext)._expr)

		}
		p.SetState(322)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext())
	}

	return localctx
}

// IMapInitializerListContext is an interface to support dynamic dispatch.
type IMapInitializerListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS29 returns the s29 token.
	GetS29() antlr.Token

	// SetS29 sets the s29 token.
	SetS29(antlr.Token)

	// GetCols returns the cols token list.
	GetCols() []antlr.Token

	// SetCols sets the cols token list.
	SetCols([]antlr.Token)

	// Get_expr returns the _expr rule contexts.
	Get_expr() IExprContext

	// Set_expr sets the _expr rule contexts.
	Set_expr(IExprContext)

	// GetKeys returns the keys rule context list.
	GetKeys() []IExprContext

	// GetValues returns the values rule context list.
	GetValues() []IExprContext

	// SetKeys sets the keys rule context list.
	SetKeys([]IExprContext)

	// SetValues sets the values rule context list.
	SetValues([]IExprContext)

	// IsMapInitializerListContext differentiates from other interfaces.
	IsMapInitializerListContext()
}

type MapInitializerListContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	_expr  IExprContext
	keys   []IExprContext
	s29    antlr.Token
	cols   []antlr.Token
	values []IExprContext
}

func NewEmptyMapInitializerListContext() *MapInitializerListContext {
	var p = new(MapInitializerListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_mapInitializerList
	return p
}

func (*MapInitializerListContext) IsMapInitializerListContext() {}

func NewMapInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapInitializerListContext {
	var p = new(MapInitializerListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_mapInitializerList

	return p
}

func (s *MapInitializerListContext) GetParser() antlr.Parser { return s.parser }

func (s *MapInitializerListContext) GetS29() antlr.Token { return s.s29 }

func (s *MapInitializerListContext) SetS29(v antlr.Token) { s.s29 = v }

func (s *MapInitializerListContext) GetCols() []antlr.Token { return s.cols }

func (s *MapInitializerListContext) SetCols(v []antlr.Token) { s.cols = v }

func (s *MapInitializerListContext) Get_expr() IExprContext { return s._expr }

func (s *MapInitializerListContext) Set_expr(v IExprContext) { s._expr = v }

func (s *MapInitializerListContext) GetKeys() []IExprContext { return s.keys }

func (s *MapInitializerListContext) GetValues() []IExprContext { return s.values }

func (s *MapInitializerListContext) SetKeys(v []IExprContext) { s.keys = v }

func (s *MapInitializerListContext) SetValues(v []IExprContext) { s.values = v }

func (s *MapInitializerListContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *MapInitializerListContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *MapInitializerListContext) AllCOLON() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOLON)
}

func (s *MapInitializerListContext) COLON(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOLON, i)
}

func (s *MapInitializerListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOMMA)
}

func (s *MapInitializerListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, i)
}

func (s *MapInitializerListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MapInitializerListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MapInitializerListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterMapInitializerList(s)
	}
}

func (s *MapInitializerListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitMapInitializerList(s)
	}
}

func (s *MapInitializerListContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitMapInitializerList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) MapInitializerList() (localctx IMapInitializerListContext) {
	this := p
	_ = this

	localctx = NewMapInitializerListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, CommandsParserRULE_mapInitializerList)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(323)

		var _x = p.Expr()

		localctx.(*MapInitializerListContext)._expr = _x
	}
	localctx.(*MapInitializerListContext).keys = append(localctx.(*MapInitializerListContext).keys, localctx.(*MapInitializerListContext)._expr)
	{
		p.SetState(324)

		var _m = p.Match(CommandsParserCOLON)

		localctx.(*MapInitializerListContext).s29 = _m
	}
	localctx.(*MapInitializerListContext).cols = append(localctx.(*MapInitializerListContext).cols, localctx.(*MapInitializerListContext).s29)
	{
		p.SetState(325)

		var _x = p.Expr()

		localctx.(*MapInitializerListContext)._expr = _x
	}
	localctx.(*MapInitializerListContext).values = append(localctx.(*MapInitializerListContext).values, localctx.(*MapInitializerListContext)._expr)
	p.SetState(333)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(326)
				p.Match(CommandsParserCOMMA)
			}
			{
				p.SetState(327)

				var _x = p.Expr()

				localctx.(*MapInitializerListContext)._expr = _x
			}
			localctx.(*MapInitializerListContext).keys = append(localctx.(*MapInitializerListContext).keys, localctx.(*MapInitializerListContext)._expr)
			{
				p.SetState(328)

				var _m = p.Match(CommandsParserCOLON)

				localctx.(*MapInitializerListContext).s29 = _m
			}
			localctx.(*MapInitializerListContext).cols = append(localctx.(*MapInitializerListContext).cols, localctx.(*MapInitializerListContext).s29)
			{
				p.SetState(329)

				var _x = p.Expr()

				localctx.(*MapInitializerListContext)._expr = _x
			}
			localctx.(*MapInitializerListContext).values = append(localctx.(*MapInitializerListContext).values, localctx.(*MapInitializerListContext)._expr)

		}
		p.SetState(335)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 41, p.GetParserRuleContext())
	}

	return localctx
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CommandsParserRULE_literal
	return p
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) CopyFrom(ctx *LiteralContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type BytesContext struct {
	*LiteralContext
	tok antlr.Token
}

func NewBytesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BytesContext {
	var p = new(BytesContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *BytesContext) GetTok() antlr.Token { return s.tok }

func (s *BytesContext) SetTok(v antlr.Token) { s.tok = v }

func (s *BytesContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BytesContext) BYTES() antlr.TerminalNode {
	return s.GetToken(CommandsParserBYTES, 0)
}

func (s *BytesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterBytes(s)
	}
}

func (s *BytesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitBytes(s)
	}
}

func (s *BytesContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitBytes(s)

	default:
		return t.VisitChildren(s)
	}
}

type UintContext struct {
	*LiteralContext
	tok antlr.Token
}

func NewUintContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UintContext {
	var p = new(UintContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *UintContext) GetTok() antlr.Token { return s.tok }

func (s *UintContext) SetTok(v antlr.Token) { s.tok = v }

func (s *UintContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UintContext) NUM_UINT() antlr.TerminalNode {
	return s.GetToken(CommandsParserNUM_UINT, 0)
}

func (s *UintContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterUint(s)
	}
}

func (s *UintContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitUint(s)
	}
}

func (s *UintContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitUint(s)

	default:
		return t.VisitChildren(s)
	}
}

type NullContext struct {
	*LiteralContext
	tok antlr.Token
}

func NewNullContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullContext {
	var p = new(NullContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *NullContext) GetTok() antlr.Token { return s.tok }

func (s *NullContext) SetTok(v antlr.Token) { s.tok = v }

func (s *NullContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NullContext) NUL() antlr.TerminalNode {
	return s.GetToken(CommandsParserNUL, 0)
}

func (s *NullContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterNull(s)
	}
}

func (s *NullContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitNull(s)
	}
}

func (s *NullContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitNull(s)

	default:
		return t.VisitChildren(s)
	}
}

type BoolFalseContext struct {
	*LiteralContext
	tok antlr.Token
}

func NewBoolFalseContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoolFalseContext {
	var p = new(BoolFalseContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *BoolFalseContext) GetTok() antlr.Token { return s.tok }

func (s *BoolFalseContext) SetTok(v antlr.Token) { s.tok = v }

func (s *BoolFalseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolFalseContext) CEL_FALSE() antlr.TerminalNode {
	return s.GetToken(CommandsParserCEL_FALSE, 0)
}

func (s *BoolFalseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterBoolFalse(s)
	}
}

func (s *BoolFalseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitBoolFalse(s)
	}
}

func (s *BoolFalseContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitBoolFalse(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringContext struct {
	*LiteralContext
	tok antlr.Token
}

func NewStringContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringContext {
	var p = new(StringContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *StringContext) GetTok() antlr.Token { return s.tok }

func (s *StringContext) SetTok(v antlr.Token) { s.tok = v }

func (s *StringContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringContext) STRING() antlr.TerminalNode {
	return s.GetToken(CommandsParserSTRING, 0)
}

func (s *StringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterString(s)
	}
}

func (s *StringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitString(s)
	}
}

func (s *StringContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitString(s)

	default:
		return t.VisitChildren(s)
	}
}

type DoubleContext struct {
	*LiteralContext
	sign antlr.Token
	tok  antlr.Token
}

func NewDoubleContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DoubleContext {
	var p = new(DoubleContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *DoubleContext) GetSign() antlr.Token { return s.sign }

func (s *DoubleContext) GetTok() antlr.Token { return s.tok }

func (s *DoubleContext) SetSign(v antlr.Token) { s.sign = v }

func (s *DoubleContext) SetTok(v antlr.Token) { s.tok = v }

func (s *DoubleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DoubleContext) NUM_FLOAT() antlr.TerminalNode {
	return s.GetToken(CommandsParserNUM_FLOAT, 0)
}

func (s *DoubleContext) MINUS() antlr.TerminalNode {
	return s.GetToken(CommandsParserMINUS, 0)
}

func (s *DoubleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterDouble(s)
	}
}

func (s *DoubleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitDouble(s)
	}
}

func (s *DoubleContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitDouble(s)

	default:
		return t.VisitChildren(s)
	}
}

type BoolTrueContext struct {
	*LiteralContext
	tok antlr.Token
}

func NewBoolTrueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoolTrueContext {
	var p = new(BoolTrueContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *BoolTrueContext) GetTok() antlr.Token { return s.tok }

func (s *BoolTrueContext) SetTok(v antlr.Token) { s.tok = v }

func (s *BoolTrueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolTrueContext) CEL_TRUE() antlr.TerminalNode {
	return s.GetToken(CommandsParserCEL_TRUE, 0)
}

func (s *BoolTrueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterBoolTrue(s)
	}
}

func (s *BoolTrueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitBoolTrue(s)
	}
}

func (s *BoolTrueContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitBoolTrue(s)

	default:
		return t.VisitChildren(s)
	}
}

type IntContext struct {
	*LiteralContext
	sign antlr.Token
	tok  antlr.Token
}

func NewIntContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntContext {
	var p = new(IntContext)

	p.LiteralContext = NewEmptyLiteralContext()
	p.parser = parser
	p.CopyFrom(ctx.(*LiteralContext))

	return p
}

func (s *IntContext) GetSign() antlr.Token { return s.sign }

func (s *IntContext) GetTok() antlr.Token { return s.tok }

func (s *IntContext) SetSign(v antlr.Token) { s.sign = v }

func (s *IntContext) SetTok(v antlr.Token) { s.tok = v }

func (s *IntContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntContext) NUM_INT() antlr.TerminalNode {
	return s.GetToken(CommandsParserNUM_INT, 0)
}

func (s *IntContext) MINUS() antlr.TerminalNode {
	return s.GetToken(CommandsParserMINUS, 0)
}

func (s *IntContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterInt(s)
	}
}

func (s *IntContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitInt(s)
	}
}

func (s *IntContext) Accept(visitor antlr.ParseTreeVisitor) any {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitInt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Literal() (localctx ILiteralContext) {
	this := p
	_ = this

	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, CommandsParserRULE_literal)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(350)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 44, p.GetParserRuleContext()) {
	case 1:
		localctx = NewIntContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(337)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserMINUS {
			{
				p.SetState(336)

				var _m = p.Match(CommandsParserMINUS)

				localctx.(*IntContext).sign = _m
			}

		}
		{
			p.SetState(339)

			var _m = p.Match(CommandsParserNUM_INT)

			localctx.(*IntContext).tok = _m
		}

	case 2:
		localctx = NewUintContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(340)

			var _m = p.Match(CommandsParserNUM_UINT)

			localctx.(*UintContext).tok = _m
		}

	case 3:
		localctx = NewDoubleContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		p.SetState(342)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserMINUS {
			{
				p.SetState(341)

				var _m = p.Match(CommandsParserMINUS)

				localctx.(*DoubleContext).sign = _m
			}

		}
		{
			p.SetState(344)

			var _m = p.Match(CommandsParserNUM_FLOAT)

			localctx.(*DoubleContext).tok = _m
		}

	case 4:
		localctx = NewStringContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(345)

			var _m = p.Match(CommandsParserSTRING)

			localctx.(*StringContext).tok = _m
		}

	case 5:
		localctx = NewBytesContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(346)

			var _m = p.Match(CommandsParserBYTES)

			localctx.(*BytesContext).tok = _m
		}

	case 6:
		localctx = NewBoolTrueContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(347)

			var _m = p.Match(CommandsParserCEL_TRUE)

			localctx.(*BoolTrueContext).tok = _m
		}

	case 7:
		localctx = NewBoolFalseContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(348)

			var _m = p.Match(CommandsParserCEL_FALSE)

			localctx.(*BoolFalseContext).tok = _m
		}

	case 8:
		localctx = NewNullContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(349)

			var _m = p.Match(CommandsParserNUL)

			localctx.(*NullContext).tok = _m
		}

	}

	return localctx
}

func (p *CommandsParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 20:
		var t *RelationContext = nil
		if localctx != nil {
			t = localctx.(*RelationContext)
		}
		return p.Relation_Sempred(t, predIndex)

	case 21:
		var t *CalcContext = nil
		if localctx != nil {
			t = localctx.(*CalcContext)
		}
		return p.Calc_Sempred(t, predIndex)

	case 23:
		var t *MemberContext = nil
		if localctx != nil {
			t = localctx.(*MemberContext)
		}
		return p.Member_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *CommandsParser) Relation_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *CommandsParser) Calc_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 1:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *CommandsParser) Member_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 3:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
