package ast

type Kind int

const (
	KindNumber Kind = iota
	KindBool
	KindSymbol
	KindTriple
	KindFunc1
	KindBinOp
	KindIfStmt
)

type Node interface {
	Kind() Kind
}

// Number is a literal number.
type Number float64

func (Number) Kind() Kind { return KindNumber }

var _ Node = (*Number)(nil)

// Bool is a literal boolean value.
type Bool bool

func (Bool) Kind() Kind { return KindBool }

var _ Node = (*Bool)(nil)

// Symbol is one of "x", "y" or "rnd".
//
// TODO: Or make a separate Node for each of these?
type Symbol string

func (Symbol) Kind() Kind { return KindSymbol }

var _ Node = (*Symbol)(nil)

// Triple is a three-element Node tuple.
type Triple struct {
	A, B, C Node
}

func (Triple) Kind() Kind { return KindTriple }

var _ Node = (*Triple)(nil)

// Func1 is a function with arity of 1.
type Func1 struct {
	Name string
	Arg  Node
}

func (Func1) Kind() Kind { return KindFunc1 }

var _ Node = (*Func1)(nil)

// BinOp is a binary operation.
type BinOp struct {
	Op       string
	Lhs, Rhs Node
}

func (BinOp) Kind() Kind { return KindBinOp }

var _ Node = (*BinOp)(nil)

// IfStmt represents an if statement.
// Cond must evaluate to Bool.
type IfStmt struct {
	Cond Node
	Then Node
	Else Node
}

func (IfStmt) Kind() Kind { return KindIfStmt }

var _ Node = (*IfStmt)(nil)
