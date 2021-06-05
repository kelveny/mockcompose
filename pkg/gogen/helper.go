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

	"github.com/kelveny/mockcompose/pkg/gosyntax"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

type ImportSpec struct {
	Name string
	Path string
}

func (s *ImportSpec) IsNameDefault() bool {
	return s.Name == s.Path || strings.HasSuffix(s.Path, "/"+s.Name)
}

func GetPackageImports(pkg *packages.Package) []ImportSpec {
	var specs []ImportSpec

	for _, p := range pkg.Imports {
		specs = append(specs, ImportSpec{
			Name: p.Name,
			Path: p.PkgPath,
		})
	}

	return specs
}

func GetFileImports(file *ast.File) []ImportSpec {
	var specs []ImportSpec

	for _, s := range file.Imports {
		n := ""
		if s.Name != nil {
			n = s.Name.Name
		}

		specs = append(specs, ImportSpec{
			Name: n,
			Path: strings.Trim(s.Path.Value, "\""),
		})
	}

	return specs
}

func LookupInPackageScope(pkg *packages.Package, name string) types.Object {
	return pkg.Types.Scope().Lookup(name)
}

func AppendImportSpec(specs []ImportSpec, name, p string) []ImportSpec {
	if strings.HasSuffix(p, "/"+name) || name == p {
		name = ""
	}

	appearsAsNew := true
	if len(specs) > 0 {
		for _, spec := range specs {
			if spec.Name == name && spec.Path == p {
				appearsAsNew = false
				break
			}
		}
	}

	if appearsAsNew {
		return append(specs, ImportSpec{
			Name: name,
			Path: p,
		})
	}

	return specs
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

func CleanImports(f *ast.File, cleanedImports []ImportSpec) []ImportSpec {
	if len(f.Imports) > 0 {
		for _, imp := range f.Imports {
			// Path.Value has been quoted, remove it
			if astutil.UsesImport(f, strings.Trim(imp.Path.Value, "\"")) {
				var name, p string
				if imp.Name != nil {
					name = imp.Name.Name
				}
				p = strings.Trim(imp.Path.Value, "\"")
				cleanedImports = AppendImportSpec(cleanedImports, name, p)
			}
		}
	}

	return cleanedImports
}

func WriteImportDecls(writer io.Writer, imports []ImportSpec) {
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
			fmt.Fprintf(
				writer,
				"func %s",
				fnName,
			)
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
