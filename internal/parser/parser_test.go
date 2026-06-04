package parser

import (
	"fmt"
	"testing"

	"monkey/internal/ast"
	"monkey/internal/lexer"
)

func TestParsingLetStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedValue any
	}{
		{
			"let x = 5;", "x", 5,
		},
		{
			"let y = true;", "y", true,
		},
		{
			"let foobar = y;", "foobar", "y",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain 1 statements. got=%d",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("stmt not ast.LetStatement. got=%T", program.Statements[0])
		}

		if !testLetStatement(t, stmt, tt.expectedIdent) {
			return
		}

		val := stmt.Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not let. got=%q", s.TokenLiteral())

		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement. got=%T", s)

		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)

		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())

		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func TestParsingReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"return 5;",
			5,
		},
		{
			"return 10;",
			10,
		},
		{
			"return 993322;",
			993322,
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain 1 statements. got=%d",
				len(program.Statements),
			)
		}

		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)

			if !ok {
				t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)

				continue
			}

			if returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
			}

			if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expected) {
				return
			}
		}
	}
}

func TestParsingIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestParsingIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"+15;", "+", 15},
		{"!false;", "!", false},
		{"!true;", "!", true},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d,%q\n",
				1, len(program.Statements), program.Statements,
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf(
				"stmt is not ast.PrefixExpression. got=%T", stmt.Expression,
			)
		}

		if exp.Operator != tt.operator {
			t.Fatalf(
				"exp.Operator is not %s. got=%s", exp.Operator, tt.operator,
			)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(
	t *testing.T, il ast.Expression, value int64,
) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)

		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integer.Value)

		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf(
			"integer.TokenLiteral() %d. got=%s",
			value, integer.TokenLiteral(),
		)

		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)

		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)

		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf(
			"ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral(),
		)

		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected any,
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))

	case int64:
		return testIntegerLiteral(t, exp, v)

	case string:
		return testIdentifier(t, exp, v)

	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)

	return false
}

func testInfixExpression(
	t *testing.T, exp ast.Expression,
	left any, operator string, right any,
) bool {
	expected, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)

		return false
	}

	if !testLiteralExpression(t, expected.Left, left) {
		return false
	}

	if expected.Operator != operator {
		t.Errorf(
			"exp.Operator is not '%s'. got=%q", operator, expected.Operator,
		)

		return false
	}

	if !testLiteralExpression(t, expected.Right, right) {
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"6 + 5;", 6, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"1 < 5;", 1, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 7;", 5, "!=", 7},
		{"true == true;", true, "==", true},
		{"true == false;", true, "==", false},
		{"false == false;", false, "==", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		testInfixExpression(
			t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue,
		)
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
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
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == false",
			"((3 < 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a,b,1,2 * 3,4 + 5,add(6,7 * 8))",
			"add(a,b,1,(2 * 3),(4 + 5),add(6,(7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])),(b[1]),(2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("for %s, expected=%q, got=%q", tt.input, tt.expected, actual)
		}
	}
}

func testBooleanLiteral(
	t *testing.T, il ast.Expression, value bool,
) bool {
	boolean, ok := il.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("il not *ast.BooleanLiteral. got=%T", il)

		return false
	}

	if boolean.Value != value {
		t.Errorf("integer.Value not %v. got=%v", value, boolean.Value)

		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%v", value) {
		t.Errorf(
			"integer.TokenLiteral() %v. got=%s",
			value, boolean.TokenLiteral(),
		)

		return false
	}

	return true
}

func TestParsingBooleanLiteralExpression(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{
			"true;",
			true,
		},
		{
			"false;",
			false,
		},
	}

	for _, tc := range cases {
		l := lexer.New(tc.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		testBooleanLiteral(t, stmt.Expression, tc.expected)
	}
}

