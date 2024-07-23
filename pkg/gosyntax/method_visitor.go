package gosyntax

import (
	"go/ast"
)

type MethodCalleeVisitor struct {
	clzMethods map[string]*ReceiverSpec
	receiver   string
	name       string

	// peer callees (excluding the method myselfm but including peer methodd and callees via field function pointers)
	thisClassCallees []string

	// functions within the same package
	thisPkgCallees []string

	// functions from imported packages
	// map package name -> functions
	otherPkgcallees map[string][]string
}

func NewMethodCalleeVisitor(
	clzMethods map[string]*ReceiverSpec,
	receiver, name string,
) *MethodCalleeVisitor {
	return &MethodCalleeVisitor{
		clzMethods: clzMethods,
		receiver:   receiver,
		name:       name,
	}
}

func (v *MethodCalleeVisitor) AppendPeerCallee(calleeName string) {
	v.thisClassCallees = append(v.thisClassCallees, calleeName)
}

func (v *MethodCalleeVisitor) GetPeerCallees() []string {
	return v.thisClassCallees
}

func (v *MethodCalleeVisitor) AppendThisPackageCallee(calleeName string) {
	v.thisPkgCallees = append(v.thisPkgCallees, calleeName)
}

func (v *MethodCalleeVisitor) GetThisPackageCallees() []string {
	return v.thisPkgCallees
}

func (v *MethodCalleeVisitor) AppendOtherPackageCallee(pkgName, calleeName string) {
	if v.otherPkgcallees == nil {
		v.otherPkgcallees = make(map[string][]string)
	}

	v.otherPkgcallees[pkgName] = append(v.otherPkgcallees[pkgName], calleeName)
}

func (v *MethodCalleeVisitor) GetOtherPackageCallees() map[string][]string {
	return v.otherPkgcallees
}

func (v *MethodCalleeVisitor) isSelf(x, sel string) bool {
	return x == v.receiver && sel == v.name
}

func (v *MethodCalleeVisitor) isPeer(x, sel string) bool {
	return x == v.receiver && sel != v.name
}

func (v *MethodCalleeVisitor) isMethod(x, sel string) bool {
	if len(v.clzMethods) > 0 {
		if spec, ok := v.clzMethods[sel]; ok && spec.Name == x {
			return true
		}
	}
	return false
}

func (v *MethodCalleeVisitor) Visit(node ast.Node) ast.Visitor {
	if callExpr, ok := node.(*ast.CallExpr); ok {
		if fun, ok := callExpr.Fun.(*ast.Ident); ok {
			v.AppendThisPackageCallee(fun.Name)
		} else if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if _, ok := sel.X.(*ast.Ident); ok {
				x := sel.X.(*ast.Ident).Name
				n := sel.Sel.Name

				if !v.isSelf(x, n) {
					if v.isPeer(x, n) {
						if v.isMethod(x, n) {
							v.AppendPeerCallee(n)
						}
					} else {
						v.AppendOtherPackageCallee(x, n)
					}
				}
			}
		}
	}
	return v
}
