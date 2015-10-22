package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strings"
)

var builtinTypes = [...]string{"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64",
	"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr"}

func check(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

var filename = flag.String("file", "/Users/AnimotoOverstreet/go/src/github.com/firegoblin/goActorBinaryTree", "file to parse on")
var includeTestFiles = flag.Bool("test", true, "whether or not to include test files in the graph")

var structList = make([]*Struct, 0)
var funcList = make([]*Function, 0)
var interfaceList = make([]*Interface, 0)

func processTypeDecl(obj *ast.Object, typ *Type) {
	decl, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		panic("unexpected type in processTypeDecl")
	}

	switch decl.Type.(type) {
	case *ast.StructType:
		structList = append(structList, makeStruct(decl, typ.base))
	case *ast.InterfaceType:
		interfaceList = append(interfaceList, makeInterface(decl, typ.base))
	default:
		panic("unexpected type of s.Type")
	}
}

func main() {
	//initialize master map with builtin types
	fset := token.NewFileSet()

	var pkgs map[string]*ast.Package
	var err error

	if *includeTestFiles {
		pkgs, err = parser.ParseDir(fset, *filename, nil, 0)
	} else {
		pkgs, err = parser.ParseDir(fset, *filename, func(f os.FileInfo) bool {
			return !strings.Contains(f.Name(), "_test.go")
		}, 0)
	}
	if err != nil {
		panic("could not read original directory")
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			//add all types to master list before processing delcarations
			//minimizes creation of unknown types
			for key := range pkg.Scope.Objects {
				typeMap.lookupOrAdd(key)
			}

			//processes all structs, interfaces, and embedded types
			for key, scope := range pkg.Scope.Objects {
				typ := typeMap.lookupOrAdd(key)
				processTypeDecl(scope, typ)
			}

			//processes all the function declarations
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					funcMap.lookupOrAddFromExpr(d.Name.Name, d.Type)
				}
			}
		}
	}

	for _, i := range interfaceList {
		implementingStructs := i.implementedBy(structList)

		fmt.Println("Interface", i, "is implemented by the following types:")
		for _, s := range implementingStructs {
			fmt.Println("   ", s)
		}
	}
}
