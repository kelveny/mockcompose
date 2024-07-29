package gotype

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

// for debugging purpose only
func TestFieldInfoFromSignature(t *testing.T) {
	assert := require.New(t)

	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}

	// return packages.Package
	pkgs, err := packages.Load(cfg, "github.com/kelveny/mockcompose/test/libfn")
	if err != nil {
		t.FailNow()
	}

	assert.True(len(pkgs) == 1)

	fn := FindFuncSignature(pkgs[0], "GetSecrets")
	assert.True(fn != nil)

	paramInfos := GetFuncParamInfosFromSignature(fn, "")
	assert.True(paramInfos != nil)

	returnInfos := GetFuncReturnInfosFromSignature(fn, "")
	assert.True(returnInfos != nil)

	intf := FindInterfaceMethodSignature(pkgs[0], "SecretsInterface", "GetSecrets")

	assert.True(intf != nil)
}

func TestLoadDotPackage(t *testing.T) {
	assert := require.New(t)

	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}

	pkgs, err := packages.Load(cfg, ".")
	assert.NoError(err)
	assert.Equal(1, len(pkgs))

	fn := FindFuncSignature(pkgs[0], "RenderTypeDeclString")
	assert.NotNil(fn)

	paramInfos := GetFuncParamInfosFromSignature(fn, "")
	assert.True(paramInfos != nil)

	returnInfos := GetFuncReturnInfosFromSignature(fn, "")
	assert.True(returnInfos != nil)
}

func TestVariadicFieldInfoFromSignature(t *testing.T) {
	assert := require.New(t)

	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}

	// return packages.Package
	pkgs, err := packages.Load(cfg, "fmt")
	if err != nil {
		t.FailNow()
	}

	assert.True(len(pkgs) == 1)

	fn := FindFuncSignature(pkgs[0], "Sprintf")
	assert.True(fn != nil)

	paramInfos := GetFuncParamInfosFromSignature(fn, "")
	assert.True(paramInfos != nil)

	returnInfos := GetFuncReturnInfosFromSignature(fn, "")
	assert.True(returnInfos != nil)
}
