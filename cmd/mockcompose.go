package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

const (
	MATCH_NONE matchType = iota
	MATCH_CLONE
	MATCH_MOCK
)

type matchType int

type generator struct {
	clzName     string // name of the class that implements class interface
	mockClzName string // name of the class that implements the same class interface in mocking
	mockPkgName string // package name that mocking class resides
	mockName    string // the mocking composite class name

	mothodsToClone []string // method function names that need to be cloned in mocking class
	methodsToMock  []string // method function names that need to be mocked
}

type importSpec struct {
	name string
	path string
}

type stringSlice []string

func (ss *stringSlice) String() string {
	return strings.Join(*ss, ",")
}

func (ss *stringSlice) Set(val string) error {
	*ss = append(*ss, val)
	return nil
}

var (
	header = `//
// CODE GENERATED AUTOMATICALLY WITH github.com/kelveny/mockcompose
// THIS FILE SHOULD NOT BE EDITED BY HAND
//
package %s
`
	compositeClzTemplate = `type %s struct {
	%s
	%s
}

`
)

// match checks if a FuncDecl matches condition, if it is matched,
// change and prepare receiver type to the mocking target class
func (g *generator) match(fnSpec *ast.FuncDecl) (bool, matchType) {
	if fnSpec.Recv != nil {
		t := fnSpec.Recv.List[0].Type

		if expr, ok := t.(*ast.StarExpr); ok {
			recvrClzName := expr.X.(*ast.Ident).Name
			if g.clzName == recvrClzName {
				if matchType := g.matchMethod(fnSpec.Name.Name); matchType != MATCH_NONE {
					expr.X.(*ast.Ident).Name = g.mockName
					return true, matchType
				}
			}
		} else {
			recvrClzName := t.(*ast.Ident).Name
			if g.clzName == recvrClzName {
				if matchType := g.matchMethod(fnSpec.Name.Name); matchType != MATCH_NONE {
					t.(*ast.Ident).Name = g.mockName
					return true, matchType
				}
			}
		}
	}

	return false, MATCH_NONE
}

func (g *generator) matchMethod(fnName string) matchType {
	if len(g.mothodsToClone) > 0 {
		for _, name := range g.mothodsToClone {
			if name == fnName {
				return MATCH_CLONE
			}
		}
	}

	if len(g.methodsToMock) > 0 {
		for _, name := range g.methodsToMock {
			if name == fnName {
				return MATCH_MOCK
			}
		}
	}

	return MATCH_NONE
}

func (g *generator) composeMock(fnSpec *ast.FuncDecl, writer io.Writer) {
	// TODO we may directly generate mocking implementation later
	// to remove dependencies to mockery tool
}

func addImportSpec(specs []importSpec, name, p string) []importSpec {
	if strings.HasSuffix(p, "/"+name) || name == p {
		name = ""
	}

	appearsAsNew := true
	if len(specs) > 0 {
		for _, spec := range specs {
			if spec.name == name && spec.path == p {
				appearsAsNew = false
				break
			}
		}
	}

	if appearsAsNew {
		return append(specs, importSpec{
			name: name,
			path: p,
		})
	}

	return specs
}

func (g *generator) generate(
	file *ast.File,
	mockFile *ast.File,
	writer io.Writer) error {

	var buf bytes.Buffer

	fset := token.NewFileSet()
	if g.generateInternal(fset, file, &buf) {
		// reload generated content
		f, err := parser.ParseFile(fset, "", buf.Bytes(), parser.ParseComments)
		if err != nil {
			return err
		}

		// remove unused imports
		var cleanedImports []importSpec
		if len(f.Imports) > 0 {
			for _, imp := range f.Imports {
				// Path.Value has been quoted, remove it
				if astutil.UsesImport(f, strings.Trim(imp.Path.Value, "\"")) {
					var name, p string
					if imp.Name != nil {
						name = imp.Name.Name
					}
					p = strings.Trim(imp.Path.Value, "\"")
					cleanedImports = addImportSpec(cleanedImports, name, p)
				}
			}
		}

		//
		// Generate mocking methods using an already-generated mock implementation
		// class with a cloned generator that has different configuration
		//
		// TODO we may directly generate mocking methods later
		//
		var (
			g2    *generator
			f2    *ast.File
			buf2  bytes.Buffer
			fset2 *token.FileSet
		)

		if mockFile != nil {
			fset2 = token.NewFileSet()
			g2 = &generator{
				clzName:        g.mockClzName,
				mockClzName:    "",
				mockPkgName:    g.mockPkgName,
				mockName:       g.mockName,
				mothodsToClone: g.methodsToMock,
				methodsToMock:  nil,
			}
			g2.generateInternal(fset2, mockFile, &buf2)
			f2, _ = parser.ParseFile(fset2, "", buf2.Bytes(), parser.ParseComments)

			if len(f2.Imports) > 0 {
				for _, imp := range f2.Imports {
					// Path.Value has been quoted, remove it
					if astutil.UsesImport(f2, strings.Trim(imp.Path.Value, "\"")) {
						var name, p string
						if imp.Name != nil {
							name = imp.Name.Name
						}
						p = strings.Trim(imp.Path.Value, "\"")
						cleanedImports = addImportSpec(cleanedImports, name, p)
					}
				}
			}
		}

		// compose final output
		fmt.Fprintf(writer, header, g.mockPkgName)

		writeImportDecls(writer, cleanedImports)
		fmt.Fprintf(writer, compositeClzTemplate, g.mockName, g.clzName, g.mockClzName)

		g.writeFuncDecls(fset, f, writer)

		if mockFile != nil {
			g2.writeFuncDecls(fset2, f2, writer)
		}
	}

	return nil
}

