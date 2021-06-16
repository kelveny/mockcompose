package gofile

import (
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kelveny/mockcompose/pkg/logger"
)

func DerivePackage(anchor bool) string {
	path, err := filepath.Abs("")
	if err != nil {
		logger.Log(logger.ERROR, "Error in accessing file system. error: %s\n", err)
		os.Exit(1)
	}

	if path, ok := StripGopath(path); ok {
		if anchor {
			return path
		}
		return filepath.Base(path)
	}

	return filepath.Base(path)
}

// StripGopath strips GOPATH roots and normalizes it to Go package path format
func StripGopath(p string) (string, bool) {
	gopathConfig := GetGoPathConfig()

	for _, gopath := range strings.Split(gopathConfig, string(filepath.ListSeparator)) {
		logger.Log(logger.VERBOSE, "check gopath: %s\n", gopath)

		if strings.HasPrefix(p, gopath+string(filepath.Separator)) {
			p = strings.Replace(p, gopath+string(filepath.Separator), "", 1)

			// normalize to Go package path format
			if filepath.Separator != '/' {
				p = strings.ReplaceAll(p, string(filepath.Separator), "/")
			}

			p = strings.Replace(p, "src/", "", 1)
			p = strings.Replace(p, "pkg/", "", 1)

			return p, true
		}
	}

	return p, false
}

func GetGoPathConfig() string {
	gopathConfig := os.Getenv("GOPATH")

	pkgRoot, _ := filepath.Abs("")

	var detected string
	if pkgPos := strings.LastIndex(pkgRoot, string(filepath.Separator)+"pkg"); pkgPos >= 0 {
		detected = pkgRoot[0:pkgPos]

		if strings.Index(gopathConfig, detected) < 0 {
			gopathConfig = gopathConfig + string(filepath.ListSeparator) + detected
		}
	}

	if pkgPos := strings.LastIndex(pkgRoot, string(filepath.Separator)+"src"); pkgPos >= 0 {
		detected2 := pkgRoot[0:pkgPos]

		if strings.Index(gopathConfig, detected2) < 0 {
			return gopathConfig + string(filepath.ListSeparator) + detected2
		}
	}

	return gopathConfig
}

func FormatGoFile(filePath string) {
	b, err := ioutil.ReadFile(filePath)

	if err != nil {
		logger.Log(logger.ERROR, "Error in reading file %s, error: %s\n", filePath, err)
		return
	}

	bb, err := format.Source(b)
	if err != nil {
		logger.Log(logger.ERROR, "Error in formatting Go source %s, error: %s\n", filePath, err)
		return
	}

	ioutil.WriteFile(filePath, bb, 0644)
}
