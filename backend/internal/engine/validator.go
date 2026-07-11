package engine

import (
	"github.com/vimal/sams/internal/models"
)

type Validator struct {
	rules []models.Rule
}

func NewValidator(rules []models.Rule) *Validator {
	return &Validator{rules: rules}
}

// RuleCount returns the number of rules the validator checks.
func (v *Validator) RuleCount() int {
	return len(v.rules)
}

func (v *Validator) Validate(req models.ValidationRequest) models.ValidationResult {
	if len(req.Nodes) == 0 {
		return models.ValidationResult{
			Errors:   []models.Issue{},
			Warnings: []models.Issue{},
			Passed:   []models.PassedRule{},
			Score:    100,
		}
	}

	graph := BuildGraph(req.Nodes, req.Edges)
	ctx := req.Context

	errors := []models.Issue{}
	warnings := []models.Issue{}
	rulesFired := make(map[string]bool)

	addError := func(rule models.Rule, nodes []string, edges []string) {
		rulesFired[rule.ID] = true
		errors = append(errors, models.Issue{
			RuleID:        rule.ID,
			Category:      rule.Category,
			Title:         rule.Title,
			Description:   rule.Description,
			AffectedNodes: nodes,
			AffectedEdges: edges,
			Suggestion:    rule.Suggestion,
		})
	}

	addWarning := func(rule models.Rule, nodes []string, edges []string) {
		rulesFired[rule.ID] = true
		warnings = append(warnings, models.Issue{
			RuleID:        rule.ID,
			Category:      rule.Category,
			Title:         rule.Title,
			Description:   rule.Description,
			AffectedNodes: nodes,
			AffectedEdges: edges,
			Suggestion:    rule.Suggestion,
		})
	}

	ruleMap := make(map[string]models.Rule)
	for _, r := range v.rules {
		ruleMap[r.ID] = r
	}

	// ── Performance ──────────────────────────────────────────────────────────

	// rule_001: Missing Cache Layer
	if r, ok := ruleMap["rule_001"]; ok {
		highTraffic := ctx.TrafficLevel == "high" || ctx.TrafficLevel == "massive"
		backendToDb := len(edgesBetween(graph, "backend", "database")) > 0 ||
			len(edgesBetween(graph, "microservice", "database")) > 0
		noCache := !hasNodeOfType(graph, "cache")
		if highTraffic && backendToDb && noCache {
			nodes := append(nodeIDsOfType(graph, "backend"), nodeIDsOfType(graph, "database")...)
			addWarning(r, nodes, nil)
		}
	}

	// rule_005: Synchronous Operations — multiple backends, all sync, no queue
	if r, ok := ruleMap["rule_005"]; ok {
		multiBackend := countNodesByType(graph, "backend")+countNodesByType(graph, "microservice") > 1
		noQueue := !hasNodeOfType(graph, "queue")
		if multiBackend && allEdgesSync(graph) && noQueue {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_018: N+1 Query Problem
	if r, ok := ruleMap["rule_018"]; ok {
		backendToDb := len(edgesBetween(graph, "backend", "database")) > 0
		noCache := !hasNodeOfType(graph, "cache")
		readHeavy := ctx.ReadWriteRatio == "read_heavy"
		highTraffic := ctx.TrafficLevel == "high" || ctx.TrafficLevel == "massive"
		if backendToDb && noCache && readHeavy && highTraffic {
			nodes := append(nodeIDsOfType(graph, "backend"), nodeIDsOfType(graph, "database")...)
			addWarning(r, nodes, nil)
		}
	}

	// rule_019: Missing Database Indexes
	if r, ok := ruleMap["rule_019"]; ok {
		highTraffic := ctx.TrafficLevel == "high" || ctx.TrafficLevel == "massive"
		hasDB := hasNodeOfType(graph, "database")
		if highTraffic && hasDB {
			addWarning(r, nodeIDsOfType(graph, "database"), nil)
		}
	}

	// rule_020: No CDN for Static Assets
	if r, ok := ruleMap["rule_020"]; ok {
		hasFrontend := hasNodeOfType(graph, "frontend")
		globalBase := ctx.UserBase == "global"
		noCDN := !hasNodeOfType(graph, "cdn")
		if hasFrontend && globalBase && noCDN {
			addWarning(r, nodeIDsOfType(graph, "frontend"), nil)
		}
	}

	// rule_021: No Computed Result Caching
	if r, ok := ruleMap["rule_021"]; ok {
		hasBackend := hasNodeOfType(graph, "backend") || hasNodeOfType(graph, "microservice")
		hasDB := hasNodeOfType(graph, "database")
		readHeavy := ctx.ReadWriteRatio == "read_heavy"
		highTraffic := ctx.TrafficLevel == "high" || ctx.TrafficLevel == "massive"
		noCache := !hasNodeOfType(graph, "cache")
		if hasBackend && hasDB && readHeavy && highTraffic && noCache {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// ── Scalability ───────────────────────────────────────────────────────────

	// rule_003: No Queue for High Write Traffic
	if r, ok := ruleMap["rule_003"]; ok {
		highTraffic := ctx.TrafficLevel == "high" || ctx.TrafficLevel == "massive"
		writeHeavy := ctx.ReadWriteRatio == "write_heavy"
		directToDb := len(edgesBetween(graph, "backend", "database")) > 0
		noQueue := !hasNodeOfType(graph, "queue")
		if highTraffic && writeHeavy && directToDb && noQueue {
			nodes := append(nodeIDsOfType(graph, "backend"), nodeIDsOfType(graph, "database")...)
			addWarning(r, nodes, nil)
		}
	}

	// rule_004: Wrong Database for Use Case (MongoDB + relational context)
	if r, ok := ruleMap["rule_004"]; ok {
		for _, n := range graph.Nodes {
			if n.Type == "database" && n.Data.Subtype == "mongodb" {
				// Heuristic: relational signals = high team, scale stage, complex edge graph
				complexRelations := countNodesByType(graph, "backend")+countNodesByType(graph, "microservice") > 2
				scaleTeam := ctx.TeamSize == "large" || ctx.TeamSize == "medium"
				if complexRelations && scaleTeam {
					addWarning(r, []string{n.ID}, nil)
				}
			}
		}
	}

	// rule_022: No Centralized Session Storage
	if r, ok := ruleMap["rule_022"]; ok {
		multiBackend := countNodesByType(graph, "backend")+countNodesByType(graph, "microservice") > 1
		hasGateway := hasNodeOfType(graph, "apigateway")
		noRedis := true
		for _, n := range graph.Nodes {
			if n.Type == "cache" && n.Data.Subtype == "redis" {
				noRedis = false
				break
			}
		}
		if multiBackend && hasGateway && noRedis {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_023: Stateful Servers
	if r, ok := ruleMap["rule_023"]; ok {
		multiBackend := countNodesByType(graph, "backend")+countNodesByType(graph, "microservice") > 1
		if multiBackend {
			for _, n := range graph.Nodes {
				if (n.Type == "backend" || n.Type == "microservice") &&
					n.Data.Properties["stateless"] == false {
					addWarning(r, []string{n.ID}, nil)
				}
			}
		}
	}

	// rule_024: Monolith at Large Scale
	if r, ok := ruleMap["rule_024"]; ok {
		singleBackend := countNodesByType(graph, "backend") == 1 &&
			countNodesByType(graph, "microservice") == 0
		massiveTraffic := ctx.TrafficLevel == "massive" || ctx.TrafficLevel == "high"
		largeTeam := ctx.TeamSize == "large"
		if singleBackend && massiveTraffic && largeTeam {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_025: Premature Microservices
	if r, ok := ruleMap["rule_025"]; ok {
		manyMicroservices := countNodesByType(graph, "microservice") > 2
		smallTeam := ctx.TeamSize == "solo" || ctx.TeamSize == "small"
		early := ctx.Stage == "early"
		if manyMicroservices && smallTeam && early {
			addWarning(r, nodeIDsOfType(graph, "microservice"), nil)
		}
	}

	// rule_026: No Circuit Breaker
	if r, ok := ruleMap["rule_026"]; ok {
		hasMicroservices := hasNodeOfType(graph, "microservice")
		if hasMicroservices && allEdgesSync(graph) {
			for _, n := range graph.Nodes {
				if n.Type == "microservice" && !n.Data.Properties["circuitBreaker"] {
					addWarning(r, []string{n.ID}, nil)
				}
			}
		}
	}

	// rule_027: No Auto Scaling
	if r, ok := ruleMap["rule_027"]; ok {
		highTraffic := ctx.TrafficLevel == "high" || ctx.TrafficLevel == "massive"
		noLB := !hasNodeOfType(graph, "loadbalancer")
		singleBackend := countNodesByType(graph, "backend")+countNodesByType(graph, "microservice") <= 1
		if highTraffic && (noLB || singleBackend) {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_028: Redis Overkill for Simple Cache (suggestion only)
	if r, ok := ruleMap["rule_028"]; ok {
		for _, n := range graph.Nodes {
			if n.Type == "cache" && n.Data.Subtype == "redis" {
				hasGateway := hasNodeOfType(graph, "apigateway")
				multiBackend := countNodesByType(graph, "backend")+countNodesByType(graph, "microservice") > 1
				// Overkill only when simple single-backend setup with no auth implied
				if !hasGateway && !multiBackend {
					addWarning(r, []string{n.ID}, nil)
				}
			}
		}
	}

	// rule_029: Memcached Cannot Handle Sessions
	if r, ok := ruleMap["rule_029"]; ok {
		for _, n := range graph.Nodes {
			if n.Type == "cache" && n.Data.Subtype == "memcached" {
				hasGateway := hasNodeOfType(graph, "apigateway")
				if hasGateway {
					addError(r, []string{n.ID}, nil)
				}
			}
		}
	}

	// ── Reliability ───────────────────────────────────────────────────────────

	// rule_002: Files Stored in Wrong Place
	if r, ok := ruleMap["rule_002"]; ok {
		backendToDb := len(edgesBetween(graph, "backend", "database")) > 0
		noStorage := !hasNodeOfType(graph, "storage")
		// Heuristic: if there's a DB but no object storage and high traffic → likely storing files in DB
		highTraffic := ctx.TrafficLevel == "high" || ctx.TrafficLevel == "massive"
		if backendToDb && noStorage && highTraffic {
			nodes := append(nodeIDsOfType(graph, "backend"), nodeIDsOfType(graph, "database")...)
			addWarning(r, nodes, nil)
		}
	}

	// rule_013: Single Point of Failure (Server)
	if r, ok := ruleMap["rule_013"]; ok {
		singleBackend := countNodesByType(graph, "backend") == 1 &&
			countNodesByType(graph, "microservice") == 0
		noLB := !hasNodeOfType(graph, "loadbalancer")
		if singleBackend && noLB {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_014: No Database Replication
	if r, ok := ruleMap["rule_014"]; ok {
		highTraffic := ctx.TrafficLevel == "high" || ctx.TrafficLevel == "massive"
		scale := ctx.Stage == "scale" || ctx.Stage == "growing"
		for _, n := range graph.Nodes {
			if n.Type == "database" && !n.Data.Properties["replication"] &&
				(highTraffic || scale) {
				addWarning(r, []string{n.ID}, nil)
			}
		}
	}

	// rule_015: No Backup Strategy
	if r, ok := ruleMap["rule_015"]; ok {
		for _, n := range graph.Nodes {
			if n.Type == "database" && !n.Data.Properties["backups"] {
				addWarning(r, []string{n.ID}, nil)
			}
		}
	}

	// rule_016: No Monitoring
	if r, ok := ruleMap["rule_016"]; ok {
		hasBackend := hasNodeOfType(graph, "backend") || hasNodeOfType(graph, "microservice")
		growing := ctx.Stage == "growing" || ctx.Stage == "scale"
		noMonitoring := !hasNodeOfType(graph, "monitoring")
		if hasBackend && growing && noMonitoring {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_017: No Deployment Strategy
	if r, ok := ruleMap["rule_017"]; ok {
		multiBackend := countNodesByType(graph, "backend")+countNodesByType(graph, "microservice") > 1
		scale := ctx.Stage == "scale"
		if multiBackend && scale {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// ── Security ──────────────────────────────────────────────────────────────

	// rule_006: Frontend directly to Database
	if r, ok := ruleMap["rule_006"]; ok {
		for _, n := range graph.Nodes {
			if n.Type == "frontend" {
				for _, connID := range graph.Adjacency[n.ID] {
					conn, ok := graph.Nodes[connID]
					if ok && conn.Type == "database" {
						// Find the edge
						var edgeIDs []string
						for _, e := range graph.Edges {
							if e.Source == n.ID && e.Target == connID {
								edgeIDs = append(edgeIDs, e.ID)
							}
						}
						addError(r, []string{n.ID, connID}, edgeIDs)
					}
				}
			}
		}
	}

	// rule_007: Missing Authentication
	if r, ok := ruleMap["rule_007"]; ok {
		hasBackend := hasNodeOfType(graph, "backend") || hasNodeOfType(graph, "microservice")
		noGateway := !hasNodeOfType(graph, "apigateway")
		if hasBackend && noGateway {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_008: Missing Authorization
	if r, ok := ruleMap["rule_008"]; ok {
		// Fires when there's an API gateway (auth implied) but no authorization hints
		// In v1 we can't detect RBAC, so we fire for scale/growing stage with gateway
		hasGateway := hasNodeOfType(graph, "apigateway")
		growing := ctx.Stage == "growing" || ctx.Stage == "scale"
		if hasGateway && growing {
			addWarning(r, nodeIDsOfType(graph, "apigateway"), nil)
		}
	}

	// rule_009: SQL Injection Risk
	if r, ok := ruleMap["rule_009"]; ok {
		hasPostgres := false
		for _, n := range graph.Nodes {
			if n.Type == "database" && n.Data.Subtype == "postgresql" {
				hasPostgres = true
				break
			}
		}
		hasBackend := hasNodeOfType(graph, "backend") || hasNodeOfType(graph, "microservice")
		noGateway := !hasNodeOfType(graph, "apigateway")
		if hasPostgres && hasBackend && noGateway {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_010: Verbose Error Messages
	if r, ok := ruleMap["rule_010"]; ok {
		hasBackend := hasNodeOfType(graph, "backend") || hasNodeOfType(graph, "microservice")
		// In v1 there's no "error handling" node; fire for early stage public APIs
		noGateway := !hasNodeOfType(graph, "apigateway")
		if hasBackend && noGateway {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_011: No Rate Limiting
	if r, ok := ruleMap["rule_011"]; ok {
		hasFrontend := hasNodeOfType(graph, "frontend")
		noGateway := !hasNodeOfType(graph, "apigateway")
		noLB := !hasNodeOfType(graph, "loadbalancer")
		if hasFrontend && noGateway && noLB {
			addWarning(r, nodeIDsOfType(graph, "backend"), nil)
		}
	}

	// rule_012: No HTTPS/TLS
	if r, ok := ruleMap["rule_012"]; ok {
		for _, e := range graph.Edges {
			if e.Protocol == "http" {
				src := graph.Nodes[e.Source]
				tgt := graph.Nodes[e.Target]
				// Only flag external-facing edges (frontend or gateway involved)
				if src.Type == "frontend" || tgt.Type == "frontend" ||
					src.Type == "apigateway" || tgt.Type == "apigateway" {
					addError(r, []string{e.Source, e.Target}, []string{e.ID})
				}
			}
		}
	}

	// ── Build passed rules ────────────────────────────────────────────────────
	passed := []models.PassedRule{}
	for _, r := range v.rules {
		if !rulesFired[r.ID] {
			passed = append(passed, models.PassedRule{
				RuleID:   r.ID,
				Category: r.Category,
				Title:    r.Title,
			})
		}
	}

	score := calculateScore(errors, warnings, v.rules, ctx)

	return models.ValidationResult{
		Errors:   errors,
		Warnings: warnings,
		Passed:   passed,
		Score:    score,
	}
}
