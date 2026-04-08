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
			ident, ok := field.Type.(*ast.Ident)
			if !ok || !scalars[ident.Name] {
				continue
			}
			if !isBuiltinType(pass.TypesInfo.Uses[ident]) {
				continue
			}
			for _, name := range field.Names {
				pass.Reportf(name.Pos(), "safe-go-types/no-scalar: field %q has raw scalar type %q", name.Name, ident.Name)
			}
		}
	})

	checkNoConstructor(pass)

	return nil, nil
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

// isBuiltinType reports whether obj is a built-in type (no package).
func isBuiltinType(obj types.Object) bool {
	_, ok := obj.(*types.TypeName)
	return ok && obj.Pkg() == nil
}
