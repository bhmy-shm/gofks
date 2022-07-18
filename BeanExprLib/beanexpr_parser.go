// Code generated from E:/Go/gofks/BeanExprLib\BeanExpr.g4 by ANTLR 4.10.1. DO NOT EDIT.

package BeanExprLib // BeanExpr

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

type BeanExprParser struct {
	*antlr.BaseParser
}

var beanexprParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func beanexprParserInit() {
	staticData := &beanexprParserStaticData
	staticData.literalNames = []string{
		"", "'('", "')'", "','", "", "", "", "", "", "'.'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "StringArg", "FloatArg", "IntArg", "FuncName", "MethodName",
		"Dot", "Float", "WHITESPACE",
	}
	staticData.ruleNames = []string{
		"start", "methodCall", "functionCall", "functionArgs",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 11, 37, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 1, 0, 1,
		0, 1, 0, 1, 0, 3, 0, 13, 8, 0, 1, 1, 1, 1, 1, 1, 3, 1, 18, 8, 1, 1, 1,
		1, 1, 1, 2, 1, 2, 1, 2, 3, 2, 25, 8, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 5,
		3, 32, 8, 3, 10, 3, 12, 3, 35, 9, 3, 1, 3, 0, 0, 4, 0, 2, 4, 6, 0, 1, 1,
		0, 4, 6, 36, 0, 12, 1, 0, 0, 0, 2, 14, 1, 0, 0, 0, 4, 21, 1, 0, 0, 0, 6,
		28, 1, 0, 0, 0, 8, 13, 3, 4, 2, 0, 9, 10, 3, 2, 1, 0, 10, 11, 5, 0, 0,
		1, 11, 13, 1, 0, 0, 0, 12, 8, 1, 0, 0, 0, 12, 9, 1, 0, 0, 0, 13, 1, 1,
		0, 0, 0, 14, 15, 5, 8, 0, 0, 15, 17, 5, 1, 0, 0, 16, 18, 3, 6, 3, 0, 17,
		16, 1, 0, 0, 0, 17, 18, 1, 0, 0, 0, 18, 19, 1, 0, 0, 0, 19, 20, 5, 2, 0,
		0, 20, 3, 1, 0, 0, 0, 21, 22, 5, 7, 0, 0, 22, 24, 5, 1, 0, 0, 23, 25, 3,
		6, 3, 0, 24, 23, 1, 0, 0, 0, 24, 25, 1, 0, 0, 0, 25, 26, 1, 0, 0, 0, 26,
		27, 5, 2, 0, 0, 27, 5, 1, 0, 0, 0, 28, 33, 7, 0, 0, 0, 29, 30, 5, 3, 0,
		0, 30, 32, 7, 0, 0, 0, 31, 29, 1, 0, 0, 0, 32, 35, 1, 0, 0, 0, 33, 31,
		1, 0, 0, 0, 33, 34, 1, 0, 0, 0, 34, 7, 1, 0, 0, 0, 35, 33, 1, 0, 0, 0,
		4, 12, 17, 24, 33,
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

// BeanExprParserInit initializes any static state used to implement BeanExprParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewBeanExprParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func BeanExprParserInit() {
	staticData := &beanexprParserStaticData
	staticData.once.Do(beanexprParserInit)
}

// NewBeanExprParser produces a new parser instance for the optional input antlr.TokenStream.
func NewBeanExprParser(input antlr.TokenStream) *BeanExprParser {
	BeanExprParserInit()
	this := new(BeanExprParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &beanexprParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "BeanExpr.g4"

	return this
}

// BeanExprParser tokens.
const (
	BeanExprParserEOF        = antlr.TokenEOF
	BeanExprParserT__0       = 1
	BeanExprParserT__1       = 2
	BeanExprParserT__2       = 3
	BeanExprParserStringArg  = 4
	BeanExprParserFloatArg   = 5
	BeanExprParserIntArg     = 6
	BeanExprParserFuncName   = 7
	BeanExprParserMethodName = 8
	BeanExprParserDot        = 9
	BeanExprParserFloat      = 10
	BeanExprParserWHITESPACE = 11
)

// BeanExprParser rules.
const (
	BeanExprParserRULE_start        = 0
	BeanExprParserRULE_methodCall   = 1
	BeanExprParserRULE_functionCall = 2
	BeanExprParserRULE_functionArgs = 3
)

// IStartContext is an interface to support dynamic dispatch.
type IStartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStartContext differentiates from other interfaces.
	IsStartContext()
}

type StartContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStartContext() *StartContext {
	var p = new(StartContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = BeanExprParserRULE_start
	return p
}

func (*StartContext) IsStartContext() {}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	var p = new(StartContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = BeanExprParserRULE_start

	return p
}

func (s *StartContext) GetParser() antlr.Parser { return s.parser }

func (s *StartContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *StartContext) MethodCall() IMethodCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMethodCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMethodCallContext)
}

