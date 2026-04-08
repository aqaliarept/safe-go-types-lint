package safegotypes

import (
	"go/ast"
	"go/token"
	"go/types"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "safe_go_types",
	Doc:      "Flags raw scalar types used as struct fields.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

var scalars = map[string]bool{
	"string": true, "bool": true,
	"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
	"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
	"float32": true, "float64": true,
	"complex64": true, "complex128": true,
	"byte": true, "rune": true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	insp.Preorder([]ast.Node{(*ast.StructType)(nil)}, func(n ast.Node) {
		for _, field := range n.(*ast.StructType).Fields.List {
			if !containsScalar(pass, field.Type) {
				continue
			}
			for _, name := range field.Names {
				pass.Reportf(name.Pos(), "safe-go-types/no-scalar: field %q has raw scalar type", name.Name)
			}
		}
	})

	checkNoScalarLocalVar(pass)

	customTypes := collectCustomTypes(pass)
	checkNoConstructor(pass)
	checkNoZeroValue(pass, customTypes)
	checkNoCast(pass, customTypes)
	checkUntypedLiteral(pass, customTypes)

	return nil, nil
}

// checkNoScalarLocalVar flags local variable declarations with explicit raw scalar types.
func checkNoScalarLocalVar(pass *analysis.Pass) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}
			// Walk inside function bodies only.
			ast.Inspect(fn.Body, func(inner ast.Node) bool {
				switch node := inner.(type) {
				case *ast.GenDecl:
					if node.Tok != token.VAR {
						return true
					}
					for _, spec := range node.Specs {
						vs, ok := spec.(*ast.ValueSpec)
						if !ok || vs.Type == nil {
							continue
						}
						if !containsScalar(pass, vs.Type) {
							continue
						}
						for _, name := range vs.Names {
							pass.Reportf(name.Pos(), "safe-go-types/no-scalar: variable %q has raw scalar type", name.Name)
						}
					}
				case *ast.AssignStmt:
					if node.Tok != token.DEFINE {
						return true
					}
					// Flag only when RHS is a literal (not a function call or identifier).
					for i, rhs := range node.Rhs {
						if !isScalarLiteralExpr(rhs) {
							continue
						}
						if i >= len(node.Lhs) {
							break
						}
						lhsIdent, ok := node.Lhs[i].(*ast.Ident)
						if !ok {
							continue
						}
						obj := pass.TypesInfo.Defs[lhsIdent]
						if obj == nil {
							continue
						}
						typeName := scalarTypeName(obj.Type())
						if typeName == "" {
							continue
						}
						pass.Reportf(lhsIdent.Pos(), "safe-go-types/no-scalar: variable %q has raw scalar type %q", lhsIdent.Name, typeName)
					}
				}
				return true
			})
			return false
		})
	}
}

// checkNoConstructor flags custom types defined over scalars (or other custom
// types in the same package) that lack a valid constructor.
func checkNoConstructor(pass *analysis.Pass) {
	// Pass 1: collect custom types that need a constructor.
	// typePos maps each qualifying type name to its declaration position.
	typePos := map[string]token.Pos{} // type name → position to report

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				ident, ok := ts.Type.(*ast.Ident)
				if !ok {
					continue
				}
				obj := pass.TypesInfo.Uses[ident]
				// Scalar builtin (e.g. type Foo string)
				if scalars[ident.Name] && isBuiltinType(obj) {
					typePos[ts.Name.Name] = ts.Name.Pos()
					continue
				}
				// Custom type in the same package (e.g. type Foo Bar)
				if obj != nil && obj.Pkg() != nil && obj.Pkg() == pass.Pkg {
					typePos[ts.Name.Name] = ts.Name.Pos()
				}
			}
		}
	}

	// Pass 2: collect valid constructors in this package.
	// A valid constructor for Foo is a func named NewFoo or newFoo
	// returning exactly (Foo, error).
	hasConstructor := map[string]bool{}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}
			typeName := constructorTarget(fn.Name.Name)
			if _, needed := typePos[typeName]; !needed {
				continue
			}
			if isValidConstructorSignature(pass, fn, typeName) {
				hasConstructor[typeName] = true
			}
		}
	}

	// Emit diagnostics for types without a valid constructor.
	for typeName, pos := range typePos {
		if !hasConstructor[typeName] {
			pass.Reportf(pos, "safe-go-types/no-constructor: type %q has no valid constructor", typeName)
		}
	}
}

