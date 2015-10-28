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

var filename = flag.String("file", "github.com/firegoblin/golangTypeGraph", "file to parse on, relative to $GOPATH/src")
var includeTestFiles = flag.Bool("test", false, "whether or not to include test files in the graph")
var defaultPackageName = flag.String("pkg", "main", "the package that will not have its types prefiexed with package name")
var onlyExports = flag.Bool("exports", false, "marks whether only exported nodes are shown")
var implementMax = flag.Int("imax", 9, "the maximum number of structs implementing an interface before edges are not drawn")
var envVar = flag.String("env", "GOPATH", "environment variable to use instead of GOPATH")
var maxDepth = flag.Int("depth", 1, "maximum depth of recursively searching imports")

func processTypeDecl(obj *ast.Object, typ *Type, structList *[]*structNode, interfaceList *[]*interfaceNode) {
	decl, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		panic("unexpected type in processTypeDecl")
	}

	node := typ.base.node

	switch decl.Type.(type) {
	case *ast.InterfaceType:
		if node == nil {
			*interfaceList = append(*interfaceList, newInterface(decl, typ.base))
		} else {
			switch n := node.(type) {
			case *interfaceNode:
				*interfaceList = append(*interfaceList, n.remakeInterface(decl))
			case *unknownNode:
				*interfaceList = append(*interfaceList, n.remakeInterface(decl))
			}
		}
	//case StructType or redefinied type
	default:
		if node == nil {
			*structList = append(*structList, newStruct(decl, typ.base))
		} else {
			*structList = append(*structList, node.(*unknownNode).remakeStruct(decl))
		}
	}
}

func main() {
	flag.Parse()

	//initialize master map with builtin types
	fset := token.NewFileSet()

	var pkgs map[string]*ast.Package
	var err error

	gopath := os.Getenv(*envVar) + "/src/"

	var structList []*structNode
	var interfaceList []*interfaceNode
	//var funcList []*function

	var directories []string
	var depth []int
	directories = append(directories, *filename)
	depth = append(depth, 0)

	var searchedDirectories []string

	//adding os as it does not work with what seems like duplicate types
	searchedDirectories = append(searchedDirectories, "os")

	for len(directories) > 0 {
		if *includeTestFiles {
			pkgs, err = parser.ParseDir(fset, gopath+directories[len(directories)-1], nil, 0)
		} else {
			pkgs, err = parser.ParseDir(fset, gopath+directories[len(directories)-1], func(f os.FileInfo) bool {
				return !strings.Contains(f.Name(), "_test.go")
			}, 0)
		}

		dep := depth[len(depth)-1]
		searchedDirectories = append(searchedDirectories, directories[len(directories)-1])

		depth = depth[:len(depth)-1]
		directories = directories[:len(directories)-1]

		if err != nil {
			continue
		}

		//TODO: fix for multiple packages
		for _, pkg := range pkgs {
			if *onlyExports {
				hasExports := ast.PackageExports(pkg)
				if !hasExports {
					continue
				}
			}
			typeMap.currentPkg = pkg.Name
			for _, file := range pkg.Files {
				//add imports to directories to check
				if dep < *maxDepth {
					for _, impor := range file.Imports {
						importName := strings.Trim(impor.Path.Value, "\"")
						found := false
						for _, v := range searchedDirectories {
							if v == importName {
								found = true
								break
							}
						}
						for _, v := range directories {
							if v == importName || found {
								found = true
								break
							}
						}
						if !found {
							depth = append(depth, dep+1)
							directories = append(directories, importName)
						}
					}
				}

				//add all types to master list before processing delcarations
				//minimizes creation of unknown types
				for key := range file.Scope.Objects {
					typeMap.lookupOrAdd(key)
				}

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
						f := funcMap.lookupOrAddFromExpr(d.Name.Name, d.Type)
						//funcList = append(funcList, f)
						if d.Recv != nil {
							recv := typeMap.lookupOrAddFromExpr(d.Recv.List[0].Type).base.node
							if recv != nil {
								recv.(*structNode).addFunction(f, d.Recv.List[0])
							}
						}
					}
				}
			}
		}
	}

	for _, i := range interfaceList {
		i.setImplementedBy(structList)
	}

	g := gographviz.NewGraph()
	g.SetName((*filename)[strings.LastIndex(*filename, "/")+1:])
	g.SetDir(true)
	for _, s := range structList {
		g.AddGraphableNode("G", s)
	}
	for _, i := range interfaceList {
		g.AddGraphableNode("G", i)
	}
	s := g.String()
	fmt.Println(s)
}
