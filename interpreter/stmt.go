package interpreter

type StmtVisitor interface {
	VisitBlock(stmt StmtBlock) error
	VisitVarDecl(stmt StmtVarDecl) error
	VisitFuncDecl(stmt StmtFuncDecl) error
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtBlock struct {
	Statements []Stmt
}

func (s StmtBlock) Accept(visitor StmtVisitor) error {
	return visitor.VisitBlock(s)
}

type StmtVarDecl struct {
	Name Token
	Expr Expr
}

func (s StmtVarDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarDecl(s)
}

type StmtFuncDecl struct {
	Name Token
	Body Stmt
}

func (s StmtFuncDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitFuncDecl(s)
}