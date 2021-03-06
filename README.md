<h1>golangTypeGraph</h1>

<h3>Description</h3>

	Outputs a graphviz compatible .dot file that is a graph of the connections of
	different types in the target dir. Uses the fork github.com/firegoblin/gographviz
	of github.com/awalterschulze/gographviz to handle making the dot file.
	The fork adds an interface for GraphableNode and functions upon it to work with
	the program.  Uses go/ast and go/parser extensively for parsing the target dir.

<h3>Basic Installation</h3>

	brew install graphviz (or visit http://www.graphviz.org/ for other options)

	go get github.com/firegoblin/golangTypeGraph
	go install github.com/firegoblin/golangTypeGraph

<h3>Basic Usage</h3>

	golangTypeGraph github.com/UserName/TargetDir > output.dot
	dot -Tpng output.dot > output.png

<h3>Connections</h3>

	Struct -> Struct:
		parent
		field
		inherited
	Struct -> Interface:
		implments
	Interface -> Struct:
		inherited
		field

<h3>Flags for golangTypeGraph</h3>

  	-depth int
    	maximum depth of recursively searching imports (default 1)
  	-edgeless
    	include nodes that have no edges connected to them (default true)
  	-env string
    	environment variable to use instead of GOPATH (default "GOPATH")
  	-exports
    	marks whether only exported nodes are shown
  	-file string
    	file to parse on, relative to $GOPATH/src (default "github.com/firegoblin/golangTypeGraph")
  	-imax int
    	the maximum number of structs implementing an interface before edges are not drawn (default 9)
    -json
    	include tags for struct fields in print out
  	-pkg string
    	the package that will not have its types prefiexed with package name (default "main")
  	-test
    	whether or not to include test files in the graph


<h3>assumptions for target dir</h3>

	Compiles
	Expects dir to be in $GOPATH/src (unless -env is set, then checks in in $<env>/src)
	Assumes default golang style import folders
	Does not use import . or renaming

<h3>workarounds</h3>

	Repeat definitions are ignored (temporary workaround for OS specific files)
	Error handling added to protect against crashes on unexpected cases


<h3>flag ideas</h3>

	Output file flag
	Perform the dot command for conversion to graphics file types through the program
	Modify verbosity of nodes


<h3>future improvements</h3>

	Add connections to functions, partiularly interfaces to functions they're used in
	Add tests to allow safer changes.
	Print the string to a file instead of piping
	Better attributes, especially to ensure a less cluttered layout.
	Seperate ast_helper into on package as contains some useful general AST functions.