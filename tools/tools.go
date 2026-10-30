package tools

import (
	"context"
	"fmt"

	"mcp-system-monitor/monitor"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools adds all health-monitoring tools to the MCP server
func RegisterTools(s *server.MCPServer) {

	// 1. CPU Tool
	cpuTool := mcp.NewTool("get_cpu_usage",
		mcp.WithDescription("Get current CPU usage percentage of the host machine"),
	)
	s.AddTool(cpuTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fmt.Printf("[Tool Handler] Executing 'get_cpu_usage'...\n")
		val, err := monitor.GetCPUUsage()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get CPU: %v", err)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("CPU Usage: %.2f%%", val)), nil
	})

	// 2. Memory Tool
	memTool := mcp.NewTool("get_memory_usage",
		mcp.WithDescription("Get RAM allocation and usage statistics"),
	)
	s.AddTool(memTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fmt.Printf("[Tool Handler] Executing 'get_memory_usage'...\n")
		v, err := monitor.GetMemoryUsage()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get memory: %v", err)), nil
		}
		report := fmt.Sprintf("Total RAM: %d MB\nUsed RAM: %d MB (%.2f%%)\nFree RAM: %d MB",
			v.Total/1024/1024, v.Used/1024/1024, v.UsedPercent, v.Free/1024/1024)
		return mcp.NewToolResultText(report), nil
	})

	// 3. Disk Tool
	diskTool := mcp.NewTool("get_disk_usage",
		mcp.WithDescription("Get disk partition storage stats"),
	)
	s.AddTool(diskTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fmt.Printf("[Tool Handler] Executing 'get_disk_usage'...\n")
		d, err := monitor.GetDiskUsage()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get disk: %v", err)), nil
		}
		report := fmt.Sprintf("Total Disk: %d GB\nUsed Disk: %d GB (%.2f%%)\nFree Disk: %d GB",
			d.Total/1024/1024/1024, d.Used/1024/1024/1024, d.UsedPercent, d.Free/1024/1024/1024)
		return mcp.NewToolResultText(report), nil
	})

	// 4. Network Activity Tool
	netTool := mcp.NewTool("get_network_activity",
		mcp.WithDescription("Get current network bandwidth speed (bytes sent/received per second)"),
	)
	s.AddTool(netTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fmt.Printf("[Tool Handler] Executing 'get_network_activity' (sampling over 1 second)...\n")
		sent, recv, err := monitor.GetNetworkBandwidth()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get network stats: %v", err)), nil
		}
		report := fmt.Sprintf("Network Bandwidth:\nUpload Speed: %d bytes/sec\nDownload Speed: %d bytes/sec", sent, recv)
		return mcp.NewToolResultText(report), nil
	})
}