package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"reflect"
)

func main() {
	path := flag.String("path", "../example", "The source path of a Go package")
	flag.Parse()

	fset, astf := parsePackage(*path)

	info := inferTypes(*path, fset, astf)

	for _, f := range astf {
		ast.Walk(&PrintASTVisitor{info}, f)
	}
}

func parsePackage(path string) (*token.FileSet, []*ast.File) {
	fset := token.NewFileSet()

	pkgs, e := parser.ParseDir(fset, path, nil, 0)
	if e != nil {
		log.Fatal(e)
	}

	astf := make([]*ast.File, 0)
	for _, pkg := range pkgs {
		fmt.Printf("Parsed package %v\n", pkg.Name)
		for fn, f := range pkg.Files {
			fmt.Printf("\t%v\n", fn)
			astf = append(astf, f)
		}
	}

	return fset, astf
}

func inferTypes(path string, fset *token.FileSet, astf []*ast.File) *types.Info {
	config := &types.Config{
		Error: func(e error) {
			fmt.Println(e)
		},
		Importer: importer.Default(),
	}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	pkg, e := config.Check(path, fset, astf, info)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Printf("Type inference done: %v\n", pkg.String())

	return info
}

type PrintASTVisitor struct {
	info *types.Info
}

func (v *PrintASTVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		fmt.Printf("%s", reflect.TypeOf(node).String())
		switch node.(type) {
		case ast.Expr:
			t := v.info.TypeOf(node.(ast.Expr))
			if t != nil {
				fmt.Printf(" : %s", t.String())
			}
		}
		fmt.Println()
	}
	return v
}
