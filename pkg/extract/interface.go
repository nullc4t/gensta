package extract

import (
	"fmt"
	"go/ast"
	"go/token"
)

func Interfaces(file *ast.File) []Interface {
	var ifaces []Interface

	importMap := make(map[Import]struct{})

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			iface, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			i := Interface{Name: typeSpec.Name.Name}

			for _, field := range iface.Methods.List {
				funcType, ok := field.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				i.Methods = append(i.Methods, Method{
					Name:    field.Names[0].Name,
					Args:    GetArgs(file, funcType.Params),
					Results: Results{GetArgs(file, funcType.Results)},
				})
			}

			for _, method := range i.Methods {
				for _, arg := range method.Args {
					fmt.Println(i.Name, "agr :", arg.Type.String())
					if arg.Type.IsImported() {
						importMap[Import{arg.Type.Package, arg.Type.ImportPath}] = struct{}{}
					}
				}
				for _, arg := range method.Results.Args {
					fmt.Println(i.Name, "agr :", arg.Type.String())
					if arg.Type.IsImported() {
						importMap[Import{arg.Type.Package, arg.Type.ImportPath}] = struct{}{}
					}
				}
			}

			for imp := range importMap {
				i.Imports = append(i.Imports, imp)
				fmt.Println(i.Name, "import map:", imp.Name, imp.Path)
			}

			ifaces = append(ifaces, i)
		}
	}

	return ifaces
}
