package gosyntax

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

// ForEachDeclInPackage iterates all AST Decl objects in a AST syntactic package
func ForEachDeclInPackage(
	p *packages.Package,
	do func(ast.Decl),
) {
	for _, astFile := range p.Syntax {
		for _, d := range astFile.Decls {
			do(d)
		}
	}
}

// ForEachFuncDeclInPackage iterates all AST FuncDecl objects in a AST syntactic package
func ForEachFuncDeclInPackage(
	p *packages.Package,
	do func(*ast.FuncDecl),
) {
	for _, astFile := range p.Syntax {
		for _, d := range astFile.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok {
				do(fn)
			}
		}
	}
}

func ForEachInterfaceDeclInPackage(
	p *packages.Package,
	do func(name string, methods []*ast.Field),
) {
	for _, astFile := range p.Syntax {
		ForEachInterfaceDeclInFile(astFile, do)
	}
}

func ForEachFuncDeclInFile(
	file *ast.File,
	do func(*ast.FuncDecl),
) {
	for _, d := range file.Decls {
		if fn, ok := d.(*ast.FuncDecl); ok {
			do(fn)
		}
	}
}

func ForEachInterfaceDeclInFile(file *ast.File,
	do func(name string, methods []*ast.Field),
) {
	for _, d := range file.Decls {
		if gd, ok := d.(*ast.GenDecl); ok {
			for _, spec := range gd.Specs {
				if tspec, ok := spec.(*ast.TypeSpec); ok {
					if intf, ok := tspec.Type.(*ast.InterfaceType); ok {
						do(tspec.Name.Name, intf.Methods.List)
					}
				}
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
				names := []string{}
				for _, n := range field.Names {
					names = append(names, n.Name)
				}

				name = strings.Join(names, ",") + " " + ExprDeclString(fset, field.Type)
			} else {
				name = ExprDeclString(fset, field.Type)
			}

			names = append(names, name)
		}
	}
	return strings.Join(names, ", ")
}

