package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

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

func loadConfig() *Config {
	pkgDir, err := filepath.Abs("")
	logger.Log(logger.VERBOSE, "Check directory %s for YAML configuration\n", pkgDir)
	if err != nil {
		logger.Log(logger.ERROR, "Error in accessing file system. error: %s\n", err)
		os.Exit(1)
	}

	if cfg := loadYamlConfig(filepath.Join(pkgDir, ".mockcompose.yaml")); cfg != nil {
		return cfg
	}

	if cfg := loadYamlConfig(filepath.Join(pkgDir, ".mockcompose.yml")); cfg != nil {
		return cfg
	}

	return nil
}

func loadYamlConfig(yamlFile string) *Config {
	yamlConfig, err := ioutil.ReadFile(yamlFile)
	if err == nil {
		cfg := Config{}

		err := yaml.Unmarshal(yamlConfig, &cfg)
		if err != nil {
			logger.Log(logger.ERROR, "Failed to load YAML config: %s\n", err)
		}
		return &cfg
	}

	return nil
}

func executeOptions(options *CommandOptions) {
	var g parsedFileGenerator

	if options.ClzName != "" {
		if len(options.MethodsToClone) == 0 {
			logger.Log(logger.ERROR, "Please specify at least one real method name with -real option\n")
			os.Exit(1)
		}

		g = &classMethodGenerator{
			clzName:        options.ClzName,
			mockPkgName:    options.MockPkg,
			mockName:       options.MockName,
			methodsToClone: options.MethodsToClone,
			methodsToMock:  options.MethodsToMock,
		}
	} else if options.IntfName != "" {
		g = &interfaceMockGenerator{
			mockPkgName: options.MockPkg,
			mockName:    options.MockName,
			intfName:    options.IntfName,
		}

		if options.SrcPkg != "" {
			scanPackageToGenerate(g.(loadedPackageGenerator), options)
			return
		}
	} else {
		if len(options.MethodsToMock) == 0 && len(options.MethodsToClone) == 0 {
			logger.Log(logger.ERROR, "no function to mock or clone\n")
			os.Exit(1)
		}

		if len(options.MethodsToMock) > 0 && len(options.MethodsToClone) > 0 {
			logger.Log(logger.ERROR, "option -real and option -mock are exclusive in function clone generation\n")
			os.Exit(1)
		}

		if len(options.MethodsToClone) > 0 {
			if options.SrcPkg != "" {
				logger.Log(logger.PROMPT,
					"No source package support in function clone generation, ignore source package %s\n",
					options.SrcPkg)
			}
			g = &functionCloneGenerator{
				mockPkgName:    options.MockPkg,
				mockName:       options.MockName,
				methodsToClone: options.MethodsToClone,
			}
		} else {
			g = &functionMockGenerator{
				mockPkgName:   options.MockPkg,
				mockName:      options.MockName,
				methodsToMock: options.MethodsToMock,
			}

			if options.SrcPkg != "" {
				scanPackageToGenerate(g.(loadedPackageGenerator), options)
				return
			}
		}
	}

	scanCWDToGenerate(g, options)
}

func Execute() {
	var methodsToClone stringSlice
	var methodsToMock stringSlice

	vb := flag.Bool("v", false, "if set, print verbose logging messages")
	testOnly := flag.Bool("testonly", true, "if set, append _test to generated file name")
	prtVersion := flag.Bool("version", false, "if set, print version information")
	help := flag.Bool("help", false, "if set, print usage information")
	mockName := flag.String("n", "", "name of the generated class")
	mockPkg := flag.String("pkg", "", "name of the package that the generated class resides")
	clzName := flag.String("c", "", "name of the source class to generate against")
	srcPkg := flag.String("p", "", "path of the source package in which to search interfaces and functions")
	intfName := flag.String("i", "", "name of the source interface to generate against")
	flag.Var(&methodsToClone, "real", "name of the method function to be cloned from source class or source function")
	flag.Var(&methodsToMock, "mock", "name of the function to be mocked")

	flag.Parse()

	if *prtVersion {
		fmt.Println(GetSemverInfo())
		os.Exit(0)
	}

	if *help {
		usage()
		os.Exit(0)
	}

	if *vb {
		logger.LogLevel = int(logger.VERBOSE)

		logger.Log(logger.VERBOSE, "Set logging to verbose mode\n")
	}

	if cfg := loadConfig(); cfg != nil {
		logger.Log(logger.VERBOSE, "Found mockcompose YAML configuration, ignore command line options\n")

		derivedPkg := gofile.DerivePackage(false)
		for _, options := range cfg.Mockcompose {
			if options.MockPkg == "" {
				options.MockPkg = derivedPkg
			}

			executeOptions(&options)
		}

		return
	}

	if *mockPkg == "" {
		*mockPkg = gofile.DerivePackage(false)

		logger.Log(logger.VERBOSE, "Derive package name as: %s\n", *mockPkg)
	}
	fmt.Println()

	if *mockName == "" {
		usage()
		os.Exit(1)
	}

	options := &CommandOptions{
		MockName:       *mockName,
		MockPkg:        *mockPkg,
		ClzName:        *clzName,
		IntfName:       *intfName,
		SrcPkg:         *srcPkg,
		TestOnly:       *testOnly,
		MethodsToClone: methodsToClone,
		MethodsToMock:  methodsToMock,
	}

	executeOptions(options)
}
