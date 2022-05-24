// Code generated from ./Types.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Types
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

type TypesParser struct {
	*antlr.BaseParser
}

var typesParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func typesParserInit() {
	staticData := &typesParserStaticData
	staticData.literalNames = []string{
		"", "'=='", "'!='", "'in'", "'<'", "'<='", "'>='", "'>'", "'&&'", "'||'",
		"'['", "']'", "'{'", "'}'", "'('", "')'", "'.'", "','", "'-'", "'!'",
		"'?'", "':'", "'+'", "'*'", "'/'", "'%'", "'true'", "'false'", "'null'",
	}
	staticData.symbolicNames = []string{
		"", "EQUALS", "NOT_EQUALS", "IN", "LESS", "LESS_EQUALS", "GREATER_EQUALS",
		"GREATER", "LOGICAL_AND", "LOGICAL_OR", "LBRACKET", "RPRACKET", "LBRACE",
		"RBRACE", "LPAREN", "RPAREN", "DOT", "COMMA", "MINUS", "EXCLAM", "QUESTIONMARK",
		"COLON", "PLUS", "STAR", "SLASH", "PERCENT", "CEL_TRUE", "CEL_FALSE",
		"NUL", "WHITESPACE", "COMMENT", "NUM_FLOAT", "NUM_INT", "NUM_UINT",
		"STRING", "BYTES", "IDENTIFIER",
	}
	staticData.ruleNames = []string{
		"start", "type", "typeId", "typeParamList", "expr", "conditionalOr",
		"conditionalAnd", "relation", "calc", "unary", "member", "primary",
		"exprList", "fieldInitializerList", "mapInitializerList", "literal",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 36, 241, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 3, 1, 38, 8, 1, 1, 2, 3, 2, 41, 8, 2, 1,
		2, 1, 2, 1, 2, 5, 2, 46, 8, 2, 10, 2, 12, 2, 49, 9, 2, 1, 3, 1, 3, 1, 3,
		1, 3, 5, 3, 55, 8, 3, 10, 3, 12, 3, 58, 9, 3, 1, 3, 1, 3, 1, 4, 1, 4, 1,
		4, 1, 4, 1, 4, 1, 4, 3, 4, 68, 8, 4, 1, 5, 1, 5, 1, 5, 5, 5, 73, 8, 5,
		10, 5, 12, 5, 76, 9, 5, 1, 6, 1, 6, 1, 6, 5, 6, 81, 8, 6, 10, 6, 12, 6,
		84, 9, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 5, 7, 92, 8, 7, 10, 7, 12,
		7, 95, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 5, 8,
		106, 8, 8, 10, 8, 12, 8, 109, 9, 8, 1, 9, 1, 9, 4, 9, 113, 8, 9, 11, 9,
		12, 9, 114, 1, 9, 1, 9, 4, 9, 119, 8, 9, 11, 9, 12, 9, 120, 1, 9, 3, 9,
		124, 8, 9, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10, 3, 10,
		134, 8, 10, 1, 10, 3, 10, 137, 8, 10, 1, 10, 1, 10, 1, 10, 1, 10, 1, 10,
		1, 10, 1, 10, 1, 10, 3, 10, 147, 8, 10, 1, 10, 3, 10, 150, 8, 10, 1, 10,
		5, 10, 153, 8, 10, 10, 10, 12, 10, 156, 9, 10, 1, 11, 3, 11, 159, 8, 11,
		1, 11, 1, 11, 1, 11, 3, 11, 164, 8, 11, 1, 11, 3, 11, 167, 8, 11, 1, 11,
		1, 11, 1, 11, 1, 11, 1, 11, 1, 11, 3, 11, 175, 8, 11, 1, 11, 3, 11, 178,
		8, 11, 1, 11, 1, 11, 1, 11, 3, 11, 183, 8, 11, 1, 11, 3, 11, 186, 8, 11,
		1, 11, 1, 11, 3, 11, 190, 8, 11, 1, 12, 1, 12, 1, 12, 5, 12, 195, 8, 12,
		10, 12, 12, 12, 198, 9, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1,
		13, 5, 13, 207, 8, 13, 10, 13, 12, 13, 210, 9, 13, 1, 14, 1, 14, 1, 14,
		1, 14, 1, 14, 1, 14, 1, 14, 1, 14, 5, 14, 220, 8, 14, 10, 14, 12, 14, 223,
		9, 14, 1, 15, 3, 15, 226, 8, 15, 1, 15, 1, 15, 1, 15, 3, 15, 231, 8, 15,
		1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 1, 15, 3, 15, 239, 8, 15, 1, 15, 0,
		3, 14, 16, 20, 16, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28,
		30, 0, 4, 2, 0, 28, 28, 36, 36, 1, 0, 1, 7, 1, 0, 23, 25, 2, 0, 18, 18,
		22, 22, 268, 0, 32, 1, 0, 0, 0, 2, 35, 1, 0, 0, 0, 4, 40, 1, 0, 0, 0, 6,
		50, 1, 0, 0, 0, 8, 61, 1, 0, 0, 0, 10, 69, 1, 0, 0, 0, 12, 77, 1, 0, 0,
		0, 14, 85, 1, 0, 0, 0, 16, 96, 1, 0, 0, 0, 18, 123, 1, 0, 0, 0, 20, 125,
		1, 0, 0, 0, 22, 189, 1, 0, 0, 0, 24, 191, 1, 0, 0, 0, 26, 199, 1, 0, 0,
		0, 28, 211, 1, 0, 0, 0, 30, 238, 1, 0, 0, 0, 32, 33, 3, 2, 1, 0, 33, 34,
		5, 0, 0, 1, 34, 1, 1, 0, 0, 0, 35, 37, 3, 4, 2, 0, 36, 38, 3, 6, 3, 0,
		37, 36, 1, 0, 0, 0, 37, 38, 1, 0, 0, 0, 38, 3, 1, 0, 0, 0, 39, 41, 5, 16,
		0, 0, 40, 39, 1, 0, 0, 0, 40, 41, 1, 0, 0, 0, 41, 42, 1, 0, 0, 0, 42, 47,
		7, 0, 0, 0, 43, 44, 5, 16, 0, 0, 44, 46, 5, 36, 0, 0, 45, 43, 1, 0, 0,
		0, 46, 49, 1, 0, 0, 0, 47, 45, 1, 0, 0, 0, 47, 48, 1, 0, 0, 0, 48, 5, 1,
		0, 0, 0, 49, 47, 1, 0, 0, 0, 50, 51, 5, 14, 0, 0, 51, 56, 3, 2, 1, 0, 52,
		53, 5, 17, 0, 0, 53, 55, 3, 2, 1, 0, 54, 52, 1, 0, 0, 0, 55, 58, 1, 0,
		0, 0, 56, 54, 1, 0, 0, 0, 56, 57, 1, 0, 0, 0, 57, 59, 1, 0, 0, 0, 58, 56,
		1, 0, 0, 0, 59, 60, 5, 15, 0, 0, 60, 7, 1, 0, 0, 0, 61, 67, 3, 10, 5, 0,
		62, 63, 5, 20, 0, 0, 63, 64, 3, 10, 5, 0, 64, 65, 5, 21, 0, 0, 65, 66,
		3, 8, 4, 0, 66, 68, 1, 0, 0, 0, 67, 62, 1, 0, 0, 0, 67, 68, 1, 0, 0, 0,
		68, 9, 1, 0, 0, 0, 69, 74, 3, 12, 6, 0, 70, 71, 5, 9, 0, 0, 71, 73, 3,
		12, 6, 0, 72, 70, 1, 0, 0, 0, 73, 76, 1, 0, 0, 0, 74, 72, 1, 0, 0, 0, 74,
		75, 1, 0, 0, 0, 75, 11, 1, 0, 0, 0, 76, 74, 1, 0, 0, 0, 77, 82, 3, 14,
		7, 0, 78, 79, 5, 8, 0, 0, 79, 81, 3, 14, 7, 0, 80, 78, 1, 0, 0, 0, 81,
		84, 1, 0, 0, 0, 82, 80, 1, 0, 0, 0, 82, 83, 1, 0, 0, 0, 83, 13, 1, 0, 0,
		0, 84, 82, 1, 0, 0, 0, 85, 86, 6, 7, -1, 0, 86, 87, 3, 16, 8, 0, 87, 93,
		1, 0, 0, 0, 88, 89, 10, 1, 0, 0, 89, 90, 7, 1, 0, 0, 90, 92, 3, 14, 7,
		2, 91, 88, 1, 0, 0, 0, 92, 95, 1, 0, 0, 0, 93, 91, 1, 0, 0, 0, 93, 94,
		1, 0, 0, 0, 94, 15, 1, 0, 0, 0, 95, 93, 1, 0, 0, 0, 96, 97, 6, 8, -1, 0,
		97, 98, 3, 18, 9, 0, 98, 107, 1, 0, 0, 0, 99, 100, 10, 2, 0, 0, 100, 101,
		7, 2, 0, 0, 101, 106, 3, 16, 8, 3, 102, 103, 10, 1, 0, 0, 103, 104, 7,
		3, 0, 0, 104, 106, 3, 16, 8, 2, 105, 99, 1, 0, 0, 0, 105, 102, 1, 0, 0,
		0, 106, 109, 1, 0, 0, 0, 107, 105, 1, 0, 0, 0, 107, 108, 1, 0, 0, 0, 108,
		17, 1, 0, 0, 0, 109, 107, 1, 0, 0, 0, 110, 124, 3, 20, 10, 0, 111, 113,
		5, 19, 0, 0, 112, 111, 1, 0, 0, 0, 113, 114, 1, 0, 0, 0, 114, 112, 1, 0,
		0, 0, 114, 115, 1, 0, 0, 0, 115, 116, 1, 0, 0, 0, 116, 124, 3, 20, 10,
		0, 117, 119, 5, 18, 0, 0, 118, 117, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120,
		118, 1, 0, 0, 0, 120, 121, 1, 0, 0, 0, 121, 122, 1, 0, 0, 0, 122, 124,
		3, 20, 10, 0, 123, 110, 1, 0, 0, 0, 123, 112, 1, 0, 0, 0, 123, 118, 1,
		0, 0, 0, 124, 19, 1, 0, 0, 0, 125, 126, 6, 10, -1, 0, 126, 127, 3, 22,
		11, 0, 127, 154, 1, 0, 0, 0, 128, 129, 10, 3, 0, 0, 129, 130, 5, 16, 0,
		0, 130, 136, 5, 36, 0, 0, 131, 133, 5, 14, 0, 0, 132, 134, 3, 24, 12, 0,
		133, 132, 1, 0, 0, 0, 133, 134, 1, 0, 0, 0, 134, 135, 1, 0, 0, 0, 135,
		137, 5, 15, 0, 0, 136, 131, 1, 0, 0, 0, 136, 137, 1, 0, 0, 0, 137, 153,
		1, 0, 0, 0, 138, 139, 10, 2, 0, 0, 139, 140, 5, 10, 0, 0, 140, 141, 3,
		8, 4, 0, 141, 142, 5, 11, 0, 0, 142, 153, 1, 0, 0, 0, 143, 144, 10, 1,
		0, 0, 144, 146, 5, 12, 0, 0, 145, 147, 3, 26, 13, 0, 146, 145, 1, 0, 0,
		0, 146, 147, 1, 0, 0, 0, 147, 149, 1, 0, 0, 0, 148, 150, 5, 17, 0, 0, 149,
		148, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150, 151, 1, 0, 0, 0, 151, 153,
		5, 13, 0, 0, 152, 128, 1, 0, 0, 0, 152, 138, 1, 0, 0, 0, 152, 143, 1, 0,
		0, 0, 153, 156, 1, 0, 0, 0, 154, 152, 1, 0, 0, 0, 154, 155, 1, 0, 0, 0,
		155, 21, 1, 0, 0, 0, 156, 154, 1, 0, 0, 0, 157, 159, 5, 16, 0, 0, 158,
		157, 1, 0, 0, 0, 158, 159, 1, 0, 0, 0, 159, 160, 1, 0, 0, 0, 160, 166,
		5, 36, 0, 0, 161, 163, 5, 14, 0, 0, 162, 164, 3, 24, 12, 0, 163, 162, 1,
		0, 0, 0, 163, 164, 1, 0, 0, 0, 164, 165, 1, 0, 0, 0, 165, 167, 5, 15, 0,
		0, 166, 161, 1, 0, 0, 0, 166, 167, 1, 0, 0, 0, 167, 190, 1, 0, 0, 0, 168,
		169, 5, 14, 0, 0, 169, 170, 3, 8, 4, 0, 170, 171, 5, 15, 0, 0, 171, 190,
		1, 0, 0, 0, 172, 174, 5, 10, 0, 0, 173, 175, 3, 24, 12, 0, 174, 173, 1,
		0, 0, 0, 174, 175, 1, 0, 0, 0, 175, 177, 1, 0, 0, 0, 176, 178, 5, 17, 0,
		0, 177, 176, 1, 0, 0, 0, 177, 178, 1, 0, 0, 0, 178, 179, 1, 0, 0, 0, 179,
		190, 5, 11, 0, 0, 180, 182, 5, 12, 0, 0, 181, 183, 3, 28, 14, 0, 182, 181,
		1, 0, 0, 0, 182, 183, 1, 0, 0, 0, 183, 185, 1, 0, 0, 0, 184, 186, 5, 17,
		0, 0, 185, 184, 1, 0, 0, 0, 185, 186, 1, 0, 0, 0, 186, 187, 1, 0, 0, 0,
		187, 190, 5, 13, 0, 0, 188, 190, 3, 30, 15, 0, 189, 158, 1, 0, 0, 0, 189,
		168, 1, 0, 0, 0, 189, 172, 1, 0, 0, 0, 189, 180, 1, 0, 0, 0, 189, 188,
		1, 0, 0, 0, 190, 23, 1, 0, 0, 0, 191, 196, 3, 8, 4, 0, 192, 193, 5, 17,
		0, 0, 193, 195, 3, 8, 4, 0, 194, 192, 1, 0, 0, 0, 195, 198, 1, 0, 0, 0,
		196, 194, 1, 0, 0, 0, 196, 197, 1, 0, 0, 0, 197, 25, 1, 0, 0, 0, 198, 196,
		1, 0, 0, 0, 199, 200, 5, 36, 0, 0, 200, 201, 5, 21, 0, 0, 201, 208, 3,
		8, 4, 0, 202, 203, 5, 17, 0, 0, 203, 204, 5, 36, 0, 0, 204, 205, 5, 21,
		0, 0, 205, 207, 3, 8, 4, 0, 206, 202, 1, 0, 0, 0, 207, 210, 1, 0, 0, 0,
		208, 206, 1, 0, 0, 0, 208, 209, 1, 0, 0, 0, 209, 27, 1, 0, 0, 0, 210, 208,
		1, 0, 0, 0, 211, 212, 3, 8, 4, 0, 212, 213, 5, 21, 0, 0, 213, 221, 3, 8,
		4, 0, 214, 215, 5, 17, 0, 0, 215, 216, 3, 8, 4, 0, 216, 217, 5, 21, 0,
		0, 217, 218, 3, 8, 4, 0, 218, 220, 1, 0, 0, 0, 219, 214, 1, 0, 0, 0, 220,
		223, 1, 0, 0, 0, 221, 219, 1, 0, 0, 0, 221, 222, 1, 0, 0, 0, 222, 29, 1,
		0, 0, 0, 223, 221, 1, 0, 0, 0, 224, 226, 5, 18, 0, 0, 225, 224, 1, 0, 0,
		0, 225, 226, 1, 0, 0, 0, 226, 227, 1, 0, 0, 0, 227, 239, 5, 32, 0, 0, 228,
		239, 5, 33, 0, 0, 229, 231, 5, 18, 0, 0, 230, 229, 1, 0, 0, 0, 230, 231,
		1, 0, 0, 0, 231, 232, 1, 0, 0, 0, 232, 239, 5, 31, 0, 0, 233, 239, 5, 34,
		0, 0, 234, 239, 5, 35, 0, 0, 235, 239, 5, 26, 0, 0, 236, 239, 5, 27, 0,
		0, 237, 239, 5, 28, 0, 0, 238, 225, 1, 0, 0, 0, 238, 228, 1, 0, 0, 0, 238,
		230, 1, 0, 0, 0, 238, 233, 1, 0, 0, 0, 238, 234, 1, 0, 0, 0, 238, 235,
		1, 0, 0, 0, 238, 236, 1, 0, 0, 0, 238, 237, 1, 0, 0, 0, 239, 31, 1, 0,
		0, 0, 33, 37, 40, 47, 56, 67, 74, 82, 93, 105, 107, 114, 120, 123, 133,
		136, 146, 149, 152, 154, 158, 163, 166, 174, 177, 182, 185, 189, 196, 208,
		221, 225, 230, 238,
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

// TypesParserInit initializes any static state used to implement TypesParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewTypesParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func TypesParserInit() {
	staticData := &typesParserStaticData
	staticData.once.Do(typesParserInit)
}

// NewTypesParser produces a new parser instance for the optional input antlr.TokenStream.
func NewTypesParser(input antlr.TokenStream) *TypesParser {
	TypesParserInit()
	this := new(TypesParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &typesParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "Types.g4"

	return this
}

// TypesParser tokens.
const (
	TypesParserEOF            = antlr.TokenEOF
	TypesParserEQUALS         = 1
	TypesParserNOT_EQUALS     = 2
	TypesParserIN             = 3
	TypesParserLESS           = 4
	TypesParserLESS_EQUALS    = 5
	TypesParserGREATER_EQUALS = 6
	TypesParserGREATER        = 7
	TypesParserLOGICAL_AND    = 8
	TypesParserLOGICAL_OR     = 9
	TypesParserLBRACKET       = 10
	TypesParserRPRACKET       = 11
	TypesParserLBRACE         = 12
	TypesParserRBRACE         = 13
	TypesParserLPAREN         = 14
	TypesParserRPAREN         = 15
	TypesParserDOT            = 16
	TypesParserCOMMA          = 17
	TypesParserMINUS          = 18
	TypesParserEXCLAM         = 19
	TypesParserQUESTIONMARK   = 20
	TypesParserCOLON          = 21
	TypesParserPLUS           = 22
	TypesParserSTAR           = 23
	TypesParserSLASH          = 24
	TypesParserPERCENT        = 25
	TypesParserCEL_TRUE       = 26
	TypesParserCEL_FALSE      = 27
	TypesParserNUL            = 28
	TypesParserWHITESPACE     = 29
	TypesParserCOMMENT        = 30
	TypesParserNUM_FLOAT      = 31
	TypesParserNUM_INT        = 32
	TypesParserNUM_UINT       = 33
	TypesParserSTRING         = 34
	TypesParserBYTES          = 35
	TypesParserIDENTIFIER     = 36
)

// TypesParser rules.
const (
	TypesParserRULE_start                = 0
	TypesParserRULE_type                 = 1
	TypesParserRULE_typeId               = 2
	TypesParserRULE_typeParamList        = 3
	TypesParserRULE_expr                 = 4
	TypesParserRULE_conditionalOr        = 5
	TypesParserRULE_conditionalAnd       = 6
	TypesParserRULE_relation             = 7
	TypesParserRULE_calc                 = 8
	TypesParserRULE_unary                = 9
	TypesParserRULE_member               = 10
	TypesParserRULE_primary              = 11
	TypesParserRULE_exprList             = 12
	TypesParserRULE_fieldInitializerList = 13
	TypesParserRULE_mapInitializerList   = 14
	TypesParserRULE_literal              = 15
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetT returns the t rule contexts.
	GetT() ITypeContext

	// SetT sets the t rule contexts.
	SetT(ITypeContext)

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
	t      ITypeContext
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = TypesParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) GetT() ITypeContext { return s.t }

func (s *StartContext) SetT(v ITypeContext) { s.t = v }

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(TypesParserEOF, 0)
}

func (s *StartContext) Type() ITypeContext {
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

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitStart(s)
	}
}

