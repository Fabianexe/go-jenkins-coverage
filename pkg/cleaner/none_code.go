package cleaner

import (
	"go/ast"
	"go/token"

	"github.com/Fabianexe/go-jenkins-coverage/pkg/entity"
)

func cleanNoneCodeLines(project *entity.Project) *entity.Project {
	for _, p := range project.Packages {
		for _, f := range p.Files {
			for _, method := range f.Methods {
				noneCodeVisitor := &noneCodeVisitor{
					validLines: make(map[int]struct{}, len(method.Lines)),
					fset:       p.Fset,
				}

				ast.Walk(noneCodeVisitor, method.Body)

				i := 0
				for i < len(method.Lines) {
					line := method.Lines[i]
					if _, ok := noneCodeVisitor.validLines[line.Number]; !ok {
						method.Lines = append(method.Lines[:i], method.Lines[i+1:]...)
						continue
					}

					i++
				}
			}
		}
	}

	return project
}

type noneCodeVisitor struct {
	validLines map[int]struct{}
	fset       *token.FileSet
}

func (n *noneCodeVisitor) Visit(node ast.Node) (w ast.Visitor) {
	if node == nil {
		return n
	}

	lineNUmber := n.fset.Position(node.Pos()).Line
	n.validLines[lineNUmber] = struct{}{}

	return n
}
