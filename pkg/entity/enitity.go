// Package entity contains all entities that are used in every component of the application
package entity

import (
	"go/ast"
	"go/token"
)

type Package struct {
	Name  string
	Files []*File
	Fset  *token.FileSet
}

type File struct {
	Name    string
	Ast     *ast.File
	Methods []*Method
}

type Method struct {
	Name  string
	Body  *ast.BlockStmt
	Lines []*Line
}

type Line struct {
	Number int
}