func (s *StartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitStart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Start() (localctx IStartContext) {
	this := p
	_ = this

	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, TypesParserRULE_start)

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
		p.SetState(32)

		var _x = p.Type()

		localctx.(*StartContext).t = _x
	}
	{
		p.SetState(33)
		p.Match(TypesParserEOF)
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
	p.RuleIndex = TypesParserRULE_type
	return p
}

func (*TypeContext) IsTypeContext() {}

func NewTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeContext {
	var p = new(TypeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_type

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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterType(s)
	}
}

func (s *TypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitType(s)
	}
}

func (s *TypeContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitType(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Type() (localctx ITypeContext) {
	this := p
	_ = this

	localctx = NewTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, TypesParserRULE_type)
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
		p.SetState(35)

		var _x = p.TypeId()

		localctx.(*TypeContext).id = _x
	}
	p.SetState(37)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == TypesParserLPAREN {
		{
			p.SetState(36)

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
	p.RuleIndex = TypesParserRULE_typeId
	return p
}

func (*TypeIdContext) IsTypeIdContext() {}

func NewTypeIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeIdContext {
	var p = new(TypeIdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_typeId

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
	return s.GetTokens(TypesParserIDENTIFIER)
}

func (s *TypeIdContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserIDENTIFIER, i)
}

func (s *TypeIdContext) NUL() antlr.TerminalNode {
	return s.GetToken(TypesParserNUL, 0)
}

func (s *TypeIdContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(TypesParserDOT)
}

func (s *TypeIdContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserDOT, i)
}

func (s *TypeIdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeIdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeIdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterTypeId(s)
	}
}

