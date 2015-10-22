package main

import "github.com/firegoblin/gographviz"

// EdgeSet does not implement gographviz.EdgesInterface
// instead the parent GraphableNode uses it to implement EdgesInterface
// all edges must share attributes and must use parent as a source or destination
type EdgeSet struct {
	neighbors     []gographviz.GraphableNode
	attrs         gographviz.Attrs
	isDestination bool
}

func (e *EdgeSet) Add(node gographviz.GraphableNode) {
	e.neighbors = append(e.neighbors, node)
}

func (e *EdgeSet) Size() int {
	return len(e.neighbors)
}

func (e *EdgeSet) Empty() bool {
	return e.neighbors == nil
}

func (e *EdgeSet) SetNeighbors(nodes []gographviz.GraphableNode) {
	e.neighbors = make([]gographviz.GraphableNode, 0, len(nodes))
	copy(e.neighbors, nodes)
}

func (e *EdgeSet) Edges(parent gographviz.NodeInterface) []*gographviz.Edge {
	retval := make([]*gographviz.Edge, 0, len(e.neighbors))

	for _, v := range e.neighbors {
		if e.isDestination {
			retval = append(retval, &gographviz.Edge{parent.Name(), "", v.Name(), "", true, e.attrs})
		} else {
			retval = append(retval, &gographviz.Edge{v.Name(), "", parent.Name(), "", true, e.attrs})
		}
	}

	return retval
}
