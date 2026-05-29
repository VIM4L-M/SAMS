package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vimal/sams/internal/engine"
	"github.com/vimal/sams/internal/models"
)

const (
	maxNodes     = 100
	maxEdges     = 200
	maxStringLen = 256
)

var validNodeTypes = map[string]bool{
	"frontend": true, "backend": true, "microservice": true,
	"database": true, "cache": true, "queue": true,
	"loadbalancer": true, "apigateway": true, "cdn": true, "storage": true,
}

var validConnectionTypes = map[string]bool{"sync": true, "async": true}
var validProtocols = map[string]bool{"http": true, "https": true, "grpc": true, "amqp": true}
var validTrafficLevels = map[string]bool{"low": true, "medium": true, "high": true, "massive": true}
var validRWRatios = map[string]bool{"read_heavy": true, "balanced": true, "write_heavy": true}
var validUserBases = map[string]bool{"local": true, "regional": true, "global": true}
var validTeamSizes = map[string]bool{"solo": true, "small": true, "medium": true, "large": true}
var validStages = map[string]bool{"early": true, "growing": true, "scale": true}

func safeString(s string) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) > maxStringLen {
		runes := []rune(s)
		s = string(runes[:maxStringLen])
	}
	return s
}

func validateRequest(req *models.ValidationRequest) string {
	if len(req.Nodes) > maxNodes {
		return "too many nodes"
	}
	if len(req.Edges) > maxEdges {
		return "too many edges"
	}

	nodeIDs := make(map[string]bool)
	for i, n := range req.Nodes {
		n.ID = safeString(n.ID)
		n.Type = strings.ToLower(safeString(n.Type))
		if n.ID == "" {
			return "node missing id"
		}
		if nodeIDs[n.ID] {
			return "duplicate node id: " + n.ID
		}
		nodeIDs[n.ID] = true
		if !validNodeTypes[n.Type] {
			return "invalid node type: " + n.Type
		}
		n.Data.Label = safeString(n.Data.Label)
		n.Data.Subtype = strings.ToLower(safeString(n.Data.Subtype))
		req.Nodes[i] = n
	}

	edgeIDs := make(map[string]bool)
	for i, e := range req.Edges {
		e.ID = safeString(e.ID)
		e.Source = safeString(e.Source)
		e.Target = safeString(e.Target)
		e.ConnectionType = strings.ToLower(safeString(e.ConnectionType))
		e.Protocol = strings.ToLower(safeString(e.Protocol))
		if e.ID == "" {
			return "edge missing id"
		}
		if edgeIDs[e.ID] {
			return "duplicate edge id: " + e.ID
		}
		edgeIDs[e.ID] = true
		if !nodeIDs[e.Source] {
			return "edge references unknown source node"
		}
		if !nodeIDs[e.Target] {
			return "edge references unknown target node"
		}
		if e.ConnectionType != "" && !validConnectionTypes[e.ConnectionType] {
			return "invalid connectionType: " + e.ConnectionType
		}
		if e.Protocol != "" && !validProtocols[e.Protocol] {
			return "invalid protocol: " + e.Protocol
		}
		req.Edges[i] = e
	}

	ctx := &req.Context
	if ctx.TrafficLevel != "" && !validTrafficLevels[ctx.TrafficLevel] {
		return "invalid trafficLevel"
	}
	if ctx.ReadWriteRatio != "" && !validRWRatios[ctx.ReadWriteRatio] {
		return "invalid readWriteRatio"
	}
	if ctx.UserBase != "" && !validUserBases[ctx.UserBase] {
		return "invalid userBase"
	}
	if ctx.TeamSize != "" && !validTeamSizes[ctx.TeamSize] {
		return "invalid teamSize"
	}
	if ctx.Stage != "" && !validStages[ctx.Stage] {
		return "invalid stage"
	}

	return ""
}

func Validate(v *engine.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.ContentLength > 1<<20 { // 1 MB limit
			http.Error(w, `{"error":"request too large"}`, http.StatusRequestEntityTooLarge)
			return
		}

		var req models.ValidationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if errMsg := validateRequest(&req); errMsg != "" {
			http.Error(w, `{"error":"invalid request: `+errMsg+`"}`, http.StatusBadRequest)
			return
		}

		start := time.Now()
		result := v.Validate(req)
		elapsed := time.Since(start).Milliseconds()

		resp := models.ValidationResponse{
			Results: result,
			Metadata: models.Metadata{
				RulesChecked: 29,
				TimeMs:       elapsed,
			},
		}
		json.NewEncoder(w).Encode(resp)
	}
}
