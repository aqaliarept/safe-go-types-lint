package safegotypes

import (
	"go/ast"
	"go/types"

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

	return nil, nil
}

// isBuiltinType reports whether obj is a built-in type (no package).
func isBuiltinType(obj types.Object) bool {
	_, ok := obj.(*types.TypeName)
	return ok && obj.Pkg() == nil
}
