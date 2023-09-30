package main

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"testing"
)

func TestNewImport(t *testing.T) {
	// Create a new file set
	fset := token.NewFileSet()

	// Create an empty AST
	f := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{},
	}

	// Add the "go/token" import
	f.Decls = append(f.Decls, &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"go/token\"",
				},
			},
		},
	})

	// Add the "os" import after the "go/token" import
	f.Decls = append(f.Decls, &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"os\"",
				},
			},
		},
	})

	// Format the imports with multi-line format
	ast.SortImports(fset, f)
	cfg := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	cfg.Fprint(os.Stdout, fset, f)
}

func TestCreateInterface(t *testing.T) {
	// Create the User struct type
	userStructType := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{
						ast.NewIdent("Name"),
					},
					Type: ast.NewIdent("string"),
				},
				&ast.Field{
					Names: []*ast.Ident{
						ast.NewIdent("Age"),
					},
					Type: ast.NewIdent("int"),
				},
			},
		},
	}

	// Create the User pointer type
	userPointerType := &ast.StarExpr{
		X: ast.NewIdent("User"),
	}

	// Create the GetUsers method
	getUsersMethod := &ast.FuncDecl{
		Name: ast.NewIdent("GetUsers"),
		Type: &ast.FuncType{
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Type: &ast.ArrayType{
							Elt: userPointerType,
						},
					},
				},
			},
		},
	}

	// Create the AddUser method
	addUserMethod := &ast.FuncDecl{
		Name: ast.NewIdent("AddUser"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
							ast.NewIdent("user"),
						},
						Type: userPointerType,
					},
				},
			},
		},
	}

	// Create the UserService interface
	userServiceInterface := &ast.InterfaceType{
		Methods: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Names: []*ast.Ident{
						ast.NewIdent("GetUsers"),
					},
					Type: getUsersMethod.Type,
				},
				&ast.Field{
					Names: []*ast.Ident{
						ast.NewIdent("AddUser"),
					},
					Type: addUserMethod.Type,
				},
			},
		},
	}

	// Create the UserService declaration
	userServiceDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("UserService"),
				Type: userServiceInterface,
			},
		},
	}

	// Create the module
	module := &ast.File{
		Name: ast.NewIdent("main"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("User"),
						Type: userStructType,
					},
				},
			},
			userServiceDecl,
		},
	}

	// Print the AST to stdout
	printer.Fprint(os.Stdout, token.NewFileSet(), module)
}

func TestCreateFunction(t *testing.T) {
	xParam := &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent("x"),
		},
		Type: ast.NewIdent("int"),
	}
	yParam := &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent("y"),
		},
		Type: ast.NewIdent("int"),
	}

	// Create the sum expression
	sumExpr := &ast.BinaryExpr{
		X:  ast.NewIdent("x"),
		Op: token.ADD,
		Y:  ast.NewIdent("y"),
	}

	// Create the return statement
	returnStmt := &ast.ReturnStmt{
		Results: []ast.Expr{
			sumExpr,
		},
	}

	// Create the function type
	funcType := &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				xParam,
				yParam,
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				&ast.Field{
					Type: ast.NewIdent("int"),
				},
			},
		},
	}

	// Create the function declaration
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent("sum"),
		Type: funcType,
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				returnStmt,
			},
		},
	}

	// Create the AST file
	file := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{funcDecl},
	}

	// Print the AST to stdout
	printer.Fprint(os.Stdout, token.NewFileSet(), file)
}

func TestCreateStruct(t *testing.T) {
	// Create the Age field with a comment
	ageField := &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent("Age"),
		},
		Type: ast.NewIdent("int"),
		Comment: &ast.CommentGroup{
			List: []*ast.Comment{
				&ast.Comment{
					Text:  "// you years",
					Slash: token.Pos(1),
				},
			},
		},
	}

	// Create the Name field with a comment
	nameField := &ast.Field{
		Names: []*ast.Ident{
			ast.NewIdent("Name"),
		},
		Type: ast.NewIdent("string"),
		Comment: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text:  "// your name",
					Slash: token.Pos(1),
				},
			},
		},
	}

	// Create the Struct type
	structType := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{
				ageField,
				nameField,
			},
		},
	}

	// Create the Struct declaration
	structDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("Struct"),
				Type: structType,
			},
		},
	}

	// Create the AST file
	file := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{structDecl},
	}

	// Print the AST to stdout
	printer.Fprint(os.Stdout, token.NewFileSet(), file)
}

