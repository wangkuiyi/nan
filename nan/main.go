package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	path := flag.String("path", "../example", "The source path of a Go package")
	check := flag.Bool("check", true, "Add or remove NaN checks")
	flag.Parse()

	fset, astf := parsePackage(*path)

	info := inferTypes(*path, fset, astf)

	for _, f := range astf {
		ast.Walk(&Printer{info, *check}, f)
		printFile(fset, f)
	}
}

func parsePackage(path string) (*token.FileSet, []*ast.File) {
	fset := token.NewFileSet()

	pkgs, e := parser.ParseDir(fset, path, nil, parser.ParseComments)
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

type Printer struct {
	info  *types.Info
	check bool
}

func (v *Printer) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		fmt.Printf("%s", reflect.TypeOf(node).String())
		switch node.(type) {
		case *ast.CallExpr:
			fmt.Printf(" : %#v", spew.Sdump(node.(*ast.CallExpr)))
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

func printFile(fset *token.FileSet, astf *ast.File) {
	f, e := os.Create(fset.Position(astf.Package).Filename)
	if e != nil {
		log.Fatal(e)
	}

	if e := printer.Fprint(f, fset, astf); e != nil {
		log.Fatal(e)
	}
}
