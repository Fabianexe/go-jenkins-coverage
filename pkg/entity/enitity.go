// Package entity contains all entities that are used in every component of the application
package entity

import (
	"go/ast"
	"go/token"
)

type Project struct {
	Packages       []*Package
	LineCoverage   LineCounter
	BranchCoverage BranchCounter
}

type Package struct {
	Name           string
	Files          []*File
	Fset           *token.FileSet
	LineCoverage   LineCounter
	BranchCoverage BranchCounter
}

type File struct {
	Name           string
	FilePath       string
	Ast            *ast.File
	Methods        []*Method
	LineCoverage   LineCounter
	BranchCoverage BranchCounter
}

type Method struct {
	Name           string
	Body           *ast.BlockStmt
	Lines          []*Line
	Branches       []*Branch
	LineCoverage   LineCounter
	BranchCoverage BranchCounter
}

type Line struct {
	Number        int
	CoverageCount int
}

type Branch struct {
	StartLine int
	EndeLine  int
	Covered   bool
}
