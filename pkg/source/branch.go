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
		branch := &entity.Branch{}
		branch.StartLine = b.fset.Position(v.Colon).Line + 1
		branch.EndeLine = b.fset.Position(v.End()).Line
		b.branches = append(b.branches, branch)
	case *ast.CommClause:
		branch := &entity.Branch{}
		branch.StartLine = b.fset.Position(v.Colon).Line + 1
		branch.EndeLine = b.fset.Position(v.End()).Line
		b.branches = append(b.branches, branch)
	case *ast.IfStmt:
		branch := &entity.Branch{}
		branch.StartLine = b.fset.Position(v.Body.Pos()).Line + 1
		branch.EndeLine = b.fset.Position(v.Body.End()).Line
		b.branches = append(b.branches, branch)
	case *ast.ForStmt:
		branch := &entity.Branch{}
		branch.StartLine = b.fset.Position(v.Body.Pos()).Line + 1
		branch.EndeLine = b.fset.Position(v.Body.End()).Line
		b.branches = append(b.branches, branch)
	case *ast.RangeStmt:
		branch := &entity.Branch{}
		branch.StartLine = b.fset.Position(v.Body.Pos()).Line + 1
		branch.EndeLine = b.fset.Position(v.Body.End()).Line
		b.branches = append(b.branches, branch)
	case *ast.FuncDecl:
		branch := &entity.Branch{}
		branch.StartLine = b.fset.Position(v.Body.Pos()).Line + 1
		branch.EndeLine = b.fset.Position(v.Body.End()).Line
		b.branches = append(b.branches, branch)
	}

	return b
}
