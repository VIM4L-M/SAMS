package engine

import "github.com/vimal/sams/internal/models"

func BuildGraph(nodes []models.Node, edges []models.Edge) models.Graph {
	graph := models.Graph{
		Nodes:     make(map[string]models.Node),
		Adjacency: make(map[string][]string),
	}
	for _, node := range nodes {
		graph.Nodes[node.ID] = node
		graph.Adjacency[node.ID] = []string{}
	}
	for _, edge := range edges {
		graph.Edges = append(graph.Edges, edge)
		graph.Adjacency[edge.Source] = append(
			graph.Adjacency[edge.Source], edge.Target,
		)
	}
	return graph
}

func countNodesByType(graph models.Graph, nodeType string) int {
	count := 0
	for _, n := range graph.Nodes {
		if n.Type == nodeType {
			count++
		}
	}
	return count
}

func hasNodeOfType(graph models.Graph, nodeType string) bool {
	for _, n := range graph.Nodes {
		if n.Type == nodeType {
			return true
		}
	}
	return false
}

func nodesOfType(graph models.Graph, nodeType string) []models.Node {
	var result []models.Node
	for _, n := range graph.Nodes {
		if n.Type == nodeType {
			result = append(result, n)
		}
	}
	return result
}

func nodeIDsOfType(graph models.Graph, nodeType string) []string {
	var ids []string
	for _, n := range graph.Nodes {
		if n.Type == nodeType {
			ids = append(ids, n.ID)
		}
	}
	return ids
}

// edgesBetween returns edges where source type → target type match.
func edgesBetween(graph models.Graph, srcType, tgtType string) []models.Edge {
	var result []models.Edge
	for _, e := range graph.Edges {
		src, srcOk := graph.Nodes[e.Source]
		tgt, tgtOk := graph.Nodes[e.Target]
		if srcOk && tgtOk && src.Type == srcType && tgt.Type == tgtType {
			result = append(result, e)
		}
	}
	return result
}

// directlyConnected returns true if source node directly connects to target node.
func directlyConnected(graph models.Graph, srcID, tgtID string) bool {
	for _, id := range graph.Adjacency[srcID] {
		if id == tgtID {
			return true
		}
	}
	return false
}

// allEdgesSync returns true when every edge in the graph is synchronous.
func allEdgesSync(graph models.Graph) bool {
	for _, e := range graph.Edges {
		if e.ConnectionType == "async" {
			return false
		}
	}
	return true
}
