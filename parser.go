package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Condition interface {
	Score() bool
}

type AndExpr struct {
	Lhs Condition
	Rhs Condition
}

func (a AndExpr) Score() bool {
	return a.Lhs.Score() && a.Rhs.Score()
}

type OrExpr struct {
	Lhs Condition
	Rhs Condition
}

func (a OrExpr) Score() bool {
	return a.Lhs.Score() || a.Rhs.Score()
}

type NotFunc struct {
	Func Condition
}

func (n NotFunc) Score() bool {
	return !n.Func.Score()
}

type Parser struct {
	lexer        *Lexer
	currentToken Token

	// permit substitution of custom conditions
	// only allowed when parsing check conditions (not
	// custom conditions)
	customConds   bool
	customCondMap map[string]string
}

func (p *Parser) Fatal(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func tokenTypesToStr(tokenTypes []TokenType) string {
	var strTypeArr []string
	for _, num := range tokenTypes {
		strTypeArr = append(strTypeArr, strconv.Itoa(int(num)))
	}
	return strings.Join(strTypeArr, ", ")
}

func (p *Parser) ExpectAny(tokenTypes []TokenType) Token {
	tok := p.lexer.Next()
	if tok.Type == TokenEOF {
		p.Fatal(fmt.Errorf("unexpected end of token stream, expected one of: [%s]",
			tokenTypesToStr(tokenTypes)))
	}

	if !slices.Contains(tokenTypes, tok.Type) {
		p.Fatal(fmt.Errorf("unexpected token %s (type %d), expected one of: [%s] at pos: %d",
			tok.Value, tok.Type, tokenTypesToStr(tokenTypes), p.lexer.pos))
	}

	return tok
}

func (p *Parser) Expect(tokenType TokenType) Token {
	return p.ExpectAny([]TokenType{tokenType})
}

func (p *Parser) ExpectCheckArg() Token {
	return p.ExpectAny([]TokenType{TokenString, TokenIdent})
}

func (p *Parser) NextToken() Token {
	p.currentToken = p.lexer.Next()
	return p.currentToken
}

// OR is top-level expr (has highest op precedence)
func (p *Parser) ParseExpression() Condition {
	lhs := p.ParseAnd()
	for p.NextToken().Type == TokenOr {
		rhs := p.ParseAnd()
		lhs = OrExpr{Lhs: lhs, Rhs: rhs}
	}
	return lhs
}

func (p *Parser) ParseAnd() Condition {
	lhs := p.ParseFactor()
	for p.NextToken().Type == TokenAnd {
		rhs := p.ParseFactor()
		lhs = AndExpr{Lhs: lhs, Rhs: rhs}
	}
	return lhs
}

func (p *Parser) ParseFactor() Condition {
	if p.NextToken().Type == TokenLParen {
		expr := p.ParseExpression()
		return expr
	} else if p.currentToken.Type == TokenIdent {
		if p.customConds && p.currentToken.Value[0] == '$' {
			return p.ParseCondCall()
		}
		return p.ParseFunc()
	}

	p.Fatal(fmt.Errorf("unhandled token: %s (type %d) at pos %d", p.currentToken.Value, p.currentToken.Type, p.lexer.pos))
	return nil // unreachable
}

func (p *Parser) ParseCondCall() Condition {
	customCondName := p.currentToken.Value[1:]
	customCondStr := p.customCondMap[customCondName]

	var args []string
	for p.NextToken().IsLiteral() {
		args = append(args, p.currentToken.Value)
	}

	for argPos, arg := range args {
		customCondStr = strings.ReplaceAll(customCondStr, fmt.Sprintf("$%d", argPos+1), arg)
	}

	return ParseConditionFromString(customCondStr)
}

func (p *Parser) ParseFunc() Condition {
	funcName := p.currentToken.Value
	if len(funcName) <= len("Not") {
		p.Fatal(fmt.Errorf("invalid check name: %s", funcName))
	}

	notFunc := false
	notIdx := len(funcName) - 3
	if funcName[notIdx:] == "Not" {
		notFunc = true
		funcName = funcName[:notIdx]
	}

	var fun Condition
	switch funcName {
	case "PathExists":
		path := p.ExpectCheckArg()
		fun = PathExists{Path: path.Value}
	case "FileContains":
		file := p.ExpectCheckArg()
		value := p.ExpectCheckArg()
		fun = FileContains{File: file.Value, Value: value.Value}
	default:
		p.Fatal(fmt.Errorf("unrecognized check type: %s", funcName))
	}

	if notFunc {
		return NotFunc{Func: fun}
	}
	return fun
}

func ParseConditionFromString(source string) Condition {
	lexer := Lexer{source: source, pos: 0}
	parser := Parser{lexer: &lexer, customConds: false}
	return parser.ParseExpression()
}

func ParseConditionFromStringWith(source string, customCondMap map[string]string) Condition {
	lexer := Lexer{source: source, pos: 0}
	parser := Parser{lexer: &lexer, customConds: true, customCondMap: customCondMap}
	return parser.ParseExpression()
}