func (s *StartContext) EOF() antlr.TerminalNode {
	return s.GetToken(BeanExprParserEOF, 0)
}

func (s *StartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.EnterStart(s)
	}
}

func (s *StartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.ExitStart(s)
	}
}

func (s *StartContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case BeanExprVisitor:
		return t.VisitStart(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *BeanExprParser) Start() (localctx IStartContext) {
	this := p
	_ = this

	localctx = NewStartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, BeanExprParserRULE_start)

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

	p.SetState(12)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case BeanExprParserFuncName:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(8)
			p.FunctionCall()
		}

	case BeanExprParserMethodName:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(9)
			p.MethodCall()
		}
		{
			p.SetState(10)
			p.Match(BeanExprParserEOF)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IMethodCallContext is an interface to support dynamic dispatch.
type IMethodCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMethodCallContext differentiates from other interfaces.
	IsMethodCallContext()
}

type MethodCallContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMethodCallContext() *MethodCallContext {
	var p = new(MethodCallContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = BeanExprParserRULE_methodCall
	return p
}

func (*MethodCallContext) IsMethodCallContext() {}

func NewMethodCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MethodCallContext {
	var p = new(MethodCallContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = BeanExprParserRULE_methodCall

	return p
}

func (s *MethodCallContext) GetParser() antlr.Parser { return s.parser }

func (s *MethodCallContext) MethodName() antlr.TerminalNode {
	return s.GetToken(BeanExprParserMethodName, 0)
}

func (s *MethodCallContext) FunctionArgs() IFunctionArgsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionArgsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionArgsContext)
}

func (s *MethodCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MethodCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MethodCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.EnterMethodCall(s)
	}
}

func (s *MethodCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.ExitMethodCall(s)
	}
}

func (s *MethodCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case BeanExprVisitor:
		return t.VisitMethodCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *BeanExprParser) MethodCall() (localctx IMethodCallContext) {
	this := p
	_ = this

	localctx = NewMethodCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, BeanExprParserRULE_methodCall)
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
		p.SetState(14)
		p.Match(BeanExprParserMethodName)
	}
	{
		p.SetState(15)
		p.Match(BeanExprParserT__0)
	}
	p.SetState(17)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<BeanExprParserStringArg)|(1<<BeanExprParserFloatArg)|(1<<BeanExprParserIntArg))) != 0 {
		{
			p.SetState(16)
			p.FunctionArgs()
		}

	}
	{
		p.SetState(19)
		p.Match(BeanExprParserT__1)
	}

	return localctx
}

// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = BeanExprParserRULE_functionCall
	return p
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = BeanExprParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) CopyFrom(ctx *FunctionCallContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type FuncCallContext struct {
	*FunctionCallContext
}

func NewFuncCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FuncCallContext {
	var p = new(FuncCallContext)

	p.FunctionCallContext = NewEmptyFunctionCallContext()
	p.parser = parser
	p.CopyFrom(ctx.(*FunctionCallContext))

	return p
}

func (s *FuncCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncCallContext) FuncName() antlr.TerminalNode {
	return s.GetToken(BeanExprParserFuncName, 0)
}

func (s *FuncCallContext) FunctionArgs() IFunctionArgsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionArgsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionArgsContext)
}

