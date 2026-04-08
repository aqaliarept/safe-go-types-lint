package safegotypes

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"
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
	// A type needs a constructor if it is defined as:
	//   type Foo <scalar>
	//   type Foo Bar  (where Bar is a custom type in the same package)
	needsConstructor := map[string]bool{}  // type name → true
	typePos := map[string]token.Pos{}      // type name → position to report

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
				typeName := ts.Name.Name
				// Check if underlying type is a scalar builtin
				if scalars[ident.Name] && isBuiltinType(pass.TypesInfo.Uses[ident]) {
					needsConstructor[typeName] = true
					typePos[typeName] = ts.Name.Pos()
					continue
				}
				// Check if underlying type is another custom type in this package
				obj := pass.TypesInfo.Uses[ident]
				if obj != nil && obj.Pkg() != nil && obj.Pkg() == pass.Pkg {
					needsConstructor[typeName] = true
					typePos[typeName] = ts.Name.Pos()
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
			name := fn.Name.Name
			// Derive the type name this constructor targets.
			typeName := constructorTarget(name)
			if typeName == "" {
				continue
			}
			if !needsConstructor[typeName] {
				continue
			}
			// Check return signature: exactly (TypeName, error)
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
// For "NewFoo" or "newFoo" it returns "Foo" or "foo".
// Returns "" if the name doesn't match the pattern.
func constructorTarget(funcName string) string {
	if strings.HasPrefix(funcName, "New") && len(funcName) > 3 {
		suffix := funcName[3:]
		// The type name matches the suffix directly (exported constructor → exported type)
		return suffix
	}
	if strings.HasPrefix(funcName, "new") && len(funcName) > 3 {
		suffix := funcName[3:]
		// Unexported constructor: newFoo → foo (lowercase first letter)
		if len(suffix) > 0 {
			runes := []rune(suffix)
			runes[0] = unicode.ToLower(runes[0])
			return string(runes)
		}
	}
	return ""
}

// isValidConstructorSignature checks that fn returns exactly (TypeName, error).
func isValidConstructorSignature(pass *analysis.Pass, fn *ast.FuncDecl, typeName string) bool {
	if fn.Type.Results == nil {
		return false
	}
	results := fn.Type.Results.List
	// Count total return values (fields can have multiple names, but for
	// return values that's unusual; still handle it).
	total := 0
	for _, r := range results {
		if len(r.Names) == 0 {
			total++
		} else {
			total += len(r.Names)
		}
	}
	if total != 2 {
		return false
	}

	// First return must be the type itself (unqualified ident).
	first := results[0]
	firstIdent, ok := first.Type.(*ast.Ident)
	if !ok || firstIdent.Name != typeName {
		return false
	}
	// The object must be in the same package.
	obj := pass.TypesInfo.Uses[firstIdent]
	if obj == nil || obj.Pkg() != pass.Pkg {
		return false
	}

	// Second return must be the builtin error interface.
	second := results[1]
	secondIdent, ok := second.Type.(*ast.Ident)
	if !ok || secondIdent.Name != "error" {
		return false
	}
	// error is a builtin — no package.
	errObj := pass.TypesInfo.Uses[secondIdent]
	if errObj == nil || errObj.Pkg() != nil {
		return false
	}

	return true
}

// isBuiltinType reports whether obj is a built-in type (no package).
func isBuiltinType(obj types.Object) bool {
	_, ok := obj.(*types.TypeName)
	return ok && obj.Pkg() == nil
}