func writeImportDecls(writer io.Writer, imports []importSpec) {
	if len(imports) > 0 {
		fmt.Fprintln(writer, "import (")

		for _, spec := range imports {
			if spec.name != "" {
				fmt.Fprintf(writer, "\t%s \"%s\"\n", spec.name, spec.path)
			} else {
				fmt.Fprintf(writer, "\t\"%s\"\n", spec.path)
			}
		}
		fmt.Fprintln(writer, ")")
	}

	fmt.Fprintln(writer)
}

func (g *generator) writeFuncDecls(
	fset *token.FileSet,
	file *ast.File,
	writer io.Writer,
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

// setImportSpecs sets ImportSpecs in f's Decls section using the passed ImportSpecs
func setImportSpecs(f *ast.File, imports []*ast.ImportSpec) {
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

func (g *generator) generateInternal(
	fset *token.FileSet,
	file *ast.File,
	writer io.Writer,
) bool {
	found := false

	writer.Write([]byte(fmt.Sprintf("package %s\n\n", g.mockPkgName)))

	if len(file.Decls) > 0 {
		for _, d := range file.Decls {
			if fnSpec, ok := d.(*ast.FuncDecl); ok {
				matched, matchType := g.match(fnSpec)
				if matched {
					found = true
				}

				if matchType == MATCH_CLONE {
					// clone receiver-modified method
					format.Node(writer, fset, fnSpec)
					writer.Write([]byte("\n\n"))
				} else if matchType == MATCH_MOCK {
					// generate mocked method
					g.composeMock(fnSpec, writer)
				}
			} else {
				// for any non-function declaration, export only imports
				if dd, ok := d.(*ast.GenDecl); ok && dd.Tok == token.IMPORT {
					format.Node(writer, fset, d)
					writer.Write([]byte("\n\n"))
				}
			}
		}
	}

	return found
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: %s [-help] [options]

mockcompose generate a composite mocking class implementation.
`, os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func derivePackage() string {
	path, err := filepath.Abs("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in accessing file system. error: %s\n", err)
		os.Exit(1)
	}
	path = stripGopath(path)

	return path
}

func stripGopath(p string) string {
	gopathConfig := getGoPathConfig()

	for _, gopath := range strings.Split(gopathConfig, string(filepath.ListSeparator)) {
		fmt.Printf("check gopath: %s\n", gopath)
		p = strings.Replace(p, gopath+"/", "", 1)
	}

	p = strings.Replace(p, "src/", "", 1)
	p = strings.Replace(p, "pkg/", "", 1)

	return p
}

func getGoPathConfig() string {
	gopathConfig := os.Getenv("GOPATH")

	pkgRoot, _ := filepath.Abs("")

	var detected string
	if pkgPos := strings.LastIndex(pkgRoot, "/pkg"); pkgPos >= 0 {
		detected = pkgRoot[0:pkgPos]

		if strings.Index(gopathConfig, detected) < 0 {
			gopathConfig = gopathConfig + string(filepath.ListSeparator) + detected
		}
	}

	if pkgPos := strings.LastIndex(pkgRoot, "/src"); pkgPos >= 0 {
		detected2 := pkgRoot[0:pkgPos]

		if strings.Index(gopathConfig, detected2) < 0 {
			return gopathConfig + string(filepath.ListSeparator) + detected2
		}
	}

	return gopathConfig
}

func formatGoFile(filePath string) {
	b, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Printf("Error in reading file %s, error: %s", filePath, err)
		return
	}

	bb, err := format.Source(b)
	if err != nil {
		fmt.Printf("Error in formatting Go source %s, error: %s", filePath, err)
		return
	}

	ioutil.WriteFile(filePath, bb, 0644)
}

func Execute() {
	var methodsToClone stringSlice
	var methodsToMock stringSlice

	pkg := flag.String("pkg", "", "Name of the package for which to generate composite mocking implementation")
	clzName := flag.String("c", "", "Name of the class for which to generate composite mocking implementation")
	mockClzName := flag.String("m", "", "Name of the mocking class that implements the same interface as -c specified class")
	mockClzFileName := flag.String("mm", "", "Name of the file that implements -m specified mocking class")
	mockName := flag.String("n", "", "Name of the class(struct) with composite mocking implementation")
	flag.Var(&methodsToClone, "real", "Name of the method function name to be cloned from class implementation")
	flag.Var(&methodsToMock, "mock", "Name of the method function name to be mocked")

	flag.Parse()

	if *pkg == "" {
		*pkg = derivePackage()

		fmt.Printf("Derive package name as: %s\n", *pkg)
	}
	fmt.Println()

	if *mockName == "" {
		fmt.Fprintf(os.Stderr, "Please specify composite mocking class name with -n option\n")
		os.Exit(1)
	}

	if *clzName == "" {
		fmt.Fprintf(os.Stderr, "Please specify the class name for which to generate composite mocking implementation with -c option\n")
		os.Exit(1)
	}

	if *mockClzName == "" {
		fmt.Fprintf(os.Stderr, "Please specify the mocking class name that implements the same interface as -c specified class with -m option\n")
		os.Exit(1)
	}

	if len(methodsToClone) == 0 {
		fmt.Fprintf(os.Stderr, "Please specify at least one real method name with -real option\n")
		os.Exit(1)
	}

	if len(methodsToMock) > 0 && *mockClzFileName == "" {
		fmt.Fprintf(os.Stderr, "Please specify the name of the file that implements -m specified mocking class\n")
		os.Exit(1)
	}

	// iterate candidates from package directory
	gopathConfig := getGoPathConfig()

	for _, gopath := range strings.Split(gopathConfig, string(filepath.ListSeparator)) {
		// support scanning of subfolder src/ and pkg/
		for _, subFolder := range []string{"src", "pkg"} {
			pkgDir, err := filepath.Abs(path.Join(gopath, subFolder, *pkg))
			fmt.Printf("Check directory %s for code generation\n", pkgDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error in accessing file system. error: %s\n", err)
				os.Exit(1)
			}

			if dir, err := os.Stat(pkgDir); err == nil && dir.IsDir() {
				fileInfos, err := ioutil.ReadDir(pkgDir)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error in accessing file system. error: %s\n", err)
					os.Exit(1)
				}

				for _, fileInfo := range fileInfos {
					if strings.HasSuffix(fileInfo.Name(), ".go") &&
						!strings.HasSuffix(fileInfo.Name(), "_test.go") {

						fmt.Printf("Scan %s...\n", filepath.Join(pkgDir, fileInfo.Name()))

						fset := token.NewFileSet()
						file, err := parser.ParseFile(
							fset,
							filepath.Join(pkgDir, fileInfo.Name()),
							nil,
							parser.ParseComments)

						if err != nil {
							fmt.Fprintf(
								os.Stderr, "Error in parsing %s, error: %s\n",
								filepath.Join(pkgDir, fileInfo.Name()), err,
							)
							continue
						}

						var mockFile *ast.File
						if len(methodsToMock) > 0 {
							fset2 := token.NewFileSet()
							mockFile, err = parser.ParseFile(
								fset2,
								filepath.Join(pkgDir, *mockClzFileName),
								nil,
								parser.ParseComments)

							if err != nil {
								fmt.Fprintf(
									os.Stderr, "Error in parsing %s, error: %s\n",
									filepath.Join(pkgDir, *mockClzFileName), err,
								)
							}
						}

						g := &generator{
							clzName:        *clzName,
							mockClzName:    *mockClzName,
							mockPkgName:    *pkg,
							mockName:       *mockName,
							mothodsToClone: methodsToClone,
							methodsToMock:  methodsToMock,
						}

						outputFileName := fmt.Sprintf("mockc_%s_test.go", *mockName)
						output, err := os.OpenFile(
							filepath.Join(pkgDir, outputFileName),
							os.O_CREATE|os.O_RDWR,
							0644)
						if err != nil {
							fmt.Fprintf(
								os.Stderr, "Error in creating %s, error: %s\n",
								outputFileName, err,
							)

							continue
						}

						g.generate(file, mockFile, output)

						offset, err := output.Seek(0, io.SeekCurrent)
						if err != nil {
							fmt.Printf("Error in file operation on %s, error: %s", outputFileName, err)
						} else {
							fi, _ := output.Stat()
							if offset > 0 && offset < fi.Size() {
								output.Truncate(offset)
							}
						}
						output.Close()

						formatGoFile(filepath.Join(pkgDir, outputFileName))

						fmt.Printf("Done scan with %s\n\n", filepath.Join(pkgDir, fileInfo.Name()))
					}
				}
			}
		}
	}
}
