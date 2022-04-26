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
	VisitReturn(stmt *StmtReturn) error
	VisitThrow(stmt *StmtThrow) error
	VisitTry(stmt *StmtTry) error
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
	Operator Token
	Names    []Token
	Expr     Expr
}

func (s *StmtVarDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarDecl(s)
}

type StmtFuncDecl struct {
	Name             Token
	Body             Stmt
	Parameters       []string
	ReturnValueCount int
	Throws           bool
}

func (s *StmtFuncDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitFuncDecl(s)
}

type StmtIf struct {
	Keyword   Token
	Condition Expr
	Body      Stmt
	ElseBody  Stmt
}

func (s *StmtIf) Accept(visitor StmtVisitor) error {
	return visitor.VisitIf(s)
}

type StmtWhile struct {
	Keyword   Token
	Condition Expr
	Body      Stmt
}

func (s *StmtWhile) Accept(visitor StmtVisitor) error {
	return visitor.VisitWhile(s)
}

type StmtFor struct {
	Keyword     Token
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

type StmtReturn struct {
	Keyword Token
	Values  []Expr
}

func (s *StmtReturn) Accept(visitor StmtVisitor) error {
	return visitor.VisitReturn(s)
}

type StmtThrow struct {
	Keyword Token
	Value   Expr
}

func (s *StmtThrow) Accept(visitor StmtVisitor) error {
	return visitor.VisitThrow(s)
}

type StmtTry struct {
	Keyword       Token
	Body          Stmt
	CatchBody     Stmt
	ExceptionName Token
}

func (s *StmtTry) Accept(visitor StmtVisitor) error {
	return visitor.VisitTry(s)
}
