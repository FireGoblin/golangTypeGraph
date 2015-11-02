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
var edgelessNodes = flag.Bool("edgeless", true, "include nodes that have no edges connected to them")

//var operatingSystem = flag.String("os", "linux", "define the os to use when choosing os specific files")
//var operatingArchitecture = flag.String("arch", "amd64", "define the architecture to use when choosing os specific files")

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
			//fmt.Println(node)
			switch n := node.(type) {
			case *unknownNode:
				*structList = append(*structList, n.remakeStruct(decl))
			default:
				fmt.Fprintln(os.Stderr, "attempt to recreate struct:", n)
			}
		}
	}
}

func processFuncDecl(decl ast.Decl) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		f := funcMap.lookupOrAddFromExpr(d.Name.Name, d.Type)
		//funcList = append(funcList, f)
		if d.Recv != nil {
			recv := typeMap.lookupOrAddFromExpr(d.Recv.List[0].Type).base.node
			if recv != nil {
				//fmt.Println(d.Recv.List[0])
				switch r := recv.(type) {
				case *structNode:
					r.addFunction(f, d.Recv.List[0])
				case *unknownNode:
					r.addFunction(f, d.Recv.List[0])
				default:
					panic("trying to add receiver to interface")
				}
			}
		}
	}
}

//var osVar = [...]string{"freebsd", "windows", "linux", "dragonfly", "openbsd", "netbsd", "darwin", "plan9", "solaris", "nacl"}
//var arch = [...]string{"386", "amd64", "arm"}

func legalFile(f os.FileInfo) bool {
	return *includeTestFiles || !strings.Contains(f.Name(), "_test.go")
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()

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

	//loop until directories to search is empty
	for len(directories) > 0 {
		pkgs, err = parser.ParseDir(fset, gopath+directories[len(directories)-1], legalFile, 0)

		//dep is used to
		currentDepth := depth[len(depth)-1]
		searchedDirectories = append(searchedDirectories, directories[len(directories)-1])

		depth = depth[:len(depth)-1]
		directories = directories[:len(directories)-1]

		//skip this folder if there was an error parsing, usually meaning the directory is not found
		if err != nil {
			continue
		}

		for _, pkg := range pkgs {
			//remove unexported types/functions if onlyExports
			if *onlyExports {
				hasExports := ast.PackageExports(pkg)
				if !hasExports {
					continue
				}
			}
			typeMap.currentPkg = pkg.Name

			for _, file := range pkg.Files {
				//add imports to directories to check if not at maxDepth yet
				if currentDepth < *maxDepth {
					for _, impor := range file.Imports {
						importName := strings.Trim(impor.Path.Value, "\"")

						if !containsString(searchedDirectories, importName) && !containsString(directories, importName) {
							depth = append(depth, currentDepth+1)
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
					processFuncDecl(decl)
				}
			}
		}
	}

	//set all interfaces' implementedByCache
	for _, i := range interfaceList {
		i.setImplementedBy(structList)
	}

	//Create the graph and print it out
	g := gographviz.NewGraph()
	g.SetName((*filename)[strings.LastIndex(*filename, "/")+1:])
	g.SetDir(true)
	for _, s := range structList {
		g.AddGraphableNode("G", s)
	}
	for _, i := range interfaceList {
		g.AddGraphableNode("G", i)
	}
	if !*edgelessNodes {
		fmt.Fprintln(os.Stderr, "removing edgeless nodes")
		g.RemoveEdgelessNodes("G")
	}
	s := g.String()
	fmt.Println(s)
}