// constructorTarget returns the type name that a constructor function targets.
// "NewFoo" → "Foo", "newFoo" → "foo". Returns "" if the pattern doesn't match.
func constructorTarget(funcName string) string {
	if len(funcName) <= 3 {
		return ""
	}
	suffix := funcName[3:]
	switch funcName[:3] {
	case "New":
		return suffix
	case "new":
		runes := []rune(suffix)
		runes[0] = unicode.ToLower(runes[0])
		return string(runes)
	}
	return ""
}

// isValidConstructorSignature checks that fn returns exactly (TypeName, error).
func isValidConstructorSignature(pass *analysis.Pass, fn *ast.FuncDecl, typeName string) bool {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 2 {
		return false
	}
	results := fn.Type.Results.List

	// First return must be an unqualified ident naming the type in this package.
	firstIdent, ok := results[0].Type.(*ast.Ident)
	if !ok || firstIdent.Name != typeName {
		return false
	}
	obj := pass.TypesInfo.Uses[firstIdent]
	if obj == nil || obj.Pkg() != pass.Pkg {
		return false
	}

	// Second return must be the builtin error interface.
	secondIdent, ok := results[1].Type.(*ast.Ident)
	if !ok || secondIdent.Name != "error" {
		return false
	}
	errObj := pass.TypesInfo.Uses[secondIdent]
	return errObj != nil && errObj.Pkg() == nil
}

// containsScalar reports whether the type expression expr contains a raw scalar
// builtin type anywhere in its structure (e.g. []string, map[string]int, *bool, chan int).
func containsScalar(pass *analysis.Pass, expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		if !scalars[e.Name] {
			return false
		}
		return isBuiltinType(pass.TypesInfo.Uses[e])
	case *ast.ArrayType:
		return containsScalar(pass, e.Elt)
	case *ast.MapType:
		return containsScalar(pass, e.Key) || containsScalar(pass, e.Value)
	case *ast.StarExpr:
		return containsScalar(pass, e.X)
	case *ast.ChanType:
		return containsScalar(pass, e.Value)
	}
	return false
}

// isScalarLiteralExpr reports whether expr is a scalar literal (BasicLit),
// composite literal, or explicit type conversion — not a function call or identifier.
func isScalarLiteralExpr(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.CompositeLit:
		return true
	}
	return false
}

// scalarTypeName returns the scalar type name for a types.Type if it is a builtin scalar,
// otherwise returns "".
func scalarTypeName(t types.Type) string {
	basic, ok := t.(*types.Basic)
	if !ok {
		return ""
	}
	name := basic.Name()
	if scalars[name] {
		return name
	}
	return ""
}

// isBuiltinType reports whether obj is a built-in type (no package).
func isBuiltinType(obj types.Object) bool {
	_, ok := obj.(*types.TypeName)
	return ok && obj.Pkg() == nil
}

// collectCustomTypes returns the set of custom type names in this package.
// A custom type is: type T scalar, type T U where U is also a custom type,
// or a struct type that contains at least one field of a custom type.
func collectCustomTypes(pass *analysis.Pass) map[string]bool {
	customTypes := map[string]bool{}

	// First pass: scalar-backed and custom-type-backed types.
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				ident, ok := ts.Type.(*ast.Ident)
				if !ok {
					continue
				}
				obj := pass.TypesInfo.Uses[ident]
				if scalars[ident.Name] && isBuiltinType(obj) {
					customTypes[ts.Name.Name] = true
					continue
				}
				if obj != nil && obj.Pkg() != nil && obj.Pkg() == pass.Pkg {
					customTypes[ts.Name.Name] = true
				}
			}
		}
	}

	// Second pass: struct types that have at least one field of a custom type.
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genDecl.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range st.Fields.List {
					fieldIdent, ok := field.Type.(*ast.Ident)
					if !ok {
						continue
					}
					if customTypes[fieldIdent.Name] {
						customTypes[ts.Name.Name] = true
						break
					}
				}
			}
		}
	}

	return customTypes
}

