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

	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		st := n.(*ast.StructType)
		for _, field := range st.Fields.List {
			if ident, ok := field.Type.(*ast.Ident); ok {
				if scalars[ident.Name] {
					// Confirm the type resolves to a built-in (not a named type from the package)
					if obj := pass.TypesInfo.Uses[ident]; obj != nil {
						if _, isBuiltin := obj.(*types.TypeName); isBuiltin {
							if obj.Pkg() == nil { // built-in types have no package
								for _, name := range field.Names {
									pass.Reportf(name.Pos(), "safe-go-types/no-scalar: field %q has raw scalar type %q", name.Name, ident.Name)
								}
							}
						}
					}
				}
			}
		}
	})

	return nil, nil
}
