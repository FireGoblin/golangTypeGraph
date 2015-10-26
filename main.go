package main

import (
	"flag"
	"fmt"
	"github.com/firegoblin/gographviz"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// var builtinTypes = [...]string{"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
// 	"int", "int8", "int16", "int32", "int64",
// 	"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr"}

var filename = flag.String("file", "github.com/firegoblin/golangTypeGraph", "file to parse on, relative to $GOPATH/src")
var includeTestFiles = flag.Bool("test", true, "whether or not to include test files in the graph")

func processTypeDecl(obj *ast.Object, typ *Type, structList *[]*Struct, interfaceList *[]*Interface) {
	fmt.Println(obj.Decl)
	decl, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		panic("unexpected type in processTypeDecl")
	}

	node := typ.base.node

	switch decl.Type.(type) {
	case *ast.StructType:
		if node == nil {
			*structList = append(*structList, makeStruct(decl, typ.base))
		} else {
			*structList = append(*structList, node.(*Unknown).remakeStruct(decl))
		}
	case *ast.InterfaceType:
		if node == nil {
			*interfaceList = append(*interfaceList, makeInterface(decl, typ.base))
		} else {
			switch n := node.(type) {
			case *Interface:
				*interfaceList = append(*interfaceList, n.remakeInterface(decl))
			case *Unknown:
				*interfaceList = append(*interfaceList, n.remakeInterface(decl))
			}
		}
	//case *ast.Ident:
	default:
		if node == nil {
			*structList = append(*structList, makeStruct(decl, typ.base))
		} else {
			*structList = append(*structList, node.(*Unknown).remakeStruct(decl))
		}
		//default:
		//	panic("unexpected type of s.Type")
	}
}

func main() {
	flag.Parse()

	//initialize master map with builtin types
	fset := token.NewFileSet()

	var pkgs map[string]*ast.Package
	var err error

	gopath := os.Getenv("GOPATH") + "/src/"

	if *includeTestFiles {
		pkgs, err = parser.ParseDir(fset, gopath+*filename, nil, 0)
	} else {
		pkgs, err = parser.ParseDir(fset, gopath+*filename, func(f os.FileInfo) bool {
			return !strings.Contains(f.Name(), "_test.go")
		}, 0)
	}
	if err != nil {
		panic(err)
	}

	structList := make([]*Struct, 0)
	interfaceList := make([]*Interface, 0)
	funcList := make([]*Function, 0)

	//ast.Print(fset, pkgs)

	//TODO: fix for multiple packages
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			//add all types to master list before processing delcarations
			//minimizes creation of unknown types
			//fmt.Println("file.Scope.Objects:")
			//fmt.Println(file.Scope.Objects)
			for key := range file.Scope.Objects {
				typeMap.lookupOrAdd(key)
			}

			//fmt.Println("typeMap:")
			//fmt.Println(typeMap)

			//processes all structs, interfaces, and embedded types
			for key, scope := range file.Scope.Objects {
				//non-receiver functions are found in scope
				if scope.Kind == ast.Typ {
					typ := typeMap.lookupOrAdd(key)
					processTypeDecl(scope, typ, &structList, &interfaceList)
				}
			}

			//processes all the function declarations
			for _, decl := range file.Decls {
				switch d := decl.(type) {
				case *ast.FuncDecl:
					//fmt.Println("funcDecl")
					f := funcMap.lookupOrAddFromExpr(d.Name.Name, d.Type)
					funcList = append(funcList, f)
					if d.Recv != nil {
						//fmt.Println("non-nil recv")
						recv := typeMap.lookupOrAddFromExpr(d.Recv.List[0].Type).base.node
						if recv != nil {
							f.addReceiver(recv.(*Struct))
							//fmt.Println(recv.(*Struct))
						}
					}
				}
			}

			//fmt.Println(interfaceList[0].requiredFunctions[0] == structList[1].receiverFunctions[0])

			fmt.Println("funcMap:")
			fmt.Println(funcMap)
		}
	}

	fmt.Println("interface:", len(interfaceList))
	fmt.Println("struct:", len(structList))
	fmt.Println("func:", len(funcList))

	for _, i := range interfaceList {
		i.setImplementedBy(structList)

		//fmt.Println("Interface", i, "is implemented by the following types:")
		//for _, s := range implementingStructs {
		//	fmt.Println("   ", s)
		//}
	}

	g := gographviz.NewGraph()
	g.SetName((*filename)[strings.LastIndex(*filename, "/")+1:])
	g.SetDir(true)
	for _, s := range structList {
		fmt.Println("Add Struct Node:", s.Name())
		g.AddGraphableNode("G", s)
	}
	for _, i := range interfaceList {
		fmt.Println("Add Interface Node:", i.Name())
		g.AddGraphableNode("G", i)
	}
	s := g.String()
	fmt.Println(s)
}
