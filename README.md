<h1>golangTypeGraph</h1>

Description:

	Outputs a graphviz compatible .dot file that is a graph of the connections of
	different types in the target dir. Uses the fork github.com/firegoblin/gographviz
	of https://github.com/awalterschulze/gographviz to handle making the dot file.  
	The fork adds an interface for GraphableNode and functions upon it to work with
	the program.  Uses go/ast and go/parser extensively for parsing the target dir.

Basic Installation:

	brew install graphviz (or visit http://www.graphviz.org/ for other options)

	go get github.com/firegoblin/golangTypeGraph
	go get github.com/firegoblin/gographviz
	go install github.com/firegoblin/golangTypeGraph

Basic Usage:

	golangTypeGraph github.com/UserName/TargetDir > output.dot
	dot -Tpng output.dot > output.png

Connections:

	Struct -> Struct:
		parent
		field
		inherited
	Struct -> Interface:
		implments
	Interface -> Struct:
		inherited
		field

Flags for golangTypeGraph:

	-exports
    	marks whether only exported nodes are shown
  	-file string
    	file to parse on, relative to $GOPATH/src
    	(default "github.com/firegoblin/golangTypeGraph")
  	-imax int
    	the maximum number of structs implementing an interface before edges 
    	are not drawn (default 9)
  	-imports
    	whether or not to parse import directories recrusively (default true)
  	-pkg string
    	the package that will not have its types prefiexed with package name 
    	(default "main")
  	-test
    	whether or not to include test files in the graph (default true)


assumptions for target dir:

	Compiles
	Receiver functions are in same file as the struct
	Expects dir to be in $GOPATH/src (unless -goroot=true, then its in $GOROOT/src)
	Assumes default golang style import folders
	Does not use import .


flag ideas:

	Flag for particular level of recursion
	Output file flag
	Perform the dot command for conversion to graphics file types through the program
	Modify verbosity of nodes


future improvements:

	Add connections to functions, partiularly interfaces to functions they're used in
	Add tests to allow safer changes.


known bug:

	Functions may display params/results with or without pkg name indeterminately.


possible untested issues:

	May have issues with elaborate types such as complex func types