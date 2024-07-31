package gosyntax

import (
	"go/ast"
	"go/types"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
)

type CalleeVisitor struct {
	imports map[string]string

	// if caller has a receiver, set receiver and clzMethods accordingly,
	// otherwise, set to zero values
	receiver   string // receiver variable name
	clzMethods map[string]*ReceiverSpec

	name string // function name

	// peer callees (excluding the method myself and method-like callees invoked from function pointers)
	thisClassCallees []string

	// functions within the same package
	thisPkgCallees []string

	// functions from imported packages
	// map package name -> functions
	otherPkgcallees map[string][]string
}

func NewCalleeVisitor(
	imports map[string]string,
	clzMethods map[string]*ReceiverSpec,
	receiver, name string,
) *CalleeVisitor {
	return &CalleeVisitor{
		imports:    imports,
		clzMethods: clzMethods,
		receiver:   receiver,
		name:       name,
	}
}

func (v *CalleeVisitor) ReceiverName() string {
	return v.receiver
}

func (v *CalleeVisitor) MethodFuncName() string {
	return v.name
}

func (v *CalleeVisitor) AppendPeerCallee(calleeName string) {
	v.thisClassCallees = append(v.thisClassCallees, calleeName)
}

func (v *CalleeVisitor) GetPeerCallees() []string {
	return v.thisClassCallees
}

func (v *CalleeVisitor) AppendThisPackageCallee(calleeName string) {
	if !slices.Contains(v.thisPkgCallees, calleeName) {
		v.thisPkgCallees = append(v.thisPkgCallees, calleeName)
	}
}

func (v *CalleeVisitor) GetThisPackageCallees() []string {
	return v.thisPkgCallees
}

func (v *CalleeVisitor) AppendOtherPackageCallee(pkgName, calleeName string) {
	if _, ok := v.imports[pkgName]; ok {
		if v.otherPkgcallees == nil {
			v.otherPkgcallees = make(map[string][]string)
		}

		if !slices.Contains(v.otherPkgcallees[pkgName], calleeName) {
			v.otherPkgcallees[pkgName] = append(v.otherPkgcallees[pkgName], calleeName)
		}
	}
}

func (v *CalleeVisitor) GetOtherPackageCallees() map[string][]string {
	return v.otherPkgcallees
}

func (v *CalleeVisitor) isSelf(x, sel string) bool {
	return x == v.receiver && sel == v.name
}

func (v *CalleeVisitor) isPeer(x, sel string) bool {
	return x == v.receiver && sel != v.name
}

func (v *CalleeVisitor) isMethod(x, sel string) bool {
	if len(v.clzMethods) > 0 {
		if spec, ok := v.clzMethods[sel]; ok && spec.Name == x {
			return true
		}
	}
	return false
}

func (v *CalleeVisitor) Visit(node ast.Node) ast.Visitor {
	if callExpr, ok := node.(*ast.CallExpr); ok {
		if fun, ok := callExpr.Fun.(*ast.Ident); ok {
			if len(v.receiver) > 0 {
				v.AppendThisPackageCallee(fun.Name)
			} else {
				if !v.isSelf("", fun.Name) {
					v.AppendThisPackageCallee(fun.Name)
				}
			}
		} else if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if _, ok := sel.X.(*ast.Ident); ok {
				x := sel.X.(*ast.Ident).Name
				n := sel.Sel.Name

				if len(v.receiver) > 0 {
					if !v.isSelf(x, n) {
						if v.isPeer(x, n) {
							if v.isMethod(x, n) {
								v.AppendPeerCallee(n)
							}
						} else {
							v.AppendOtherPackageCallee(x, n)
						}
					}
				} else {
					v.AppendOtherPackageCallee(x, n)
				}
			}
		}
	}
	return v
}

func (v *CalleeVisitor) SanitizeCallees(
	imports map[string]string,
) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedSyntax}

	if len(v.thisPkgCallees) > 0 {
		filteredCallees := []string{}

		pkgs, _ := packages.Load(cfg, ".")
		for _, callee := range v.thisPkgCallees {
			if findFuncSignature(pkgs[0], callee) != nil {
				filteredCallees = append(filteredCallees, callee)
			}
		}
		v.thisPkgCallees = filteredCallees
	}

	if len(v.otherPkgcallees) > 0 {
		for pkgName, callees := range v.otherPkgcallees {
			if _, ok := imports[pkgName]; ok {
				pkgs, err := packages.Load(cfg, imports[pkgName])
				if err == nil {
					filteredCallees := []string{}

					for _, callee := range callees {
						if findFuncSignature(pkgs[0], callee) != nil {
							filteredCallees = append(filteredCallees, callee)
						}
					}
					if len(filteredCallees) > 0 {
						v.otherPkgcallees[pkgName] = filteredCallees
					} else {
						delete(v.otherPkgcallees, pkgName)
					}
				} else {
					delete(v.otherPkgcallees, pkgName)
				}
			} else {
				delete(v.otherPkgcallees, pkgName)
			}
		}
	}
}

func findFuncSignature(p *packages.Package, fnName string) *types.Signature {
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