func TestCreatePipeline_2(t *testing.T) {
	fset := token.NewFileSet()

	// Create an empty AST
	f := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{},
	}

	// Create the BasePipeline struct
	pl := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{Name: "pl"},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: &ast.Ident{Name: "BasePipeline"},
					Elts: []ast.Expr{
						&ast.KeyValueExpr{
							Key: &ast.Ident{Name: "MetaInfo"},
							Value: &ast.UnaryExpr{
								Op: token.AND,
								X: &ast.CompositeLit{
									Type: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "label"},
										Sel: &ast.Ident{Name: "MetaInfo"},
									},
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key:   &ast.Ident{Name: "Name"},
											Value: &ast.BasicLit{Kind: token.STRING, Value: "\"stage-1\""},
										},
										&ast.KeyValueExpr{
											Key:   &ast.Ident{Name: "Required"},
											Value: &ast.Ident{Name: "true"},
										},
									},
								},
							},
						},
						&ast.KeyValueExpr{
							Key:   &ast.Ident{Name: "Fn"},
							Value: createFnAST(),
						},
					},
				},
			},
		},
	}

	// Add the BasePipeline struct to the main function body
	f.Decls = append(f.Decls, &ast.FuncDecl{
		Name: &ast.Ident{Name: "main"},
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				pl,
			},
		},
	})
	ast.SortImports(fset, f)
	cfg := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	cfg.Fprint(os.Stdout, fset, f)

}

func TestCreatePipeline(t *testing.T) {
	// Create a new file set
	fset := token.NewFileSet()

	// Create an empty AST
	f := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{},
	}

	// Create the BasePipeline struct
	pl := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{Name: "pl"},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "BasePipeline"},
					Sel: &ast.Ident{Name: ""},
				},
				Elts: []ast.Expr{
					&ast.KeyValueExpr{
						Key: &ast.Ident{Name: "MetaInfo"},
						Value: &ast.CompositeLit{
							Type: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "label"},
								Sel: &ast.Ident{Name: "MetaInfo"},
							},
							Elts: []ast.Expr{
								&ast.KeyValueExpr{
									Key:   &ast.Ident{Name: "Name"},
									Value: &ast.BasicLit{Kind: token.STRING, Value: "\"stage-1\""},
								},
								&ast.KeyValueExpr{
									Key:   &ast.Ident{Name: "Required"},
									Value: &ast.Ident{Name: "true"},
								},
							},
						},
					},
					&ast.KeyValueExpr{
						Key:   &ast.Ident{Name: "Fn"},
						Value: createFnAST(),
					},
				},
			},
		},
	}

	// Add the BasePipeline struct to the main function body
	f.Decls = append(f.Decls, &ast.FuncDecl{
		Name: &ast.Ident{Name: "main"},
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				pl,
			},
		},
	})

	// Format the code
	ast.SortImports(fset, f)
	cfg := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	cfg.Fprint(os.Stdout, fset, f)
}

func createFnAST() ast.Expr {
	return &ast.FuncLit{
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							&ast.Ident{Name: "ctx"},
						},
						Type: &ast.Ident{Name: "context.Context"},
					},
					{
						Names: []*ast.Ident{
							&ast.Ident{Name: "logger"},
						},
						Type: &ast.Ident{Name: "interfaces.LoggerInterface"},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{Name: "error"},
					},
					{
						Type: &ast.Ident{Name: "context.Context"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "logger"},
							Sel: &ast.Ident{Name: "Info"},
						},
						Args: []ast.Expr{
							&ast.BasicLit{Kind: token.STRING, Value: "\"run stage 1\""},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.Ident{Name: "nil"},
						&ast.Ident{Name: "ctx"},
					},
				},
			},
		},
	}
}

func TestNewPbErrorIntoMap(t *testing.T) {
	// Create a new file set
	fset := token.NewFileSet()

	// Create an empty AST
	f := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{},
	}

	// Create the ReasonToCode map
	reasonToCode := &ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{Name: "ReasonToCode"},
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			&ast.CompositeLit{
				Type: &ast.MapType{
					Key:   &ast.SelectorExpr{X: &ast.Ident{Name: "pb"}, Sel: &ast.Ident{Name: "ErrorReason"}},
					Value: &ast.SelectorExpr{X: &ast.Ident{Name: "codes"}, Sel: &ast.Ident{Name: "Code"}},
				},
				Elts: []ast.Expr{
					&ast.KeyValueExpr{
						Key:   &ast.SelectorExpr{X: &ast.Ident{Name: "pb"}, Sel: &ast.Ident{Name: "ErrorReason_TOKEN_REQUIRED"}},
						Value: &ast.SelectorExpr{X: &ast.Ident{Name: "codes"}, Sel: &ast.Ident{Name: "Unauthenticated"}},
					},
					&ast.KeyValueExpr{
						Key:   &ast.SelectorExpr{X: &ast.Ident{Name: "pb"}, Sel: &ast.Ident{Name: "ErrorReason_TOKEN_AUTH_INVALID"}},
						Value: &ast.SelectorExpr{X: &ast.Ident{Name: "codes"}, Sel: &ast.Ident{Name: "Unauthenticated"}},
					},
				},
			},
		},
	}
	// Add the ReasonToCode map to the main function body
	f.Decls = append(f.Decls, &ast.FuncDecl{
		Name: &ast.Ident{Name: "main"},
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				reasonToCode,
			},
		},
	})

	// Format the code
	ast.SortImports(fset, f)
	cfg := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	cfg.Fprint(os.Stdout, fset, f)

	// Format the code

}
