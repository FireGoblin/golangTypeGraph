package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	. "regexp"
)

var builtinTypes = [...]string{"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64",
	"rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr"}
var typeMap = MasterTypeMap(make(map[string]*Type))
var funcMap = MasterFuncMap(make(map[string]*Function))

var filename = flag.String("file", ".", "file to parse on")

var interfac = MustCompile(`^type ([^ ]+) interface \{$`)

//this is for function declarations in interface, so unnamed params
var FunctionParser = MustCompile(`^([\w]+)(\(.*?\) .*)$`)

//Made for parsing a function type, such as func(int, bool) (int, error)
var FuncTypeParser = MustCompile(`^func\((.*?)\) (.*)$`)

var functionLine = MustCompile(`^func (.*)$`)
var receiverFunction = MustCompile(`^\([^ ]+ \**([^ ]+)\) (.*)$`)

var struc = MustCompile(`^type ([^ ]+) struct(.*)$`)
var strucMulti = MustCompile(`^ \{$`)
var structOneLiner = MustCompile(`^\{ ([^ ]+(?:, [^ ]+)) (.+) \}$`)

//may be removable
var AnonymousStructParser = MustCompile(`^[^ ]+$`)

//pair of struct name and struct type
var NamedTypeParser = MustCompile(`^(.+?)[ ]+((?:(?:<-)?chan.+)|(?:func.+)|(?:[^ ]+))$`)

func check(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

type State struct {
	inStruct    bool
	inInterface bool

	structList    []*Struct
	interfaceList []*Interface
}

func processLine(line string, state *State, fileReader *bufio.Reader) {
	if x := struc.FindStringSubmatch(line); x != nil {
		fmt.Println("struc:", x[0])
		name := x[1]
		if y := strucMulti.FindStringSubmatch(x[2]); y = nil {

		} else if y := strucOneLiner.FindStringSubmatch(x[2]); y = nil {

			} else {
				panic("unrecognized struct")
			}
		}
	} else if x := functionLine.FindStringSubmatch(line); x != nil {
		fmt.Println("function:", x[0])
	} else if x := interfac.FindStringSubmatch(line); x != nil {
		fmt.Println("interfac:", x[0])
	}
}

func main() {
	//initialize master map with builtin types
	for _, v := range builtinTypes {
		typeMap.lookupOrAdd(v)
	}

	file, err := os.Open(*filename)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(bufio.NewReader(file))

	state := &State{false, false, make([]*Struct, 0), make([]*Interface, 0)}

	for scanner.Scan() {
		processLine(scanner.Text(), state, scanner)
	}

	for _, i := range state.interfaceList {
		implementingStructs := i.implementedBy(state.structList)

		fmt.Println("Interface", i, "is implemented by the following types:")
		for _, s := range implementingStructs {
			fmt.Println("   ", s)
		}
	}
}
