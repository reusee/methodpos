package main

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/reusee/e/v2"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
)

var (
	ce, he = e.New(e.Default.WithStack())
	pt     = fmt.Printf
)

func main() {
	pkgs, err := packages.Load(
		&packages.Config{
			Mode: 0 |
				packages.NeedImports |
				packages.NeedDeps |
				packages.NeedTypes |
				packages.NeedSyntax |
				packages.NeedTypesInfo,
		},
		"bytes",
	)
	ce(err)
	if packages.PrintErrors(pkgs) > 0 {
		return
	}
	pkg := pkgs[0]

	globalScope := pkg.Types.Scope()
	namedType := globalScope.Lookup("Buffer").(*types.TypeName).Type().(*types.Named)
	method := func() *types.Func {
		for i := 0; i < namedType.NumMethods(); i++ {
			method := namedType.Method(i)
			if method.Name() == "Read" {
				return method
			}
		}
		panic("no such method")
	}()
	beginPos := pkg.Fset.Position(method.Pos())
	pt("begin at %v\n", beginPos)
	endPos := pkg.Fset.Position(method.Scope().End())
	pt("end at %v\n", endPos)

	comments := func() *ast.CommentGroup {
		for _, file := range pkg.Syntax {
			path, _ := astutil.PathEnclosingInterval(file, method.Pos(), method.Pos())
			for _, node := range path {
				decl, ok := node.(*ast.FuncDecl)
				if ok {
					return decl.Doc
				}
			}
		}
		return nil
	}()
	pt("comments begin at %v\n", pkg.Fset.Position(comments.Pos()))

}
