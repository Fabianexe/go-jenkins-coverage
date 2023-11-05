// Package source reads source files from a given path and create entities from this
package source

import (
	"go/ast"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"

	"github.com/Fabianexe/go-jenkins-coverage/pkg/entity"
)

func LoadSources(path string) (*entity.Project, error) {
	goPath := make(map[string]struct{}, 1000)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" {
			goPath[filepath.Dir(path)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	gofiles := make([]string, 0, len(goPath))
	for pack := range goPath {
		gofiles = append(gofiles, pack)
	}

	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedModule |
			packages.NeedTypes,
		Dir: path,
	}

	pkgs, err := packages.Load(cfg, gofiles...)
	if err != nil {
		return nil, err
	}

	allPackages := make([]*entity.Package, 0, len(pkgs))
	for _, pkg := range pkgs {
		pack := &entity.Package{
			Name:  pkg.PkgPath,
			Files: make([]*entity.File, 0, len(pkg.Syntax)),
			Fset:  pkg.Fset,
		}
		for i, fileAst := range pkg.Syntax {
			file := &entity.File{
				Name:     filepath.Base(pkg.GoFiles[i]),
				FilePath: pkg.GoFiles[i],
				Ast:      fileAst,
				Methods:  make([]*entity.Method, 0, len(fileAst.Decls)),
			}
			for _, decl := range fileAst.Decls {
				if fun, ok := decl.(*ast.FuncDecl); ok {

					method := &entity.Method{
						Name: fun.Name.Name,
						Body: fun.Body,
					}

					// start after the function declaration
					startLine := pkg.Fset.Position(fun.Pos()).Line + 1
					endLine := pkg.Fset.Position(fun.End()).Line
					if startLine >= endLine {
						continue
					}

					lines := make([]*entity.Line, 0, endLine-startLine)
					for i := startLine; i < endLine; i++ {
						lines = append(lines, &entity.Line{Number: i})
					}
					method.Lines = lines
					bV := &branchVisitor{
						fset: pkg.Fset,
					}

					ast.Walk(bV, fun)

					method.Branches = bV.branches

					file.Methods = append(file.Methods, method)
				}
			}

			pack.Files = append(pack.Files, file)
		}

		allPackages = append(allPackages, pack)
	}

	return &entity.Project{Packages: allPackages}, nil
}
