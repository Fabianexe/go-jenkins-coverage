package utility

import (
	"go/ast"
	"go/token"
)

func IsErrorIf(node ast.Node) bool {
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

	if _, ok := v.Body.List[0].(*ast.ReturnStmt); !ok { // body is not a return statement
		return false
	}

	return isErrorVar(cond)
}

func isErrorVar(cond *ast.BinaryExpr) bool {
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
	}

	return true
}
