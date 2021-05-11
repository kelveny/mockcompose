package syntax

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"io"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

//
// Quick reference to GO AST syntactic objects that are used in this package
//
//	Node
//		Expr
//		Stmt
//		Decl
//			BadDecl
//			GenDecl (Tok filed has value of token.IMPORT, token.CONST, token.TYPE, token.VAR)
//				Spec(A Spec node represents a single (non-parenthesized) import, constant, type, or variable declaration)
//					ImportSpec
//					ValueSpec (const/var)
//					TypeSpec
//						*Ident, *ParenExpr, *SelectorExpr, *StarExpr
//						ArrayType
//						StructType
//						FuncType
//						InterfaceType
//						MapType
//						ChanType
//			FuncDecl
//
//
// 	Expr
// 		BadExpr
//		Ident
//		Ellipsis
// 		BasicLit (basic literal)
//		FuncLit
//		CompositeLit
//		ParenExpr
//		SelectorExpr
//		IndexExpr
//		SliceExpr
//		TypeAssertExpr
//		CallExpr
//		StarExpr
//		UnaryExpr
//		BinaryExpr
//		KeyValueExpr
//
/*
	// FieldList contains list of fields
	// A FieldList represents a list of Fields, enclosed by parentheses or braces.
	//
	// struct fields (Names/Type)
	// interface fields (Names -> method name)
	// unamed embedded field (Names -> embedded Type name)
	// fields of param/return list (Names nil means unamed, otherwise name of the parameter)
	type Field struct {
		Doc     *CommentGroup // associated documentation; or nil
		Names   []*Ident      // field/method/parameter names; or nil
		Type    Expr          // field/method/parameter type
		Tag     *BasicLit     // field tag; or nil
		Comment *CommentGroup // line comments; or nil
	}
*/

const (
	returnFieldTemplate = ` 
	_mc_ret := {{ .MockCallExpr }}
	{{ range $index, $f := .Fields }}
	var _r{{ $index }} {{ $f.Typ }}

	if _rfn, ok := _mc_ret.Get({{ $index }}).({{ $f.TypeFuncDecl }}); ok {
		_r{{ $index }} = _rfn({{ $.FuncInvokeParamsExpr }})
	} else {	
	{{- if isErrorType $f }}
		_r{{ $index }} = _mc_ret.Error({{ $index }})
	{{- else }}
		if _mc_ret.Get({{ $index }}) != nil {
			_r{{ $index }} = _mc_ret.Get({{ $index }}).({{ $f.Typ }})
		}
	{{- end }}
	}
	{{ end }}
	return {{ join . }}
`
)

// ForEachDecl iterates all AST Decl objects in a AST syntactic package
func ForEachDecl(
	p *packages.Package,
	do func(*packages.Package, ast.Decl),
) {
	for _, astFile := range p.Syntax {
		for _, d := range astFile.Decls {
			do(p, d)
		}
	}
}

// ForEachFuncDecl iterates all AST FuncDecl objects in a AST syntactic package
func ForEachFuncDecl(
	p *packages.Package,
	do func(*packages.Package, *ast.FuncDecl),
) {
	for _, astFile := range p.Syntax {
		for _, d := range astFile.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok {
				do(p, fn)
			}
		}
	}
}

// ExprDeclString returns declarative string for AST Expr object
func ExprDeclString(fset *token.FileSet, e ast.Expr) string {
	var b bytes.Buffer
	format.Node(&b, fset, e)
	return b.String()
}

// ParamListDeclString returns declarative string for AST FieldList object
func ParamListDeclString(fset *token.FileSet, fl *ast.FieldList) string {
	names := []string{}

	if fl.NumFields() > 0 {
		for _, field := range fl.List {
			name := ""
			if field.Names != nil {
				name = field.Names[0].Name
				name = name + " " + ExprDeclString(fset, field.Type)
			} else {
				name = ExprDeclString(fset, field.Type)
			}

			names = append(names, name)
		}
	}
	return strings.Join(names, ", ")
}

// ParamListTypeOnlyDeclString returns declarative string for AST FieldList object
// compared with ParamListDeclString, the returned declarative string contains only type
// declarations
func ParamListTypeOnlyDeclString(fset *token.FileSet, fl *ast.FieldList) string {
	names := []string{}

	if fl.NumFields() > 0 {
		for _, field := range fl.List {
			names = append(names, ExprDeclString(fset, field.Type))
		}
	}
	return strings.Join(names, ", ")
}

type FieldDeclInfo struct {
	Name     string
	Typ      string
	Variadic bool
}

func ParamListDeclInfo(fset *token.FileSet, fl *ast.FieldList) []FieldDeclInfo {
	infos := []FieldDeclInfo{}
	if fl != nil && fl.NumFields() > 0 {
		for _, f := range fl.List {
			name := ""
			if f.Names != nil {
				name = f.Names[0].Name
			}

			typ := ExprDeclString(fset, f.Type)

			variadic := false
			if strings.HasPrefix(typ, "...") {
				variadic = true
			}

			infos = append(infos, FieldDeclInfo{
				Name:     name,
				Typ:      typ,
				Variadic: variadic,
			})
		}
	}

	return infos
}

func ParamListInvokeString(params []FieldDeclInfo) string {
	s := []string{}
	for _, paramInfo := range params {
		if paramInfo.Variadic {
			s = append(s, paramInfo.Name+"...")
		} else {
			s = append(s, paramInfo.Name)
		}
	}

	return strings.Join(s, ", ")
}