func (s *TypeIdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitTypeId(s)
	}
}

func (s *TypeIdContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitTypeId(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) TypeId() (localctx ITypeIdContext) {
	this := p
	_ = this

	localctx = NewTypeIdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, TypesParserRULE_typeId)
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
	p.SetState(40)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == TypesParserDOT {
		{
			p.SetState(39)

			var _m = p.Match(TypesParserDOT)

			localctx.(*TypeIdContext).leadingDot = _m
		}

	}
	{
		p.SetState(42)

		var _lt = p.GetTokenStream().LT(1)

		localctx.(*TypeIdContext).id = _lt

		_la = p.GetTokenStream().LA(1)

		if !(_la == TypesParserNUL || _la == TypesParserIDENTIFIER) {
			var _ri = p.GetErrorHandler().RecoverInline(p)

			localctx.(*TypeIdContext).id = _ri
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(47)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == TypesParserDOT {
		{
			p.SetState(43)
			p.Match(TypesParserDOT)
		}
		{
			p.SetState(44)

			var _m = p.Match(TypesParserIDENTIFIER)

			localctx.(*TypeIdContext)._IDENTIFIER = _m
		}
		localctx.(*TypeIdContext).qualifiers = append(localctx.(*TypeIdContext).qualifiers, localctx.(*TypeIdContext)._IDENTIFIER)

		p.SetState(49)
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
	p.RuleIndex = TypesParserRULE_typeParamList
	return p
}

func (*TypeParamListContext) IsTypeParamListContext() {}

func NewTypeParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamListContext {
	var p = new(TypeParamListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_typeParamList

	return p
}

func (s *TypeParamListContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeParamListContext) Get_type() ITypeContext { return s._type }

func (s *TypeParamListContext) Set_type(v ITypeContext) { s._type = v }

func (s *TypeParamListContext) GetTypes() []ITypeContext { return s.types }

func (s *TypeParamListContext) SetTypes(v []ITypeContext) { s.types = v }

func (s *TypeParamListContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(TypesParserLPAREN, 0)
}

func (s *TypeParamListContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(TypesParserRPAREN, 0)
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
	return s.GetTokens(TypesParserCOMMA)
}

func (s *TypeParamListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserCOMMA, i)
}

func (s *TypeParamListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeParamListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterTypeParamList(s)
	}
}

func (s *TypeParamListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitTypeParamList(s)
	}
}

func (s *TypeParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitTypeParamList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) TypeParamList() (localctx ITypeParamListContext) {
	this := p
	_ = this

	localctx = NewTypeParamListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, TypesParserRULE_typeParamList)
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
		p.SetState(50)
		p.Match(TypesParserLPAREN)
	}
	{
		p.SetState(51)

		var _x = p.Type()

		localctx.(*TypeParamListContext)._type = _x
	}
	localctx.(*TypeParamListContext).types = append(localctx.(*TypeParamListContext).types, localctx.(*TypeParamListContext)._type)
	p.SetState(56)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == TypesParserCOMMA {
		{
			p.SetState(52)
			p.Match(TypesParserCOMMA)
		}
		{
			p.SetState(53)

			var _x = p.Type()

			localctx.(*TypeParamListContext)._type = _x
		}
		localctx.(*TypeParamListContext).types = append(localctx.(*TypeParamListContext).types, localctx.(*TypeParamListContext)._type)

		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(59)
		p.Match(TypesParserRPAREN)
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
	p.RuleIndex = TypesParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_expr

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
	return s.GetToken(TypesParserCOLON, 0)
}

