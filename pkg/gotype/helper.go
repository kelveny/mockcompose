package gotype

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/kelveny/mockcompose/pkg/gosyntax"
	"golang.org/x/tools/go/packages"
)

//
//  Quick reference to GO type objects that are used in this package
//
//	Package {
//		Name() string
//		Path() string
//		Imports() []*Package
//		Scope() *Scope
//			elems (string -> types.Object)
//  }
//
//  types.Object (interface) {
//		Pkg() *Package  // package to which this object belongs; nil for labels and objects in the Universe scope
//		Name() string   // package local object name
//		Type() Type     // object type
//		Exported() bool // reports whether the name starts with a capital letter
//		Id() string     // object name if exported, qualified name if not exported (see func Id)
//		String() string
//  }
//		PkgName
//		TypeName
//		Const
//		Var
//		Func
//		Label
//		Builtin
//		Nil
//
//	types.Type (interface) {
//		Underlying() Type
//		String() string
//  }
//
//		Basic
//		Slice
//		Array
//		Struct
//		Interface
//		Map
//		Chan
//		Pointer
//		Signature (A Signature represents a (non-builtin) function or method type)
//		Tuple
//		Named
//
func FindFuncSignature(p *packages.Package, fnName string) *types.Signature {
	if p != nil && p.Types != nil {
		ret := p.Types.Scope().Lookup(fnName)
		if ret != nil {
			if fn, ok := ret.(*types.Func); ok {
				return fn.Type().(*types.Signature)
			}
		}
	}
	return nil
}

func FindInterfaceMethodSignature(
	p *packages.Package,
	intfName, methodName string,
) *types.Signature {
	if p != nil && p.Types != nil {
		ret := p.Types.Scope().Lookup(intfName)
		if ret != nil {
			if named, ok := ret.Type().(*types.Named); ok {
				underlying := named.Underlying()
				if intf, ok := underlying.(*types.Interface); ok {
					for i := 0; i < intf.NumMethods(); i++ {
						if intf.Method(i).Name() == methodName {
							return intf.Method(i).Type().(*types.Signature)
						}
					}
				}
			}
		}
	}

	return nil
}

func GetFuncParamInfosFromSignature(fn *types.Signature, mockPkg string) []*gosyntax.FieldDeclInfo {
	paramInfos := []*gosyntax.FieldDeclInfo{}

	tuples := fn.Params()
	if tuples != nil && tuples.Len() > 0 {
		for i := 0; i < tuples.Len(); i++ {
			v := tuples.At(i)

			variadic := fn.Variadic() && i == tuples.Len()-1
			paramInfos = append(paramInfos, &gosyntax.FieldDeclInfo{
				Name:     v.Name(),
				Typ:      RenderTypeDeclString(v.Type(), variadic, mockPkg),
				Variadic: variadic,
			})
		}
	}

	return paramInfos
}

func GetFuncReturnInfosFromSignature(fn *types.Signature, mockPkg string) []*gosyntax.FieldDeclInfo {
	paramInfos := []*gosyntax.FieldDeclInfo{}

	tuples := fn.Results()
	if tuples != nil && tuples.Len() > 0 {
		for i := 0; i < tuples.Len(); i++ {
			v := tuples.At(i)

			paramInfos = append(paramInfos, &gosyntax.FieldDeclInfo{
				Name:     v.Name(),
				Typ:      RenderTypeDeclString(v.Type(), false, mockPkg),
				Variadic: false,
			})
		}
	}

	return paramInfos
}

func RenderTypeDeclString(t types.Type, variadic bool, mockPkg string) string {
	switch tt := t.(type) {
	case *types.Basic:
		return tt.Name()
	case *types.Slice:
		if variadic {
			return "..." + RenderTypeDeclString(tt.Elem(), false, mockPkg)
		}
		return "[]" + RenderTypeDeclString(tt.Elem(), false, mockPkg)
	case *types.Array:
		return fmt.Sprintf("[%d]%s", tt.Len(), RenderTypeDeclString(tt.Elem(), false, mockPkg))
	case *types.Struct:
		var fields []string

		for i := 0; i < tt.NumFields(); i++ {
			field := tt.Field(i)

			if field.Anonymous() {
				fields = append(fields, RenderTypeDeclString(field.Type(), false, mockPkg))
			} else {
				fields = append(fields,
					fmt.Sprintf("%s %s", field.Name(), RenderTypeDeclString(field.Type(), false, mockPkg)),
				)
			}
		}
		return fmt.Sprintf("struct{%s}", strings.Join(fields, ";"))
	case *types.Interface:
		if tt.NumMethods() != 0 {
			panic("Empty interface")
		}
		return "interface{}"

	case *types.Map:
		key := RenderTypeDeclString(tt.Key(), false, mockPkg)
		val := RenderTypeDeclString(tt.Elem(), false, mockPkg)

		return fmt.Sprintf("map[%s]%s", key, val)
	case *types.Chan:
		switch tt.Dir() {
		case types.SendRecv:
			return "chan " + RenderTypeDeclString(tt.Elem(), false, mockPkg)
		case types.RecvOnly:
			return "<-chan " + RenderTypeDeclString(tt.Elem(), false, mockPkg)
		default:
			return "chan<- " + RenderTypeDeclString(tt.Elem(), false, mockPkg)
		}

	case *types.Pointer:
		return "*" + RenderTypeDeclString(tt.Elem(), false, mockPkg)

	case *types.Signature:
		switch tt.Results().Len() {
		case 0:
			return fmt.Sprintf(
				"func(%s)",
				RenderTypeDeclString(tt.Params(), false, mockPkg),
			)
		case 1:
			return fmt.Sprintf(
				"func(%s) %s",
				RenderTypeDeclString(tt.Params(), false, mockPkg),
				RenderTypeDeclString(tt.Results().At(0).Type(), false, mockPkg),
			)
		default:
			return fmt.Sprintf(
				"func(%s)(%s)",
				RenderTypeDeclString(tt.Params(), false, mockPkg),
				RenderTypeDeclString(tt.Results(), false, mockPkg),
			)
		}
	case *types.Tuple:
		var parts []string

		for i := 0; i < tt.Len(); i++ {
			part := tt.At(i)
			parts = append(parts, RenderTypeDeclString(part.Type(), false, mockPkg))
		}

		return strings.Join(parts, ", ")

	case *types.Named:
		o := tt.Obj()
		if o.Pkg() == nil || o.Pkg().Name() == "main" || o.Pkg().Name() == mockPkg {
			return o.Name()
		}
		return o.Pkg().Name() + "." + o.Name()

	default:
		panic(fmt.Sprintf("Unsupported type: %#v (%T)", t, tt))
	}
}
