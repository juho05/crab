package interpreter

import (
	"fmt"
)

type parser struct {
	tokens  []Token
	current int
	lines   [][]rune
	errors  []error
}

func Parse(tokens []Token, lines [][]rune) ([]Stmt, []error) {
	parser := &parser{
		tokens: tokens,
		lines:  lines,
		errors: make([]error, 0),
	}
	return parser.parse()
}

func (p *parser) parse() ([]Stmt, []error) {
	statements := make([]Stmt, 0)
	for p.peek().Type != EOF {
		statements = append(statements, p.declaration(false))
	}
	return statements, p.errors
}

func (p *parser) declaration(allowNonDeclarationStatements bool) Stmt {
	var stmt Stmt
	var err error
	if p.match(VAR) {
		stmt, err = p.varDecl()
	} else if p.match(FUNC) {
		stmt, err = p.funcDecl()
	} else if allowNonDeclarationStatements {
		stmt, err = p.statement()
	}

	if stmt == nil && err == nil {
		err = p.newError(fmt.Sprintf("Unexpected token '%s'", p.peek().Lexeme))
	}

	if err != nil {
		p.errors = append(p.errors, err)
		p.synchronize()
	}
	return stmt
}

func (p *parser) varDecl() (Stmt, error) {
	if !p.match(IDENTIFIER) {
		return nil, p.newError("Expect identifier after 'var' keyword.")
	}

	names := make([]Token, 1)
	names[0] = p.previous()

	for p.match(COMMA) {
		if !p.match(IDENTIFIER) {
			return nil, p.newError("Expect identifier after ','.")
		}
		names = append(names, p.previous())
	}

	var expr Expr
	var err error
	var operator Token
	if p.match(EQUAL) {
		operator = p.previous()
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if !p.match(SEMICOLON) {
		return nil, p.newError("Missing semicolon.")
	}

	return &StmtVarDecl{
		Operator: operator,
		Names:    names,
		Expr:     expr,
	}, nil
}

func (p *parser) funcDecl() (Stmt, error) {
	if !p.match(IDENTIFIER) {
		return nil, p.newError("Expect identifier after 'func' keyword.")
	}

	name := p.previous()

	if !p.match(OPEN_PAREN) {
		return nil, p.newError("Expect '(' after function name.")
	}

	parameters := make([]string, 0)
	for p.peek().Type != CLOSE_PAREN {
		if !p.match(IDENTIFIER) {
			return nil, p.newError("Invalid parameter name.")
		}
		parameters = append(parameters, p.previous().Lexeme)
		if p.peek().Type == CLOSE_PAREN {
			break
		}
		if !p.match(COMMA) {
			return nil, p.newError("Expect ',' between parameters.")
		}
	}

	if !p.match(CLOSE_PAREN) {
		return nil, p.newError("Expect ')' after function parameter list.")
	}

	returnValueCount := 0
	if p.peek().Type == NUMBER {
		num := p.peek()
		if num.Lexeme != "0" && num.Lexeme != "1" && num.Lexeme != "2" && num.Lexeme != "3" && num.Lexeme != "4" {
			return nil, p.newError("Only 0, 1, 2, 3 or 4 return values are allowed.")
		}
		returnValueCount = int(num.Literal.(float64))
		p.current++
	}

	throws := false
	if p.match(THROWS) {
		throws = true
	}

	if !p.match(OPEN_BRACE) {
		return nil, p.newError("Expect block after function signature.")
	}

	block, err := p.block()
	if err != nil {
		return nil, err
	}

	return &StmtFuncDecl{
		Name:             name,
		Body:             block,
		Parameters:       parameters,
		ReturnValueCount: returnValueCount,
		Throws:           throws,
	}, nil
}

func (p *parser) statement() (Stmt, error) {
	if p.match(OPEN_BRACE) {
		return p.block()
	}
	if p.match(IF) {
		return p.ifStmt()
	}
	if p.match(WHILE) {
		return p.whileLoop()
	}
	if p.match(FOR) {
		return p.forLoop()
	}
	if p.match(BREAK, CONTINUE) {
		return p.loopControl()
	}
	if p.match(RETURN) {
		return p.returnStmt()
	}
	if p.match(THROW) {
		return p.throwStmt()
	}
	if p.match(TRY) {
		return p.tryStmt()
	}
	return p.expressionStmt()
}

func (p *parser) ifStmt() (Stmt, error) {
	keyword := p.previous()
	if !p.match(OPEN_PAREN) {
		return nil, p.newError("Expect '(' after 'if'.")
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(CLOSE_PAREN) {
		return nil, p.newError("Expect ')' after if condition.")
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBody Stmt
	if p.match(ELSE) {
		elseBody, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &StmtIf{
		Keyword:   keyword,
		Condition: condition,
		Body:      body,
		ElseBody:  elseBody,
	}, nil
}

func (p *parser) whileLoop() (Stmt, error) {
	keyword := p.previous()
	if !p.match(OPEN_PAREN) {
		return nil, p.newError("Expect '(' after 'while'.")
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(CLOSE_PAREN) {
		return nil, p.newError("Expect ')' after while condition.")
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &StmtWhile{
		Keyword:   keyword,
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *parser) forLoop() (Stmt, error) {
	keyword := p.previous()
	if !p.match(OPEN_PAREN) {
		return nil, p.newError("Expect '(' after 'while'.")
	}

	var initializer Stmt
	var err error
	if !p.match(SEMICOLON) {
		if p.match(VAR) {
			initializer, err = p.varDecl()
		} else {
			initializer, err = p.expressionStmt()
		}
	}
	if err != nil {
		return nil, err
	}
	if initializer == nil {
		initializer = &StmtExpression{
			Expr: &ExprLiteral{},
		}
	}

	var condition Expr
	if p.peek().Type != SEMICOLON {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if condition == nil {
		condition = &ExprLiteral{
			Value: true,
		}
	}

	if !p.match(SEMICOLON) {
		return nil, p.newError("Expect ';' after for loop condition.")
	}

	var increment Expr
	if p.peek().Type != CLOSE_PAREN {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if increment == nil {
		increment = &ExprLiteral{}
	}

	if !p.match(CLOSE_PAREN) {
		return nil, p.newError("Expect ')' after for loop clauses.")
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &StmtBlock{
		Statements: []Stmt{&StmtFor{
			Keyword:     keyword,
			Initializer: initializer,
			Condition:   condition,
			Increment:   increment,
			Body:        body,
		}},
	}, nil
}

func (p *parser) loopControl() (Stmt, error) {
	keyword := p.previous()
	if !p.match(SEMICOLON) {
		return nil, p.newError("Missing semicolon.")
	}
	return &StmtLoopControl{
		Keyword: keyword,
	}, nil
}

func (p *parser) returnStmt() (Stmt, error) {
	keyword := p.previous()

	values := make([]Expr, 0)

	for p.peek().Type != SEMICOLON {
		expr, err := p.conditional()
		if err != nil {
			return nil, err
		}
		values = append(values, expr)
		if !p.match(COMMA) {
			break
		}
	}

	if !p.match(SEMICOLON) {
		return nil, p.newError("Missing semicolon.")
	}

	return &StmtReturn{
		Keyword: keyword,
		Values:  values,
	}, nil
}

func (p *parser) throwStmt() (Stmt, error) {
	keyword := p.previous()

	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(SEMICOLON) {
		return nil, p.newError("Missing semicolon.")
	}

	return &StmtThrow{
		Keyword: keyword,
		Value:   expr,
	}, nil
}

func (p *parser) tryStmt() (Stmt, error) {
	keyword := p.previous()
	if !p.match(OPEN_BRACE) {
		return nil, p.newError("Expect '{' after 'try'.")
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	if !p.match(CATCH) {
		return nil, p.newError("Expect 'catch' after try body.")
	}

	var exceptionName Token
	if p.match(OPEN_PAREN) {
		if !p.match(IDENTIFIER) {
			return nil, p.newError("Expect exception name.")
		}
		exceptionName = p.previous()
		if !p.match(CLOSE_PAREN) {
			return nil, p.newError("Expect ')' after exception name.")
		}
	}

	if !p.match(OPEN_BRACE) {
		return nil, p.newError("Expect '{' after 'catch'.")
	}
	catchBody, err := p.block()
	if err != nil {
		return nil, err
	}

	return &StmtTry{
		Keyword:       keyword,
		Body:          body,
		CatchBody:     catchBody,
		ExceptionName: exceptionName,
	}, nil
}

func (p *parser) expressionStmt() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(SEMICOLON) {
		return nil, p.newError("Missing semicolon.")
	}

	return &StmtExpression{
		Expr: expr,
	}, nil
}

func (p *parser) block() (Stmt, error) {
	openBrace := p.previous()
	statements := make([]Stmt, 0)

	for p.peek().Type != CLOSE_BRACE && p.peek().Type != EOF {
		stmt := p.declaration(true)
		statements = append(statements, stmt)
	}

	if !p.match(CLOSE_BRACE) {
		return nil, p.newErrorAt("Block never closed.", openBrace)
	}

	return &StmtBlock{
		Statements: statements,
	}, nil
}

func (p *parser) expression() (Expr, error) {
	return p.assign()
}

func (p *parser) assign() (Expr, error) {
	exprs := make([]Expr, 1)
	var err error
	exprs[0], err = p.conditional()
	if err != nil {
		return nil, err
	}

	isAssign := false
	if p.match(COMMA) {
		isAssign = true
		expr, err := p.conditional()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}
	if isAssign && !p.match(EQUAL, PLUS_EQUAL, MINUS_EQUAL, ASTERISK_EQUAL, ASTERISK_ASTERISK_EQUAL, SLASH_EQUAL, PERCENT_EQUAL) {
		return nil, p.newError("Expect assignment operator after identifier list.")
	}

	if isAssign || p.match(EQUAL, PLUS_EQUAL, MINUS_EQUAL, ASTERISK_EQUAL, ASTERISK_ASTERISK_EQUAL, SLASH_EQUAL, PERCENT_EQUAL) {
		operator := p.previous()
		assignees := make([]Expr, 0)
		for _, expr := range exprs {
			if v, ok := expr.(*ExprVariable); ok {
				assignees = append(assignees, v)
			} else if s, ok := expr.(*ExprSubscript); ok {
				assignees = append(assignees, s)
			} else {
				return nil, p.newErrorAt("Can only assign to variables.", operator)
			}
		}

		right, err := p.conditional()
		if err != nil {
			return nil, err
		}

		tokenType := EQUAL
		switch operator.Type {
		case PLUS_EQUAL:
			tokenType = PLUS
		case MINUS_EQUAL:
			tokenType = MINUS
		case ASTERISK_EQUAL:
			tokenType = ASTERISK
		case ASTERISK_ASTERISK_EQUAL:
			tokenType = ASTERISK_ASTERISK
		case SLASH_EQUAL:
			tokenType = SLASH
		case PERCENT_EQUAL:
			tokenType = PERCENT
		}

		if tokenType != EQUAL {
			if len(assignees) > 1 {
				return nil, p.newErrorAt("Multi value assignment only allowed for '=' operator.", operator)
			}
			right = &ExprBinary{
				Operator: Token{
					Line:   operator.Line,
					Type:   tokenType,
					Column: operator.Column,
					Lexeme: operator.Lexeme,
				},
				Left:  exprs[0],
				Right: right,
			}
		}

		return &ExprAssign{
			Operator:  operator,
			Assignees: assignees,
			Expr:      right,
		}, nil
	}

	return exprs[0], nil
}

func (p *parser) conditional() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(QUESTION_MARK) {
		operator1 := p.previous()
		center, err := p.conditional()
		if err != nil {
			return nil, err
		}
		if !p.match(COLON) {
			return nil, p.newError("Expect ':' after '?'.")
		}
		operator2 := p.previous()
		right, err := p.conditional()
		if err != nil {
			return nil, err
		}
		expr = &ExprTernary{
			Left:      expr,
			Operator1: operator1,
			Center:    center,
			Operator2: operator2,
			Right:     right,
		}
	}

	return expr, nil
}

func (p *parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR, XOR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &ExprLogical{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) and() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &ExprLogical{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(EQUAL_EQUAL, BANG_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(LESS, LESS_EQUAL, GREATER, GREATER_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(PLUS, MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) factor() (Expr, error) {
	expr, err := p.power()
	if err != nil {
		return nil, err
	}

	for p.match(ASTERISK, SLASH, PERCENT) {
		operator := p.previous()
		right, err := p.power()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) power() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(ASTERISK_ASTERISK) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ExprBinary{
			Operator: operator,
			Left:     expr,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ExprUnary{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.postfix()
}

func (p *parser) postfix() (Expr, error) {
	expr, err := p.subscriptOrCall()
	if err != nil {
		return nil, err
	}
	if p.match(PLUS_PLUS, MINUS_MINUS) {
		operator := p.previous()
		if _, ok := expr.(*ExprVariable); !ok {
			if _, ok := expr.(*ExprSubscript); !ok {
				return nil, p.newErrorAt("Can only increment/decrement variables.", operator)
			}
		}
		tokenType := PLUS
		if operator.Type == MINUS_MINUS {
			tokenType = MINUS
		}
		expr = &ExprAssign{
			Assignees: []Expr{expr},
			Expr: &ExprBinary{
				Operator: Token{
					Line:   operator.Line,
					Type:   tokenType,
					Lexeme: operator.Lexeme,
					Column: operator.Column,
				},
				Left: expr,
				Right: &ExprLiteral{
					Value: 1.0,
				},
			},
		}
	}
	return expr, nil
}

func (p *parser) subscriptOrCall() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	for p.match(OPEN_BRACKET, OPEN_PAREN) {
		token := p.previous()

		if token.Type == OPEN_BRACKET {
			subscript, err := p.expression()
			if err != nil {
				return nil, err
			}

			if !p.match(CLOSE_BRACKET) {
				return nil, p.newError("Expect ')' after argument list.")
			}

			expr = &ExprSubscript{
				OpenBracket: token,
				Object:      expr,
				Subscript:   subscript,
			}
		} else if token.Type == OPEN_PAREN {
			args := make([]Expr, 0)
			for p.peek().Type != CLOSE_PAREN {
				arg, err := p.conditional()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				if p.peek().Type == CLOSE_PAREN {
					break
				}
				if !p.match(COMMA) {
					return nil, p.newError("Expect ',' between arguments.")
				}
			}
			if !p.match(CLOSE_PAREN) {
				return nil, p.newError("Expect ')' after argument list.")
			}
			expr = &ExprCall{
				OpenParen: token,
				Callee:    expr,
				Args:      args,
			}
		}
	}
	return expr, nil
}

func (p *parser) primary() (Expr, error) {
	if p.match(NUMBER, STRING, TRUE, FALSE) {
		return &ExprLiteral{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(IDENTIFIER) {
		return &ExprVariable{
			Name: p.previous(),
		}, nil
	}

	if p.match(OPEN_PAREN) {
		openingParen := p.previous()
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.match(CLOSE_PAREN) {
			return nil, p.newErrorAt("Parenthesis never closed.", openingParen)
		}
		return &ExprGrouping{
			Expr: expr,
		}, nil
	}

	if p.match(OPEN_BRACKET) {
		return p.list()
	}

	return nil, p.newError(fmt.Sprintf("Unexpected token '%s'", p.peek().Lexeme))
}

func (p *parser) list() (Expr, error) {
	openingBracket := p.previous()

	values := make([]Expr, 0)

	for p.peek().Type != CLOSE_BRACKET {
		expr, err := p.conditional()
		if err != nil {
			return nil, err
		}
		values = append(values, expr)

		if !p.match(COMMA) {
			break
		}
	}

	if !p.match(CLOSE_BRACKET) {
		return nil, p.newErrorAt("Bracket never closed.", openingBracket)
	}

	return &ExprList{
		OpenBracket: openingBracket,
		Values:      values,
	}, nil
}

func (p *parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.peek().Type == t {
			p.current++
			return true
		}
	}
	return false
}

func (p *parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *parser) peek() Token {
	return p.tokens[p.current]
}

func (p *parser) synchronize() {
	if p.peek().Type == EOF {
		return
	}
	p.current++
	for p.peek().Type != EOF {
		switch p.peek().Type {
		case SEMICOLON:
			p.current++
			return
		case VAR, FUNC, IF, WHILE, FOR:
			return
		}
		p.current++
	}
}

type ParseError struct {
	Token   Token
	Message string
	Line    []rune
}

func (p ParseError) Error() string {
	return generateErrorText(p.Message, p.Line, p.Token.Line, p.Token.Column, p.Token.Column+len([]rune(p.Token.Lexeme)))
}

func (p *parser) newError(message string) error {
	return ParseError{
		Token:   p.peek(),
		Message: message,
		Line:    p.lines[p.peek().Line],
	}
}

func (p *parser) newErrorAt(message string, token Token) error {
	return ParseError{
		Token:   token,
		Message: message,
		Line:    p.lines[token.Line],
	}
}