func (s *ExprContext) QUESTIONMARK() antlr.TerminalNode {
	return s.GetToken(TypesParserQUESTIONMARK, 0)
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitExpr(s)
	}
}

func (s *ExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Expr() (localctx IExprContext) {
	this := p
	_ = this

	localctx = NewExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, TypesParserRULE_expr)
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
		p.SetState(61)

		var _x = p.ConditionalOr()

		localctx.(*ExprContext).e = _x
	}
	p.SetState(67)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == TypesParserQUESTIONMARK {
		{
			p.SetState(62)

			var _m = p.Match(TypesParserQUESTIONMARK)

			localctx.(*ExprContext).op = _m
		}
		{
			p.SetState(63)

			var _x = p.ConditionalOr()

			localctx.(*ExprContext).e1 = _x
		}
		{
			p.SetState(64)
			p.Match(TypesParserCOLON)
		}
		{
			p.SetState(65)

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

	// GetS9 returns the s9 token.
	GetS9() antlr.Token

	// SetS9 sets the s9 token.
	SetS9(antlr.Token)

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
	s9              antlr.Token
	ops             []antlr.Token
	_conditionalAnd IConditionalAndContext
	e1              []IConditionalAndContext
}

func NewEmptyConditionalOrContext() *ConditionalOrContext {
	var p = new(ConditionalOrContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = TypesParserRULE_conditionalOr
	return p
}

func (*ConditionalOrContext) IsConditionalOrContext() {}

func NewConditionalOrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalOrContext {
	var p = new(ConditionalOrContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_conditionalOr

	return p
}

func (s *ConditionalOrContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionalOrContext) GetS9() antlr.Token { return s.s9 }

func (s *ConditionalOrContext) SetS9(v antlr.Token) { s.s9 = v }

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
	return s.GetTokens(TypesParserLOGICAL_OR)
}

func (s *ConditionalOrContext) LOGICAL_OR(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserLOGICAL_OR, i)
}

func (s *ConditionalOrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionalOrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConditionalOrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterConditionalOr(s)
	}
}

func (s *ConditionalOrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitConditionalOr(s)
	}
}

func (s *ConditionalOrContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitConditionalOr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) ConditionalOr() (localctx IConditionalOrContext) {
	this := p
	_ = this

	localctx = NewConditionalOrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, TypesParserRULE_conditionalOr)
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
		p.SetState(69)

		var _x = p.ConditionalAnd()

		localctx.(*ConditionalOrContext).e = _x
	}
	p.SetState(74)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == TypesParserLOGICAL_OR {
		{
			p.SetState(70)

			var _m = p.Match(TypesParserLOGICAL_OR)

			localctx.(*ConditionalOrContext).s9 = _m
		}
		localctx.(*ConditionalOrContext).ops = append(localctx.(*ConditionalOrContext).ops, localctx.(*ConditionalOrContext).s9)
		{
			p.SetState(71)

			var _x = p.ConditionalAnd()

			localctx.(*ConditionalOrContext)._conditionalAnd = _x
		}
		localctx.(*ConditionalOrContext).e1 = append(localctx.(*ConditionalOrContext).e1, localctx.(*ConditionalOrContext)._conditionalAnd)

		p.SetState(76)
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

	// GetS8 returns the s8 token.
	GetS8() antlr.Token

	// SetS8 sets the s8 token.
	SetS8(antlr.Token)

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
	s8        antlr.Token
	ops       []antlr.Token
	_relation IRelationContext
	e1        []IRelationContext
}

func NewEmptyConditionalAndContext() *ConditionalAndContext {
	var p = new(ConditionalAndContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = TypesParserRULE_conditionalAnd
	return p
}

func (*ConditionalAndContext) IsConditionalAndContext() {}

func NewConditionalAndContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalAndContext {
	var p = new(ConditionalAndContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_conditionalAnd

	return p
}

func (s *ConditionalAndContext) GetParser() antlr.Parser { return s.parser }

func (s *ConditionalAndContext) GetS8() antlr.Token { return s.s8 }

func (s *ConditionalAndContext) SetS8(v antlr.Token) { s.s8 = v }

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
	return s.GetTokens(TypesParserLOGICAL_AND)
}

func (s *ConditionalAndContext) LOGICAL_AND(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserLOGICAL_AND, i)
}

func (s *ConditionalAndContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConditionalAndContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConditionalAndContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterConditionalAnd(s)
	}
}

func (s *ConditionalAndContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitConditionalAnd(s)
	}
}

func (s *ConditionalAndContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitConditionalAnd(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) ConditionalAnd() (localctx IConditionalAndContext) {
	this := p
	_ = this

	localctx = NewConditionalAndContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, TypesParserRULE_conditionalAnd)
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
		p.SetState(77)

		var _x = p.relation(0)

		localctx.(*ConditionalAndContext).e = _x
	}
	p.SetState(82)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == TypesParserLOGICAL_AND {
		{
			p.SetState(78)

			var _m = p.Match(TypesParserLOGICAL_AND)

			localctx.(*ConditionalAndContext).s8 = _m
		}
		localctx.(*ConditionalAndContext).ops = append(localctx.(*ConditionalAndContext).ops, localctx.(*ConditionalAndContext).s8)
		{
			p.SetState(79)

			var _x = p.relation(0)

			localctx.(*ConditionalAndContext)._relation = _x
		}
		localctx.(*ConditionalAndContext).e1 = append(localctx.(*ConditionalAndContext).e1, localctx.(*ConditionalAndContext)._relation)

		p.SetState(84)
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
	p.RuleIndex = TypesParserRULE_relation
	return p
}

func (*RelationContext) IsRelationContext() {}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	var p = new(RelationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_relation

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
	return s.GetToken(TypesParserLESS, 0)
}

func (s *RelationContext) LESS_EQUALS() antlr.TerminalNode {
	return s.GetToken(TypesParserLESS_EQUALS, 0)
}

func (s *RelationContext) GREATER_EQUALS() antlr.TerminalNode {
	return s.GetToken(TypesParserGREATER_EQUALS, 0)
}

func (s *RelationContext) GREATER() antlr.TerminalNode {
	return s.GetToken(TypesParserGREATER, 0)
}

func (s *RelationContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(TypesParserEQUALS, 0)
}

func (s *RelationContext) NOT_EQUALS() antlr.TerminalNode {
	return s.GetToken(TypesParserNOT_EQUALS, 0)
}

func (s *RelationContext) IN() antlr.TerminalNode {
	return s.GetToken(TypesParserIN, 0)
}

func (s *RelationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RelationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RelationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterRelation(s)
	}
}

func (s *RelationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitRelation(s)
	}
}

