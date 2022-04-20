package interpreter

type ExprVisitor interface {
	VisitLiteral(expr ExprLiteral) (any, error)
	VisitGrouping(expr ExprGrouping) (any, error)
	VisitUnary(expr ExprUnary) (any, error)
	VisitBinary(expr ExprBinary) (any, error)
}

type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

type ExprLiteral struct {
	Value any
}

func (e ExprLiteral) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteral(e)
}

type ExprGrouping struct {
	Expr Expr
}

func (e ExprGrouping) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGrouping(e)
}

type ExprUnary struct {
	Operator Token
	Right    Expr
}

func (e ExprUnary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnary(e)
}

type ExprBinary struct {
	Operator Token
	Left     Expr
	Right    Expr
}

func (e ExprBinary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinary(e)
}
