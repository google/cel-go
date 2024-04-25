// Code generated from ./Commands.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Commands
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type CommandsParser struct {
	*antlr.BaseParser
}

var CommandsParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func commandsParserInit() {
	staticData := &CommandsParserStaticData
	staticData.LiteralNames = []string{
		"", "'%help'", "'%?'", "'%let'", "'%declare'", "'%delete'", "'%compile'",
		"'%parse'", "'%eval'", "", "", "'->'", "'='", "'=='", "'!='", "'in'",
		"'<'", "'<='", "'>='", "'>'", "'&&'", "'||'", "'['", "']'", "'{'", "'}'",
		"'('", "')'", "'.'", "','", "'-'", "'!'", "'?'", "':'", "'+'", "'*'",
		"'/'", "'%'", "'true'", "'false'", "'null'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "COMMAND", "FLAG", "ARROW", "EQUAL_ASSIGN",
		"EQUALS", "NOT_EQUALS", "IN", "LESS", "LESS_EQUALS", "GREATER_EQUALS",
		"GREATER", "LOGICAL_AND", "LOGICAL_OR", "LBRACKET", "RPRACKET", "LBRACE",
		"RBRACE", "LPAREN", "RPAREN", "DOT", "COMMA", "MINUS", "EXCLAM", "QUESTIONMARK",
		"COLON", "PLUS", "STAR", "SLASH", "PERCENT", "CEL_TRUE", "CEL_FALSE",
		"NUL", "WHITESPACE", "COMMENT", "NUM_FLOAT", "NUM_INT", "NUM_UINT",
		"STRING", "BYTES", "IDENTIFIER",
	}
	staticData.RuleNames = []string{
		"startCommand", "command", "help", "let", "declare", "varDecl", "fnDecl",
		"param", "delete", "simple", "empty", "compile", "parse", "exprCmd",
		"qualId", "startType", "type", "typeId", "typeParamList", "start", "expr",
		"conditionalOr", "conditionalAnd", "relation", "calc", "unary", "member",
		"primary", "exprList", "listInit", "fieldInitializerList", "optField",
		"mapInitializerList", "optExpr", "literal",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 48, 412, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 1, 0, 1, 0, 1, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 83, 8, 1, 1, 2, 1, 2,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 94, 8, 3, 1, 3, 1, 3, 1,
		4, 1, 4, 1, 4, 3, 4, 101, 8, 4, 1, 5, 1, 5, 1, 5, 3, 5, 106, 8, 5, 1, 6,
		1, 6, 1, 6, 1, 6, 1, 6, 5, 6, 113, 8, 6, 10, 6, 12, 6, 116, 9, 6, 3, 6,
		118, 8, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8,
		1, 8, 3, 8, 131, 8, 8, 1, 9, 1, 9, 1, 9, 5, 9, 136, 8, 9, 10, 9, 12, 9,
		139, 9, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 13,
		3, 13, 150, 8, 13, 1, 13, 1, 13, 1, 14, 3, 14, 155, 8, 14, 1, 14, 1, 14,
		1, 14, 5, 14, 160, 8, 14, 10, 14, 12, 14, 163, 9, 14, 1, 15, 1, 15, 1,
		15, 1, 16, 1, 16, 3, 16, 170, 8, 16, 1, 17, 3, 17, 173, 8, 17, 1, 17, 1,
		17, 1, 17, 5, 17, 178, 8, 17, 10, 17, 12, 17, 181, 9, 17, 1, 18, 1, 18,
		1, 18, 1, 18, 5, 18, 187, 8, 18, 10, 18, 12, 18, 190, 9, 18, 1, 18, 1,
		18, 1, 19, 1, 19, 1, 19, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 1, 20, 3, 20,
		203, 8, 20, 1, 21, 1, 21, 1, 21, 5, 21, 208, 8, 21, 10, 21, 12, 21, 211,
		9, 21, 1, 22, 1, 22, 1, 22, 5, 22, 216, 8, 22, 10, 22, 12, 22, 219, 9,
		22, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 1, 23, 5, 23, 227, 8, 23, 10, 23,
		12, 23, 230, 9, 23, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1, 24, 1,
		24, 1, 24, 5, 24, 241, 8, 24, 10, 24, 12, 24, 244, 9, 24, 1, 25, 1, 25,
		4, 25, 248, 8, 25, 11, 25, 12, 25, 249, 1, 25, 1, 25, 4, 25, 254, 8, 25,
		11, 25, 12, 25, 255, 1, 25, 3, 25, 259, 8, 25, 1, 26, 1, 26, 1, 26, 1,
		26, 1, 26, 1, 26, 3, 26, 267, 8, 26, 1, 26, 1, 26, 1, 26, 1, 26, 1, 26,
		1, 26, 3, 26, 275, 8, 26, 1, 26, 1, 26, 1, 26, 1, 26, 3, 26, 281, 8, 26,
		1, 26, 1, 26, 1, 26, 5, 26, 286, 8, 26, 10, 26, 12, 26, 289, 9, 26, 1,
		27, 3, 27, 292, 8, 27, 1, 27, 1, 27, 1, 27, 3, 27, 297, 8, 27, 1, 27, 3,
		27, 300, 8, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 1, 27, 3, 27, 308, 8,
		27, 1, 27, 3, 27, 311, 8, 27, 1, 27, 1, 27, 1, 27, 3, 27, 316, 8, 27, 1,
		27, 3, 27, 319, 8, 27, 1, 27, 1, 27, 3, 27, 323, 8, 27, 1, 27, 1, 27, 1,
		27, 5, 27, 328, 8, 27, 10, 27, 12, 27, 331, 9, 27, 1, 27, 1, 27, 3, 27,
		335, 8, 27, 1, 27, 3, 27, 338, 8, 27, 1, 27, 1, 27, 3, 27, 342, 8, 27,
		1, 28, 1, 28, 1, 28, 5, 28, 347, 8, 28, 10, 28, 12, 28, 350, 9, 28, 1,
		29, 1, 29, 1, 29, 5, 29, 355, 8, 29, 10, 29, 12, 29, 358, 9, 29, 1, 30,
		1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 5, 30, 368, 8, 30, 10,
		30, 12, 30, 371, 9, 30, 1, 31, 3, 31, 374, 8, 31, 1, 31, 1, 31, 1, 32,
		1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 1, 32, 5, 32, 386, 8, 32, 10,
		32, 12, 32, 389, 9, 32, 1, 33, 3, 33, 392, 8, 33, 1, 33, 1, 33, 1, 34,
		3, 34, 397, 8, 34, 1, 34, 1, 34, 1, 34, 3, 34, 402, 8, 34, 1, 34, 1, 34,
		1, 34, 1, 34, 1, 34, 1, 34, 3, 34, 410, 8, 34, 1, 34, 0, 3, 46, 48, 52,
		35, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34,
		36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 0,
		5, 1, 0, 1, 2, 2, 0, 40, 40, 48, 48, 1, 0, 13, 19, 1, 0, 35, 37, 2, 0,
		30, 30, 34, 34, 446, 0, 70, 1, 0, 0, 0, 2, 82, 1, 0, 0, 0, 4, 84, 1, 0,
		0, 0, 6, 86, 1, 0, 0, 0, 8, 97, 1, 0, 0, 0, 10, 102, 1, 0, 0, 0, 12, 107,
		1, 0, 0, 0, 14, 123, 1, 0, 0, 0, 16, 127, 1, 0, 0, 0, 18, 132, 1, 0, 0,
		0, 20, 140, 1, 0, 0, 0, 22, 142, 1, 0, 0, 0, 24, 145, 1, 0, 0, 0, 26, 149,
		1, 0, 0, 0, 28, 154, 1, 0, 0, 0, 30, 164, 1, 0, 0, 0, 32, 167, 1, 0, 0,
		0, 34, 172, 1, 0, 0, 0, 36, 182, 1, 0, 0, 0, 38, 193, 1, 0, 0, 0, 40, 196,
		1, 0, 0, 0, 42, 204, 1, 0, 0, 0, 44, 212, 1, 0, 0, 0, 46, 220, 1, 0, 0,
		0, 48, 231, 1, 0, 0, 0, 50, 258, 1, 0, 0, 0, 52, 260, 1, 0, 0, 0, 54, 341,
		1, 0, 0, 0, 56, 343, 1, 0, 0, 0, 58, 351, 1, 0, 0, 0, 60, 359, 1, 0, 0,
		0, 62, 373, 1, 0, 0, 0, 64, 377, 1, 0, 0, 0, 66, 391, 1, 0, 0, 0, 68, 409,
		1, 0, 0, 0, 70, 71, 3, 2, 1, 0, 71, 72, 5, 0, 0, 1, 72, 1, 1, 0, 0, 0,
		73, 83, 3, 4, 2, 0, 74, 83, 3, 6, 3, 0, 75, 83, 3, 8, 4, 0, 76, 83, 3,
		16, 8, 0, 77, 83, 3, 18, 9, 0, 78, 83, 3, 22, 11, 0, 79, 83, 3, 24, 12,
		0, 80, 83, 3, 26, 13, 0, 81, 83, 3, 20, 10, 0, 82, 73, 1, 0, 0, 0, 82,
		74, 1, 0, 0, 0, 82, 75, 1, 0, 0, 0, 82, 76, 1, 0, 0, 0, 82, 77, 1, 0, 0,
		0, 82, 78, 1, 0, 0, 0, 82, 79, 1, 0, 0, 0, 82, 80, 1, 0, 0, 0, 82, 81,
		1, 0, 0, 0, 83, 3, 1, 0, 0, 0, 84, 85, 7, 0, 0, 0, 85, 5, 1, 0, 0, 0, 86,
		93, 5, 3, 0, 0, 87, 88, 3, 10, 5, 0, 88, 89, 5, 12, 0, 0, 89, 94, 1, 0,
		0, 0, 90, 91, 3, 12, 6, 0, 91, 92, 5, 11, 0, 0, 92, 94, 1, 0, 0, 0, 93,
		87, 1, 0, 0, 0, 93, 90, 1, 0, 0, 0, 94, 95, 1, 0, 0, 0, 95, 96, 3, 40,
		20, 0, 96, 7, 1, 0, 0, 0, 97, 100, 5, 4, 0, 0, 98, 101, 3, 10, 5, 0, 99,
		101, 3, 12, 6, 0, 100, 98, 1, 0, 0, 0, 100, 99, 1, 0, 0, 0, 101, 9, 1,
		0, 0, 0, 102, 105, 3, 28, 14, 0, 103, 104, 5, 33, 0, 0, 104, 106, 3, 32,
		16, 0, 105, 103, 1, 0, 0, 0, 105, 106, 1, 0, 0, 0, 106, 11, 1, 0, 0, 0,
		107, 108, 3, 28, 14, 0, 108, 117, 5, 26, 0, 0, 109, 114, 3, 14, 7, 0, 110,
		111, 5, 29, 0, 0, 111, 113, 3, 14, 7, 0, 112, 110, 1, 0, 0, 0, 113, 116,
		1, 0, 0, 0, 114, 112, 1, 0, 0, 0, 114, 115, 1, 0, 0, 0, 115, 118, 1, 0,
		0, 0, 116, 114, 1, 0, 0, 0, 117, 109, 1, 0, 0, 0, 117, 118, 1, 0, 0, 0,
		118, 119, 1, 0, 0, 0, 119, 120, 5, 27, 0, 0, 120, 121, 5, 33, 0, 0, 121,
		122, 3, 32, 16, 0, 122, 13, 1, 0, 0, 0, 123, 124, 5, 48, 0, 0, 124, 125,
		5, 33, 0, 0, 125, 126, 3, 32, 16, 0, 126, 15, 1, 0, 0, 0, 127, 130, 5,
		5, 0, 0, 128, 131, 3, 10, 5, 0, 129, 131, 3, 12, 6, 0, 130, 128, 1, 0,
		0, 0, 130, 129, 1, 0, 0, 0, 131, 17, 1, 0, 0, 0, 132, 137, 5, 9, 0, 0,
		133, 136, 5, 10, 0, 0, 134, 136, 5, 46, 0, 0, 135, 133, 1, 0, 0, 0, 135,
		134, 1, 0, 0, 0, 136, 139, 1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 137, 138,
		1, 0, 0, 0, 138, 19, 1, 0, 0, 0, 139, 137, 1, 0, 0, 0, 140, 141, 1, 0,
		0, 0, 141, 21, 1, 0, 0, 0, 142, 143, 5, 6, 0, 0, 143, 144, 3, 40, 20, 0,
		144, 23, 1, 0, 0, 0, 145, 146, 5, 7, 0, 0, 146, 147, 3, 40, 20, 0, 147,
		25, 1, 0, 0, 0, 148, 150, 5, 8, 0, 0, 149, 148, 1, 0, 0, 0, 149, 150, 1,
		0, 0, 0, 150, 151, 1, 0, 0, 0, 151, 152, 3, 40, 20, 0, 152, 27, 1, 0, 0,
		0, 153, 155, 5, 28, 0, 0, 154, 153, 1, 0, 0, 0, 154, 155, 1, 0, 0, 0, 155,
		156, 1, 0, 0, 0, 156, 161, 5, 48, 0, 0, 157, 158, 5, 28, 0, 0, 158, 160,
		5, 48, 0, 0, 159, 157, 1, 0, 0, 0, 160, 163, 1, 0, 0, 0, 161, 159, 1, 0,
		0, 0, 161, 162, 1, 0, 0, 0, 162, 29, 1, 0, 0, 0, 163, 161, 1, 0, 0, 0,
		164, 165, 3, 32, 16, 0, 165, 166, 5, 0, 0, 1, 166, 31, 1, 0, 0, 0, 167,
		169, 3, 34, 17, 0, 168, 170, 3, 36, 18, 0, 169, 168, 1, 0, 0, 0, 169, 170,
		1, 0, 0, 0, 170, 33, 1, 0, 0, 0, 171, 173, 5, 28, 0, 0, 172, 171, 1, 0,
		0, 0, 172, 173, 1, 0, 0, 0, 173, 174, 1, 0, 0, 0, 174, 179, 7, 1, 0, 0,
		175, 176, 5, 28, 0, 0, 176, 178, 5, 48, 0, 0, 177, 175, 1, 0, 0, 0, 178,
		181, 1, 0, 0, 0, 179, 177, 1, 0, 0, 0, 179, 180, 1, 0, 0, 0, 180, 35, 1,
		0, 0, 0, 181, 179, 1, 0, 0, 0, 182, 183, 5, 26, 0, 0, 183, 188, 3, 32,
		16, 0, 184, 185, 5, 29, 0, 0, 185, 187, 3, 32, 16, 0, 186, 184, 1, 0, 0,
		0, 187, 190, 1, 0, 0, 0, 188, 186, 1, 0, 0, 0, 188, 189, 1, 0, 0, 0, 189,
		191, 1, 0, 0, 0, 190, 188, 1, 0, 0, 0, 191, 192, 5, 27, 0, 0, 192, 37,
		1, 0, 0, 0, 193, 194, 3, 40, 20, 0, 194, 195, 5, 0, 0, 1, 195, 39, 1, 0,
		0, 0, 196, 202, 3, 42, 21, 0, 197, 198, 5, 32, 0, 0, 198, 199, 3, 42, 21,
		0, 199, 200, 5, 33, 0, 0, 200, 201, 3, 40, 20, 0, 201, 203, 1, 0, 0, 0,
		202, 197, 1, 0, 0, 0, 202, 203, 1, 0, 0, 0, 203, 41, 1, 0, 0, 0, 204, 209,
		3, 44, 22, 0, 205, 206, 5, 21, 0, 0, 206, 208, 3, 44, 22, 0, 207, 205,
		1, 0, 0, 0, 208, 211, 1, 0, 0, 0, 209, 207, 1, 0, 0, 0, 209, 210, 1, 0,
		0, 0, 210, 43, 1, 0, 0, 0, 211, 209, 1, 0, 0, 0, 212, 217, 3, 46, 23, 0,
		213, 214, 5, 20, 0, 0, 214, 216, 3, 46, 23, 0, 215, 213, 1, 0, 0, 0, 216,
		219, 1, 0, 0, 0, 217, 215, 1, 0, 0, 0, 217, 218, 1, 0, 0, 0, 218, 45, 1,
		0, 0, 0, 219, 217, 1, 0, 0, 0, 220, 221, 6, 23, -1, 0, 221, 222, 3, 48,
		24, 0, 222, 228, 1, 0, 0, 0, 223, 224, 10, 1, 0, 0, 224, 225, 7, 2, 0,
		0, 225, 227, 3, 46, 23, 2, 226, 223, 1, 0, 0, 0, 227, 230, 1, 0, 0, 0,
		228, 226, 1, 0, 0, 0, 228, 229, 1, 0, 0, 0, 229, 47, 1, 0, 0, 0, 230, 228,
		1, 0, 0, 0, 231, 232, 6, 24, -1, 0, 232, 233, 3, 50, 25, 0, 233, 242, 1,
		0, 0, 0, 234, 235, 10, 2, 0, 0, 235, 236, 7, 3, 0, 0, 236, 241, 3, 48,
		24, 3, 237, 238, 10, 1, 0, 0, 238, 239, 7, 4, 0, 0, 239, 241, 3, 48, 24,
		2, 240, 234, 1, 0, 0, 0, 240, 237, 1, 0, 0, 0, 241, 244, 1, 0, 0, 0, 242,
		240, 1, 0, 0, 0, 242, 243, 1, 0, 0, 0, 243, 49, 1, 0, 0, 0, 244, 242, 1,
		0, 0, 0, 245, 259, 3, 52, 26, 0, 246, 248, 5, 31, 0, 0, 247, 246, 1, 0,
		0, 0, 248, 249, 1, 0, 0, 0, 249, 247, 1, 0, 0, 0, 249, 250, 1, 0, 0, 0,
		250, 251, 1, 0, 0, 0, 251, 259, 3, 52, 26, 0, 252, 254, 5, 30, 0, 0, 253,
		252, 1, 0, 0, 0, 254, 255, 1, 0, 0, 0, 255, 253, 1, 0, 0, 0, 255, 256,
		1, 0, 0, 0, 256, 257, 1, 0, 0, 0, 257, 259, 3, 52, 26, 0, 258, 245, 1,
		0, 0, 0, 258, 247, 1, 0, 0, 0, 258, 253, 1, 0, 0, 0, 259, 51, 1, 0, 0,
		0, 260, 261, 6, 26, -1, 0, 261, 262, 3, 54, 27, 0, 262, 287, 1, 0, 0, 0,
		263, 264, 10, 3, 0, 0, 264, 266, 5, 28, 0, 0, 265, 267, 5, 32, 0, 0, 266,
		265, 1, 0, 0, 0, 266, 267, 1, 0, 0, 0, 267, 268, 1, 0, 0, 0, 268, 286,
		5, 48, 0, 0, 269, 270, 10, 2, 0, 0, 270, 271, 5, 28, 0, 0, 271, 272, 5,
		48, 0, 0, 272, 274, 5, 26, 0, 0, 273, 275, 3, 56, 28, 0, 274, 273, 1, 0,
		0, 0, 274, 275, 1, 0, 0, 0, 275, 276, 1, 0, 0, 0, 276, 286, 5, 27, 0, 0,
		277, 278, 10, 1, 0, 0, 278, 280, 5, 22, 0, 0, 279, 281, 5, 32, 0, 0, 280,
		279, 1, 0, 0, 0, 280, 281, 1, 0, 0, 0, 281, 282, 1, 0, 0, 0, 282, 283,
		3, 40, 20, 0, 283, 284, 5, 23, 0, 0, 284, 286, 1, 0, 0, 0, 285, 263, 1,
		0, 0, 0, 285, 269, 1, 0, 0, 0, 285, 277, 1, 0, 0, 0, 286, 289, 1, 0, 0,
		0, 287, 285, 1, 0, 0, 0, 287, 288, 1, 0, 0, 0, 288, 53, 1, 0, 0, 0, 289,
		287, 1, 0, 0, 0, 290, 292, 5, 28, 0, 0, 291, 290, 1, 0, 0, 0, 291, 292,
		1, 0, 0, 0, 292, 293, 1, 0, 0, 0, 293, 299, 5, 48, 0, 0, 294, 296, 5, 26,
		0, 0, 295, 297, 3, 56, 28, 0, 296, 295, 1, 0, 0, 0, 296, 297, 1, 0, 0,
		0, 297, 298, 1, 0, 0, 0, 298, 300, 5, 27, 0, 0, 299, 294, 1, 0, 0, 0, 299,
		300, 1, 0, 0, 0, 300, 342, 1, 0, 0, 0, 301, 302, 5, 26, 0, 0, 302, 303,
		3, 40, 20, 0, 303, 304, 5, 27, 0, 0, 304, 342, 1, 0, 0, 0, 305, 307, 5,
		22, 0, 0, 306, 308, 3, 58, 29, 0, 307, 306, 1, 0, 0, 0, 307, 308, 1, 0,
		0, 0, 308, 310, 1, 0, 0, 0, 309, 311, 5, 29, 0, 0, 310, 309, 1, 0, 0, 0,
		310, 311, 1, 0, 0, 0, 311, 312, 1, 0, 0, 0, 312, 342, 5, 23, 0, 0, 313,
		315, 5, 24, 0, 0, 314, 316, 3, 64, 32, 0, 315, 314, 1, 0, 0, 0, 315, 316,
		1, 0, 0, 0, 316, 318, 1, 0, 0, 0, 317, 319, 5, 29, 0, 0, 318, 317, 1, 0,
		0, 0, 318, 319, 1, 0, 0, 0, 319, 320, 1, 0, 0, 0, 320, 342, 5, 25, 0, 0,
		321, 323, 5, 28, 0, 0, 322, 321, 1, 0, 0, 0, 322, 323, 1, 0, 0, 0, 323,
		324, 1, 0, 0, 0, 324, 329, 5, 48, 0, 0, 325, 326, 5, 28, 0, 0, 326, 328,
		5, 48, 0, 0, 327, 325, 1, 0, 0, 0, 328, 331, 1, 0, 0, 0, 329, 327, 1, 0,
		0, 0, 329, 330, 1, 0, 0, 0, 330, 332, 1, 0, 0, 0, 331, 329, 1, 0, 0, 0,
		332, 334, 5, 24, 0, 0, 333, 335, 3, 60, 30, 0, 334, 333, 1, 0, 0, 0, 334,
		335, 1, 0, 0, 0, 335, 337, 1, 0, 0, 0, 336, 338, 5, 29, 0, 0, 337, 336,
		1, 0, 0, 0, 337, 338, 1, 0, 0, 0, 338, 339, 1, 0, 0, 0, 339, 342, 5, 25,
		0, 0, 340, 342, 3, 68, 34, 0, 341, 291, 1, 0, 0, 0, 341, 301, 1, 0, 0,
		0, 341, 305, 1, 0, 0, 0, 341, 313, 1, 0, 0, 0, 341, 322, 1, 0, 0, 0, 341,
		340, 1, 0, 0, 0, 342, 55, 1, 0, 0, 0, 343, 348, 3, 40, 20, 0, 344, 345,
		5, 29, 0, 0, 345, 347, 3, 40, 20, 0, 346, 344, 1, 0, 0, 0, 347, 350, 1,
		0, 0, 0, 348, 346, 1, 0, 0, 0, 348, 349, 1, 0, 0, 0, 349, 57, 1, 0, 0,
		0, 350, 348, 1, 0, 0, 0, 351, 356, 3, 66, 33, 0, 352, 353, 5, 29, 0, 0,
		353, 355, 3, 66, 33, 0, 354, 352, 1, 0, 0, 0, 355, 358, 1, 0, 0, 0, 356,
		354, 1, 0, 0, 0, 356, 357, 1, 0, 0, 0, 357, 59, 1, 0, 0, 0, 358, 356, 1,
		0, 0, 0, 359, 360, 3, 62, 31, 0, 360, 361, 5, 33, 0, 0, 361, 369, 3, 40,
		20, 0, 362, 363, 5, 29, 0, 0, 363, 364, 3, 62, 31, 0, 364, 365, 5, 33,
		0, 0, 365, 366, 3, 40, 20, 0, 366, 368, 1, 0, 0, 0, 367, 362, 1, 0, 0,
		0, 368, 371, 1, 0, 0, 0, 369, 367, 1, 0, 0, 0, 369, 370, 1, 0, 0, 0, 370,
		61, 1, 0, 0, 0, 371, 369, 1, 0, 0, 0, 372, 374, 5, 32, 0, 0, 373, 372,
		1, 0, 0, 0, 373, 374, 1, 0, 0, 0, 374, 375, 1, 0, 0, 0, 375, 376, 5, 48,
		0, 0, 376, 63, 1, 0, 0, 0, 377, 378, 3, 66, 33, 0, 378, 379, 5, 33, 0,
		0, 379, 387, 3, 40, 20, 0, 380, 381, 5, 29, 0, 0, 381, 382, 3, 66, 33,
		0, 382, 383, 5, 33, 0, 0, 383, 384, 3, 40, 20, 0, 384, 386, 1, 0, 0, 0,
		385, 380, 1, 0, 0, 0, 386, 389, 1, 0, 0, 0, 387, 385, 1, 0, 0, 0, 387,
		388, 1, 0, 0, 0, 388, 65, 1, 0, 0, 0, 389, 387, 1, 0, 0, 0, 390, 392, 5,
		32, 0, 0, 391, 390, 1, 0, 0, 0, 391, 392, 1, 0, 0, 0, 392, 393, 1, 0, 0,
		0, 393, 394, 3, 40, 20, 0, 394, 67, 1, 0, 0, 0, 395, 397, 5, 30, 0, 0,
		396, 395, 1, 0, 0, 0, 396, 397, 1, 0, 0, 0, 397, 398, 1, 0, 0, 0, 398,
		410, 5, 44, 0, 0, 399, 410, 5, 45, 0, 0, 400, 402, 5, 30, 0, 0, 401, 400,
		1, 0, 0, 0, 401, 402, 1, 0, 0, 0, 402, 403, 1, 0, 0, 0, 403, 410, 5, 43,
		0, 0, 404, 410, 5, 46, 0, 0, 405, 410, 5, 47, 0, 0, 406, 410, 5, 38, 0,
		0, 407, 410, 5, 39, 0, 0, 408, 410, 5, 40, 0, 0, 409, 396, 1, 0, 0, 0,
		409, 399, 1, 0, 0, 0, 409, 401, 1, 0, 0, 0, 409, 404, 1, 0, 0, 0, 409,
		405, 1, 0, 0, 0, 409, 406, 1, 0, 0, 0, 409, 407, 1, 0, 0, 0, 409, 408,
		1, 0, 0, 0, 410, 69, 1, 0, 0, 0, 51, 82, 93, 100, 105, 114, 117, 130, 135,
		137, 149, 154, 161, 169, 172, 179, 188, 202, 209, 217, 228, 240, 242, 249,
		255, 258, 266, 274, 280, 285, 287, 291, 296, 299, 307, 310, 315, 318, 322,
		329, 334, 337, 341, 348, 356, 369, 373, 387, 391, 396, 401, 409,
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
	staticData := &CommandsParserStaticData
	staticData.once.Do(commandsParserInit)
}

// NewCommandsParser produces a new parser instance for the optional input antlr.TokenStream.
func NewCommandsParser(input antlr.TokenStream) *CommandsParser {
	CommandsParserInit()
	this := new(CommandsParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &CommandsParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
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
	CommandsParserT__4           = 5
	CommandsParserT__5           = 6
	CommandsParserT__6           = 7
	CommandsParserT__7           = 8
	CommandsParserCOMMAND        = 9
	CommandsParserFLAG           = 10
	CommandsParserARROW          = 11
	CommandsParserEQUAL_ASSIGN   = 12
	CommandsParserEQUALS         = 13
	CommandsParserNOT_EQUALS     = 14
	CommandsParserIN             = 15
	CommandsParserLESS           = 16
	CommandsParserLESS_EQUALS    = 17
	CommandsParserGREATER_EQUALS = 18
	CommandsParserGREATER        = 19
	CommandsParserLOGICAL_AND    = 20
	CommandsParserLOGICAL_OR     = 21
	CommandsParserLBRACKET       = 22
	CommandsParserRPRACKET       = 23
	CommandsParserLBRACE         = 24
	CommandsParserRBRACE         = 25
	CommandsParserLPAREN         = 26
	CommandsParserRPAREN         = 27
	CommandsParserDOT            = 28
	CommandsParserCOMMA          = 29
	CommandsParserMINUS          = 30
	CommandsParserEXCLAM         = 31
	CommandsParserQUESTIONMARK   = 32
	CommandsParserCOLON          = 33
	CommandsParserPLUS           = 34
	CommandsParserSTAR           = 35
	CommandsParserSLASH          = 36
	CommandsParserPERCENT        = 37
	CommandsParserCEL_TRUE       = 38
	CommandsParserCEL_FALSE      = 39
	CommandsParserNUL            = 40
	CommandsParserWHITESPACE     = 41
	CommandsParserCOMMENT        = 42
	CommandsParserNUM_FLOAT      = 43
	CommandsParserNUM_INT        = 44
	CommandsParserNUM_UINT       = 45
	CommandsParserSTRING         = 46
	CommandsParserBYTES          = 47
	CommandsParserIDENTIFIER     = 48
)

// CommandsParser rules.
const (
	CommandsParserRULE_startCommand         = 0
	CommandsParserRULE_command              = 1
	CommandsParserRULE_help                 = 2
	CommandsParserRULE_let                  = 3
	CommandsParserRULE_declare              = 4
	CommandsParserRULE_varDecl              = 5
	CommandsParserRULE_fnDecl               = 6
	CommandsParserRULE_param                = 7
	CommandsParserRULE_delete               = 8
	CommandsParserRULE_simple               = 9
	CommandsParserRULE_empty                = 10
	CommandsParserRULE_compile              = 11
	CommandsParserRULE_parse                = 12
	CommandsParserRULE_exprCmd              = 13
	CommandsParserRULE_qualId               = 14
	CommandsParserRULE_startType            = 15
	CommandsParserRULE_type                 = 16
	CommandsParserRULE_typeId               = 17
	CommandsParserRULE_typeParamList        = 18
	CommandsParserRULE_start                = 19
	CommandsParserRULE_expr                 = 20
	CommandsParserRULE_conditionalOr        = 21
	CommandsParserRULE_conditionalAnd       = 22
	CommandsParserRULE_relation             = 23
	CommandsParserRULE_calc                 = 24
	CommandsParserRULE_unary                = 25
	CommandsParserRULE_member               = 26
	CommandsParserRULE_primary              = 27
	CommandsParserRULE_exprList             = 28
	CommandsParserRULE_listInit             = 29
	CommandsParserRULE_fieldInitializerList = 30
	CommandsParserRULE_optField             = 31
	CommandsParserRULE_mapInitializerList   = 32
	CommandsParserRULE_optExpr              = 33
	CommandsParserRULE_literal              = 34
)

// IStartCommandContext is an interface to support dynamic dispatch.
type IStartCommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Command() ICommandContext
	EOF() antlr.TerminalNode

	// IsStartCommandContext differentiates from other interfaces.
	IsStartCommandContext()
}

type StartCommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartCommandContext() *StartCommandContext {
	var p = new(StartCommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_startCommand
	return p
}

func InitEmptyStartCommandContext(p *StartCommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_startCommand
}

func (*StartCommandContext) IsStartCommandContext() {}

func NewStartCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartCommandContext {
	var p = new(StartCommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *StartCommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitStartCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) StartCommand() (localctx IStartCommandContext) {
	localctx = NewStartCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, CommandsParserRULE_startCommand)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(70)
		p.Command()
	}
	{
		p.SetState(71)
		p.Match(CommandsParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICommandContext is an interface to support dynamic dispatch.
type ICommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Help() IHelpContext
	Let() ILetContext
	Declare() IDeclareContext
	Delete_() IDeleteContext
	Simple() ISimpleContext
	Compile() ICompileContext
	Parse() IParseContext
	ExprCmd() IExprCmdContext
	Empty() IEmptyContext

	// IsCommandContext differentiates from other interfaces.
	IsCommandContext()
}

type CommandContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandContext() *CommandContext {
	var p = new(CommandContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_command
	return p
}

func InitEmptyCommandContext(p *CommandContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_command
}

func (*CommandContext) IsCommandContext() {}

func NewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandContext {
	var p = new(CommandContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_command

	return p
}

func (s *CommandContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandContext) Help() IHelpContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IHelpContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IHelpContext)
}

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

func (s *CommandContext) Delete_() IDeleteContext {
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

func (s *CommandContext) Compile() ICompileContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompileContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompileContext)
}

func (s *CommandContext) Parse() IParseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParseContext)
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

func (s *CommandContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCommand(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Command() (localctx ICommandContext) {
	localctx = NewCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, CommandsParserRULE_command)
	p.SetState(82)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case CommandsParserT__0, CommandsParserT__1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(73)
			p.Help()
		}

	case CommandsParserT__2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(74)
			p.Let()
		}

	case CommandsParserT__3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(75)
			p.Declare()
		}

	case CommandsParserT__4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(76)
			p.Delete_()
		}

	case CommandsParserCOMMAND:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(77)
			p.Simple()
		}

	case CommandsParserT__5:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(78)
			p.Compile()
		}

	case CommandsParserT__6:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(79)
			p.Parse()
		}

	case CommandsParserT__7, CommandsParserLBRACKET, CommandsParserLBRACE, CommandsParserLPAREN, CommandsParserDOT, CommandsParserMINUS, CommandsParserEXCLAM, CommandsParserCEL_TRUE, CommandsParserCEL_FALSE, CommandsParserNUL, CommandsParserNUM_FLOAT, CommandsParserNUM_INT, CommandsParserNUM_UINT, CommandsParserSTRING, CommandsParserBYTES, CommandsParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(80)
			p.ExprCmd()
		}

	case CommandsParserEOF:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(81)
			p.Empty()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IHelpContext is an interface to support dynamic dispatch.
type IHelpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsHelpContext differentiates from other interfaces.
	IsHelpContext()
}

type HelpContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHelpContext() *HelpContext {
	var p = new(HelpContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_help
	return p
}

func InitEmptyHelpContext(p *HelpContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_help
}

func (*HelpContext) IsHelpContext() {}

func NewHelpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HelpContext {
	var p = new(HelpContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_help

	return p
}

func (s *HelpContext) GetParser() antlr.Parser { return s.parser }
func (s *HelpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HelpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HelpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterHelp(s)
	}
}

func (s *HelpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitHelp(s)
	}
}

func (s *HelpContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitHelp(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Help() (localctx IHelpContext) {
	localctx = NewHelpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, CommandsParserRULE_help)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(84)
		_la = p.GetTokenStream().LA(1)

		if !(_la == CommandsParserT__0 || _la == CommandsParserT__1) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILetContext is an interface to support dynamic dispatch.
type ILetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVar_ returns the var_ rule contexts.
	GetVar_() IVarDeclContext

	// GetFn returns the fn rule contexts.
	GetFn() IFnDeclContext

	// GetE returns the e rule contexts.
	GetE() IExprContext

	// SetVar_ sets the var_ rule contexts.
	SetVar_(IVarDeclContext)

	// SetFn sets the fn rule contexts.
	SetFn(IFnDeclContext)

	// SetE sets the e rule contexts.
	SetE(IExprContext)

	// Getter signatures
	Expr() IExprContext
	EQUAL_ASSIGN() antlr.TerminalNode
	ARROW() antlr.TerminalNode
	VarDecl() IVarDeclContext
	FnDecl() IFnDeclContext

	// IsLetContext differentiates from other interfaces.
	IsLetContext()
}

type LetContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	var_   IVarDeclContext
	fn     IFnDeclContext
	e      IExprContext
}