func (s *RelationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitRelation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Relation() (localctx IRelationContext) {
	return p.relation(0)
}

func (p *TypesParser) relation(_p int) (localctx IRelationContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewRelationContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IRelationContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 14
	p.EnterRecursionRule(localctx, 14, TypesParserRULE_relation, _p)
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
		p.SetState(86)
		p.calc(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(93)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewRelationContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, TypesParserRULE_relation)
			p.SetState(88)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(89)

				var _lt = p.GetTokenStream().LT(1)

				localctx.(*RelationContext).op = _lt

				_la = p.GetTokenStream().LA(1)

				if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<TypesParserEQUALS)|(1<<TypesParserNOT_EQUALS)|(1<<TypesParserIN)|(1<<TypesParserLESS)|(1<<TypesParserLESS_EQUALS)|(1<<TypesParserGREATER_EQUALS)|(1<<TypesParserGREATER))) != 0) {
					var _ri = p.GetErrorHandler().RecoverInline(p)

					localctx.(*RelationContext).op = _ri
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(90)
				p.relation(2)
			}

		}
		p.SetState(95)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext())
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
	p.RuleIndex = TypesParserRULE_calc
	return p
}

func (*CalcContext) IsCalcContext() {}

func NewCalcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CalcContext {
	var p = new(CalcContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_calc

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
	return s.GetToken(TypesParserSTAR, 0)
}

func (s *CalcContext) SLASH() antlr.TerminalNode {
	return s.GetToken(TypesParserSLASH, 0)
}

func (s *CalcContext) PERCENT() antlr.TerminalNode {
	return s.GetToken(TypesParserPERCENT, 0)
}

func (s *CalcContext) PLUS() antlr.TerminalNode {
	return s.GetToken(TypesParserPLUS, 0)
}

func (s *CalcContext) MINUS() antlr.TerminalNode {
	return s.GetToken(TypesParserMINUS, 0)
}

func (s *CalcContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CalcContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CalcContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterCalc(s)
	}
}

func (s *CalcContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitCalc(s)
	}
}

func (s *CalcContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitCalc(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Calc() (localctx ICalcContext) {
	return p.calc(0)
}

func (p *TypesParser) calc(_p int) (localctx ICalcContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewCalcContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ICalcContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 16
	p.EnterRecursionRule(localctx, 16, TypesParserRULE_calc, _p)
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
		p.SetState(97)
		p.Unary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(107)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(105)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
			case 1:
				localctx = NewCalcContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, TypesParserRULE_calc)
				p.SetState(99)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(100)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*CalcContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<TypesParserSTAR)|(1<<TypesParserSLASH)|(1<<TypesParserPERCENT))) != 0) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*CalcContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(101)
					p.calc(3)
				}

			case 2:
				localctx = NewCalcContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, TypesParserRULE_calc)
				p.SetState(102)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(103)

					var _lt = p.GetTokenStream().LT(1)

					localctx.(*CalcContext).op = _lt

					_la = p.GetTokenStream().LA(1)

					if !(_la == TypesParserMINUS || _la == TypesParserPLUS) {
						var _ri = p.GetErrorHandler().RecoverInline(p)

						localctx.(*CalcContext).op = _ri
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(104)
					p.calc(2)
				}

			}

		}
		p.SetState(109)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
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
	p.RuleIndex = TypesParserRULE_unary
	return p
}

func (*UnaryContext) IsUnaryContext() {}

func NewUnaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryContext {
	var p = new(UnaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_unary

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
	s19 antlr.Token
	ops []antlr.Token
}

func NewLogicalNotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalNotContext {
	var p = new(LogicalNotContext)

	p.UnaryContext = NewEmptyUnaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*UnaryContext))

	return p
}

func (s *LogicalNotContext) GetS19() antlr.Token { return s.s19 }

func (s *LogicalNotContext) SetS19(v antlr.Token) { s.s19 = v }

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
	return s.GetTokens(TypesParserEXCLAM)
}

func (s *LogicalNotContext) EXCLAM(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserEXCLAM, i)
}

func (s *LogicalNotContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterLogicalNot(s)
	}
}

func (s *LogicalNotContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitLogicalNot(s)
	}
}

func (s *LogicalNotContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterMemberExpr(s)
	}
}

func (s *MemberExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitMemberExpr(s)
	}
}

func (s *MemberExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitMemberExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

type NegateContext struct {
	*UnaryContext
	s18 antlr.Token
	ops []antlr.Token
}

func NewNegateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NegateContext {
	var p = new(NegateContext)

	p.UnaryContext = NewEmptyUnaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*UnaryContext))

	return p
}

func (s *NegateContext) GetS18() antlr.Token { return s.s18 }

func (s *NegateContext) SetS18(v antlr.Token) { s.s18 = v }

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
	return s.GetTokens(TypesParserMINUS)
}

func (s *NegateContext) MINUS(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserMINUS, i)
}

func (s *NegateContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterNegate(s)
	}
}

func (s *NegateContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitNegate(s)
	}
}

func (s *NegateContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitNegate(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Unary() (localctx IUnaryContext) {
	this := p
	_ = this

	localctx = NewUnaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, TypesParserRULE_unary)
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

	p.SetState(123)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		localctx = NewMemberExprContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(110)
			p.member(0)
		}

	case 2:
		localctx = NewLogicalNotContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		p.SetState(112)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == TypesParserEXCLAM {
			{
				p.SetState(111)

				var _m = p.Match(TypesParserEXCLAM)

				localctx.(*LogicalNotContext).s19 = _m
			}
			localctx.(*LogicalNotContext).ops = append(localctx.(*LogicalNotContext).ops, localctx.(*LogicalNotContext).s19)

			p.SetState(114)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(116)
			p.member(0)
		}

	case 3:
		localctx = NewNegateContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		p.SetState(118)
		p.GetErrorHandler().Sync(p)
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(117)

					var _m = p.Match(TypesParserMINUS)

					localctx.(*NegateContext).s18 = _m
				}
				localctx.(*NegateContext).ops = append(localctx.(*NegateContext).ops, localctx.(*NegateContext).s18)

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(120)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())
		}
		{
			p.SetState(122)
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
	p.RuleIndex = TypesParserRULE_member
	return p
}

func (*MemberContext) IsMemberContext() {}

func NewMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberContext {
	var p = new(MemberContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_member

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
	return s.GetToken(TypesParserDOT, 0)
}

func (s *SelectOrCallContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(TypesParserIDENTIFIER, 0)
}

func (s *SelectOrCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(TypesParserRPAREN, 0)
}

func (s *SelectOrCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(TypesParserLPAREN, 0)
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterSelectOrCall(s)
	}
}

func (s *SelectOrCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitSelectOrCall(s)
	}
}

func (s *SelectOrCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterPrimaryExpr(s)
	}
}

func (s *PrimaryExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitPrimaryExpr(s)
	}
}

func (s *PrimaryExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserRPRACKET, 0)
}

func (s *IndexContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(TypesParserLBRACKET, 0)
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterIndex(s)
	}
}

func (s *IndexContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitIndex(s)
	}
}

func (s *IndexContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserRBRACE, 0)
}

func (s *CreateMessageContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(TypesParserLBRACE, 0)
}

func (s *CreateMessageContext) COMMA() antlr.TerminalNode {
	return s.GetToken(TypesParserCOMMA, 0)
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterCreateMessage(s)
	}
}

func (s *CreateMessageContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitCreateMessage(s)
	}
}

