package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"
	"os"
	"testing"
)

func TestGenerateCommentForPKG(t *testing.T) {
	fset := token.NewFileSet()

	// Create the AST for the given code
	file := &ast.File{
		Name: ast.NewIdent("grpc"),
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "// Package grpc ...\n"},
			},
		},
	}

	// Generate the code from the AST
	err := printer.Fprint(os.Stdout, fset, file)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

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

func TestCreateSpaceBK(t *testing.T) {
	// Create the AST for the given code.
	fset := token.NewFileSet()
	file := &ast.File{
		Name: ast.NewIdent("main"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.VAR,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ReasonToCode"),
						},
						Type: &ast.MapType{
							Key:   ast.NewIdent("pb.ErrorReason"),
							Value: ast.NewIdent("codes.Code"),
						},
						Values: []ast.Expr{
							&ast.CompositeLit{
								Type: &ast.MapType{
									Key:   ast.NewIdent("pb.ErrorReason"),
									Value: ast.NewIdent("codes.Code"),
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("pb.ErrorReason_TOKEN_REQUIRED"),
										Value: ast.NewIdent("codes.Unauthenticated"),
									},
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("pb.ErrorReason_TOKEN_AUTH_INVALID"),
										Value: ast.NewIdent("codes.Unauthenticated"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Format the AST with separate line declarations for key-value pairs.
	var output []byte
	ast.Inspect(file, func(n ast.Node) bool {
		if n != nil {
			if pos := fset.Position(n.Pos()); pos.Line > 0 {
				output = append(output, []byte(fmt.Sprintf("\n// Line %d\n", pos.Line))...)
			}
		}
		return true
	})

	println(fmt.Sprintf(">>>>> %s", output))
	err := format.Node(os.Stdout, fset, file)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestUsersOfSlice(t *testing.T) {

	//tests := []*User{&User{Name: "vortex", Age: 12}}
	// Create the AST by parsing the source code.
	// Create an empty AST
	fset := token.NewFileSet()
	f := &ast.File{
		Name:  ast.NewIdent("main"),
		Decls: []ast.Decl{},
	}

	// Add the User struct type to the AST.
	userStruct := &ast.TypeSpec{
		Name: ast.NewIdent("User"),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("Name")},
						Type:  ast.NewIdent("string"),
					},
					{
						Names: []*ast.Ident{ast.NewIdent("Age")},
						Type:  ast.NewIdent("int"),
					},
				},
			},
		},
	}
	f.Decls = append(f.Decls, &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{userStruct},
	})

	// Add the Test function to the AST.
	testFunc := &ast.FuncDecl{
		Name: ast.NewIdent("Test"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("users")},
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("*User"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{ast.NewIdent("nil")},
				},
			},
		},
	}
	f.Decls = append(f.Decls, testFunc)

	// Add the main function to the AST.
	mainFunc := &ast.FuncDecl{
		Name: ast.NewIdent("main"),
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("tests")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.ArrayType{
								Elt: ast.NewIdent("*User"),
							},
							Elts: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: ast.NewIdent("User"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Name"),
												Value: &ast.BasicLit{Kind: token.STRING, Value: "\"vortex\""},
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Age"),
												Value: &ast.BasicLit{Kind: token.INT, Value: "12"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	f.Decls = append(f.Decls, mainFunc)

	ast.SortImports(fset, f)
	cfg := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	cfg.Fprint(os.Stdout, fset, f)

}

func TestCommentForConst(t *testing.T) {
	// Create the AST for the given code
	fset := token.NewFileSet()
	f := &ast.File{
		Name: ast.NewIdent("main"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.CONST,
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "// my const"},
					},
				},
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent("A")},
						Doc: &ast.CommentGroup{
							List: []*ast.Comment{
								{Text: "\n// data \n"},
							},
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.INT,
								Value: "1",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{ast.NewIdent("B")},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.INT,
								Value: "2",
							},
						},
					},
				},
			},
		},
	}

	cfg := &printer.Config{Mode: printer.UseSpaces, Tabwidth: 4}
	cfg.Fprint(os.Stdout, fset, f)
}

func TestCreateStructWithMethods(t *testing.T) {

	fset := token.NewFileSet()

	// Create the main package
	pkg := &ast.File{
		Name: ast.NewIdent("main"),
	}

	// Create the handler struct
	handlerType := &ast.TypeSpec{
		Name: ast.NewIdent("handler"),
		Type: &ast.StructType{
			Fields: &ast.FieldList{},
		},
	}

	// Add the handler struct to the main package
	pkg.Decls = append(pkg.Decls, &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			handlerType,
		},
	})

	// Create the MyMethod1 method
	myMethod1 := &ast.FuncDecl{
		Name: ast.NewIdent("MyMethod1"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("h")},
					Type:  &ast.StarExpr{X: ast.NewIdent("handler")},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{},
	}

	// Add the MyMethod1 method to the main package
	pkg.Decls = append(pkg.Decls, myMethod1)

	// Create the MyMethod2 method
	myMethod2 := &ast.FuncDecl{
		Name: ast.NewIdent("MyMethod2"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("h")},
					Type:  &ast.StarExpr{X: ast.NewIdent("handler")},
				},
			},
		},
		Type: &ast.FuncType{
			Params:  &ast.FieldList{},
			Results: nil,
		},
		Body: &ast.BlockStmt{},
	}

	// Add the MyMethod2 method to the main package
	pkg.Decls = append(pkg.Decls, myMethod2)

	// Print the AST to Go source code
	format.Node(os.Stdout, fset, pkg)

}
