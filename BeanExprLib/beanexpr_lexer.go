// Code generated from E:/Go/gofks/BeanExprLib\BeanExpr.g4 by ANTLR 4.10.1. DO NOT EDIT.

package BeanExprLib

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type BeanExprLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var beanexprlexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func beanexprlexerLexerInit() {
	staticData := &beanexprlexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.literalNames = []string{
		"", "'('", "')'", "','", "", "", "", "", "", "'.'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "StringArg", "FloatArg", "IntArg", "FuncName", "MethodName",
		"Dot", "Float", "WHITESPACE",
	}
	staticData.ruleNames = []string{
		"T__0", "T__1", "T__2", "DIGIT", "StringArg", "FloatArg", "IntArg",
		"FuncName", "MethodName", "Dot", "Float", "WHITESPACE",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 11, 96, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3,
		1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 5, 4, 40, 8, 4, 10, 4, 12, 4, 43, 9,
		4, 1, 4, 1, 4, 1, 5, 3, 5, 48, 8, 5, 1, 5, 1, 5, 1, 6, 3, 6, 53, 8, 6,
		1, 6, 4, 6, 56, 8, 6, 11, 6, 12, 6, 57, 1, 7, 1, 7, 5, 7, 62, 8, 7, 10,
		7, 12, 7, 65, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 4, 8, 71, 8, 8, 11, 8, 12,
		8, 72, 1, 9, 1, 9, 1, 10, 4, 10, 78, 8, 10, 11, 10, 12, 10, 79, 3, 10,
		82, 8, 10, 1, 10, 1, 10, 4, 10, 86, 8, 10, 11, 10, 12, 10, 87, 1, 11, 4,
		11, 91, 8, 11, 11, 11, 12, 11, 92, 1, 11, 1, 11, 0, 0, 12, 1, 1, 3, 2,
		5, 3, 7, 0, 9, 4, 11, 5, 13, 6, 15, 7, 17, 8, 19, 9, 21, 10, 23, 11, 1,
		0, 5, 1, 0, 48, 57, 2, 0, 39, 39, 92, 92, 2, 0, 65, 90, 97, 122, 3, 0,
		48, 57, 65, 90, 97, 122, 3, 0, 9, 10, 13, 13, 32, 32, 106, 0, 1, 1, 0,
		0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 11, 1, 0,
		0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0, 0, 17, 1, 0, 0, 0, 0, 19, 1,
		0, 0, 0, 0, 21, 1, 0, 0, 0, 0, 23, 1, 0, 0, 0, 1, 25, 1, 0, 0, 0, 3, 27,
		1, 0, 0, 0, 5, 29, 1, 0, 0, 0, 7, 31, 1, 0, 0, 0, 9, 33, 1, 0, 0, 0, 11,
		47, 1, 0, 0, 0, 13, 52, 1, 0, 0, 0, 15, 59, 1, 0, 0, 0, 17, 66, 1, 0, 0,
		0, 19, 74, 1, 0, 0, 0, 21, 81, 1, 0, 0, 0, 23, 90, 1, 0, 0, 0, 25, 26,
		5, 40, 0, 0, 26, 2, 1, 0, 0, 0, 27, 28, 5, 41, 0, 0, 28, 4, 1, 0, 0, 0,
		29, 30, 5, 44, 0, 0, 30, 6, 1, 0, 0, 0, 31, 32, 7, 0, 0, 0, 32, 8, 1, 0,
		0, 0, 33, 41, 5, 39, 0, 0, 34, 35, 5, 92, 0, 0, 35, 40, 9, 0, 0, 0, 36,
		37, 5, 39, 0, 0, 37, 40, 5, 39, 0, 0, 38, 40, 8, 1, 0, 0, 39, 34, 1, 0,
		0, 0, 39, 36, 1, 0, 0, 0, 39, 38, 1, 0, 0, 0, 40, 43, 1, 0, 0, 0, 41, 39,
		1, 0, 0, 0, 41, 42, 1, 0, 0, 0, 42, 44, 1, 0, 0, 0, 43, 41, 1, 0, 0, 0,
		44, 45, 5, 39, 0, 0, 45, 10, 1, 0, 0, 0, 46, 48, 5, 45, 0, 0, 47, 46, 1,
		0, 0, 0, 47, 48, 1, 0, 0, 0, 48, 49, 1, 0, 0, 0, 49, 50, 3, 21, 10, 0,
		50, 12, 1, 0, 0, 0, 51, 53, 5, 45, 0, 0, 52, 51, 1, 0, 0, 0, 52, 53, 1,
		0, 0, 0, 53, 55, 1, 0, 0, 0, 54, 56, 3, 7, 3, 0, 55, 54, 1, 0, 0, 0, 56,
		57, 1, 0, 0, 0, 57, 55, 1, 0, 0, 0, 57, 58, 1, 0, 0, 0, 58, 14, 1, 0, 0,
		0, 59, 63, 7, 2, 0, 0, 60, 62, 7, 3, 0, 0, 61, 60, 1, 0, 0, 0, 62, 65,
		1, 0, 0, 0, 63, 61, 1, 0, 0, 0, 63, 64, 1, 0, 0, 0, 64, 16, 1, 0, 0, 0,
		65, 63, 1, 0, 0, 0, 66, 70, 3, 15, 7, 0, 67, 68, 3, 19, 9, 0, 68, 69, 3,
		15, 7, 0, 69, 71, 1, 0, 0, 0, 70, 67, 1, 0, 0, 0, 71, 72, 1, 0, 0, 0, 72,
		70, 1, 0, 0, 0, 72, 73, 1, 0, 0, 0, 73, 18, 1, 0, 0, 0, 74, 75, 5, 46,
		0, 0, 75, 20, 1, 0, 0, 0, 76, 78, 3, 7, 3, 0, 77, 76, 1, 0, 0, 0, 78, 79,
		1, 0, 0, 0, 79, 77, 1, 0, 0, 0, 79, 80, 1, 0, 0, 0, 80, 82, 1, 0, 0, 0,
		81, 77, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82, 83, 1, 0, 0, 0, 83, 85, 5,
		46, 0, 0, 84, 86, 3, 7, 3, 0, 85, 84, 1, 0, 0, 0, 86, 87, 1, 0, 0, 0, 87,
		85, 1, 0, 0, 0, 87, 88, 1, 0, 0, 0, 88, 22, 1, 0, 0, 0, 89, 91, 7, 4, 0,
		0, 90, 89, 1, 0, 0, 0, 91, 92, 1, 0, 0, 0, 92, 90, 1, 0, 0, 0, 92, 93,
		1, 0, 0, 0, 93, 94, 1, 0, 0, 0, 94, 95, 6, 11, 0, 0, 95, 24, 1, 0, 0, 0,
		12, 0, 39, 41, 47, 52, 57, 63, 72, 79, 81, 87, 92, 1, 6, 0, 0,
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

// BeanExprLexerInit initializes any static state used to implement BeanExprLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewBeanExprLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func BeanExprLexerInit() {
	staticData := &beanexprlexerLexerStaticData
	staticData.once.Do(beanexprlexerLexerInit)
}

// NewBeanExprLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewBeanExprLexer(input antlr.CharStream) *BeanExprLexer {
	BeanExprLexerInit()
	l := new(BeanExprLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &beanexprlexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "BeanExpr.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// BeanExprLexer tokens.
const (
	BeanExprLexerT__0       = 1
	BeanExprLexerT__1       = 2
	BeanExprLexerT__2       = 3
	BeanExprLexerStringArg  = 4
	BeanExprLexerFloatArg   = 5
	BeanExprLexerIntArg     = 6
	BeanExprLexerFuncName   = 7
	BeanExprLexerMethodName = 8
	BeanExprLexerDot        = 9
	BeanExprLexerFloat      = 10
	BeanExprLexerWHITESPACE = 11
)
