package interpreter

type StmtVisitor interface {
	VisitVarDecl(stmt StmtVarDecl) error
}

type Stmt interface {
	Accept(visitor StmtVisitor) error
}

type StmtVarDecl struct {
	Name Token
	Expr Expr
}

func (s StmtVarDecl) Accept(visitor StmtVisitor) error {
	return visitor.VisitVarDecl(s)
}