func NewEmptyLetContext() *LetContext {
	var p = new(LetContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_let
	return p
}

func InitEmptyLetContext(p *LetContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_let
}

func (*LetContext) IsLetContext() {}

func NewLetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LetContext {
	var p = new(LetContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_let

	return p
}

func (s *LetContext) GetParser() antlr.Parser { return s.parser }

func (s *LetContext) GetVar_() IVarDeclContext { return s.var_ }

func (s *LetContext) GetFn() IFnDeclContext { return s.fn }

func (s *LetContext) GetE() IExprContext { return s.e }

func (s *LetContext) SetVar_(v IVarDeclContext) { s.var_ = v }

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

func (s *LetContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitLet(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Let() (localctx ILetContext) {
	localctx = NewLetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, CommandsParserRULE_let)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(86)
		p.Match(CommandsParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(93)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(87)

			var _x = p.VarDecl()

			localctx.(*LetContext).var_ = _x
		}
		{
			p.SetState(88)
			p.Match(CommandsParserEQUAL_ASSIGN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		{
			p.SetState(90)

			var _x = p.FnDecl()

			localctx.(*LetContext).fn = _x
		}
		{
			p.SetState(91)
			p.Match(CommandsParserARROW)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	{
		p.SetState(95)

		var _x = p.Expr()

		localctx.(*LetContext).e = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDeclareContext is an interface to support dynamic dispatch.
type IDeclareContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVar_ returns the var_ rule contexts.
	GetVar_() IVarDeclContext

	// GetFn returns the fn rule contexts.
	GetFn() IFnDeclContext

	// SetVar_ sets the var_ rule contexts.
	SetVar_(IVarDeclContext)

	// SetFn sets the fn rule contexts.
	SetFn(IFnDeclContext)

	// Getter signatures
	VarDecl() IVarDeclContext
	FnDecl() IFnDeclContext

	// IsDeclareContext differentiates from other interfaces.
	IsDeclareContext()
}

type DeclareContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	var_   IVarDeclContext
	fn     IFnDeclContext
}

func NewEmptyDeclareContext() *DeclareContext {
	var p = new(DeclareContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_declare
	return p
}

func InitEmptyDeclareContext(p *DeclareContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_declare
}

func (*DeclareContext) IsDeclareContext() {}

func NewDeclareContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclareContext {
	var p = new(DeclareContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_declare

	return p
}

func (s *DeclareContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclareContext) GetVar_() IVarDeclContext { return s.var_ }

func (s *DeclareContext) GetFn() IFnDeclContext { return s.fn }

func (s *DeclareContext) SetVar_(v IVarDeclContext) { s.var_ = v }

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

func (s *DeclareContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitDeclare(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Declare() (localctx IDeclareContext) {
	localctx = NewDeclareContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, CommandsParserRULE_declare)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(97)
		p.Match(CommandsParserT__3)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(98)

			var _x = p.VarDecl()

			localctx.(*DeclareContext).var_ = _x
		}

	case 2:
		{
			p.SetState(99)

			var _x = p.FnDecl()

			localctx.(*DeclareContext).fn = _x
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	QualId() IQualIdContext
	COLON() antlr.TerminalNode
	Type_() ITypeContext

	// IsVarDeclContext differentiates from other interfaces.
	IsVarDeclContext()
}

type VarDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	id     IQualIdContext
	t      ITypeContext
}

func NewEmptyVarDeclContext() *VarDeclContext {
	var p = new(VarDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_varDecl
	return p
}

func InitEmptyVarDeclContext(p *VarDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_varDecl
}

func (*VarDeclContext) IsVarDeclContext() {}

func NewVarDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarDeclContext {
	var p = new(VarDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *VarDeclContext) Type_() ITypeContext {
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

func (s *VarDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitVarDecl(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) VarDecl() (localctx IVarDeclContext) {
	localctx = NewVarDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, CommandsParserRULE_varDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(102)

		var _x = p.QualId()

		localctx.(*VarDeclContext).id = _x
	}
	p.SetState(105)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserCOLON {
		{
			p.SetState(103)
			p.Match(CommandsParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(104)

			var _x = p.Type_()

			localctx.(*VarDeclContext).t = _x
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	COLON() antlr.TerminalNode
	QualId() IQualIdContext
	Type_() ITypeContext
	AllParam() []IParamContext
	Param(i int) IParamContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFnDeclContext differentiates from other interfaces.
	IsFnDeclContext()
}

type FnDeclContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	id     IQualIdContext
	_param IParamContext
	params []IParamContext
	rType  ITypeContext
}

func NewEmptyFnDeclContext() *FnDeclContext {
	var p = new(FnDeclContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_fnDecl
	return p
}

func InitEmptyFnDeclContext(p *FnDeclContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_fnDecl
}

func (*FnDeclContext) IsFnDeclContext() {}

func NewFnDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FnDeclContext {
	var p = new(FnDeclContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *FnDeclContext) Type_() ITypeContext {
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

func (s *FnDeclContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitFnDecl(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) FnDecl() (localctx IFnDeclContext) {
	localctx = NewFnDeclContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, CommandsParserRULE_fnDecl)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(107)

		var _x = p.QualId()

		localctx.(*FnDeclContext).id = _x
	}
	{
		p.SetState(108)
		p.Match(CommandsParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(117)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserIDENTIFIER {
		{
			p.SetState(109)

			var _x = p.Param()

			localctx.(*FnDeclContext)._param = _x
		}
		localctx.(*FnDeclContext).params = append(localctx.(*FnDeclContext).params, localctx.(*FnDeclContext)._param)
		p.SetState(114)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == CommandsParserCOMMA {
			{
				p.SetState(110)
				p.Match(CommandsParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(111)

				var _x = p.Param()

				localctx.(*FnDeclContext)._param = _x
			}
			localctx.(*FnDeclContext).params = append(localctx.(*FnDeclContext).params, localctx.(*FnDeclContext)._param)

			p.SetState(116)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(119)
		p.Match(CommandsParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(120)
		p.Match(CommandsParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(121)

		var _x = p.Type_()

		localctx.(*FnDeclContext).rType = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	COLON() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	Type_() ITypeContext

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}

type ParamContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	pid    antlr.Token
	t      ITypeContext
}

func NewEmptyParamContext() *ParamContext {
	var p = new(ParamContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_param
	return p
}

func InitEmptyParamContext(p *ParamContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_param
}

func (*ParamContext) IsParamContext() {}

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext {
	var p = new(ParamContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *ParamContext) Type_() ITypeContext {
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

func (s *ParamContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitParam(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Param() (localctx IParamContext) {
	localctx = NewParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, CommandsParserRULE_param)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(123)

		var _m = p.Match(CommandsParserIDENTIFIER)

		localctx.(*ParamContext).pid = _m
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(124)
		p.Match(CommandsParserCOLON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(125)

		var _x = p.Type_()

		localctx.(*ParamContext).t = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDeleteContext is an interface to support dynamic dispatch.
type IDeleteContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetVar_ returns the var_ rule contexts.
	GetVar_() IVarDeclContext

	// GetFn returns the fn rule contexts.
	GetFn() IFnDeclContext

	// SetVar_ sets the var_ rule contexts.
	SetVar_(IVarDeclContext)

	// SetFn sets the fn rule contexts.
	SetFn(IFnDeclContext)

	// Getter signatures
	VarDecl() IVarDeclContext
	FnDecl() IFnDeclContext

	// IsDeleteContext differentiates from other interfaces.
	IsDeleteContext()
}

type DeleteContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	var_   IVarDeclContext
	fn     IFnDeclContext
}

func NewEmptyDeleteContext() *DeleteContext {
	var p = new(DeleteContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_delete
	return p
}

func InitEmptyDeleteContext(p *DeleteContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_delete
}

func (*DeleteContext) IsDeleteContext() {}

func NewDeleteContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeleteContext {
	var p = new(DeleteContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_delete

	return p
}

func (s *DeleteContext) GetParser() antlr.Parser { return s.parser }

func (s *DeleteContext) GetVar_() IVarDeclContext { return s.var_ }

func (s *DeleteContext) GetFn() IFnDeclContext { return s.fn }

func (s *DeleteContext) SetVar_(v IVarDeclContext) { s.var_ = v }

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

func (s *DeleteContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitDelete(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Delete_() (localctx IDeleteContext) {
	localctx = NewDeleteContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, CommandsParserRULE_delete)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(127)
		p.Match(CommandsParserT__4)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(130)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(128)

			var _x = p.VarDecl()

			localctx.(*DeleteContext).var_ = _x
		}

	case 2:
		{
			p.SetState(129)

			var _x = p.FnDecl()

			localctx.(*DeleteContext).fn = _x
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	COMMAND() antlr.TerminalNode
	AllFLAG() []antlr.TerminalNode
	FLAG(i int) antlr.TerminalNode
	AllSTRING() []antlr.TerminalNode
	STRING(i int) antlr.TerminalNode

	// IsSimpleContext differentiates from other interfaces.
	IsSimpleContext()
}

type SimpleContext struct {
	antlr.BaseParserRuleContext
	parser  antlr.Parser
	cmd     antlr.Token
	_FLAG   antlr.Token
	args    []antlr.Token
	_STRING antlr.Token
}

func NewEmptySimpleContext() *SimpleContext {
	var p = new(SimpleContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_simple
	return p
}

func InitEmptySimpleContext(p *SimpleContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_simple
}

func (*SimpleContext) IsSimpleContext() {}

func NewSimpleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SimpleContext {
	var p = new(SimpleContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *SimpleContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitSimple(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Simple() (localctx ISimpleContext) {
	localctx = NewSimpleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, CommandsParserRULE_simple)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(132)

		var _m = p.Match(CommandsParserCOMMAND)

		localctx.(*SimpleContext).cmd = _m
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserFLAG || _la == CommandsParserSTRING {
		p.SetState(135)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case CommandsParserFLAG:
			{
				p.SetState(133)

				var _m = p.Match(CommandsParserFLAG)

				localctx.(*SimpleContext)._FLAG = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*SimpleContext).args = append(localctx.(*SimpleContext).args, localctx.(*SimpleContext)._FLAG)

		case CommandsParserSTRING:
			{
				p.SetState(134)

				var _m = p.Match(CommandsParserSTRING)

				localctx.(*SimpleContext)._STRING = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*SimpleContext).args = append(localctx.(*SimpleContext).args, localctx.(*SimpleContext)._STRING)

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(139)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEmptyContext() *EmptyContext {
	var p = new(EmptyContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_empty
	return p
}

func InitEmptyEmptyContext(p *EmptyContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_empty
}

func (*EmptyContext) IsEmptyContext() {}

func NewEmptyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EmptyContext {
	var p = new(EmptyContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *EmptyContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitEmpty(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Empty() (localctx IEmptyContext) {
	localctx = NewEmptyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, CommandsParserRULE_empty)
	p.EnterOuterAlt(localctx, 1)

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICompileContext is an interface to support dynamic dispatch.
type ICompileContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetE returns the e rule contexts.
	GetE() IExprContext

	// SetE sets the e rule contexts.
	SetE(IExprContext)

	// Getter signatures
	Expr() IExprContext

	// IsCompileContext differentiates from other interfaces.
	IsCompileContext()
}

type CompileContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IExprContext
}

func NewEmptyCompileContext() *CompileContext {
	var p = new(CompileContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_compile
	return p
}

func InitEmptyCompileContext(p *CompileContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_compile
}

func (*CompileContext) IsCompileContext() {}

func NewCompileContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompileContext {
	var p = new(CompileContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_compile

	return p
}

func (s *CompileContext) GetParser() antlr.Parser { return s.parser }

func (s *CompileContext) GetE() IExprContext { return s.e }

func (s *CompileContext) SetE(v IExprContext) { s.e = v }

func (s *CompileContext) Expr() IExprContext {
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

func (s *CompileContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompileContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompileContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterCompile(s)
	}
}

func (s *CompileContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitCompile(s)
	}
}

func (s *CompileContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCompile(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Compile() (localctx ICompileContext) {
	localctx = NewCompileContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, CommandsParserRULE_compile)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(142)
		p.Match(CommandsParserT__5)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(143)

		var _x = p.Expr()

		localctx.(*CompileContext).e = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IParseContext is an interface to support dynamic dispatch.
type IParseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetE returns the e rule contexts.
	GetE() IExprContext

	// SetE sets the e rule contexts.
	SetE(IExprContext)

	// Getter signatures
	Expr() IExprContext

	// IsParseContext differentiates from other interfaces.
	IsParseContext()
}

type ParseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IExprContext
}

func NewEmptyParseContext() *ParseContext {
	var p = new(ParseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_parse
	return p
}

func InitEmptyParseContext(p *ParseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_parse
}

func (*ParseContext) IsParseContext() {}

func NewParseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseContext {
	var p = new(ParseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_parse

	return p
}

func (s *ParseContext) GetParser() antlr.Parser { return s.parser }

func (s *ParseContext) GetE() IExprContext { return s.e }

func (s *ParseContext) SetE(v IExprContext) { s.e = v }

func (s *ParseContext) Expr() IExprContext {
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

func (s *ParseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterParse(s)
	}
}

func (s *ParseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitParse(s)
	}
}

func (s *ParseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitParse(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Parse() (localctx IParseContext) {
	localctx = NewParseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, CommandsParserRULE_parse)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(145)
		p.Match(CommandsParserT__6)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(146)

		var _x = p.Expr()

		localctx.(*ParseContext).e = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	Expr() IExprContext

	// IsExprCmdContext differentiates from other interfaces.
	IsExprCmdContext()
}

type ExprCmdContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IExprContext
}

func NewEmptyExprCmdContext() *ExprCmdContext {
	var p = new(ExprCmdContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_exprCmd
	return p
}

func InitEmptyExprCmdContext(p *ExprCmdContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_exprCmd
}

func (*ExprCmdContext) IsExprCmdContext() {}

func NewExprCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprCmdContext {
	var p = new(ExprCmdContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *ExprCmdContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitExprCmd(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ExprCmd() (localctx IExprCmdContext) {
	localctx = NewExprCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, CommandsParserRULE_exprCmd)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(149)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserT__7 {
		{
			p.SetState(148)
			p.Match(CommandsParserT__7)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(151)

		var _x = p.Expr()

		localctx.(*ExprCmdContext).e = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsQualIdContext differentiates from other interfaces.
	IsQualIdContext()
}

type QualIdContext struct {
	antlr.BaseParserRuleContext
	parser      antlr.Parser
	leadingDot  antlr.Token
	rid         antlr.Token
	_IDENTIFIER antlr.Token
	qualifiers  []antlr.Token
}

func NewEmptyQualIdContext() *QualIdContext {
	var p = new(QualIdContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_qualId
	return p
}

func InitEmptyQualIdContext(p *QualIdContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_qualId
}

func (*QualIdContext) IsQualIdContext() {}

func NewQualIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualIdContext {
	var p = new(QualIdContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *QualIdContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitQualId(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) QualId() (localctx IQualIdContext) {
	localctx = NewQualIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, CommandsParserRULE_qualId)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserDOT {
		{
			p.SetState(153)

			var _m = p.Match(CommandsParserDOT)

			localctx.(*QualIdContext).leadingDot = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(156)

		var _m = p.Match(CommandsParserIDENTIFIER)

		localctx.(*QualIdContext).rid = _m
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(161)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserDOT {
		{
			p.SetState(157)
			p.Match(CommandsParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(158)

			var _m = p.Match(CommandsParserIDENTIFIER)

			localctx.(*QualIdContext)._IDENTIFIER = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*QualIdContext).qualifiers = append(localctx.(*QualIdContext).qualifiers, localctx.(*QualIdContext)._IDENTIFIER)

		p.SetState(163)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	EOF() antlr.TerminalNode
	Type_() ITypeContext

	// IsStartTypeContext differentiates from other interfaces.
	IsStartTypeContext()
}

type StartTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	t      ITypeContext
}

func NewEmptyStartTypeContext() *StartTypeContext {
	var p = new(StartTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_startType
	return p
}

func InitEmptyStartTypeContext(p *StartTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_startType
}

func (*StartTypeContext) IsStartTypeContext() {}

func NewStartTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartTypeContext {
	var p = new(StartTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *StartTypeContext) Type_() ITypeContext {
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

func (s *StartTypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitStartType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) StartType() (localctx IStartTypeContext) {
	localctx = NewStartTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, CommandsParserRULE_startType)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(164)

		var _x = p.Type_()

		localctx.(*StartTypeContext).t = _x
	}
	{
		p.SetState(165)
		p.Match(CommandsParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	TypeId() ITypeIdContext
	TypeParamList() ITypeParamListContext

	// IsTypeContext differentiates from other interfaces.
	IsTypeContext()
}

type TypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	id     ITypeIdContext
	params ITypeParamListContext
}

func NewEmptyTypeContext() *TypeContext {
	var p = new(TypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_type
	return p
}

func InitEmptyTypeContext(p *TypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_type
}

func (*TypeContext) IsTypeContext() {}

func NewTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeContext {
	var p = new(TypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *TypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Type_() (localctx ITypeContext) {
	localctx = NewTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, CommandsParserRULE_type)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(167)

		var _x = p.TypeId()

		localctx.(*TypeContext).id = _x
	}
	p.SetState(169)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserLPAREN {
		{
			p.SetState(168)

			var _x = p.TypeParamList()

			localctx.(*TypeContext).params = _x
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	NUL() antlr.TerminalNode
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode

	// IsTypeIdContext differentiates from other interfaces.
	IsTypeIdContext()
}

type TypeIdContext struct {
	antlr.BaseParserRuleContext
	parser      antlr.Parser
	leadingDot  antlr.Token
	id          antlr.Token
	_IDENTIFIER antlr.Token
	qualifiers  []antlr.Token
}

func NewEmptyTypeIdContext() *TypeIdContext {
	var p = new(TypeIdContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_typeId
	return p
}

func InitEmptyTypeIdContext(p *TypeIdContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_typeId
}

func (*TypeIdContext) IsTypeIdContext() {}

func NewTypeIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeIdContext {
	var p = new(TypeIdContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *TypeIdContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitTypeId(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) TypeId() (localctx ITypeIdContext) {
	localctx = NewTypeIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, CommandsParserRULE_typeId)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(172)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserDOT {
		{
			p.SetState(171)

			var _m = p.Match(CommandsParserDOT)

			localctx.(*TypeIdContext).leadingDot = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(174)

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
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserDOT {
		{
			p.SetState(175)
			p.Match(CommandsParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(176)

			var _m = p.Match(CommandsParserIDENTIFIER)

			localctx.(*TypeIdContext)._IDENTIFIER = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*TypeIdContext).qualifiers = append(localctx.(*TypeIdContext).qualifiers, localctx.(*TypeIdContext)._IDENTIFIER)

		p.SetState(181)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllType_() []ITypeContext
	Type_(i int) ITypeContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsTypeParamListContext differentiates from other interfaces.
	IsTypeParamListContext()
}

type TypeParamListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	_type  ITypeContext
	types  []ITypeContext
}

func NewEmptyTypeParamListContext() *TypeParamListContext {
	var p = new(TypeParamListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_typeParamList
	return p
}

func InitEmptyTypeParamListContext(p *TypeParamListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_typeParamList
}

func (*TypeParamListContext) IsTypeParamListContext() {}

func NewTypeParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamListContext {
	var p = new(TypeParamListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *TypeParamListContext) AllType_() []ITypeContext {
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

func (s *TypeParamListContext) Type_(i int) ITypeContext {
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

func (s *TypeParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitTypeParamList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) TypeParamList() (localctx ITypeParamListContext) {
	localctx = NewTypeParamListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, CommandsParserRULE_typeParamList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(182)
		p.Match(CommandsParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(183)

		var _x = p.Type_()

		localctx.(*TypeParamListContext)._type = _x
	}
	localctx.(*TypeParamListContext).types = append(localctx.(*TypeParamListContext).types, localctx.(*TypeParamListContext)._type)
	p.SetState(188)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserCOMMA {
		{
			p.SetState(184)
			p.Match(CommandsParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(185)

			var _x = p.Type_()

			localctx.(*TypeParamListContext)._type = _x
		}
		localctx.(*TypeParamListContext).types = append(localctx.(*TypeParamListContext).types, localctx.(*TypeParamListContext)._type)

		p.SetState(190)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(191)
		p.Match(CommandsParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	EOF() antlr.TerminalNode
	Expr() IExprContext

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IExprContext
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_start
	return p
}

func InitEmptyStartContext(p *StartContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_start
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *StartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitStart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Start_() (localctx IStartContext) {
	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, CommandsParserRULE_start)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(193)

		var _x = p.Expr()

		localctx.(*StartContext).e = _x
	}
	{
		p.SetState(194)
		p.Match(CommandsParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	AllConditionalOr() []IConditionalOrContext
	ConditionalOr(i int) IConditionalOrContext
	COLON() antlr.TerminalNode
	QUESTIONMARK() antlr.TerminalNode
	Expr() IExprContext

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	e      IConditionalOrContext
	op     antlr.Token
	e1     IConditionalOrContext
	e2     IExprContext
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_expr
	return p
}

func InitEmptyExprContext(p *ExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_expr
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *ExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Expr() (localctx IExprContext) {
	localctx = NewExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, CommandsParserRULE_expr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(196)

		var _x = p.ConditionalOr()

		localctx.(*ExprContext).e = _x
	}
	p.SetState(202)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserQUESTIONMARK {
		{
			p.SetState(197)

			var _m = p.Match(CommandsParserQUESTIONMARK)

			localctx.(*ExprContext).op = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(198)

			var _x = p.ConditionalOr()

			localctx.(*ExprContext).e1 = _x
		}
		{
			p.SetState(199)
			p.Match(CommandsParserCOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(200)

			var _x = p.Expr()

			localctx.(*ExprContext).e2 = _x
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConditionalOrContext is an interface to support dynamic dispatch.
type IConditionalOrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS21 returns the s21 token.
	GetS21() antlr.Token

	// SetS21 sets the s21 token.
	SetS21(antlr.Token)

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

	// Getter signatures
	AllConditionalAnd() []IConditionalAndContext
	ConditionalAnd(i int) IConditionalAndContext
	AllLOGICAL_OR() []antlr.TerminalNode
	LOGICAL_OR(i int) antlr.TerminalNode

	// IsConditionalOrContext differentiates from other interfaces.
	IsConditionalOrContext()
}

type ConditionalOrContext struct {
	antlr.BaseParserRuleContext
	parser          antlr.Parser
	e               IConditionalAndContext
	s21             antlr.Token
	ops             []antlr.Token
	_conditionalAnd IConditionalAndContext
	e1              []IConditionalAndContext
}

func NewEmptyConditionalOrContext() *ConditionalOrContext {
	var p = new(ConditionalOrContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_conditionalOr
	return p
}

func InitEmptyConditionalOrContext(p *ConditionalOrContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_conditionalOr
}

func (*ConditionalOrContext) IsConditionalOrContext() {}

func NewConditionalOrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalOrContext {
	var p = new(ConditionalOrContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_conditionalOr

	return p
}

func (s *ConditionalOrContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionalOrContext) GetS21() antlr.Token { return s.s21 }

func (s *ConditionalOrContext) SetS21(v antlr.Token) { s.s21 = v }

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

func (s *ConditionalOrContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitConditionalOr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ConditionalOr() (localctx IConditionalOrContext) {
	localctx = NewConditionalOrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, CommandsParserRULE_conditionalOr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(204)

		var _x = p.ConditionalAnd()

		localctx.(*ConditionalOrContext).e = _x
	}
	p.SetState(209)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserLOGICAL_OR {
		{
			p.SetState(205)

			var _m = p.Match(CommandsParserLOGICAL_OR)

			localctx.(*ConditionalOrContext).s21 = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*ConditionalOrContext).ops = append(localctx.(*ConditionalOrContext).ops, localctx.(*ConditionalOrContext).s21)
		{
			p.SetState(206)

			var _x = p.ConditionalAnd()

			localctx.(*ConditionalOrContext)._conditionalAnd = _x
		}
		localctx.(*ConditionalOrContext).e1 = append(localctx.(*ConditionalOrContext).e1, localctx.(*ConditionalOrContext)._conditionalAnd)

		p.SetState(211)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConditionalAndContext is an interface to support dynamic dispatch.
type IConditionalAndContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS20 returns the s20 token.
	GetS20() antlr.Token

	// SetS20 sets the s20 token.
	SetS20(antlr.Token)

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

	// Getter signatures
	AllRelation() []IRelationContext
	Relation(i int) IRelationContext
	AllLOGICAL_AND() []antlr.TerminalNode
	LOGICAL_AND(i int) antlr.TerminalNode

	// IsConditionalAndContext differentiates from other interfaces.
	IsConditionalAndContext()
}

type ConditionalAndContext struct {
	antlr.BaseParserRuleContext
	parser    antlr.Parser
	e         IRelationContext
	s20       antlr.Token
	ops       []antlr.Token
	_relation IRelationContext
	e1        []IRelationContext
}

func NewEmptyConditionalAndContext() *ConditionalAndContext {
	var p = new(ConditionalAndContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_conditionalAnd
	return p
}

func InitEmptyConditionalAndContext(p *ConditionalAndContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_conditionalAnd
}

func (*ConditionalAndContext) IsConditionalAndContext() {}

func NewConditionalAndContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalAndContext {
	var p = new(ConditionalAndContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_conditionalAnd

	return p
}

func (s *ConditionalAndContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionalAndContext) GetS20() antlr.Token { return s.s20 }

func (s *ConditionalAndContext) SetS20(v antlr.Token) { s.s20 = v }

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

func (s *ConditionalAndContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitConditionalAnd(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ConditionalAnd() (localctx IConditionalAndContext) {
	localctx = NewConditionalAndContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, CommandsParserRULE_conditionalAnd)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(212)

		var _x = p.relation(0)

		localctx.(*ConditionalAndContext).e = _x
	}
	p.SetState(217)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserLOGICAL_AND {
		{
			p.SetState(213)

			var _m = p.Match(CommandsParserLOGICAL_AND)

			localctx.(*ConditionalAndContext).s20 = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*ConditionalAndContext).ops = append(localctx.(*ConditionalAndContext).ops, localctx.(*ConditionalAndContext).s20)
		{
			p.SetState(214)

			var _x = p.relation(0)

			localctx.(*ConditionalAndContext)._relation = _x
		}
		localctx.(*ConditionalAndContext).e1 = append(localctx.(*ConditionalAndContext).e1, localctx.(*ConditionalAndContext)._relation)

		p.SetState(219)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	Calc() ICalcContext
	AllRelation() []IRelationContext
	Relation(i int) IRelationContext
	LESS() antlr.TerminalNode
	LESS_EQUALS() antlr.TerminalNode
	GREATER_EQUALS() antlr.TerminalNode
	GREATER() antlr.TerminalNode
	EQUALS() antlr.TerminalNode
	NOT_EQUALS() antlr.TerminalNode
	IN() antlr.TerminalNode

	// IsRelationContext differentiates from other interfaces.
	IsRelationContext()
}

type RelationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	op     antlr.Token
}

func NewEmptyRelationContext() *RelationContext {
	var p = new(RelationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_relation
	return p
}

func InitEmptyRelationContext(p *RelationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_relation
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *RelationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
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
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewRelationContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IRelationContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 46
	p.EnterRecursionRule(localctx, 46, CommandsParserRULE_relation, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(221)
		p.calc(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(228)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewRelationContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_relation)
			p.SetState(223)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				goto errorExit
			}
			{
				p.SetState(224)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*RelationContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1040384) != 0) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*RelationContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(225)
				p.relation(2)
			}

		}
		p.SetState(230)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	Unary() IUnaryContext
	AllCalc() []ICalcContext
	Calc(i int) ICalcContext
	STAR() antlr.TerminalNode
	SLASH() antlr.TerminalNode
	PERCENT() antlr.TerminalNode
	PLUS() antlr.TerminalNode
	MINUS() antlr.TerminalNode

	// IsCalcContext differentiates from other interfaces.
	IsCalcContext()
}

type CalcContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	op     antlr.Token
}

func NewEmptyCalcContext() *CalcContext {
	var p = new(CalcContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_calc
	return p
}

func InitEmptyCalcContext(p *CalcContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_calc
}

func (*CalcContext) IsCalcContext() {}

func NewCalcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CalcContext {
	var p = new(CalcContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *CalcContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
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
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewCalcContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ICalcContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 48
	p.EnterRecursionRule(localctx, 48, CommandsParserRULE_calc, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(232)
		p.Unary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(242)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(240)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 20, p.GetParserRuleContext()) {
			case 1:
				localctx = NewCalcContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_calc)
				p.SetState(234)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(235)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*CalcContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&240518168576) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*CalcContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(236)
					p.calc(3)
				}

			case 2:
				localctx = NewCalcContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_calc)
				p.SetState(237)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(238)

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
					p.SetState(239)
					p.calc(2)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(244)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 21, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnaryContext() *UnaryContext {
	var p = new(UnaryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_unary
	return p
}

func InitEmptyUnaryContext(p *UnaryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_unary
}

func (*UnaryContext) IsUnaryContext() {}

func NewUnaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryContext {
	var p = new(UnaryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_unary

	return p
}

func (s *UnaryContext) GetParser() antlr.Parser { return s.parser }

func (s *UnaryContext) CopyAll(ctx *UnaryContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *UnaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type LogicalNotContext struct {
	UnaryContext
	s31 antlr.Token
	ops []antlr.Token
}

func NewLogicalNotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalNotContext {
	var p = new(LogicalNotContext)

	InitEmptyUnaryContext(&p.UnaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*UnaryContext))

	return p
}

func (s *LogicalNotContext) GetS31() antlr.Token { return s.s31 }

func (s *LogicalNotContext) SetS31(v antlr.Token) { s.s31 = v }

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

func (s *LogicalNotContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitLogicalNot(s)

	default:
		return t.VisitChildren(s)
	}
}

type MemberExprContext struct {
	UnaryContext
}

func NewMemberExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberExprContext {
	var p = new(MemberExprContext)

	InitEmptyUnaryContext(&p.UnaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*UnaryContext))

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

func (s *MemberExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitMemberExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type NegateContext struct {
	UnaryContext
	s30 antlr.Token
	ops []antlr.Token
}

func NewNegateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NegateContext {
	var p = new(NegateContext)

	InitEmptyUnaryContext(&p.UnaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*UnaryContext))

	return p
}

func (s *NegateContext) GetS30() antlr.Token { return s.s30 }

func (s *NegateContext) SetS30(v antlr.Token) { s.s30 = v }

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

func (s *NegateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitNegate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Unary() (localctx IUnaryContext) {
	localctx = NewUnaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, CommandsParserRULE_unary)
	var _la int

	var _alt int

	p.SetState(258)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		localctx = NewMemberExprContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(245)
			p.member(0)
		}

	case 2:
		localctx = NewLogicalNotContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		p.SetState(247)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == CommandsParserEXCLAM {
			{
				p.SetState(246)

				var _m = p.Match(CommandsParserEXCLAM)

				localctx.(*LogicalNotContext).s31 = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*LogicalNotContext).ops = append(localctx.(*LogicalNotContext).ops, localctx.(*LogicalNotContext).s31)

			p.SetState(249)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(251)
			p.member(0)
		}

	case 3:
		localctx = NewNegateContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		p.SetState(253)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(252)

					var _m = p.Match(CommandsParserMINUS)

					localctx.(*NegateContext).s30 = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				localctx.(*NegateContext).ops = append(localctx.(*NegateContext).ops, localctx.(*NegateContext).s30)

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(255)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 23, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}
		{
			p.SetState(257)
			p.member(0)
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMemberContext() *MemberContext {
	var p = new(MemberContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_member
	return p
}

func InitEmptyMemberContext(p *MemberContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_member
}

func (*MemberContext) IsMemberContext() {}

func NewMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberContext {
	var p = new(MemberContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_member

	return p
}

func (s *MemberContext) GetParser() antlr.Parser { return s.parser }

func (s *MemberContext) CopyAll(ctx *MemberContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *MemberContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type MemberCallContext struct {
	MemberContext
	op   antlr.Token
	id   antlr.Token
	open antlr.Token
	args IExprListContext
}

func NewMemberCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberCallContext {
	var p = new(MemberCallContext)

	InitEmptyMemberContext(&p.MemberContext)
	p.parser = parser
	p.CopyAll(ctx.(*MemberContext))

	return p
}

func (s *MemberCallContext) GetOp() antlr.Token { return s.op }

func (s *MemberCallContext) GetId() antlr.Token { return s.id }

func (s *MemberCallContext) GetOpen() antlr.Token { return s.open }

func (s *MemberCallContext) SetOp(v antlr.Token) { s.op = v }

func (s *MemberCallContext) SetId(v antlr.Token) { s.id = v }

func (s *MemberCallContext) SetOpen(v antlr.Token) { s.open = v }

func (s *MemberCallContext) GetArgs() IExprListContext { return s.args }

func (s *MemberCallContext) SetArgs(v IExprListContext) { s.args = v }

func (s *MemberCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberCallContext) Member() IMemberContext {
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

func (s *MemberCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserRPAREN, 0)
}

func (s *MemberCallContext) DOT() antlr.TerminalNode {
	return s.GetToken(CommandsParserDOT, 0)
}

func (s *MemberCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, 0)
}

func (s *MemberCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(CommandsParserLPAREN, 0)
}

func (s *MemberCallContext) ExprList() IExprListContext {
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

func (s *MemberCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterMemberCall(s)
	}
}

func (s *MemberCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitMemberCall(s)
	}
}

func (s *MemberCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitMemberCall(s)

	default:
		return t.VisitChildren(s)
	}
}

type SelectContext struct {
	MemberContext
	op  antlr.Token
	opt antlr.Token
	id  antlr.Token
}

func NewSelectContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectContext {
	var p = new(SelectContext)

	InitEmptyMemberContext(&p.MemberContext)
	p.parser = parser
	p.CopyAll(ctx.(*MemberContext))

	return p
}

func (s *SelectContext) GetOp() antlr.Token { return s.op }

func (s *SelectContext) GetOpt() antlr.Token { return s.opt }

func (s *SelectContext) GetId() antlr.Token { return s.id }

func (s *SelectContext) SetOp(v antlr.Token) { s.op = v }

func (s *SelectContext) SetOpt(v antlr.Token) { s.opt = v }

func (s *SelectContext) SetId(v antlr.Token) { s.id = v }

func (s *SelectContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectContext) Member() IMemberContext {
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

func (s *SelectContext) DOT() antlr.TerminalNode {
	return s.GetToken(CommandsParserDOT, 0)
}

func (s *SelectContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, 0)
}

func (s *SelectContext) QUESTIONMARK() antlr.TerminalNode {
	return s.GetToken(CommandsParserQUESTIONMARK, 0)
}

func (s *SelectContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterSelect(s)
	}
}

func (s *SelectContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitSelect(s)
	}
}

func (s *SelectContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitSelect(s)

	default:
		return t.VisitChildren(s)
	}
}

type PrimaryExprContext struct {
	MemberContext
}

func NewPrimaryExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PrimaryExprContext {
	var p = new(PrimaryExprContext)

	InitEmptyMemberContext(&p.MemberContext)
	p.parser = parser
	p.CopyAll(ctx.(*MemberContext))

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

func (s *PrimaryExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitPrimaryExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type IndexContext struct {
	MemberContext
	op    antlr.Token
	opt   antlr.Token
	index IExprContext
}

func NewIndexContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IndexContext {
	var p = new(IndexContext)

	InitEmptyMemberContext(&p.MemberContext)
	p.parser = parser
	p.CopyAll(ctx.(*MemberContext))

	return p
}

func (s *IndexContext) GetOp() antlr.Token { return s.op }

func (s *IndexContext) GetOpt() antlr.Token { return s.opt }

func (s *IndexContext) SetOp(v antlr.Token) { s.op = v }

func (s *IndexContext) SetOpt(v antlr.Token) { s.opt = v }

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

func (s *IndexContext) QUESTIONMARK() antlr.TerminalNode {
	return s.GetToken(CommandsParserQUESTIONMARK, 0)
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

func (s *IndexContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitIndex(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Member() (localctx IMemberContext) {
	return p.member(0)
}

func (p *CommandsParser) member(_p int) (localctx IMemberContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewMemberContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IMemberContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 52
	p.EnterRecursionRule(localctx, 52, CommandsParserRULE_member, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewPrimaryExprContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(261)
		p.Primary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(287)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(285)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 28, p.GetParserRuleContext()) {
			case 1:
				localctx = NewSelectContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_member)
				p.SetState(263)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(264)

					var _m = p.Match(CommandsParserDOT)

					localctx.(*SelectContext).op = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(266)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if _la == CommandsParserQUESTIONMARK {
					{
						p.SetState(265)

						var _m = p.Match(CommandsParserQUESTIONMARK)

						localctx.(*SelectContext).opt = _m
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}

				}
				{
					p.SetState(268)

					var _m = p.Match(CommandsParserIDENTIFIER)

					localctx.(*SelectContext).id = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case 2:
				localctx = NewMemberCallContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_member)
				p.SetState(269)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(270)

					var _m = p.Match(CommandsParserDOT)

					localctx.(*MemberCallContext).op = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(271)

					var _m = p.Match(CommandsParserIDENTIFIER)

					localctx.(*MemberCallContext).id = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(272)

					var _m = p.Match(CommandsParserLPAREN)

					localctx.(*MemberCallContext).open = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(274)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&556081583489024) != 0 {
					{
						p.SetState(273)

						var _x = p.ExprList()

						localctx.(*MemberCallContext).args = _x
					}

				}
				{
					p.SetState(276)
					p.Match(CommandsParserRPAREN)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case 3:
				localctx = NewIndexContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, CommandsParserRULE_member)
				p.SetState(277)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(278)

					var _m = p.Match(CommandsParserLBRACKET)

					localctx.(*IndexContext).op = _m
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				p.SetState(280)
				p.GetErrorHandler().Sync(p)
				if p.HasError() {
					goto errorExit
				}
				_la = p.GetTokenStream().LA(1)

				if _la == CommandsParserQUESTIONMARK {
					{
						p.SetState(279)

						var _m = p.Match(CommandsParserQUESTIONMARK)

						localctx.(*IndexContext).opt = _m
						if p.HasError() {
							// Recognition error - abort rule
							goto errorExit
						}
					}

				}
				{
					p.SetState(282)

					var _x = p.Expr()

					localctx.(*IndexContext).index = _x
				}
				{
					p.SetState(283)
					p.Match(CommandsParserRPRACKET)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(289)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 29, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryContext() *PrimaryContext {
	var p = new(PrimaryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_primary
	return p
}

func InitEmptyPrimaryContext(p *PrimaryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_primary
}

func (*PrimaryContext) IsPrimaryContext() {}

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext {
	var p = new(PrimaryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_primary

	return p
}

func (s *PrimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryContext) CopyAll(ctx *PrimaryContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *PrimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type CreateListContext struct {
	PrimaryContext
	op    antlr.Token
	elems IListInitContext
}

func NewCreateListContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateListContext {
	var p = new(CreateListContext)

	InitEmptyPrimaryContext(&p.PrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryContext))

	return p
}

func (s *CreateListContext) GetOp() antlr.Token { return s.op }

func (s *CreateListContext) SetOp(v antlr.Token) { s.op = v }

func (s *CreateListContext) GetElems() IListInitContext { return s.elems }

func (s *CreateListContext) SetElems(v IListInitContext) { s.elems = v }

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

func (s *CreateListContext) ListInit() IListInitContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IListInitContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IListInitContext)
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

func (s *CreateListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCreateList(s)

	default:
		return t.VisitChildren(s)
	}
}

type CreateStructContext struct {
	PrimaryContext
	op      antlr.Token
	entries IMapInitializerListContext
}

func NewCreateStructContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateStructContext {
	var p = new(CreateStructContext)

	InitEmptyPrimaryContext(&p.PrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryContext))

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

func (s *CreateStructContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCreateStruct(s)

	default:
		return t.VisitChildren(s)
	}
}

type ConstantLiteralContext struct {
	PrimaryContext
}

func NewConstantLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ConstantLiteralContext {
	var p = new(ConstantLiteralContext)

	InitEmptyPrimaryContext(&p.PrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryContext))

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

func (s *ConstantLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitConstantLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type NestedContext struct {
	PrimaryContext
	e IExprContext
}

func NewNestedContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NestedContext {
	var p = new(NestedContext)

	InitEmptyPrimaryContext(&p.PrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryContext))

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

func (s *NestedContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitNested(s)

	default:
		return t.VisitChildren(s)
	}
}

type CreateMessageContext struct {
	PrimaryContext
	leadingDot  antlr.Token
	_IDENTIFIER antlr.Token
	ids         []antlr.Token
	s28         antlr.Token
	ops         []antlr.Token
	op          antlr.Token
	entries     IFieldInitializerListContext
}

func NewCreateMessageContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateMessageContext {
	var p = new(CreateMessageContext)

	InitEmptyPrimaryContext(&p.PrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryContext))

	return p
}

func (s *CreateMessageContext) GetLeadingDot() antlr.Token { return s.leadingDot }

func (s *CreateMessageContext) Get_IDENTIFIER() antlr.Token { return s._IDENTIFIER }

func (s *CreateMessageContext) GetS28() antlr.Token { return s.s28 }

func (s *CreateMessageContext) GetOp() antlr.Token { return s.op }

func (s *CreateMessageContext) SetLeadingDot(v antlr.Token) { s.leadingDot = v }

func (s *CreateMessageContext) Set_IDENTIFIER(v antlr.Token) { s._IDENTIFIER = v }

func (s *CreateMessageContext) SetS28(v antlr.Token) { s.s28 = v }

func (s *CreateMessageContext) SetOp(v antlr.Token) { s.op = v }

func (s *CreateMessageContext) GetIds() []antlr.Token { return s.ids }

func (s *CreateMessageContext) GetOps() []antlr.Token { return s.ops }

func (s *CreateMessageContext) SetIds(v []antlr.Token) { s.ids = v }

func (s *CreateMessageContext) SetOps(v []antlr.Token) { s.ops = v }

func (s *CreateMessageContext) GetEntries() IFieldInitializerListContext { return s.entries }

func (s *CreateMessageContext) SetEntries(v IFieldInitializerListContext) { s.entries = v }

func (s *CreateMessageContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateMessageContext) RBRACE() antlr.TerminalNode {
	return s.GetToken(CommandsParserRBRACE, 0)
}

func (s *CreateMessageContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserIDENTIFIER)
}

func (s *CreateMessageContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, i)
}

func (s *CreateMessageContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(CommandsParserLBRACE, 0)
}

func (s *CreateMessageContext) COMMA() antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, 0)
}

func (s *CreateMessageContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserDOT)
}

func (s *CreateMessageContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserDOT, i)
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

func (s *CreateMessageContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitCreateMessage(s)

	default:
		return t.VisitChildren(s)
	}
}

type IdentOrGlobalCallContext struct {
	PrimaryContext
	leadingDot antlr.Token
	id         antlr.Token
	op         antlr.Token
	args       IExprListContext
}

func NewIdentOrGlobalCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentOrGlobalCallContext {
	var p = new(IdentOrGlobalCallContext)

	InitEmptyPrimaryContext(&p.PrimaryContext)
	p.parser = parser
	p.CopyAll(ctx.(*PrimaryContext))

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

func (s *IdentOrGlobalCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitIdentOrGlobalCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Primary() (localctx IPrimaryContext) {
	localctx = NewPrimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, CommandsParserRULE_primary)
	var _la int

	p.SetState(341)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 41, p.GetParserRuleContext()) {
	case 1:
		localctx = NewIdentOrGlobalCallContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(291)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserDOT {
			{
				p.SetState(290)

				var _m = p.Match(CommandsParserDOT)

				localctx.(*IdentOrGlobalCallContext).leadingDot = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(293)

			var _m = p.Match(CommandsParserIDENTIFIER)

			localctx.(*IdentOrGlobalCallContext).id = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(299)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 32, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(294)

				var _m = p.Match(CommandsParserLPAREN)

				localctx.(*IdentOrGlobalCallContext).op = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			p.SetState(296)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)

			if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&556081583489024) != 0 {
				{
					p.SetState(295)

					var _x = p.ExprList()

					localctx.(*IdentOrGlobalCallContext).args = _x
				}

			}
			{
				p.SetState(298)
				p.Match(CommandsParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		} else if p.HasError() { // JIM
			goto errorExit
		}

	case 2:
		localctx = NewNestedContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(301)
			p.Match(CommandsParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(302)

			var _x = p.Expr()

			localctx.(*NestedContext).e = _x
		}
		{
			p.SetState(303)
			p.Match(CommandsParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewCreateListContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(305)

			var _m = p.Match(CommandsParserLBRACKET)

			localctx.(*CreateListContext).op = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(307)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&556085878456320) != 0 {
			{
				p.SetState(306)

				var _x = p.ListInit()

				localctx.(*CreateListContext).elems = _x
			}

		}
		p.SetState(310)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserCOMMA {
			{
				p.SetState(309)
				p.Match(CommandsParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(312)
			p.Match(CommandsParserRPRACKET)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewCreateStructContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(313)

			var _m = p.Match(CommandsParserLBRACE)

			localctx.(*CreateStructContext).op = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(315)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&556085878456320) != 0 {
			{
				p.SetState(314)

				var _x = p.MapInitializerList()

				localctx.(*CreateStructContext).entries = _x
			}

		}
		p.SetState(318)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserCOMMA {
			{
				p.SetState(317)
				p.Match(CommandsParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(320)
			p.Match(CommandsParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewCreateMessageContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		p.SetState(322)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserDOT {
			{
				p.SetState(321)

				var _m = p.Match(CommandsParserDOT)

				localctx.(*CreateMessageContext).leadingDot = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(324)

			var _m = p.Match(CommandsParserIDENTIFIER)

			localctx.(*CreateMessageContext)._IDENTIFIER = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		localctx.(*CreateMessageContext).ids = append(localctx.(*CreateMessageContext).ids, localctx.(*CreateMessageContext)._IDENTIFIER)
		p.SetState(329)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == CommandsParserDOT {
			{
				p.SetState(325)

				var _m = p.Match(CommandsParserDOT)

				localctx.(*CreateMessageContext).s28 = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*CreateMessageContext).ops = append(localctx.(*CreateMessageContext).ops, localctx.(*CreateMessageContext).s28)
			{
				p.SetState(326)

				var _m = p.Match(CommandsParserIDENTIFIER)

				localctx.(*CreateMessageContext)._IDENTIFIER = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*CreateMessageContext).ids = append(localctx.(*CreateMessageContext).ids, localctx.(*CreateMessageContext)._IDENTIFIER)

			p.SetState(331)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(332)

			var _m = p.Match(CommandsParserLBRACE)

			localctx.(*CreateMessageContext).op = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(334)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserQUESTIONMARK || _la == CommandsParserIDENTIFIER {
			{
				p.SetState(333)

				var _x = p.FieldInitializerList()

				localctx.(*CreateMessageContext).entries = _x
			}

		}
		p.SetState(337)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserCOMMA {
			{
				p.SetState(336)
				p.Match(CommandsParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(339)
			p.Match(CommandsParserRBRACE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		localctx = NewConstantLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(340)
			p.Literal()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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

	// Getter signatures
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsExprListContext differentiates from other interfaces.
	IsExprListContext()
}

type ExprListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	_expr  IExprContext
	e      []IExprContext
}

func NewEmptyExprListContext() *ExprListContext {
	var p = new(ExprListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_exprList
	return p
}

func InitEmptyExprListContext(p *ExprListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_exprList
}

func (*ExprListContext) IsExprListContext() {}

func NewExprListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprListContext {
	var p = new(ExprListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

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

func (s *ExprListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitExprList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ExprList() (localctx IExprListContext) {
	localctx = NewExprListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, CommandsParserRULE_exprList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(343)

		var _x = p.Expr()

		localctx.(*ExprListContext)._expr = _x
	}
	localctx.(*ExprListContext).e = append(localctx.(*ExprListContext).e, localctx.(*ExprListContext)._expr)
	p.SetState(348)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == CommandsParserCOMMA {
		{
			p.SetState(344)
			p.Match(CommandsParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(345)

			var _x = p.Expr()

			localctx.(*ExprListContext)._expr = _x
		}
		localctx.(*ExprListContext).e = append(localctx.(*ExprListContext).e, localctx.(*ExprListContext)._expr)

		p.SetState(350)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IListInitContext is an interface to support dynamic dispatch.
type IListInitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Get_optExpr returns the _optExpr rule contexts.
	Get_optExpr() IOptExprContext

	// Set_optExpr sets the _optExpr rule contexts.
	Set_optExpr(IOptExprContext)

	// GetElems returns the elems rule context list.
	GetElems() []IOptExprContext

	// SetElems sets the elems rule context list.
	SetElems([]IOptExprContext)

	// Getter signatures
	AllOptExpr() []IOptExprContext
	OptExpr(i int) IOptExprContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsListInitContext differentiates from other interfaces.
	IsListInitContext()
}

type ListInitContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	_optExpr IOptExprContext
	elems    []IOptExprContext
}

func NewEmptyListInitContext() *ListInitContext {
	var p = new(ListInitContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_listInit
	return p
}

func InitEmptyListInitContext(p *ListInitContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_listInit
}

func (*ListInitContext) IsListInitContext() {}

func NewListInitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListInitContext {
	var p = new(ListInitContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_listInit

	return p
}

func (s *ListInitContext) GetParser() antlr.Parser { return s.parser }

func (s *ListInitContext) Get_optExpr() IOptExprContext { return s._optExpr }

func (s *ListInitContext) Set_optExpr(v IOptExprContext) { s._optExpr = v }

func (s *ListInitContext) GetElems() []IOptExprContext { return s.elems }

func (s *ListInitContext) SetElems(v []IOptExprContext) { s.elems = v }

func (s *ListInitContext) AllOptExpr() []IOptExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOptExprContext); ok {
			len++
		}
	}

	tst := make([]IOptExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOptExprContext); ok {
			tst[i] = t.(IOptExprContext)
			i++
		}
	}

	return tst
}

func (s *ListInitContext) OptExpr(i int) IOptExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOptExprContext); ok {
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

	return t.(IOptExprContext)
}

func (s *ListInitContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOMMA)
}

func (s *ListInitContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOMMA, i)
}

func (s *ListInitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ListInitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ListInitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterListInit(s)
	}
}

func (s *ListInitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitListInit(s)
	}
}

func (s *ListInitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitListInit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) ListInit() (localctx IListInitContext) {
	localctx = NewListInitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, CommandsParserRULE_listInit)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(351)

		var _x = p.OptExpr()

		localctx.(*ListInitContext)._optExpr = _x
	}
	localctx.(*ListInitContext).elems = append(localctx.(*ListInitContext).elems, localctx.(*ListInitContext)._optExpr)
	p.SetState(356)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(352)
				p.Match(CommandsParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(353)

				var _x = p.OptExpr()

				localctx.(*ListInitContext)._optExpr = _x
			}
			localctx.(*ListInitContext).elems = append(localctx.(*ListInitContext).elems, localctx.(*ListInitContext)._optExpr)

		}
		p.SetState(358)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 43, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFieldInitializerListContext is an interface to support dynamic dispatch.
type IFieldInitializerListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS33 returns the s33 token.
	GetS33() antlr.Token

	// SetS33 sets the s33 token.
	SetS33(antlr.Token)

	// GetCols returns the cols token list.
	GetCols() []antlr.Token

	// SetCols sets the cols token list.
	SetCols([]antlr.Token)

	// Get_optField returns the _optField rule contexts.
	Get_optField() IOptFieldContext

	// Get_expr returns the _expr rule contexts.
	Get_expr() IExprContext

	// Set_optField sets the _optField rule contexts.
	Set_optField(IOptFieldContext)

	// Set_expr sets the _expr rule contexts.
	Set_expr(IExprContext)

	// GetFields returns the fields rule context list.
	GetFields() []IOptFieldContext

	// GetValues returns the values rule context list.
	GetValues() []IExprContext

	// SetFields sets the fields rule context list.
	SetFields([]IOptFieldContext)

	// SetValues sets the values rule context list.
	SetValues([]IExprContext)

	// Getter signatures
	AllOptField() []IOptFieldContext
	OptField(i int) IOptFieldContext
	AllCOLON() []antlr.TerminalNode
	COLON(i int) antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFieldInitializerListContext differentiates from other interfaces.
	IsFieldInitializerListContext()
}

type FieldInitializerListContext struct {
	antlr.BaseParserRuleContext
	parser    antlr.Parser
	_optField IOptFieldContext
	fields    []IOptFieldContext
	s33       antlr.Token
	cols      []antlr.Token
	_expr     IExprContext
	values    []IExprContext
}

func NewEmptyFieldInitializerListContext() *FieldInitializerListContext {
	var p = new(FieldInitializerListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_fieldInitializerList
	return p
}

func InitEmptyFieldInitializerListContext(p *FieldInitializerListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_fieldInitializerList
}

func (*FieldInitializerListContext) IsFieldInitializerListContext() {}

func NewFieldInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldInitializerListContext {
	var p = new(FieldInitializerListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_fieldInitializerList

	return p
}

func (s *FieldInitializerListContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldInitializerListContext) GetS33() antlr.Token { return s.s33 }

func (s *FieldInitializerListContext) SetS33(v antlr.Token) { s.s33 = v }

func (s *FieldInitializerListContext) GetCols() []antlr.Token { return s.cols }

func (s *FieldInitializerListContext) SetCols(v []antlr.Token) { s.cols = v }

func (s *FieldInitializerListContext) Get_optField() IOptFieldContext { return s._optField }

func (s *FieldInitializerListContext) Get_expr() IExprContext { return s._expr }

func (s *FieldInitializerListContext) Set_optField(v IOptFieldContext) { s._optField = v }

func (s *FieldInitializerListContext) Set_expr(v IExprContext) { s._expr = v }

func (s *FieldInitializerListContext) GetFields() []IOptFieldContext { return s.fields }

func (s *FieldInitializerListContext) GetValues() []IExprContext { return s.values }

func (s *FieldInitializerListContext) SetFields(v []IOptFieldContext) { s.fields = v }

func (s *FieldInitializerListContext) SetValues(v []IExprContext) { s.values = v }

func (s *FieldInitializerListContext) AllOptField() []IOptFieldContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOptFieldContext); ok {
			len++
		}
	}

	tst := make([]IOptFieldContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOptFieldContext); ok {
			tst[i] = t.(IOptFieldContext)
			i++
		}
	}

	return tst
}

func (s *FieldInitializerListContext) OptField(i int) IOptFieldContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOptFieldContext); ok {
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

	return t.(IOptFieldContext)
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

func (s *FieldInitializerListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitFieldInitializerList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) FieldInitializerList() (localctx IFieldInitializerListContext) {
	localctx = NewFieldInitializerListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, CommandsParserRULE_fieldInitializerList)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(359)

		var _x = p.OptField()

		localctx.(*FieldInitializerListContext)._optField = _x
	}
	localctx.(*FieldInitializerListContext).fields = append(localctx.(*FieldInitializerListContext).fields, localctx.(*FieldInitializerListContext)._optField)
	{
		p.SetState(360)

		var _m = p.Match(CommandsParserCOLON)

		localctx.(*FieldInitializerListContext).s33 = _m
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	localctx.(*FieldInitializerListContext).cols = append(localctx.(*FieldInitializerListContext).cols, localctx.(*FieldInitializerListContext).s33)
	{
		p.SetState(361)

		var _x = p.Expr()

		localctx.(*FieldInitializerListContext)._expr = _x
	}
	localctx.(*FieldInitializerListContext).values = append(localctx.(*FieldInitializerListContext).values, localctx.(*FieldInitializerListContext)._expr)
	p.SetState(369)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(362)
				p.Match(CommandsParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(363)

				var _x = p.OptField()

				localctx.(*FieldInitializerListContext)._optField = _x
			}
			localctx.(*FieldInitializerListContext).fields = append(localctx.(*FieldInitializerListContext).fields, localctx.(*FieldInitializerListContext)._optField)
			{
				p.SetState(364)

				var _m = p.Match(CommandsParserCOLON)

				localctx.(*FieldInitializerListContext).s33 = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*FieldInitializerListContext).cols = append(localctx.(*FieldInitializerListContext).cols, localctx.(*FieldInitializerListContext).s33)
			{
				p.SetState(365)

				var _x = p.Expr()

				localctx.(*FieldInitializerListContext)._expr = _x
			}
			localctx.(*FieldInitializerListContext).values = append(localctx.(*FieldInitializerListContext).values, localctx.(*FieldInitializerListContext)._expr)

		}
		p.SetState(371)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 44, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOptFieldContext is an interface to support dynamic dispatch.
type IOptFieldContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOpt returns the opt token.
	GetOpt() antlr.Token

	// SetOpt sets the opt token.
	SetOpt(antlr.Token)

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	QUESTIONMARK() antlr.TerminalNode

	// IsOptFieldContext differentiates from other interfaces.
	IsOptFieldContext()
}

type OptFieldContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	opt    antlr.Token
}

func NewEmptyOptFieldContext() *OptFieldContext {
	var p = new(OptFieldContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_optField
	return p
}

func InitEmptyOptFieldContext(p *OptFieldContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_optField
}

func (*OptFieldContext) IsOptFieldContext() {}

func NewOptFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptFieldContext {
	var p = new(OptFieldContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_optField

	return p
}

func (s *OptFieldContext) GetParser() antlr.Parser { return s.parser }

func (s *OptFieldContext) GetOpt() antlr.Token { return s.opt }

func (s *OptFieldContext) SetOpt(v antlr.Token) { s.opt = v }

func (s *OptFieldContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(CommandsParserIDENTIFIER, 0)
}

func (s *OptFieldContext) QUESTIONMARK() antlr.TerminalNode {
	return s.GetToken(CommandsParserQUESTIONMARK, 0)
}

func (s *OptFieldContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OptFieldContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OptFieldContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterOptField(s)
	}
}

func (s *OptFieldContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitOptField(s)
	}
}

func (s *OptFieldContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitOptField(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) OptField() (localctx IOptFieldContext) {
	localctx = NewOptFieldContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, CommandsParserRULE_optField)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(373)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserQUESTIONMARK {
		{
			p.SetState(372)

			var _m = p.Match(CommandsParserQUESTIONMARK)

			localctx.(*OptFieldContext).opt = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(375)
		p.Match(CommandsParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMapInitializerListContext is an interface to support dynamic dispatch.
type IMapInitializerListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS33 returns the s33 token.
	GetS33() antlr.Token

	// SetS33 sets the s33 token.
	SetS33(antlr.Token)

	// GetCols returns the cols token list.
	GetCols() []antlr.Token

	// SetCols sets the cols token list.
	SetCols([]antlr.Token)

	// Get_optExpr returns the _optExpr rule contexts.
	Get_optExpr() IOptExprContext

	// Get_expr returns the _expr rule contexts.
	Get_expr() IExprContext

	// Set_optExpr sets the _optExpr rule contexts.
	Set_optExpr(IOptExprContext)

	// Set_expr sets the _expr rule contexts.
	Set_expr(IExprContext)

	// GetKeys returns the keys rule context list.
	GetKeys() []IOptExprContext

	// GetValues returns the values rule context list.
	GetValues() []IExprContext

	// SetKeys sets the keys rule context list.
	SetKeys([]IOptExprContext)

	// SetValues sets the values rule context list.
	SetValues([]IExprContext)

	// Getter signatures
	AllOptExpr() []IOptExprContext
	OptExpr(i int) IOptExprContext
	AllCOLON() []antlr.TerminalNode
	COLON(i int) antlr.TerminalNode
	AllExpr() []IExprContext
	Expr(i int) IExprContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsMapInitializerListContext differentiates from other interfaces.
	IsMapInitializerListContext()
}

type MapInitializerListContext struct {
	antlr.BaseParserRuleContext
	parser   antlr.Parser
	_optExpr IOptExprContext
	keys     []IOptExprContext
	s33      antlr.Token
	cols     []antlr.Token
	_expr    IExprContext
	values   []IExprContext
}

func NewEmptyMapInitializerListContext() *MapInitializerListContext {
	var p = new(MapInitializerListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_mapInitializerList
	return p
}

func InitEmptyMapInitializerListContext(p *MapInitializerListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_mapInitializerList
}

func (*MapInitializerListContext) IsMapInitializerListContext() {}

func NewMapInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapInitializerListContext {
	var p = new(MapInitializerListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_mapInitializerList

	return p
}

func (s *MapInitializerListContext) GetParser() antlr.Parser { return s.parser }

func (s *MapInitializerListContext) GetS33() antlr.Token { return s.s33 }

func (s *MapInitializerListContext) SetS33(v antlr.Token) { s.s33 = v }

func (s *MapInitializerListContext) GetCols() []antlr.Token { return s.cols }

func (s *MapInitializerListContext) SetCols(v []antlr.Token) { s.cols = v }

func (s *MapInitializerListContext) Get_optExpr() IOptExprContext { return s._optExpr }

func (s *MapInitializerListContext) Get_expr() IExprContext { return s._expr }

func (s *MapInitializerListContext) Set_optExpr(v IOptExprContext) { s._optExpr = v }

func (s *MapInitializerListContext) Set_expr(v IExprContext) { s._expr = v }

func (s *MapInitializerListContext) GetKeys() []IOptExprContext { return s.keys }

func (s *MapInitializerListContext) GetValues() []IExprContext { return s.values }

func (s *MapInitializerListContext) SetKeys(v []IOptExprContext) { s.keys = v }

func (s *MapInitializerListContext) SetValues(v []IExprContext) { s.values = v }

func (s *MapInitializerListContext) AllOptExpr() []IOptExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOptExprContext); ok {
			len++
		}
	}

	tst := make([]IOptExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOptExprContext); ok {
			tst[i] = t.(IOptExprContext)
			i++
		}
	}

	return tst
}

func (s *MapInitializerListContext) OptExpr(i int) IOptExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOptExprContext); ok {
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

	return t.(IOptExprContext)
}

func (s *MapInitializerListContext) AllCOLON() []antlr.TerminalNode {
	return s.GetTokens(CommandsParserCOLON)
}

func (s *MapInitializerListContext) COLON(i int) antlr.TerminalNode {
	return s.GetToken(CommandsParserCOLON, i)
}

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

func (s *MapInitializerListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitMapInitializerList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) MapInitializerList() (localctx IMapInitializerListContext) {
	localctx = NewMapInitializerListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, CommandsParserRULE_mapInitializerList)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(377)

		var _x = p.OptExpr()

		localctx.(*MapInitializerListContext)._optExpr = _x
	}
	localctx.(*MapInitializerListContext).keys = append(localctx.(*MapInitializerListContext).keys, localctx.(*MapInitializerListContext)._optExpr)
	{
		p.SetState(378)

		var _m = p.Match(CommandsParserCOLON)

		localctx.(*MapInitializerListContext).s33 = _m
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	localctx.(*MapInitializerListContext).cols = append(localctx.(*MapInitializerListContext).cols, localctx.(*MapInitializerListContext).s33)
	{
		p.SetState(379)

		var _x = p.Expr()

		localctx.(*MapInitializerListContext)._expr = _x
	}
	localctx.(*MapInitializerListContext).values = append(localctx.(*MapInitializerListContext).values, localctx.(*MapInitializerListContext)._expr)
	p.SetState(387)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(380)
				p.Match(CommandsParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(381)

				var _x = p.OptExpr()

				localctx.(*MapInitializerListContext)._optExpr = _x
			}
			localctx.(*MapInitializerListContext).keys = append(localctx.(*MapInitializerListContext).keys, localctx.(*MapInitializerListContext)._optExpr)
			{
				p.SetState(382)

				var _m = p.Match(CommandsParserCOLON)

				localctx.(*MapInitializerListContext).s33 = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			localctx.(*MapInitializerListContext).cols = append(localctx.(*MapInitializerListContext).cols, localctx.(*MapInitializerListContext).s33)
			{
				p.SetState(383)

				var _x = p.Expr()

				localctx.(*MapInitializerListContext)._expr = _x
			}
			localctx.(*MapInitializerListContext).values = append(localctx.(*MapInitializerListContext).values, localctx.(*MapInitializerListContext)._expr)

		}
		p.SetState(389)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 46, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOptExprContext is an interface to support dynamic dispatch.
type IOptExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetOpt returns the opt token.
	GetOpt() antlr.Token

	// SetOpt sets the opt token.
	SetOpt(antlr.Token)

	// GetE returns the e rule contexts.
	GetE() IExprContext

	// SetE sets the e rule contexts.
	SetE(IExprContext)

	// Getter signatures
	Expr() IExprContext
	QUESTIONMARK() antlr.TerminalNode

	// IsOptExprContext differentiates from other interfaces.
	IsOptExprContext()
}

type OptExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
	opt    antlr.Token
	e      IExprContext
}

func NewEmptyOptExprContext() *OptExprContext {
	var p = new(OptExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_optExpr
	return p
}

func InitEmptyOptExprContext(p *OptExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_optExpr
}

func (*OptExprContext) IsOptExprContext() {}

func NewOptExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptExprContext {
	var p = new(OptExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_optExpr

	return p
}

func (s *OptExprContext) GetParser() antlr.Parser { return s.parser }

func (s *OptExprContext) GetOpt() antlr.Token { return s.opt }

func (s *OptExprContext) SetOpt(v antlr.Token) { s.opt = v }

func (s *OptExprContext) GetE() IExprContext { return s.e }

func (s *OptExprContext) SetE(v IExprContext) { s.e = v }

func (s *OptExprContext) Expr() IExprContext {
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

func (s *OptExprContext) QUESTIONMARK() antlr.TerminalNode {
	return s.GetToken(CommandsParserQUESTIONMARK, 0)
}

func (s *OptExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OptExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OptExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.EnterOptExpr(s)
	}
}

func (s *OptExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CommandsListener); ok {
		listenerT.ExitOptExpr(s)
	}
}

func (s *OptExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitOptExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) OptExpr() (localctx IOptExprContext) {
	localctx = NewOptExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, CommandsParserRULE_optExpr)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(391)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == CommandsParserQUESTIONMARK {
		{
			p.SetState(390)

			var _m = p.Match(CommandsParserQUESTIONMARK)

			localctx.(*OptExprContext).opt = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(393)

		var _x = p.Expr()

		localctx.(*OptExprContext).e = _x
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
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
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_literal
	return p
}

func InitEmptyLiteralContext(p *LiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = CommandsParserRULE_literal
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = CommandsParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) CopyAll(ctx *LiteralContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type BytesContext struct {
	LiteralContext
	tok antlr.Token
}

func NewBytesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BytesContext {
	var p = new(BytesContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

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

func (s *BytesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitBytes(s)

	default:
		return t.VisitChildren(s)
	}
}

type UintContext struct {
	LiteralContext
	tok antlr.Token
}

func NewUintContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UintContext {
	var p = new(UintContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

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

func (s *UintContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitUint(s)

	default:
		return t.VisitChildren(s)
	}
}

type NullContext struct {
	LiteralContext
	tok antlr.Token
}

func NewNullContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullContext {
	var p = new(NullContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

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

func (s *NullContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitNull(s)

	default:
		return t.VisitChildren(s)
	}
}

type BoolFalseContext struct {
	LiteralContext
	tok antlr.Token
}

func NewBoolFalseContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoolFalseContext {
	var p = new(BoolFalseContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

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

func (s *BoolFalseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitBoolFalse(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringContext struct {
	LiteralContext
	tok antlr.Token
}

func NewStringContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringContext {
	var p = new(StringContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

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

func (s *StringContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitString(s)

	default:
		return t.VisitChildren(s)
	}
}

type DoubleContext struct {
	LiteralContext
	sign antlr.Token
	tok  antlr.Token
}

func NewDoubleContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DoubleContext {
	var p = new(DoubleContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

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

func (s *DoubleContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitDouble(s)

	default:
		return t.VisitChildren(s)
	}
}

type BoolTrueContext struct {
	LiteralContext
	tok antlr.Token
}

func NewBoolTrueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoolTrueContext {
	var p = new(BoolTrueContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

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

func (s *BoolTrueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitBoolTrue(s)

	default:
		return t.VisitChildren(s)
	}
}

type IntContext struct {
	LiteralContext
	sign antlr.Token
	tok  antlr.Token
}

func NewIntContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntContext {
	var p = new(IntContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

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

func (s *IntContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case CommandsVisitor:
		return t.VisitInt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *CommandsParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, CommandsParserRULE_literal)
	var _la int

	p.SetState(409)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) {
	case 1:
		localctx = NewIntContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(396)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserMINUS {
			{
				p.SetState(395)

				var _m = p.Match(CommandsParserMINUS)

				localctx.(*IntContext).sign = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(398)

			var _m = p.Match(CommandsParserNUM_INT)

			localctx.(*IntContext).tok = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewUintContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(399)

			var _m = p.Match(CommandsParserNUM_UINT)

			localctx.(*UintContext).tok = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		localctx = NewDoubleContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		p.SetState(401)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == CommandsParserMINUS {
			{
				p.SetState(400)

				var _m = p.Match(CommandsParserMINUS)

				localctx.(*DoubleContext).sign = _m
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(403)

			var _m = p.Match(CommandsParserNUM_FLOAT)

			localctx.(*DoubleContext).tok = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewStringContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(404)

			var _m = p.Match(CommandsParserSTRING)

			localctx.(*StringContext).tok = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewBytesContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(405)

			var _m = p.Match(CommandsParserBYTES)

			localctx.(*BytesContext).tok = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		localctx = NewBoolTrueContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(406)

			var _m = p.Match(CommandsParserCEL_TRUE)

			localctx.(*BoolTrueContext).tok = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		localctx = NewBoolFalseContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(407)

			var _m = p.Match(CommandsParserCEL_FALSE)

			localctx.(*BoolFalseContext).tok = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		localctx = NewNullContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(408)

			var _m = p.Match(CommandsParserNUL)

			localctx.(*NullContext).tok = _m
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *CommandsParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 23:
		var t *RelationContext = nil
		if localctx != nil {
			t = localctx.(*RelationContext)
		}
		return p.Relation_Sempred(t, predIndex)

	case 24:
		var t *CalcContext = nil
		if localctx != nil {
			t = localctx.(*CalcContext)
		}
		return p.Calc_Sempred(t, predIndex)

	case 26:
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
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *CommandsParser) Calc_Sempred(localctx antlr.RuleContext, predIndex int) bool {
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