func TestParsingIfExpression(t *testing.T) {
	input := "if (x < y) { x }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n", program.Statements[0])
	}

	ife, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Expression statement. got=%T\n", ife)
	}

	if !testInfixExpression(t, ife.Condition, "x", "<", "y") {
		return
	}

	if len(ife.Then.Statements) != 1 {
		t.Errorf("then is not 1 statement. got=%d\n", len(ife.Then.Statements))

		return
	}

	then, ok := ife.Then.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T\n", ife.Then.Statements[0])
	}

	if !testIdentifier(t, then.Expression, "x") {
		return
	}

	if ife.Else != nil {
		t.Errorf("Else.Statements is not nil. got=%+v", ife.Else)
	}
}

func TestParsingIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T\n", program.Statements[0])
	}

	ife, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Expression statement. got=%T\n", ife)
	}

	if !testInfixExpression(t, ife.Condition, "x", "<", "y") {
		return
	}

	if len(ife.Then.Statements) != 1 {
		t.Errorf("then is not 1 statement. got=%d\n", len(ife.Then.Statements))

		return
	}

	then, ok := ife.Then.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T\n", ife.Then.Statements[0])
	}

	if !testIdentifier(t, then.Expression, "x") {
		return
	}

	if ife.Else == nil || len(ife.Else.Statements) == 0 {
		t.Fatalf("Else.Statements is empty. got=%+v", ife.Else)
	}

	if len(ife.Else.Statements) != 1 {
		t.Errorf("Else.Statements is more than 1. got=%+v", ife.Else.Statements)
	}

	elseExp, ok := ife.Else.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf(
			"Else.Statements[0] is not ast.ExpressionStatement. got=%T", ife.Else.Statements[0],
		)
	}

	if !testIdentifier(t, elseExp.Expression, "y") {
		return
	}
}

func TestParsingFunctionLiteral(t *testing.T) {
	input := "fn(x,y){x+y}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt)
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf(
			"function literal parameters wrong. want 2, got=%d, %v", len(function.Parameters),
			function.Parameters,
		)
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf(
			"function.Body.Statements does have 1 statement. got=%d",
			len(function.Body.Statements),
		)
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", bodyStmt)
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestParsingFunctionParameters(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "fn(){}",
			expected: []string{},
		},
		{
			input:    "fn(x){}",
			expected: []string{"x"},
		},
		{
			input:    "fn(x,y,z){}",
			expected: []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt)
		}

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
		}

		if len(function.Parameters) != len(tt.expected) {
			t.Fatalf(
				"function literal parameters wrong. want %d, got=%d, %v",
				len(tt.expected), len(function.Parameters),
				function.Parameters,
			)
		}

		for i, param := range function.Parameters {
			testLiteralExpression(t, param, tt.expected[i])
		}
	}
}

func TestParsingCallExpression(t *testing.T) {
	input := "add(1,2*3,4+5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not an ExpressionStatement. got=%T", stmt)
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", exp)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestParsingCallExpressionArguments(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "add()",
			expected: []string{},
		},
		{
			input:    "add(1)",
			expected: []string{"1"},
		},
		{
			input:    "add(1, 4 + 5,3)",
			expected: []string{"1", "(4 + 5)", "3"},
		},
		{
			input:    "add(x, 1 + 2)",
			expected: []string{"x", "(1 + 2)"},
		},
		{
			input:    "add(x, 1 + 2 * (1 + 1))",
			expected: []string{"x", "(1 + (2 * (1 + 1)))"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", stmt)
		}

		function, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
		}

		if len(function.Arguments) != len(tt.expected) {
			t.Fatalf(
				"call expression arguments wrong. want %d, got=%d, %v",
				len(tt.expected), len(function.Arguments),
				function.Arguments,
			)
		}

		for i, param := range function.Arguments {
			if param.String() != tt.expected[i] {
				t.Errorf(
					"argument not matching. want=%s, got=%s",
					tt.expected[i], param.String(),
				)
			}
		}
	}
}

func TestParsingStringLiteralExpression(t *testing.T) {
	input := `"hello world!";`
	expected := "hello world!"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement not ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != expected {
		t.Errorf("literal.Value not %q. got=%q", expected, literal.Value)
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}
