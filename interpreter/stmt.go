package interpreter

type StmtVisitor interface {
	VisitExpression(stmt *StmtExpression) error
	VisitBlock(stmt *StmtBlock) error
	VisitVarDecl(stmt *StmtVarDecl) error
	VisitFuncDecl(stmt *StmtFuncDecl) error
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
