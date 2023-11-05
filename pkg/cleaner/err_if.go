package cleaner

import (
	"go/ast"
	"go/token"

	"github.com/Fabianexe/go2jenkins/pkg/entity"
)

// cleanErrorIf removes all error if statements from the package data
// An error if statement is an if statement that checks if a variable x is not nil, x is named err or is of type error and has only a return statement in the body.
func cleanErrorIf(project *entity.Project) *entity.Project {
	for _, p := range project.Packages {
		for _, f := range p.Files {
			for _, method := range f.Methods {
				cleanErrorIfVisitor := &cleanErrorIfVisitor{
					fset: p.Fset,
				}

				ast.Walk(cleanErrorIfVisitor, method.Body)

				for _, errIf := range cleanErrorIfVisitor.errorIfs {
					i := 0
					for i < len(method.Lines) {
						line := method.Lines[i]
						if line.Number >= errIf.start && line.Number <= errIf.end {
							method.Lines = append(method.Lines[:i], method.Lines[i+1:]...)
							continue
						}

						i++
					}

					i = 0
					for i < len(method.Branches) {
						branch := method.Branches[i]
						if branch.StartLine >= errIf.start && branch.EndeLine <= errIf.end {
							method.Branches = append(method.Branches[:i], method.Branches[i+1:]...)
							continue
						}

						i++
					}
				}
			}
		}
	}

	return project
}

type cleanErrorIfVisitor struct {
	errorIfs []errorIF
	fset     *token.FileSet
}

type errorIF struct {
	start int
	end   int
}

func (c *cleanErrorIfVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if !isErrorIf(node) {
		return c
	}

	c.errorIfs = append(c.errorIfs, errorIF{
		start: c.fset.Position(node.Pos()).Line,
		end:   c.fset.Position(node.End()).Line,
	})

	return c
}

func isErrorIf(node ast.Node) bool {
	if node == nil { //  nothing here
		return false
	}

	v, ok := node.(*ast.IfStmt)
	if !ok { // no if statement
		return false
	}

	if v.Cond == nil || // np condition
		v.Else != nil || // has else
		len(v.Body.List) != 1 { // more than one statement in Body
		return false
	}

	cond, ok := v.Cond.(*ast.BinaryExpr)
	if !ok { // no binary expression
		return false
	}

	if cond.Op != token.NEQ { // not !=
		return false
	}

	if compare, ok := cond.Y.(*ast.Ident); !ok || compare.Name != "nil" { // not compared against nil
		return false
	}

	object, ok := cond.X.(*ast.Ident)
	if !ok || object.Obj == nil || object.Obj.Kind != ast.Var { // not a variable
		return false
	}

	if object.Name != "err" { // not named err
		// try to determine type of object
		typ, ok := object.Obj.Decl.(*ast.ValueSpec)
		if !ok { // not a value spec
			return false
		}

		if n, ok := typ.Type.(*ast.Ident); !ok || n.Name != "error" { // not an error
			return false
		}

		if _, ok := v.Body.List[0].(*ast.ReturnStmt); !ok { // body is not a return statement
			return false
		}
	}

	return true
}