func ReceiverDeclString(fset *token.FileSet, fl *ast.FieldList) string {
	if fl != nil && fl.NumFields() > 0 {
		return fmt.Sprintf("(%s)", ParamListDeclString(fset, fl))
	}
	return ""
}

// ReturnDeclString returns declative string for the given AST field list
func ReturnDeclString(fset *token.FileSet, fl *ast.FieldList) string {
	if fl != nil && fl.NumFields() > 1 {
		return fmt.Sprintf("(%s)", ParamListDeclString(fset, fl))
	} else if fl != nil && fl.NumFields() > 0 {
		if fl.List[0].Names != nil {
			return fmt.Sprintf("(%s)", ParamListDeclString(fset, fl))
		} else {
			return ParamListDeclString(fset, fl)
		}
	}
	return ""
}

func FuncDeclString(fset *token.FileSet, fn *ast.FuncDecl) string {
	recvr := ReceiverDeclString(fset, fn.Recv)

	if recvr != "" {
		return fmt.Sprintf("func %s %s(%s) %s",
			ReceiverDeclString(fset, fn.Recv),
			fn.Name.Name,
			ParamListDeclString(fset, fn.Type.Params),
			ReturnDeclString(fset, fn.Type.Results),
		)
	} else {
		return fmt.Sprintf("func %s(%s) %s",
			fn.Name.Name,
			ParamListDeclString(fset, fn.Type.Params),
			ReturnDeclString(fset, fn.Type.Results),
		)
	}
}

// generate _m.Called() expression (calling into testify/mock.Called() method)
func generateMockDotCalledExpr(fset *token.FileSet, fn *ast.FuncDecl) (string, string) {
	paramInfos := ParamListDeclInfo(fset, fn.Type.Params)

	if len(paramInfos) == 0 {
		return "_m.Called()", ""
	}

	lastParam := paramInfos[len(paramInfos)-1]
	if !lastParam.Variadic {
		return fmt.Sprintf("_m.Called(%s)", ParamListInvokeString(paramInfos)), ""
	}

	if lastParam.Typ == "...interface{}" {
		return fmt.Sprintf("_m.Called(%s)", ParamListInvokeString(paramInfos)), ""
	}

	// testify/mock.Called() accepts ...interface{}, if underlying variadic
	// type is not interface{}, perform type-casting
	setupBlock := fmt.Sprintf(`
	_mc_VArg := make([]interface{}, len(%s))
	for _mc_i := range %s {
		_mc_VArg[_mc_i] = %s[_mc_i]
	}
`, lastParam.Name, lastParam.Name, lastParam.Name)

	lastParam.Name = "_mc_VArg"
	return fmt.Sprintf("_m.Called(%s)", ParamListInvokeString(paramInfos)), setupBlock
}

type ReturnFieldBindingSpec struct {
	Name         string
	Typ          string
	TypeFuncDecl string
}

type ReturnFieldBinding struct {
	FuncInvokeParamsExpr string
	MockCallExpr         string
	Fields               []ReturnFieldBindingSpec
}

func buildReturnFieldBinding(
	fset *token.FileSet,
	fn *ast.FuncDecl,
) *ReturnFieldBinding {
	fields := []ReturnFieldBindingSpec{}

	for _, f := range fn.Type.Results.List {
		name := ""
		if f.Names != nil {
			name = f.Names[0].Name
		}
		typ := ExprDeclString(fset, f.Type)

		fields = append(fields, ReturnFieldBindingSpec{
			Name: name,
			Typ:  typ,
			TypeFuncDecl: fmt.Sprintf("func(%s) %s",
				ParamListTypeOnlyDeclString(fset, fn.Type.Params),
				typ,
			),
		})
	}

	return &ReturnFieldBinding{
		Fields: fields,
	}
}

// MockFunc generates a mocking method on mockClz class for a give AST FuncDecl object
func MockFunc(writer io.Writer, mockClz string, fset *token.FileSet, fn *ast.FuncDecl) {
	retDecl := ReturnDeclString(fset, fn.Type.Results)
	if retDecl != "" {
		fmt.Fprintf(
			writer, "func (_m *%s) %s(%s) %s {\n",
			mockClz,
			fn.Name.Name,
			ParamListDeclString(fset, fn.Type.Params),
			retDecl,
		)
	} else {
		fmt.Fprintf(
			writer, "func (_m *%s) %s(%s) {\n",
			mockClz,
			fn.Name.Name,
			ParamListDeclString(fset, fn.Type.Params),
		)
	}

	calledExpr, calledExprSetup := generateMockDotCalledExpr(fset, fn)
	fmt.Fprintf(writer, "%s", calledExprSetup)

	if len(fn.Type.Results.List) > 0 {
		paramInfos := ParamListDeclInfo(fset, fn.Type.Params)

		binding := buildReturnFieldBinding(fset, fn)
		binding.MockCallExpr = calledExpr
		binding.FuncInvokeParamsExpr = ParamListInvokeString(paramInfos)

		t := template.Must(template.New("MockCompose").
			Funcs(template.FuncMap{
				"isErrorType": func(spec ReturnFieldBindingSpec) bool {
					return spec.Typ == "error"
				},

				"join": func(binding *ReturnFieldBinding) string {
					s := []string{}
					for i := range binding.Fields {
						s = append(s, fmt.Sprintf("_r%d", i))
					}

					return strings.Join(s, ", ")
				},
			}).
			Parse(returnFieldTemplate))
		t.Execute(writer, binding)
	} else {
		fmt.Fprintf(writer, "\n\t%s\n", calledExpr)
	}

	fmt.Fprintf(writer, "\n}\n")
}
