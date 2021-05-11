package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/kelveny/mockcompose/pkg/gofile"
	"github.com/kelveny/mockcompose/pkg/logger"
)

var SemVer = "v0.0.0-devel"

func GetSemverInfo() string {
	return SemVer
}

func usage() {
	logger.Log(logger.PROMPT, `Usage: %s [-help] [options]

mockcompose generates mocking implementation for Go classes, interfaces and functions.
`, os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func Execute() {
	var methodsToClone stringSlice
	var methodsToMock stringSlice

	vb := flag.Bool("v", false, "if set, print verbose logging messages")
	testOnly := flag.Bool("testonly", true, "if set, append _test to generated file name")
	prtVersion := flag.Bool("version", false, "if set, print version information")
	mockName := flag.String("n", "", "name of the generated class")
	mockPkg := flag.String("pkg", "", "name of the package that the generated class resides")
	clzName := flag.String("c", "", "name of the source class to generate against")
	srcPkg := flag.String("p", "", "path of the source package in which to search interfaces and functions")
	intfName := flag.String("i", "", "name of the source interface to generate against")
	flag.Var(&methodsToClone, "real", "name of the method function to be cloned from source class")
	flag.Var(&methodsToMock, "mock", "name of the function to be mocked")

	flag.Parse()

	if *prtVersion {
		fmt.Println(GetSemverInfo())
		os.Exit(0)
	}

	if *vb {
		logger.LogLevel = int(logger.VERBOSE)

		logger.Log(logger.VERBOSE, "Set logging to verbose mode\n")
	}

	if *mockPkg == "" {
		*mockPkg = gofile.DerivePackage()

		logger.Log(logger.VERBOSE, "Derive package name as: %s\n", *mockPkg)
	}
	fmt.Println()

	if *mockName == "" {
		usage()
		os.Exit(1)
	}

	options := &commandOptions{
		mockName,
		mockPkg,
		clzName,
		intfName,
		srcPkg,
		testOnly,
		methodsToClone,
		methodsToMock,
	}

	var g parsedFileGenerator
	if *clzName != "" {
		if len(methodsToClone) == 0 {
			logger.Log(logger.ERROR, "Please specify at least one real method name with -real option\n")
			os.Exit(1)
		}

		g = &classMethodGenerator{
			clzName:        *options.clzName,
			mockPkgName:    *options.mockPkg,
			mockName:       *options.mockName,
			mothodsToClone: options.methodsToClone,
			methodsToMock:  options.methodsToMock,
		}
	} else if *intfName != "" {
		g = &interfaceMockGenerator{
			mockPkgName: *options.mockPkg,
			mockName:    *options.mockName,
			intfName:    *options.intfName,
		}

		if *options.srcPkg != "" {
			scanPackageToGenerate(g.(loadedPackageGenerator), options)
			return
		}

	} else {
		if len(methodsToMock) == 0 {
			logger.Log(logger.ERROR, "Please specify at least one mock function name with -mock option\n")
			os.Exit(1)
		}

		g = &functionMockGenerator{
			mockPkgName:   *options.mockPkg,
			mockName:      *options.mockName,
			methodsToMock: *&options.methodsToMock,
		}

		if *options.srcPkg != "" {
			scanPackageToGenerate(g.(loadedPackageGenerator), options)
			return
		}
	}

	scanCWDToGenerate(g, options)
}
