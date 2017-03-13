package parser

import (
	"testing"

	"github.com/dustinlindquist/monkey/ast"
	"github.com/dustinlindquist/monkey/lexer"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestLetStatements(t *testing.T) {
	var input = `let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	var l = lexer.New(input)
	var p = New(l)

	var program = p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.statements does not contain 3 statements")
	}

	var tests = []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tc := range tests {
		var stmt = program.Statements[i]
		if !checkLetStatement(t, stmt, tc.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	var input = `
	return 5;
	return 10;
	return 993322;
	`

	var l = lexer.New(input)
	var p = New(l)

	var program = p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got: %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		var returnStmt, ok = stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got:%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got: %q", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpressions(t *testing.T) {
	var input = "foobar;"

	var l = lexer.New(input)
	var p = New(l)
	var program = p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong ammount of statements: %d", len(program.Statements))
	}
	var stmt, ok = program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] isn't ast.ExpressionStatement, instead: %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier, instead: %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not foobar, got: %s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not foobar, got: %s", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	var input = `5;`
	var l = lexer.New(input)
	var p = New(l)
	var program = p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has wrong ammount of statements: %d", len(program.Statements))
	}

	var stmt, ok = program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, go: %T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral, got: %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("literal.Value not 5, got: %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not 5, got: %s", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	var prefixTests = []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15 },
	}

	for _, tc := range prefixTests {
		var l= lexer.New(tc.input)
		var p= New(l)
		var program= p.ParseProgram()
		checkParserErrors(t, p)

		assert.Len(t, program.Statements, 1)

		var stmt, ok= program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		assert.True(t, ok)

		assert.Equal(t, exp.Operator, tc.operator)

		checkIntegerLiteral(t, exp.Right, tc.integerValue)
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	var errors = p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err)
	}
	t.FailNow()
}

func checkLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got = %q", s.TokenLiteral())
		return false
	}

	var letStmt, ok = s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got: %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not: '%s' got: '%s'", name, letStmt.Name.Value)
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not: '%s' got: '%s'", name, letStmt.Name)
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	var infixTests = []struct {
		input string
		leftValue int64
		operator string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tc := range infixTests {
		var l = lexer.New(tc.input)
		var p = New(l)
		var program = p.ParseProgram()
		checkParserErrors(t, p)

		assert.Len(t, program.Statements, 1)

		var stmt, ok = program.Statements[0].(*ast.ExpressionStatement)
		assert.True(t, ok)

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		assert.True(t, ok)

		checkIntegerLiteral(t, exp.Left, tc.leftValue)
		assert.Equal(t, tc.operator, exp.Operator)
		checkIntegerLiteral(t, exp.Right, tc.rightValue)

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	var tests = []struct {
		input string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tc := range tests {
		var l = lexer.New(tc.input)
		var p = New(l)
		var program = p.ParseProgram()
		checkParserErrors(t, p)

		var actual = program.String()
		assert.Equal(t, tc.expected, actual)
	}
}

func checkIntegerLiteral(t *testing.T, il ast.Expression, value int64) {
	var integ, ok = il.(*ast.IntegerLiteral)
	assert.True(t, ok)
	assert.Equal(t, value, integ.Value)
	assert.Equal(t, fmt.Sprintf("%d", value), integ.TokenLiteral())
}
