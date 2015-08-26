# golangTypeGraph

goals:

Only one instance for every type in the graph
Only one instance of every node type in graph

limitations:

ignore names of params in function
ignore non-receiver functions for now
ignore group type definitions

type (
	Int int
	Struct struct { x, y int }
)


ignore variadic functions like a slice of the type
ignore wraparound lines