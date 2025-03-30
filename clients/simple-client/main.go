// Package mcp provides a simple example of using the MCP client
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/modelcontextprotocol/go-sdk/pkg/mcp/client"
	"github.com/modelcontextprotocol/go-sdk/pkg/mcp/types"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a FastMCP client
	client, err := client.NewFastMCPClient(ctx, "ws://localhost:8080/mcp", 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Initialize the client
	initResult, err := client.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	fmt.Printf("Connected to server: %s %s\n", initResult.ServerInfo.Name, initResult.ServerInfo.Version)
	fmt.Printf("Protocol version: %s\n", initResult.ProtocolVersion)

	// List available tools
	toolsResult, err := client.ListTools()
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("Available tools:\n")
	for _, tool := range toolsResult.Tools {
		fmt.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	// Call a tool if available
	if len(toolsResult.Tools) > 0 {
		tool := toolsResult.Tools[0]
		fmt.Printf("Calling tool: %s\n", tool.Name)

		// Prepare arguments based on the tool's input schema
		arguments := make(map[string]interface{})

		// This is a simple example; in a real application, you would
		// populate arguments based on the tool's input schema

		result, err := client.CallTool(tool.Name, arguments)
		if err != nil {
			log.Fatalf("Failed to call tool: %v", err)
		}

		fmt.Printf("Tool result: %v\n", result.Result)
	}
}
