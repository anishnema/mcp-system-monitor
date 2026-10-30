package main

import (
	"fmt"
	"log"
	"net/http"

	"mcp-system-monitor/server"
	"mcp-system-monitor/tools"

	mcpServer "github.com/mark3labs/mcp-go/server"
)

func main() {
	fmt.Println("Initializing Mac A System Health MCP Server...")

	// 1. Create the base MCP Server instance
	// Parameters: Server Name, Version
	s := mcpServer.NewMCPServer(
		"MacA-SystemMonitor",
		"1.0.0",
		mcpServer.WithLogging(),
	)

	// 2. Register all monitoring tools
	tools.RegisterTools(s)
	fmt.Println("[Init] System monitoring tools successfully registered.")

	// 3. Setup HTTP and SSE routes using our custom server package
	handler := server.SetupHTTPServer(s)

	// 4. Start listening on all network interfaces port 8080 so Mac B can dial in
	port := ":8080"
	fmt.Printf("[Init] Server listening on 0.0.0.0%s\n", port)
	fmt.Println("[Init] SSE Endpoint available at: http://<Mac-A-IP>:8080/sse")

	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("HTTP Server terminated with error: %v", err)
	}
}