func (s *FuncCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.EnterFuncCall(s)
	}
}

func (s *FuncCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.ExitFuncCall(s)
	}
}

func (s *FuncCallContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case BeanExprVisitor:
		return t.VisitFuncCall(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *BeanExprParser) FunctionCall() (localctx IFunctionCallContext) {
	this := p
	_ = this

	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, BeanExprParserRULE_functionCall)
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

	localctx = NewFuncCallContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(21)
		p.Match(BeanExprParserFuncName)
	}
	{
		p.SetState(22)
		p.Match(BeanExprParserT__0)
	}
	p.SetState(24)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<BeanExprParserStringArg)|(1<<BeanExprParserFloatArg)|(1<<BeanExprParserIntArg))) != 0 {
		{
			p.SetState(23)
			p.FunctionArgs()
		}

	}
	{
		p.SetState(26)
		p.Match(BeanExprParserT__1)
	}

	return localctx
}

// IFunctionArgsContext is an interface to support dynamic dispatch.
type IFunctionArgsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionArgsContext differentiates from other interfaces.
	IsFunctionArgsContext()
}

type FunctionArgsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionArgsContext() *FunctionArgsContext {
	var p = new(FunctionArgsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = BeanExprParserRULE_functionArgs
	return p
}

func (*FunctionArgsContext) IsFunctionArgsContext() {}

func NewFunctionArgsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionArgsContext {
	var p = new(FunctionArgsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = BeanExprParserRULE_functionArgs

	return p
}

func (s *FunctionArgsContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionArgsContext) CopyFrom(ctx *FunctionArgsContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *FunctionArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionArgsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type FuncArgsContext struct {
	*FunctionArgsContext
}

func NewFuncArgsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FuncArgsContext {
	var p = new(FuncArgsContext)

	p.FunctionArgsContext = NewEmptyFunctionArgsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*FunctionArgsContext))

	return p
}

func (s *FuncArgsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncArgsContext) AllFloatArg() []antlr.TerminalNode {
	return s.GetTokens(BeanExprParserFloatArg)
}

func (s *FuncArgsContext) FloatArg(i int) antlr.TerminalNode {
	return s.GetToken(BeanExprParserFloatArg, i)
}

func (s *FuncArgsContext) AllStringArg() []antlr.TerminalNode {
	return s.GetTokens(BeanExprParserStringArg)
}

func (s *FuncArgsContext) StringArg(i int) antlr.TerminalNode {
	return s.GetToken(BeanExprParserStringArg, i)
}

func (s *FuncArgsContext) AllIntArg() []antlr.TerminalNode {
	return s.GetTokens(BeanExprParserIntArg)
}

func (s *FuncArgsContext) IntArg(i int) antlr.TerminalNode {
	return s.GetToken(BeanExprParserIntArg, i)
}

func (s *FuncArgsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.EnterFuncArgs(s)
	}
}

func (s *FuncArgsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(BeanExprListener); ok {
		listenerT.ExitFuncArgs(s)
	}
}

func (s *FuncArgsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case BeanExprVisitor:
		return t.VisitFuncArgs(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *BeanExprParser) FunctionArgs() (localctx IFunctionArgsContext) {
	this := p
	_ = this

	localctx = NewFunctionArgsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, BeanExprParserRULE_functionArgs)
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

	localctx = NewFuncArgsContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(28)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<BeanExprParserStringArg)|(1<<BeanExprParserFloatArg)|(1<<BeanExprParserIntArg))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	p.SetState(33)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == BeanExprParserT__2 {
		{
			p.SetState(29)
			p.Match(BeanExprParserT__2)
		}
		{
			p.SetState(30)
			_la = p.GetTokenStream().LA(1)

			if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<BeanExprParserStringArg)|(1<<BeanExprParserFloatArg)|(1<<BeanExprParserIntArg))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		p.SetState(35)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}
