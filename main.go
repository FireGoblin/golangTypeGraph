package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"go/parser"
	"go/ast"
	"go/token"
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

var structList := make([]*Struct, 10) 
var interfaceList := make([]*Interface, 10)
var funcList := make([]*Function, 10)

func processDecls(decl ast.Decl) {
	select d := decl.(type) {
	case *ast.FuncDecl :
		funcList  = append(funcList, makeFunction(d))
	case *ast.GenDecl :
		for _, spec := range file.Specs {
			select s := spec.(type) {
			case *ast.TypeSpec :
				select s.Type.(type) {
				case *ast.StructType:
					makeStruct(s)
				case *ast.InterfaceType:
					makeInterface(s)
				default:
					panic("unexpected type of s.Type")
				}
			default:
				panic("unexpected type of spec")
			}
		}
	default:
		panic("unexpected type of decl")
	}
}

func main() {
	//initialize master map with builtin types
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, *filename, nil, 0)
	if err != nil {
		panic("could not read original directory")
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				processDecls(decl)
			}
		}
	}
}