func (s *CreateMessageContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitCreateMessage(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Member() (localctx IMemberContext) {
	return p.member(0)
}

func (p *TypesParser) member(_p int) (localctx IMemberContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewMemberContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IMemberContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 20
	p.EnterRecursionRule(localctx, 20, TypesParserRULE_member, _p)
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
		p.SetState(126)
		p.Primary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 18, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(152)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext()) {
			case 1:
				localctx = NewSelectOrCallContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, TypesParserRULE_member)
				p.SetState(128)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				{
					p.SetState(129)

					var _m = p.Match(TypesParserDOT)

					localctx.(*SelectOrCallContext).op = _m
				}
				{
					p.SetState(130)

					var _m = p.Match(TypesParserIDENTIFIER)

					localctx.(*SelectOrCallContext).id = _m
				}
				p.SetState(136)
				p.GetErrorHandler().Sync(p)

				if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext()) == 1 {
					{
						p.SetState(131)

						var _m = p.Match(TypesParserLPAREN)

						localctx.(*SelectOrCallContext).open = _m
					}
					p.SetState(133)
					p.GetErrorHandler().Sync(p)
					_la = p.GetTokenStream().LA(1)

					if ((_la-10)&-(0x1f+1)) == 0 && ((1<<uint((_la-10)))&((1<<(TypesParserLBRACKET-10))|(1<<(TypesParserLBRACE-10))|(1<<(TypesParserLPAREN-10))|(1<<(TypesParserDOT-10))|(1<<(TypesParserMINUS-10))|(1<<(TypesParserEXCLAM-10))|(1<<(TypesParserCEL_TRUE-10))|(1<<(TypesParserCEL_FALSE-10))|(1<<(TypesParserNUL-10))|(1<<(TypesParserNUM_FLOAT-10))|(1<<(TypesParserNUM_INT-10))|(1<<(TypesParserNUM_UINT-10))|(1<<(TypesParserSTRING-10))|(1<<(TypesParserBYTES-10))|(1<<(TypesParserIDENTIFIER-10)))) != 0 {
						{
							p.SetState(132)

							var _x = p.ExprList()

							localctx.(*SelectOrCallContext).args = _x
						}

					}
					{
						p.SetState(135)
						p.Match(TypesParserRPAREN)
					}

				}

			case 2:
				localctx = NewIndexContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, TypesParserRULE_member)
				p.SetState(138)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(139)

					var _m = p.Match(TypesParserLBRACKET)

					localctx.(*IndexContext).op = _m
				}
				{
					p.SetState(140)

					var _x = p.Expr()

					localctx.(*IndexContext).index = _x
				}
				{
					p.SetState(141)
					p.Match(TypesParserRPRACKET)
				}

			case 3:
				localctx = NewCreateMessageContext(p, NewMemberContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, TypesParserRULE_member)
				p.SetState(143)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(144)

					var _m = p.Match(TypesParserLBRACE)

					localctx.(*CreateMessageContext).op = _m
				}
				p.SetState(146)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == TypesParserIDENTIFIER {
					{
						p.SetState(145)

						var _x = p.FieldInitializerList()

						localctx.(*CreateMessageContext).entries = _x
					}

				}
				p.SetState(149)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)

				if _la == TypesParserCOMMA {
					{
						p.SetState(148)
						p.Match(TypesParserCOMMA)
					}

				}
				{
					p.SetState(151)
					p.Match(TypesParserRBRACE)
				}

			}

		}
		p.SetState(156)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 18, p.GetParserRuleContext())
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
	p.RuleIndex = TypesParserRULE_primary
	return p
}

func (*PrimaryContext) IsPrimaryContext() {}

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext {
	var p = new(PrimaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_primary

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
	return s.GetToken(TypesParserRPRACKET, 0)
}

func (s *CreateListContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(TypesParserLBRACKET, 0)
}

func (s *CreateListContext) COMMA() antlr.TerminalNode {
	return s.GetToken(TypesParserCOMMA, 0)
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterCreateList(s)
	}
}

func (s *CreateListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitCreateList(s)
	}
}

func (s *CreateListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserRBRACE, 0)
}

func (s *CreateStructContext) LBRACE() antlr.TerminalNode {
	return s.GetToken(TypesParserLBRACE, 0)
}

func (s *CreateStructContext) COMMA() antlr.TerminalNode {
	return s.GetToken(TypesParserCOMMA, 0)
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterCreateStruct(s)
	}
}

func (s *CreateStructContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitCreateStruct(s)
	}
}

func (s *CreateStructContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterConstantLiteral(s)
	}
}

func (s *ConstantLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitConstantLiteral(s)
	}
}

func (s *ConstantLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserLPAREN, 0)
}

func (s *NestedContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(TypesParserRPAREN, 0)
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterNested(s)
	}
}

func (s *NestedContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitNested(s)
	}
}

func (s *NestedContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserIDENTIFIER, 0)
}

func (s *IdentOrGlobalCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(TypesParserRPAREN, 0)
}

func (s *IdentOrGlobalCallContext) DOT() antlr.TerminalNode {
	return s.GetToken(TypesParserDOT, 0)
}

func (s *IdentOrGlobalCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(TypesParserLPAREN, 0)
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
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterIdentOrGlobalCall(s)
	}
}

func (s *IdentOrGlobalCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitIdentOrGlobalCall(s)
	}
}

func (s *IdentOrGlobalCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitIdentOrGlobalCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Primary() (localctx IPrimaryContext) {
	this := p
	_ = this

	localctx = NewPrimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, TypesParserRULE_primary)
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

	p.SetState(189)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case TypesParserDOT, TypesParserIDENTIFIER:
		localctx = NewIdentOrGlobalCallContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == TypesParserDOT {
			{
				p.SetState(157)

				var _m = p.Match(TypesParserDOT)

				localctx.(*IdentOrGlobalCallContext).leadingDot = _m
			}

		}
		{
			p.SetState(160)

			var _m = p.Match(TypesParserIDENTIFIER)

			localctx.(*IdentOrGlobalCallContext).id = _m
		}
		p.SetState(166)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(161)

				var _m = p.Match(TypesParserLPAREN)

				localctx.(*IdentOrGlobalCallContext).op = _m
			}
			p.SetState(163)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if ((_la-10)&-(0x1f+1)) == 0 && ((1<<uint((_la-10)))&((1<<(TypesParserLBRACKET-10))|(1<<(TypesParserLBRACE-10))|(1<<(TypesParserLPAREN-10))|(1<<(TypesParserDOT-10))|(1<<(TypesParserMINUS-10))|(1<<(TypesParserEXCLAM-10))|(1<<(TypesParserCEL_TRUE-10))|(1<<(TypesParserCEL_FALSE-10))|(1<<(TypesParserNUL-10))|(1<<(TypesParserNUM_FLOAT-10))|(1<<(TypesParserNUM_INT-10))|(1<<(TypesParserNUM_UINT-10))|(1<<(TypesParserSTRING-10))|(1<<(TypesParserBYTES-10))|(1<<(TypesParserIDENTIFIER-10)))) != 0 {
				{
					p.SetState(162)

					var _x = p.ExprList()

					localctx.(*IdentOrGlobalCallContext).args = _x
				}

			}
			{
				p.SetState(165)
				p.Match(TypesParserRPAREN)
			}

		}

	case TypesParserLPAREN:
		localctx = NewNestedContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(168)
			p.Match(TypesParserLPAREN)
		}
		{
			p.SetState(169)

			var _x = p.Expr()

			localctx.(*NestedContext).e = _x
		}
		{
			p.SetState(170)
			p.Match(TypesParserRPAREN)
		}

	case TypesParserLBRACKET:
		localctx = NewCreateListContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(172)

			var _m = p.Match(TypesParserLBRACKET)

			localctx.(*CreateListContext).op = _m
		}
		p.SetState(174)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if ((_la-10)&-(0x1f+1)) == 0 && ((1<<uint((_la-10)))&((1<<(TypesParserLBRACKET-10))|(1<<(TypesParserLBRACE-10))|(1<<(TypesParserLPAREN-10))|(1<<(TypesParserDOT-10))|(1<<(TypesParserMINUS-10))|(1<<(TypesParserEXCLAM-10))|(1<<(TypesParserCEL_TRUE-10))|(1<<(TypesParserCEL_FALSE-10))|(1<<(TypesParserNUL-10))|(1<<(TypesParserNUM_FLOAT-10))|(1<<(TypesParserNUM_INT-10))|(1<<(TypesParserNUM_UINT-10))|(1<<(TypesParserSTRING-10))|(1<<(TypesParserBYTES-10))|(1<<(TypesParserIDENTIFIER-10)))) != 0 {
			{
				p.SetState(173)

				var _x = p.ExprList()

				localctx.(*CreateListContext).elems = _x
			}

		}
		p.SetState(177)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == TypesParserCOMMA {
			{
				p.SetState(176)
				p.Match(TypesParserCOMMA)
			}

		}
		{
			p.SetState(179)
			p.Match(TypesParserRPRACKET)
		}

	case TypesParserLBRACE:
		localctx = NewCreateStructContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(180)

			var _m = p.Match(TypesParserLBRACE)

			localctx.(*CreateStructContext).op = _m
		}
		p.SetState(182)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if ((_la-10)&-(0x1f+1)) == 0 && ((1<<uint((_la-10)))&((1<<(TypesParserLBRACKET-10))|(1<<(TypesParserLBRACE-10))|(1<<(TypesParserLPAREN-10))|(1<<(TypesParserDOT-10))|(1<<(TypesParserMINUS-10))|(1<<(TypesParserEXCLAM-10))|(1<<(TypesParserCEL_TRUE-10))|(1<<(TypesParserCEL_FALSE-10))|(1<<(TypesParserNUL-10))|(1<<(TypesParserNUM_FLOAT-10))|(1<<(TypesParserNUM_INT-10))|(1<<(TypesParserNUM_UINT-10))|(1<<(TypesParserSTRING-10))|(1<<(TypesParserBYTES-10))|(1<<(TypesParserIDENTIFIER-10)))) != 0 {
			{
				p.SetState(181)

				var _x = p.MapInitializerList()

				localctx.(*CreateStructContext).entries = _x
			}

		}
		p.SetState(185)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == TypesParserCOMMA {
			{
				p.SetState(184)
				p.Match(TypesParserCOMMA)
			}

		}
		{
			p.SetState(187)
			p.Match(TypesParserRBRACE)
		}

	case TypesParserMINUS, TypesParserCEL_TRUE, TypesParserCEL_FALSE, TypesParserNUL, TypesParserNUM_FLOAT, TypesParserNUM_INT, TypesParserNUM_UINT, TypesParserSTRING, TypesParserBYTES:
		localctx = NewConstantLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(188)
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
	p.RuleIndex = TypesParserRULE_exprList
	return p
}