// ParamInfoListTypeOnlyDeclString returns declarative string from parameter field declarations,
// the returned declarative string contains only type declarations
func ParamInfoListTypeOnlyDeclString(paramInfos []*FieldDeclInfo) string {
	names := []string{}

	if len(paramInfos) > 0 {
		for _, field := range paramInfos {
			names = append(names, field.Typ)
		}
	}
	return strings.Join(names, ", ")
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

// unify AST field list with simplified struct. Original thought was
// to give an common interface for syntax-based and type-based implementations
type FieldDeclInfo struct {
	Name     string
	Typ      string
	Variadic bool
}

func ParamListDeclInfo(fset *token.FileSet, fl *ast.FieldList) []*FieldDeclInfo {
	infos := []*FieldDeclInfo{}
	if fl != nil && fl.NumFields() > 0 {
		for _, f := range fl.List {
			typ := ExprDeclString(fset, f.Type)

			variadic := false
			if strings.HasPrefix(typ, "...") {
				variadic = true
			}

			if f.Names != nil {
				for _, n := range f.Names {
					infos = append(infos, &FieldDeclInfo{
						Name:     n.Name,
						Typ:      typ,
						Variadic: variadic,
					})
				}
			} else {
				infos = append(infos, &FieldDeclInfo{
					Name:     "",
					Typ:      typ,
					Variadic: variadic,
				})
			}
		}
	}

	return infos
}

func ParamInfoListNameExists(paramInfos []*FieldDeclInfo, name string) bool {
	for _, param := range paramInfos {
		if param.Name == name {
			return true
		}
	}

	return false
}

func ParamInfoListFixup(paramInfos []*FieldDeclInfo) {
	i := 0
	for _, param := range paramInfos {
		if param.Name == "" {
			for {
				name := fmt.Sprintf("_a%d", i)
				if ParamInfoListNameExists(paramInfos, name) {
					i++
				} else {
					param.Name = name
					i++
					break
				}
			}
		}
	}
}

func ParamInfoListDeclString(params []*FieldDeclInfo) string {
	names := []string{}

	if len(params) > 0 {
		for _, field := range params {
			name := ""
			if field.Name != "" {
				name = field.Name + " " + field.Typ
			} else {
				name = field.Typ
			}

			names = append(names, name)
		}
	}
	return strings.Join(names, ", ")
}

func ParamInfoListInvokeString(params []*FieldDeclInfo) string {
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

// ReturnInfoListDeclString returns declative string for the given return info list
func ReturnInfoListDeclString(returns []*FieldDeclInfo) string {
	if len(returns) > 1 {
		return fmt.Sprintf("(%s)", ParamInfoListDeclString(returns))
	} else if len(returns) > 0 {
		if returns[0].Name != "" {
			return fmt.Sprintf("(%s)", ParamInfoListDeclString(returns))
		} else {
			return ParamInfoListDeclString(returns)
		}
	}
	return ""
}

// generate m.Called() expression (calling into testify/mock.Called() method)
func generateMockDotCalledExpr(
	paramInfos []*FieldDeclInfo,
) (string, string) {

	if len(paramInfos) == 0 {
		return "m.Called()", ""
	}

	lastParam := paramInfos[len(paramInfos)-1]
	if !lastParam.Variadic {
		return fmt.Sprintf("m.Called(%s)", ParamInfoListInvokeString(paramInfos)), ""
	}

	if lastParam.Typ == "...interface{}" && len(paramInfos) == 1 {
		return fmt.Sprintf("m.Called(%s...)", lastParam.Name), ""
	}

	// testify/mock.Called() accepts ...interface{}, for variadic parameters,
	// just convert it to slice
	lines := []string{}
	lines = append(lines, fmt.Sprintf(`
	_mc_args := make([]interface{}, 0, %d+len(%s))
	`, len(paramInfos)-1, lastParam.Name))

	for i := 0; i < len(paramInfos)-1; i++ {
		lines = append(lines, fmt.Sprintf(`
	_mc_args = append(_mc_args, %s)
	`, paramInfos[i].Name))
	}

	lines = append(lines, fmt.Sprintf(`
	for _, _va := range %s {
		_mc_args = append(_mc_args, _va)
	}
	`, lastParam.Name))

	setupBlock := strings.Join(lines, "")
	return "m.Called(_mc_args...)", setupBlock
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
	paramInfos []*FieldDeclInfo,
	returnInfos []*FieldDeclInfo,
) *ReturnFieldBinding {
	fields := []ReturnFieldBindingSpec{}

	for _, f := range returnInfos {
		fields = append(fields, ReturnFieldBindingSpec{
			Name: f.Name,
			Typ:  f.Typ,
			TypeFuncDecl: fmt.Sprintf("func(%s) %s",
				ParamInfoListTypeOnlyDeclString(paramInfos),
				f.Typ,
			),
		})
	}

	return &ReturnFieldBinding{
		Fields: fields,
	}
}

// MockFunc generates a mocking method on mockClz class
// generate mockery (https://github.com/vektra/mockery) compatible mocking implementation
// from syntax based declarations
//
func MockFunc(
	writer io.Writer,
	mockClz string,
	fset *token.FileSet,
	fnName string,
	fnParams *ast.FieldList,
	fnReturns *ast.FieldList,
) {
	paramInfos := ParamListDeclInfo(fset, fnParams)
	returnInfos := ParamListDeclInfo(fset, fnReturns)

	// FuncDecl of method definition from interface may come in unnamed
	// make sure that we name these parameters before code generation
	ParamInfoListFixup(paramInfos)

	retDecl := ReturnInfoListDeclString(returnInfos)
	if retDecl != "" {
		fmt.Fprintf(
			writer, "func (m *%s) %s(%s) %s {\n",
			mockClz,
			fnName,
			ParamInfoListDeclString(paramInfos),
			retDecl,
		)
	} else {
		fmt.Fprintf(
			writer, "func (m *%s) %s(%s) {\n",
			mockClz,
			fnName,
			ParamInfoListDeclString(paramInfos),
		)
	}

	calledExpr, calledExprSetup := generateMockDotCalledExpr(paramInfos)
	fmt.Fprintf(writer, "%s", calledExprSetup)

	if fnReturns != nil && len(fnReturns.List) > 0 {
		binding := buildReturnFieldBinding(paramInfos, returnInfos)
		binding.MockCallExpr = calledExpr
		binding.FuncInvokeParamsExpr = ParamInfoListInvokeString(paramInfos)

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
