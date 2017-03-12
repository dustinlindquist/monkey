package parser

import (
	"testing"

	"github.com/dustinlindquist/monkey/ast"
	"github.com/dustinlindquist/monkey/lexer"
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
