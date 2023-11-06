package complexity

import (
	"go/ast"

	"github.com/Fabianexe/go2jenkins/pkg/utility"
)

func getCyclomaticComplexity(root ast.Node, ignoreErrorIF bool) int {
	visitor := &cyclomaticVisitor{
		complexity:    1,
		ignoreErrorIF: ignoreErrorIF,
	}

	ast.Walk(visitor, root)

	return visitor.complexity

}

type cyclomaticVisitor struct {
	complexity    int
	ignoreErrorIF bool
}

func (c *cyclomaticVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch node.(type) {
	case *ast.IfStmt:
		if !c.ignoreErrorIF || !utility.IsErrorIf(node) {
			c.complexity++
		}
	case *ast.ForStmt,
		*ast.RangeStmt,
		*ast.FuncDecl,
		*ast.SwitchStmt,
		*ast.TypeSwitchStmt,
		*ast.SelectStmt:
		c.complexity++
	}

	return c
}
