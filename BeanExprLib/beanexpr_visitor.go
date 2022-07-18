// Code generated from E:/Go/gofks/BeanExprLib\BeanExpr.g4 by ANTLR 4.10.1. DO NOT EDIT.

package BeanExprLib // BeanExpr

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by BeanExprParser.
type BeanExprVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by BeanExprParser#start.
	VisitStart(ctx *StartContext) interface{}

	// Visit a parse tree produced by BeanExprParser#methodCall.
	VisitMethodCall(ctx *MethodCallContext) interface{}

	// Visit a parse tree produced by BeanExprParser#FuncCall.
	VisitFuncCall(ctx *FuncCallContext) interface{}

	// Visit a parse tree produced by BeanExprParser#FuncArgs.
	VisitFuncArgs(ctx *FuncArgsContext) interface{}
}
