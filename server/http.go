package server

import (
	"fmt"
	"net/http"
	"time"

	mcpServer "github.com/mark3labs/mcp-go/server"
)

// DebugLoggingMiddleware intercepts HTTP requests to log details about SSE frames and JSON-RPC flow
func DebugLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		clientIP := r.RemoteAddr

		fmt.Printf("\n--------------------------------------------------\n")
		fmt.Printf("[HTTP DEBUG] Incoming Request:\n")
		fmt.Printf("  Timestamp : %s\n", start.Format("15:04:05.000"))
		fmt.Printf("  Method    : %s\n", r.Method)
		fmt.Printf("  Path      : %s\n", r.URL.Path)
		fmt.Printf("  Client IP : %s\n", clientIP)
		fmt.Printf("  User-Agent: %s\n", r.Header.Get("User-Agent"))
		fmt.Printf("--------------------------------------------------\n")

		// Call next handler
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		fmt.Printf("[HTTP DEBUG] Request Completed for %s in %v\n", r.URL.Path, duration)
	})
}

// SetupHTTPServer sets up HTTP multiplexing and SSE endpoints
func SetupHTTPServer(mcpSrv *mcpServer.MCPServer) http.Handler {
	mux := http.NewServeMux()

	// Create the standard MCP SSE Server handler. 
	// The first argument is your MCP server instance, 
	// and the second argument is the endpoint path where clients post JSON-RPC messages (e.g. "/messages").
	sseServer := mcpServer.NewSSEServer(mcpSrv, "/messages")

	// Mount the SSE handler on "/sse"
	// Under the hood, this handles:
	// 1. GET /sse -> Establishes the text/event-stream connection, sends the endpoint handshake event.
	mux.Handle("/sse", sseServer)

	// Mount the messages endpoint where clients send POST requests
	// 2. POST /messages?session_id=... -> Receives JSON-RPC commands from client.
	mux.Handle("/messages", sseServer)

	// Health check endpoint for quick troubleshooting from browser or curl
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("MCP SSE Health Server is running smoothly!"))
	})

	// Wrap mux with our detailed debug logging middleware
	return DebugLoggingMiddleware(mux)
}