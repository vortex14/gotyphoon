package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// borrowed from https://github.com/lukehoban/go-outline/blob/master/main.go#L54-L107
	fset := token.NewFileSet()
	parserMode := parser.ParseComments
	var fileAst *ast.File
	var err error

	fileAst, err = parser.ParseFile(fset, "group.go", `package main

const GroupName = "andy"

func reTRALL() {

}
`, parserMode)
	if err != nil {
		panic(err)
	}

	for _, d := range fileAst.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:

			println(decl.Name.String())

			fmt.Println()
			fmt.Println("Func ??? ")

		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.ImportSpec:
					fmt.Println("Import", spec.Path.Value)
				case *ast.TypeSpec:
					fmt.Println("Type", spec.Name.String())
				case *ast.ValueSpec:
					for _, id := range spec.Names {
						fmt.Printf("Var %s: %v", id.Name, id.Obj.Decl.(*ast.ValueSpec).Values[0].(*ast.BasicLit).Value)
					}
				default:
					fmt.Printf("Unknown token type: %s\n", decl.Tok)
				}
			}
		default:
			println(fmt.Sprintf("Unknown declaration @ %+v", decl.Pos()))
			//fmt.Printf("Unknown declaration @\n", decl.Pos())
		}
	}
}
