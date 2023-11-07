package cleaner

import (
	"go/ast"
	"go/token"
	"log/slog"

	"github.com/Fabianexe/go2jenkins/pkg/entity"
)

func cleanNoneCodeLines(project *entity.Project) *entity.Project {
	slog.Info("Clean none code lines")
	for _, p := range project.Packages {
		for _, f := range p.Files {
			var cleanedLines int
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
						cleanedLines++
						method.Lines = append(method.Lines[:i], method.Lines[i+1:]...)
						continue
					}

					i++
				}
			}
			if cleanedLines > 0 {
				slog.Debug("Cleaned lines", "File", f.FilePath, "Lines", cleanedLines)
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
