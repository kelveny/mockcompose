package cmd

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kelveny/mockcompose/pkg/gofile"
	"github.com/kelveny/mockcompose/pkg/logger"
	"golang.org/x/tools/go/packages"
)

func scanPackageToGenerate(
	g loadedPackageGenerator,
	options *CommandOptions,
) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}

	pkgs, err := packages.Load(cfg, options.SrcPkg)
	if err != nil {
		logger.Log(logger.ERROR, "Error in loading package %s, error: %s\n",
			options.SrcPkg, err,
		)
		return
	}

	logger.Log(logger.PROMPT, "Scan package %s...\n", options.SrcPkg)

	var outputFileName string
	if options.TestOnly {
		outputFileName = fmt.Sprintf("mockc_%s_test.go", options.MockName)
	} else {
		outputFileName = fmt.Sprintf("mockc_%s.go", options.MockName)
	}

	output, err := os.OpenFile(
		outputFileName,
		os.O_CREATE|os.O_RDWR,
		0644)
	if err != nil {
		logger.Log(logger.ERROR, "Error in creating %s, error: %s\n",
			outputFileName, err,
		)

		return
	}

	for _, pkg := range pkgs {
		if len(pkg.Syntax) == 0 {
			for _, err := range pkg.Errors {
				logger.Log(logger.WARN, "%s error: %s\n",
					pkg.ID, err.Msg,
				)
			}
		} else {
			g.generateViaLoadedPackage(output, pkg)
		}
	}

	offset, err := output.Seek(0, io.SeekCurrent)
	if err != nil {
		logger.Log(logger.ERROR, "Error in file operation on %s, error: %s\n", outputFileName, err)
	} else {
		fi, _ := output.Stat()
		if offset > 0 && offset < fi.Size() {
			output.Truncate(offset)
		}
	}
	output.Close()

	gofile.FormatGoFile(outputFileName)

	logger.Log(logger.PROMPT, "Done scan with package %s\n\n", options.SrcPkg)
}

// scan current working directory
func scanCWDToGenerate(
	g parsedFileGenerator,
	options *CommandOptions,
) {
	pkgDir, err := filepath.Abs("")
	logger.Log(logger.VERBOSE, "Check directory %s for code generation\n", pkgDir)
	if err != nil {
		logger.Log(logger.ERROR, "Error in accessing file system. error: %s\n", err)
		os.Exit(1)
	}

	if dir, err := os.Stat(pkgDir); err == nil && dir.IsDir() {
		fileInfos, err := ioutil.ReadDir(pkgDir)
		if err != nil {
			logger.Log(logger.ERROR, "Error in accessing file system. error: %s\n", err)
			os.Exit(1)
		}

		for _, fileInfo := range fileInfos {
			scanFileToGenerate(g, options, pkgDir, fileInfo)
		}
	}
}

// not in use
func scanGoPathToGenerate(
	g parsedFileGenerator,
	options *CommandOptions,
) {
	// iterate candidates from package directory
	gopathConfig := gofile.GetGoPathConfig()

	for _, gopath := range strings.Split(gopathConfig, string(filepath.ListSeparator)) {
		// support scanning of subfolder src/ and pkg/
		for _, subFolder := range []string{"src", "pkg"} {
			pkgDir, err := filepath.Abs(path.Join(gopath, subFolder, options.MockPkg))
			logger.Log(logger.VERBOSE, "Check directory %s for code generation\n", pkgDir)
			if err != nil {
				logger.Log(logger.ERROR, "Error in accessing file system. error: %s\n", err)
				os.Exit(1)
			}

			if dir, err := os.Stat(pkgDir); err == nil && dir.IsDir() {
				fileInfos, err := ioutil.ReadDir(pkgDir)
				if err != nil {
					logger.Log(logger.ERROR, "Error in accessing file system. error: %s\n", err)
					os.Exit(1)
				}

				for _, fileInfo := range fileInfos {
					scanFileToGenerate(g, options, pkgDir, fileInfo)
				}
			}
		}
	}
}

func scanFileToGenerate(
	g parsedFileGenerator,
	options *CommandOptions,
	pkgDir string,
	fileInfo os.FileInfo,
) {
	if strings.HasSuffix(fileInfo.Name(), ".go") &&
		!strings.HasSuffix(fileInfo.Name(), "_test.go") {

		logger.Log(logger.PROMPT, "Scan %s...\n", filepath.Join(pkgDir, fileInfo.Name()))

		fset := token.NewFileSet()
		file, err := parser.ParseFile(
			fset,
			filepath.Join(pkgDir, fileInfo.Name()),
			nil,
			parser.ParseComments)

		if err != nil {
			logger.Log(logger.ERROR, "Error in parsing %s, error: %s\n",
				filepath.Join(pkgDir, fileInfo.Name()), err,
			)
			return
		}

		var outputFileName string
		if options.TestOnly {
			outputFileName = fmt.Sprintf("mockc_%s_test.go", options.MockName)
		} else {
			outputFileName = fmt.Sprintf("mockc_%s.go", options.MockName)
		}

		output, err := os.OpenFile(
			filepath.Join(pkgDir, outputFileName),
			os.O_CREATE|os.O_RDWR,
			0644)
		if err != nil {
			logger.Log(logger.ERROR, "Error in creating %s, error: %s\n",
				outputFileName, err,
			)

			return
		}

		g.generate(output, file)

		offset, err := output.Seek(0, io.SeekCurrent)
		if err != nil {
			logger.Log(logger.ERROR, "Error in file operation on %s, error: %s\n", outputFileName, err)
		} else {
			fi, _ := output.Stat()
			if offset > 0 && offset < fi.Size() {
				output.Truncate(offset)
			}
		}
		output.Close()

		gofile.FormatGoFile(filepath.Join(pkgDir, outputFileName))

		logger.Log(logger.PROMPT, "Done scan with %s\n\n", filepath.Join(pkgDir, fileInfo.Name()))
	}
}
