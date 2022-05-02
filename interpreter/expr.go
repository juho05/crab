package interpreter

type ExprVisitor interface {
	VisitLiteral(expr *ExprLiteral) (any, error)
	VisitVariable(expr *ExprVariable) (any, error)
	VisitCall(expr *ExprCall) (any, error)
	VisitSubscript(expr *ExprSubscript) (any, error)
	VisitGrouping(expr *ExprGrouping) (any, error)
	VisitList(expr *ExprList) (any, error)
	VisitUnary(expr *ExprUnary) (any, error)
	VisitBinary(expr *ExprBinary) (any, error)
	VisitLogical(expr *ExprLogical) (any, error)
	VisitTernary(expr *ExprTernary) (any, error)
	VisitAssign(expr *ExprAssign) (any, error)
	VisitAnonymousFunction(expr *ExprAnonymousFunction) (any, error)
}

type Expr interface {
	Accept(visitor ExprVisitor) (any, error)
}

type ExprLiteral struct {
	Value any
}

func (e *ExprLiteral) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLiteral(e)
}

type ExprVariable struct {
	Name         Token
	NestingLevel int
}

func (e *ExprVariable) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitVariable(e)
}

type ExprCall struct {
	OpenParen Token
	Callee    Expr
	Args      []Expr
}

func (e *ExprCall) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitCall(e)
}

type ExprSubscript struct {
	OpenBracket Token
	Object      Expr
	Subscript   Expr
}

func (e *ExprSubscript) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitSubscript(e)
}

type ExprGrouping struct {
	Expr Expr
}

func (e *ExprGrouping) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitGrouping(e)
}

type ExprList struct {
	OpenBracket Token
	Values      []Expr
}

func (e *ExprList) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitList(e)
}

type ExprUnary struct {
	Operator Token
	Right    Expr
}

func (e *ExprUnary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitUnary(e)
}

type ExprBinary struct {
	Operator Token
	Left     Expr
	Right    Expr
}

func (e *ExprBinary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitBinary(e)
}

type ExprLogical struct {
	Operator Token
	Left     Expr
	Right    Expr
}

func (e *ExprLogical) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitLogical(e)
}

type ExprTernary struct {
	Left      Expr
	Operator1 Token
	Center    Expr
	Operator2 Token
	Right     Expr
}

func (e *ExprTernary) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitTernary(e)
}

type ExprAssign struct {
	Operator  Token
	Assignees []Expr
	Expr      Expr
}

func (e *ExprAssign) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAssign(e)
}

type ExprAnonymousFunction struct {
	Keyword          Token
	Body             Stmt
	Parameters       []string
	ReturnValueCount int
	Throws           bool
}

func (e *ExprAnonymousFunction) Accept(visitor ExprVisitor) (any, error) {
	return visitor.VisitAnonymousFunction(e)
}
