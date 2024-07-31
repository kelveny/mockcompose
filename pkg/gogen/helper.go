package gogen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"io"
	"sort"
	"strings"
	"text/template"

	"github.com/kelveny/mockcompose/pkg/gosyntax"
	"github.com/kelveny/mockcompose/pkg/gotype"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

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

func GetPackageImports(pkg *packages.Package) []gosyntax.ImportSpec {
	var specs []gosyntax.ImportSpec

	if pkg.Imports != nil {
		for _, p := range pkg.Imports {
			specs = append(specs, gosyntax.ImportSpec{
				Name: p.Name,
				Path: p.PkgPath,
			})
		}
	} else {
		if pkg.Types != nil {
			for _, p := range pkg.Types.Imports() {
				specs = append(specs, gosyntax.ImportSpec{
					Name: p.Name(),
					Path: p.Path(),
				})
			}
		}
	}

	return specs
}

func LookupInPackageScope(pkg *packages.Package, name string) types.Object {
	return pkg.Types.Scope().Lookup(name)
}

// SetImportSpecs sets ImportSpecs in f's Decls section using the passed ImportSpecs
func SetImportSpecs(f *ast.File, imports []*ast.ImportSpec) {
	if len(f.Decls) > 0 {
		for _, d := range f.Decls {
			if dd, ok := d.(*ast.GenDecl); ok && dd.Tok == token.IMPORT {
				var specs []ast.Spec
				if len(imports) > 0 {
					for _, spec := range imports {
						specs = append(specs, spec)
					}
				}
				dd.Specs = specs
			}
		}
	}
}

func CleanImports(f *ast.File, cleanedImports []gosyntax.ImportSpec) []gosyntax.ImportSpec {
	if len(f.Imports) > 0 {
		for _, imp := range f.Imports {
			// Path.Value has been quoted, remove it
			if astutil.UsesImport(f, strings.Trim(imp.Path.Value, "\"")) {
				var name, p string
				if imp.Name != nil {
					name = imp.Name.Name
				}
				p = strings.Trim(imp.Path.Value, "\"")
				cleanedImports = gosyntax.AppendImportSpec(cleanedImports, name, p)
			}
		}
	}

	return cleanedImports
}

func WriteImportDecls(writer io.Writer, imports []gosyntax.ImportSpec) {
	if len(imports) > 0 {
		fmt.Fprintln(writer, "import (")

		for _, spec := range imports {
			if spec.Name != "" && !spec.IsNameDefault() {
				fmt.Fprintf(writer, "\t%s \"%s\"\n", spec.Name, spec.Path)
			} else {
				fmt.Fprintf(writer, "\t\"%s\"\n", spec.Path)
			}
		}
		fmt.Fprintln(writer, ")")
	}

	fmt.Fprintln(writer)
}

func WriteFuncDecls(
	writer io.Writer,
	fset *token.FileSet,
	file *ast.File,
) {
	if len(file.Decls) > 0 {
		for _, d := range file.Decls {
			if fnSpec, ok := d.(*ast.FuncDecl); ok {
				format.Node(writer, fset, fnSpec)
				writer.Write([]byte("\n\n"))
			}
		}
	}
}

func WriteFuncWithLocalOverrides(
	writer io.Writer,
	fset *token.FileSet,
	fnSpec *ast.FuncDecl,
	receiverDecl string,
	fnName string,
	overrides map[string]string,
) {
	if len(overrides) == 0 {
		format.Node(writer, fset, fnSpec)
		writer.Write([]byte("\n\n"))
	} else {
		if fnSpec.Recv != nil {
			fmt.Fprintf(
				writer,
				"func (%s) %s",
				gosyntax.ParamListDeclString(fset, fnSpec.Recv),
				fnName,
			)
		} else {
			if receiverDecl != "" {
				fmt.Fprintf(
					writer,
					"func %s %s",
					receiverDecl,
					fnName,
				)
			} else {
				fmt.Fprintf(
					writer,
					"func %s",
					fnName,
				)
			}
		}

		var b bytes.Buffer
		// fnSpec.Type -> func(params) (rets)
		format.Node(&b, fset, fnSpec.Type)

		// use everything after func
		fmt.Fprint(writer, string(b.Bytes()[4:]))

		b.Reset()
		format.Node(&b, fset, fnSpec.Body)
		end := len(b.Bytes()) - 1
		body := string(b.Bytes()[1:end])

		fmt.Fprint(writer, " {")
		generateLocalOverrides(writer, overrides)
		fmt.Fprint(writer, body)
		fmt.Fprintln(writer, "}")
	}
}