func (*ExprListContext) IsExprListContext() {}

func NewExprListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprListContext {
	var p = new(ExprListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_exprList

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
	return s.GetTokens(TypesParserCOMMA)
}

func (s *ExprListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserCOMMA, i)
}

func (s *ExprListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterExprList(s)
	}
}

func (s *ExprListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitExprList(s)
	}
}

func (s *ExprListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitExprList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) ExprList() (localctx IExprListContext) {
	this := p
	_ = this

	localctx = NewExprListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, TypesParserRULE_exprList)

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
		p.SetState(191)

		var _x = p.Expr()

		localctx.(*ExprListContext)._expr = _x
	}
	localctx.(*ExprListContext).e = append(localctx.(*ExprListContext).e, localctx.(*ExprListContext)._expr)
	p.SetState(196)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(192)
				p.Match(TypesParserCOMMA)
			}
			{
				p.SetState(193)

				var _x = p.Expr()

				localctx.(*ExprListContext)._expr = _x
			}
			localctx.(*ExprListContext).e = append(localctx.(*ExprListContext).e, localctx.(*ExprListContext)._expr)

		}
		p.SetState(198)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())
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

	// GetS21 returns the s21 token.
	GetS21() antlr.Token

	// Set_IDENTIFIER sets the _IDENTIFIER token.
	Set_IDENTIFIER(antlr.Token)

	// SetS21 sets the s21 token.
	SetS21(antlr.Token)

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
	s21         antlr.Token
	cols        []antlr.Token
	_expr       IExprContext
	values      []IExprContext
}

func NewEmptyFieldInitializerListContext() *FieldInitializerListContext {
	var p = new(FieldInitializerListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = TypesParserRULE_fieldInitializerList
	return p
}

func (*FieldInitializerListContext) IsFieldInitializerListContext() {}

func NewFieldInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldInitializerListContext {
	var p = new(FieldInitializerListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_fieldInitializerList

	return p
}

func (s *FieldInitializerListContext) GetParser() antlr.Parser { return s.parser }

func (s *FieldInitializerListContext) Get_IDENTIFIER() antlr.Token { return s._IDENTIFIER }

func (s *FieldInitializerListContext) GetS21() antlr.Token { return s.s21 }

func (s *FieldInitializerListContext) Set_IDENTIFIER(v antlr.Token) { s._IDENTIFIER = v }

func (s *FieldInitializerListContext) SetS21(v antlr.Token) { s.s21 = v }

func (s *FieldInitializerListContext) GetFields() []antlr.Token { return s.fields }

func (s *FieldInitializerListContext) GetCols() []antlr.Token { return s.cols }

func (s *FieldInitializerListContext) SetFields(v []antlr.Token) { s.fields = v }

func (s *FieldInitializerListContext) SetCols(v []antlr.Token) { s.cols = v }

func (s *FieldInitializerListContext) Get_expr() IExprContext { return s._expr }

func (s *FieldInitializerListContext) Set_expr(v IExprContext) { s._expr = v }

func (s *FieldInitializerListContext) GetValues() []IExprContext { return s.values }

func (s *FieldInitializerListContext) SetValues(v []IExprContext) { s.values = v }

func (s *FieldInitializerListContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(TypesParserIDENTIFIER)
}

func (s *FieldInitializerListContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserIDENTIFIER, i)
}

func (s *FieldInitializerListContext) AllCOLON() []antlr.TerminalNode {
	return s.GetTokens(TypesParserCOLON)
}

func (s *FieldInitializerListContext) COLON(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserCOLON, i)
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
	return s.GetTokens(TypesParserCOMMA)
}

func (s *FieldInitializerListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserCOMMA, i)
}

func (s *FieldInitializerListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FieldInitializerListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FieldInitializerListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterFieldInitializerList(s)
	}
}

func (s *FieldInitializerListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitFieldInitializerList(s)
	}
}

func (s *FieldInitializerListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitFieldInitializerList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) FieldInitializerList() (localctx IFieldInitializerListContext) {
	this := p
	_ = this

	localctx = NewFieldInitializerListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, TypesParserRULE_fieldInitializerList)

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
		p.SetState(199)

		var _m = p.Match(TypesParserIDENTIFIER)

		localctx.(*FieldInitializerListContext)._IDENTIFIER = _m
	}
	localctx.(*FieldInitializerListContext).fields = append(localctx.(*FieldInitializerListContext).fields, localctx.(*FieldInitializerListContext)._IDENTIFIER)
	{
		p.SetState(200)

		var _m = p.Match(TypesParserCOLON)

		localctx.(*FieldInitializerListContext).s21 = _m
	}
	localctx.(*FieldInitializerListContext).cols = append(localctx.(*FieldInitializerListContext).cols, localctx.(*FieldInitializerListContext).s21)
	{
		p.SetState(201)

		var _x = p.Expr()

		localctx.(*FieldInitializerListContext)._expr = _x
	}
	localctx.(*FieldInitializerListContext).values = append(localctx.(*FieldInitializerListContext).values, localctx.(*FieldInitializerListContext)._expr)
	p.SetState(208)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(202)
				p.Match(TypesParserCOMMA)
			}
			{
				p.SetState(203)

				var _m = p.Match(TypesParserIDENTIFIER)

				localctx.(*FieldInitializerListContext)._IDENTIFIER = _m
			}
			localctx.(*FieldInitializerListContext).fields = append(localctx.(*FieldInitializerListContext).fields, localctx.(*FieldInitializerListContext)._IDENTIFIER)
			{
				p.SetState(204)

				var _m = p.Match(TypesParserCOLON)

				localctx.(*FieldInitializerListContext).s21 = _m
			}
			localctx.(*FieldInitializerListContext).cols = append(localctx.(*FieldInitializerListContext).cols, localctx.(*FieldInitializerListContext).s21)
			{
				p.SetState(205)

				var _x = p.Expr()

				localctx.(*FieldInitializerListContext)._expr = _x
			}
			localctx.(*FieldInitializerListContext).values = append(localctx.(*FieldInitializerListContext).values, localctx.(*FieldInitializerListContext)._expr)

		}
		p.SetState(210)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext())
	}

	return localctx
}

