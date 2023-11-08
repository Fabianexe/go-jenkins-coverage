package source

import (
	"go/ast"
	"go/token"

	"github.com/Fabianexe/go2jenkins/pkg/entity"
)

type branchVisitor struct {
	branches []*entity.Branch
	fset     *token.FileSet
}

func (b *branchVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch v := node.(type) {
	case *ast.CaseClause:
		b.branches = append(b.branches, b.createBranch(v.Pos(), v.Colon, v.End()))
	case *ast.CommClause:
		b.branches = append(b.branches, b.createBranch(v.Pos(), v.Colon, v.End()))
	case *ast.IfStmt:
		b.branches = append(b.branches, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End()))
	case *ast.ForStmt:
		b.branches = append(b.branches, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End()))
	case *ast.RangeStmt:
		b.branches = append(b.branches, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End()))
	case *ast.FuncDecl:
		b.branches = append(b.branches, b.createBranch(v.Pos(), v.Body.Pos(), v.Body.End()))
	}

	return b
}

func (b *branchVisitor) createBranch(def, start, end token.Pos) *entity.Branch {
	return &entity.Branch{
		DefLine:   b.fset.Position(def).Line,
		StartLine: b.fset.Position(start).Line + 1,
		EndLine:   b.fset.Position(end).Line,
	}
}
