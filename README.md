# golangTypeGraph

goals:

Only one instance for every type in the graph

limitations:

ignore non-receiver functions for now
ignore group type definitions

type (
	Int int
	Struct struct { x, y int }
)


treat variadric functions like a slice of the type