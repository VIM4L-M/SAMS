package models

type Graph struct {
	Nodes     map[string]Node
	Edges     []Edge
	Adjacency map[string][]string // nodeID → []connected nodeIDs (directed)
}
