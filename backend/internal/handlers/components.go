package handlers

import (
	"encoding/json"
	"net/http"
)

type ComponentItem struct {
	ID                string          `json:"id"`
	Type              string          `json:"type"`
	Label             string          `json:"label"`
	Icon              string          `json:"icon"`
	DefaultProperties map[string]bool `json:"defaultProperties"`
}

type ComponentCategory struct {
	Category      string          `json:"category"`
	CategoryLabel string          `json:"categoryLabel"`
	CategoryIcon  string          `json:"categoryIcon"`
	Items         []ComponentItem `json:"items"`
}

var componentLibrary = []ComponentCategory{
	{
		Category: "frontend", CategoryLabel: "Client", CategoryIcon: "🖥️",
		Items: []ComponentItem{
			{ID: "web-frontend", Type: "frontend", Label: "Web Frontend", Icon: "web-frontend.svg", DefaultProperties: map[string]bool{}},
			{ID: "mobile-frontend", Type: "frontend", Label: "Mobile Frontend", Icon: "mobile-frontend.svg", DefaultProperties: map[string]bool{}},
		},
	},
	{
		Category: "backend", CategoryLabel: "Backend", CategoryIcon: "⚙️",
		Items: []ComponentItem{
			{ID: "api-server", Type: "backend", Label: "API Server", Icon: "api-server.svg",
				DefaultProperties: map[string]bool{"stateless": false, "healthChecks": false, "circuitBreaker": false}},
			{ID: "microservice", Type: "microservice", Label: "Microservice", Icon: "microservice.svg",
				DefaultProperties: map[string]bool{"stateless": false, "healthChecks": false, "circuitBreaker": false}},
		},
	},
	{
		Category: "database", CategoryLabel: "Database", CategoryIcon: "🗄️",
		Items: []ComponentItem{
			{ID: "postgresql", Type: "database", Label: "PostgreSQL", Icon: "postgresql.svg",
				DefaultProperties: map[string]bool{"replication": false, "backups": false, "connectionPooling": false, "encryptionAtRest": false}},
			{ID: "mongodb", Type: "database", Label: "MongoDB", Icon: "mongodb.svg",
				DefaultProperties: map[string]bool{"replication": false, "sharding": false, "backups": false}},
		},
	},
	{
		Category: "cache", CategoryLabel: "Cache", CategoryIcon: "⚡",
		Items: []ComponentItem{
			{ID: "redis", Type: "cache", Label: "Redis", Icon: "redis.svg",
				DefaultProperties: map[string]bool{"persistence": false, "clustering": false}},
			{ID: "memcached", Type: "cache", Label: "Memcached", Icon: "memcached.svg",
				DefaultProperties: map[string]bool{"clustering": false}},
		},
	},
	{
		Category: "queue", CategoryLabel: "Queue", CategoryIcon: "📨",
		Items: []ComponentItem{
			{ID: "kafka", Type: "queue", Label: "Kafka", Icon: "kafka.svg", DefaultProperties: map[string]bool{}},
			{ID: "rabbitmq", Type: "queue", Label: "RabbitMQ", Icon: "rabbitmq.svg", DefaultProperties: map[string]bool{}},
		},
	},
	{
		Category: "network", CategoryLabel: "Network", CategoryIcon: "🌐",
		Items: []ComponentItem{
			{ID: "load-balancer", Type: "loadbalancer", Label: "Load Balancer", Icon: "load-balancer.svg",
				DefaultProperties: map[string]bool{"healthChecks": false, "sslTermination": false}},
			{ID: "api-gateway", Type: "apigateway", Label: "API Gateway", Icon: "api-gateway.svg", DefaultProperties: map[string]bool{}},
			{ID: "cdn", Type: "cdn", Label: "CDN", Icon: "cdn.svg", DefaultProperties: map[string]bool{}},
		},
	},
	{
		Category: "storage", CategoryLabel: "Storage", CategoryIcon: "💾",
		Items: []ComponentItem{
			{ID: "object-storage", Type: "storage", Label: "Object Storage", Icon: "object-storage.svg", DefaultProperties: map[string]bool{}},
		},
	},
}

func Components(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"components": componentLibrary,
	})
}