func generateLocalOverrides(writer io.Writer, overrides map[string]string) {
	if len(overrides) > 0 {
		keys := make([]string, len(overrides))
		i := 0
		for k := range overrides {
			keys[i] = k
			i++
		}

		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(writer, `
	%s := %s`, k, overrides[k])
		}

		fmt.Fprintf(writer, "\n")
	}
}

// generate m.Called() expression (calling into testify/mock.Called() method)
func generateMockDotCalledExpr(
	paramInfos []*gosyntax.FieldDeclInfo,
) (string, string) {

	if len(paramInfos) == 0 {
		return "m.Called()", ""
	}

	lastParam := paramInfos[len(paramInfos)-1]
	if !lastParam.Variadic {
		return fmt.Sprintf("m.Called(%s)", gosyntax.ParamInfoListInvokeString(paramInfos)), ""
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
	paramInfos []*gosyntax.FieldDeclInfo,
	returnInfos []*gosyntax.FieldDeclInfo,
) *ReturnFieldBinding {
	fields := []ReturnFieldBindingSpec{}

	for _, f := range returnInfos {
		fields = append(fields, ReturnFieldBindingSpec{
			Name: f.Name,
			Typ:  f.Typ,
			TypeFuncDecl: fmt.Sprintf("func(%s) %s",
				gosyntax.ParamInfoListTypeOnlyDeclString(paramInfos),
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
func MockFunc(
	writer io.Writer,
	mockPkg string,
	mockClz string,
	fset *token.FileSet,
	fnName string,
	fnParams *ast.FieldList,
	fnReturns *ast.FieldList,
	signature *types.Signature,
) {
	paramInfos := gosyntax.ParamListDeclInfo(fset, fnParams)
	returnInfos := gosyntax.ParamListDeclInfo(fset, fnReturns)

	GenerateFuncMock(writer, mockPkg, mockClz, fnName, paramInfos, returnInfos, signature)
}

// GenerateFuncMock generates function mock implementation based on FieldDeclInfo
// abstraction
func GenerateFuncMock(
	writer io.Writer,
	mockPkg string,
	mockClz string,
	fnName string,
	paramInfos []*gosyntax.FieldDeclInfo,
	returnInfos []*gosyntax.FieldDeclInfo,
	signature *types.Signature,
) {
	if signature != nil {
		// override with inferred info from type signature
		p := gotype.GetFuncParamInfosFromSignature(signature, mockPkg)
		for index, info := range p {
			if !strings.Contains(info.Typ, "invalid type") {
				paramInfos[index] = p[index]
			}
		}

		p = gotype.GetFuncReturnInfosFromSignature(signature, mockPkg)
		for index, info := range p {
			if !strings.Contains(info.Typ, "invalid type") {
				returnInfos[index] = p[index]
			}
		}
	}

	// FuncDecl of method definition from interface may come in unnamed
	// make sure that we name these parameters before code generation
	gosyntax.ParamInfoListFixup(paramInfos)

	retDecl := gosyntax.ReturnInfoListDeclString(returnInfos)
	if retDecl != "" {
		fmt.Fprintf(
			writer, "func (m *%s) %s(%s) %s {\n",
			mockClz,
			fnName,
			gosyntax.ParamInfoListDeclString(paramInfos),
			retDecl,
		)
	} else {
		fmt.Fprintf(
			writer, "func (m *%s) %s(%s) {\n",
			mockClz,
			fnName,
			gosyntax.ParamInfoListDeclString(paramInfos),
		)
	}

	calledExpr, calledExprSetup := generateMockDotCalledExpr(paramInfos)
	fmt.Fprintf(writer, "%s", calledExprSetup)

	if len(returnInfos) > 0 {
		binding := buildReturnFieldBinding(paramInfos, returnInfos)
		binding.MockCallExpr = calledExpr
		binding.FuncInvokeParamsExpr = gosyntax.ParamInfoListInvokeString(paramInfos)

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
