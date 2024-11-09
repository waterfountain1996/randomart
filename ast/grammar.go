package ast

import "math/big"

type Rule struct {
	Name     string
	Branches []*Branch
}

type Branch struct {
	Node

	Probability *big.Rat
}

func (*Rule) Kind() Kind { return KindRule }

var _ Node = (*Rule)(nil)
