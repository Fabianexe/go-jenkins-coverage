package cleaner

import (
	"go/ast"
	"go/token"
	"log/slog"

	"github.com/Fabianexe/go2jenkins/pkg/entity"
	"github.com/Fabianexe/go2jenkins/pkg/utility"
)

// cleanErrorIf removes all error if statements from the package data
// An error if statement is an if statement that checks if a variable x is not nil, x is named err or is of type error and has only a return statement in the body.
func cleanErrorIf(project *entity.Project) *entity.Project {
	slog.Info("Clean error if statements")
	for _, p := range project.Packages {
		for _, f := range p.Files {
			var countErrorIf int
			for _, method := range f.Methods {
				cleanErrorIfVisitor := &cleanErrorIfVisitor{
					fset: p.Fset,
				}

				ast.Walk(cleanErrorIfVisitor, method.Body)

				countErrorIf += len(cleanErrorIfVisitor.errorIfs)

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
						if branch.StartLine >= errIf.start && branch.EndLine <= errIf.end {
							method.Branches = append(method.Branches[:i], method.Branches[i+1:]...)
							continue
						}

						i++
					}
				}
			}
			if countErrorIf > 0 {
				slog.Debug("Cleaned error if statements", "File", f.FilePath, "ErrorIfs", countErrorIf)
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
	if !utility.IsErrorIf(node) {
		return c
	}

	c.errorIfs = append(c.errorIfs, errorIF{
		start: c.fset.Position(node.Pos()).Line,
		end:   c.fset.Position(node.End()).Line,
	})

	return c
}