// IMapInitializerListContext is an interface to support dynamic dispatch.
type IMapInitializerListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// GetS21 returns the s21 token.
	GetS21() antlr.Token

	// SetS21 sets the s21 token.
	SetS21(antlr.Token)

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
	s21    antlr.Token
	cols   []antlr.Token
	values []IExprContext
}

func NewEmptyMapInitializerListContext() *MapInitializerListContext {
	var p = new(MapInitializerListContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = TypesParserRULE_mapInitializerList
	return p
}

func (*MapInitializerListContext) IsMapInitializerListContext() {}

func NewMapInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapInitializerListContext {
	var p = new(MapInitializerListContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_mapInitializerList

	return p
}

func (s *MapInitializerListContext) GetParser() antlr.Parser { return s.parser }

func (s *MapInitializerListContext) GetS21() antlr.Token { return s.s21 }

func (s *MapInitializerListContext) SetS21(v antlr.Token) { s.s21 = v }

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
	return s.GetTokens(TypesParserCOLON)
}

func (s *MapInitializerListContext) COLON(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserCOLON, i)
}

func (s *MapInitializerListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(TypesParserCOMMA)
}

func (s *MapInitializerListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(TypesParserCOMMA, i)
}

func (s *MapInitializerListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MapInitializerListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MapInitializerListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterMapInitializerList(s)
	}
}

func (s *MapInitializerListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitMapInitializerList(s)
	}
}

func (s *MapInitializerListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitMapInitializerList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) MapInitializerList() (localctx IMapInitializerListContext) {
	this := p
	_ = this

	localctx = NewMapInitializerListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, TypesParserRULE_mapInitializerList)

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
		p.SetState(211)

		var _x = p.Expr()

		localctx.(*MapInitializerListContext)._expr = _x
	}
	localctx.(*MapInitializerListContext).keys = append(localctx.(*MapInitializerListContext).keys, localctx.(*MapInitializerListContext)._expr)
	{
		p.SetState(212)

		var _m = p.Match(TypesParserCOLON)

		localctx.(*MapInitializerListContext).s21 = _m
	}
	localctx.(*MapInitializerListContext).cols = append(localctx.(*MapInitializerListContext).cols, localctx.(*MapInitializerListContext).s21)
	{
		p.SetState(213)

		var _x = p.Expr()

		localctx.(*MapInitializerListContext)._expr = _x
	}
	localctx.(*MapInitializerListContext).values = append(localctx.(*MapInitializerListContext).values, localctx.(*MapInitializerListContext)._expr)
	p.SetState(221)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(214)
				p.Match(TypesParserCOMMA)
			}
			{
				p.SetState(215)

				var _x = p.Expr()

				localctx.(*MapInitializerListContext)._expr = _x
			}
			localctx.(*MapInitializerListContext).keys = append(localctx.(*MapInitializerListContext).keys, localctx.(*MapInitializerListContext)._expr)
			{
				p.SetState(216)

				var _m = p.Match(TypesParserCOLON)

				localctx.(*MapInitializerListContext).s21 = _m
			}
			localctx.(*MapInitializerListContext).cols = append(localctx.(*MapInitializerListContext).cols, localctx.(*MapInitializerListContext).s21)
			{
				p.SetState(217)

				var _x = p.Expr()

				localctx.(*MapInitializerListContext)._expr = _x
			}
			localctx.(*MapInitializerListContext).values = append(localctx.(*MapInitializerListContext).values, localctx.(*MapInitializerListContext)._expr)

		}
		p.SetState(223)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())
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
	p.RuleIndex = TypesParserRULE_literal
	return p
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = TypesParserRULE_literal

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
	return s.GetToken(TypesParserBYTES, 0)
}

func (s *BytesContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterBytes(s)
	}
}

func (s *BytesContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitBytes(s)
	}
}

func (s *BytesContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserNUM_UINT, 0)
}

func (s *UintContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterUint(s)
	}
}

func (s *UintContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitUint(s)
	}
}

func (s *UintContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserNUL, 0)
}

func (s *NullContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterNull(s)
	}
}

func (s *NullContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitNull(s)
	}
}

func (s *NullContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserCEL_FALSE, 0)
}

func (s *BoolFalseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterBoolFalse(s)
	}
}

func (s *BoolFalseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitBoolFalse(s)
	}
}

func (s *BoolFalseContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserSTRING, 0)
}

func (s *StringContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterString(s)
	}
}

func (s *StringContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitString(s)
	}
}

func (s *StringContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserNUM_FLOAT, 0)
}

func (s *DoubleContext) MINUS() antlr.TerminalNode {
	return s.GetToken(TypesParserMINUS, 0)
}

func (s *DoubleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterDouble(s)
	}
}

func (s *DoubleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitDouble(s)
	}
}

func (s *DoubleContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserCEL_TRUE, 0)
}

func (s *BoolTrueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterBoolTrue(s)
	}
}

func (s *BoolTrueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitBoolTrue(s)
	}
}

func (s *BoolTrueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
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
	return s.GetToken(TypesParserNUM_INT, 0)
}

func (s *IntContext) MINUS() antlr.TerminalNode {
	return s.GetToken(TypesParserMINUS, 0)
}

func (s *IntContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.EnterInt(s)
	}
}

func (s *IntContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(TypesListener); ok {
		listenerT.ExitInt(s)
	}
}

func (s *IntContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case TypesVisitor:
		return t.VisitInt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *TypesParser) Literal() (localctx ILiteralContext) {
	this := p
	_ = this

	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, TypesParserRULE_literal)
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

	p.SetState(238)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		localctx = NewIntContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(225)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == TypesParserMINUS {
			{
				p.SetState(224)

				var _m = p.Match(TypesParserMINUS)

				localctx.(*IntContext).sign = _m
			}

		}
		{
			p.SetState(227)

			var _m = p.Match(TypesParserNUM_INT)

			localctx.(*IntContext).tok = _m
		}

	case 2:
		localctx = NewUintContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(228)

			var _m = p.Match(TypesParserNUM_UINT)

			localctx.(*UintContext).tok = _m
		}

	case 3:
		localctx = NewDoubleContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		p.SetState(230)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == TypesParserMINUS {
			{
				p.SetState(229)

				var _m = p.Match(TypesParserMINUS)

				localctx.(*DoubleContext).sign = _m
			}

		}
		{
			p.SetState(232)

			var _m = p.Match(TypesParserNUM_FLOAT)

			localctx.(*DoubleContext).tok = _m
		}

	case 4:
		localctx = NewStringContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(233)

			var _m = p.Match(TypesParserSTRING)

			localctx.(*StringContext).tok = _m
		}

	case 5:
		localctx = NewBytesContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(234)

			var _m = p.Match(TypesParserBYTES)

			localctx.(*BytesContext).tok = _m
		}

	case 6:
		localctx = NewBoolTrueContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(235)

			var _m = p.Match(TypesParserCEL_TRUE)

			localctx.(*BoolTrueContext).tok = _m
		}

	case 7:
		localctx = NewBoolFalseContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(236)

			var _m = p.Match(TypesParserCEL_FALSE)

			localctx.(*BoolFalseContext).tok = _m
		}

	case 8:
		localctx = NewNullContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(237)

			var _m = p.Match(TypesParserNUL)

			localctx.(*NullContext).tok = _m
		}

	}

	return localctx
}

func (p *TypesParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 7:
		var t *RelationContext = nil
		if localctx != nil {
			t = localctx.(*RelationContext)
		}
		return p.Relation_Sempred(t, predIndex)

	case 8:
		var t *CalcContext = nil
		if localctx != nil {
			t = localctx.(*CalcContext)
		}
		return p.Calc_Sempred(t, predIndex)

	case 10:
		var t *MemberContext = nil
		if localctx != nil {
			t = localctx.(*MemberContext)
		}
		return p.Member_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *TypesParser) Relation_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *TypesParser) Calc_Sempred(localctx antlr.RuleContext, predIndex int) bool {
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

func (p *TypesParser) Member_Sempred(localctx antlr.RuleContext, predIndex int) bool {
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
