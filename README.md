# golangTypeGraph

goals:

Only one instance for every type in the graph
Only one instance of every node in graph

limitations:

ignore names of params in function
ignore non-receiver functions for now
ignore if receiver is pointer or value

ignore group type declarations
type (
	Int int
	Struct struct { x, y int }
)

//ignore <- in channel declarations
ignore vairadic functions
//treat variadic functions like a slice of the type
ignore wraparound lines


performance ideas:

may be improved by converting some []*T fields to map[string]*T for O(1) lookup
may be improved by caching functions from inherited structs/interfaces