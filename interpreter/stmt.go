package interpreter

type StmtVisitor interface {
	VisitExpression(stmt *StmtExpression) error
	VisitBlock(stmt *StmtBlock) error
	VisitVarDecl(stmt *StmtVarDecl) error
	VisitFuncDecl(stmt *StmtFuncDecl) error
	VisitIf(stmt *StmtIf) error
	VisitWhile(stmt *StmtWhile) error
	VisitFor(stmt *StmtFor) error
	VisitLoopControl(stmt *StmtLoopControl) error
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtExpression struct {
	Expr Expr
}

func (s *StmtExpression) Accept(visitor StmtVisitor) error {
	return visitor.VisitExpression(s)
}

type StmtBlock struct {
	Statements []Stmt
}

func (s *StmtBlock) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlock(s)
}

type StmtVarDecl struct {
	Name Token
	Expr Expr
}

func (s *StmtVarDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarDecl(s)
}

type StmtFuncDecl struct {
	Name       Token
	Body       Stmt
	Parameters []string
}

func (s *StmtFuncDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitFuncDecl(s)
}

type StmtIf struct {
	Condition Expr
	Body      Stmt
	ElseBody  Stmt
}

func (s *StmtIf) Accept(visitor StmtVisitor) error {
	return visitor.VisitIf(s)
}

type StmtWhile struct {
	Condition Expr
	Body      Stmt
}

func (s *StmtWhile) Accept(visitor StmtVisitor) error {
	return visitor.VisitWhile(s)
}

type StmtFor struct {
	Initializer Stmt
	Condition   Expr
	Increment   Expr
	Body        Stmt
}

func (s *StmtFor) Accept(visitor StmtVisitor) error {
	return visitor.VisitFor(s)
}

type StmtLoopControl struct {
	Keyword Token
}

func (s *StmtLoopControl) Accept(visitor StmtVisitor) error {
	return visitor.VisitLoopControl(s)
}