// checkNoCast flags explicit conversions to custom types outside their constructors.
func checkNoCast(pass *analysis.Pass, customTypes map[string]bool) {

	for _, file := range pass.Files {
		// Walk the AST tracking the current enclosing function.
		var enclosingFunc string
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				prev := enclosingFunc
				enclosingFunc = node.Name.Name
				// Walk the body explicitly.
				ast.Inspect(node.Body, func(inner ast.Node) bool {
					call, ok := inner.(*ast.CallExpr)
					if !ok {
						return true
					}
					ident, ok := call.Fun.(*ast.Ident)
					if !ok {
						return true
					}
					obj := pass.TypesInfo.Uses[ident]
					if obj == nil {
						return true
					}
					typeName, ok := obj.(*types.TypeName)
					if !ok {
						return true
					}
					if !customTypes[typeName.Name()] {
						return true
					}
					// Check if enclosing func is named as a constructor for this type
					// (New<T> or new<T>), regardless of signature validity.
					if constructorTarget(enclosingFunc) == typeName.Name() {
						return true
					}
					pass.Reportf(call.Pos(), "safe-go-types/no-cast: conversion to custom type %q outside its constructor", typeName.Name())
					return true
				})
				enclosingFunc = prev
				return false // already walked body above
			}
			return true
		})
	}
}

// checkNoZeroValue flags var declarations of custom types with no initializer.
func checkNoZeroValue(pass *analysis.Pass, customTypes map[string]bool) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				// Only flag when there's no initializer at all.
				if len(vs.Values) > 0 {
					continue
				}
				// Check if the declared type is a custom type.
				if vs.Type == nil {
					continue
				}
				ident, ok := vs.Type.(*ast.Ident)
				if !ok || !customTypes[ident.Name] {
					continue
				}
				for _, name := range vs.Names {
					pass.Reportf(name.Pos(), "safe-go-types/no-zero-value: variable %q is zero-initialized custom type", name.Name)
				}
			}
		}
	}
}

// isUntypedLiteral reports whether expr is an untyped constant literal in source.
// It checks for basic literals (string, int, float, char) which are always untyped.
func isUntypedLiteral(expr ast.Expr) bool {
	_, ok := expr.(*ast.BasicLit)
	return ok
}

// checkUntypedLiteral flags untyped constant literals assigned to or passed as custom types.
func checkUntypedLiteral(pass *analysis.Pass, customTypes map[string]bool) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			// Skip const declarations (same-package constant exemption).
			if genDecl.Tok == token.CONST {
				continue
			}
			if genDecl.Tok != token.VAR {
				continue
			}
			for _, spec := range genDecl.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				// Need explicit type annotation to know the target type.
				if vs.Type == nil {
					continue
				}
				typeIdent, ok := vs.Type.(*ast.Ident)
				if !ok || !customTypes[typeIdent.Name] {
					continue
				}
				for i, val := range vs.Values {
					if isUntypedLiteral(val) {
						pass.Reportf(vs.Names[i].Pos(), "safe-go-types/untyped-literal: untyped literal assigned to custom type %q", typeIdent.Name)
					}
				}
			}
		}
	}

	// Check function call arguments.
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			// Get the callee's type from TypesInfo.
			calleeTV, ok := pass.TypesInfo.Types[call.Fun]
			if !ok {
				return true
			}
			// Skip type conversions (callee is a type, not a function).
			if calleeTV.IsType() {
				return true
			}
			sig, ok := calleeTV.Type.Underlying().(*types.Signature)
			if !ok {
				return true
			}
			params := sig.Params()
			for i, arg := range call.Args {
				if !isUntypedLiteral(arg) {
					continue
				}
				if i >= params.Len() {
					break
				}
				paramType := params.At(i).Type()
				// Get the named type.
				named, ok := paramType.(*types.Named)
				if !ok {
					continue
				}
				if named.Obj().Pkg() != pass.Pkg {
					continue
				}
				if !customTypes[named.Obj().Name()] {
					continue
				}
				pass.Reportf(arg.Pos(), "safe-go-types/untyped-literal: untyped literal passed as custom type %q", named.Obj().Name())
			}
			return true
		})
	}
}
