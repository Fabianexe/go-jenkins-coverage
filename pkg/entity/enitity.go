// Package entity contains all entities that are used in every component of the application
package entity

import (
	"go/ast"
	"go/token"
)

type Project struct {
	Packages     []*Package
	LineCoverage LineCounter
}

type Package struct {
	Name         string
	Files        []*File
	Fset         *token.FileSet
	LineCoverage LineCounter
}

type File struct {
	Name         string
	FilePath     string
	Ast          *ast.File
	Methods      []*Method
	LineCoverage LineCounter
}

type Method struct {
	Name         string
	Body         *ast.BlockStmt
	Lines        []*Line
	LineCoverage LineCounter
}

type Line struct {
	Number        int
	CoverageCount int
}
