// Package main provides a simple example of an MCP server with a tool
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/pkg/mcp/server"
	"github.com/modelcontextprotocol/go-sdk/pkg/mcp/shared"
	"github.com/modelcontextprotocol/go-sdk/pkg/mcp/types"
)

func main() {
	// Create a FastMCP server
	mcpServer := server.NewFastMCPServer(handleSession)

	// Start the server
	fmt.Println("Starting MCP server on :8080...")
	log.Fatal(mcpServer.ListenAndServe(":8080"))
}

// handleSession handles an MCP session
func handleSession(ctx context.Context, session *server.Session) error {
	// Set server capabilities
	session.SetCapabilities(types.ServerCapabilities{
		Tools: &types.ToolsCapability{
			ListChanged: boolPtr(true),
		},
	})

	// Set server info
	session.SetServerInfo("simple-tool-server", "0.1.0")

	// Register handlers
	session.RegisterRequestHandler("tools/list", handleListTools)
	session.RegisterRequestHandler("tools/call", handleCallTool)

	// Wait for context to be done
	<-ctx.Done()
	return nil
}

// handleListTools handles a tools/list request
func handleListTools(ctx context.Context, reqCtx *shared.RequestContext, params interface{}) (interface{}, error) {
	// Create a list of tools
	tools := []types.Tool{
		{
			Name:        "echo",
			Description: "Echoes back the input message",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"message": map[string]interface{}{
						"type":        "string",
						"description": "The message to echo back",
					},
				},
				"required": []string{"message"},
			},
		},
		{
			Name:        "timestamp",
			Description: "Returns the current server timestamp",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	// Return the result
	return types.ListToolsResult{
		Tools: tools,
	}, nil
}

// handleCallTool handles a tools/call request
func handleCallTool(ctx context.Context, reqCtx *shared.RequestContext, params interface{}) (interface{}, error) {
	// Parse params
	var callParams types.CallToolRequestParams
	if err := parseParams(params, &callParams); err != nil {
		return nil, &types.ErrorData{
			Code:    types.INVALID_PARAMS,
			Message: fmt.Sprintf("Invalid parameters: %v", err),
		}
	}

	// Handle different tools
	switch callParams.Name {
	case "echo":
		// Get message from arguments
		message, ok := callParams.Arguments["message"].(string)
		if !ok {
			return nil, &types.ErrorData{
				Code:    types.INVALID_PARAMS,
				Message: "Missing or invalid 'message' parameter",
			}
		}

		// Return the echo result
		return types.CallToolResult{
			Result: map[string]interface{}{
				"message": message,
			},
		}, nil

	case "timestamp":
		// Return the current timestamp
		return types.CallToolResult{
			Result: map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
			},
		}, nil

	default:
		return nil, &types.ErrorData{
			Code:    types.METHOD_NOT_FOUND,
			Message: fmt.Sprintf("Tool not found: %s", callParams.Name),
		}
	}
}

// parseParams parses JSON-RPC params into a Go struct
func parseParams(params interface{}, dest interface{}) error {
	// Convert params to JSON
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	// Unmarshal into destination
	if err := json.Unmarshal(paramsJSON, dest); err != nil {
		return fmt.Errorf("failed to unmarshal params: %w", err)
	}

	return nil
}

// boolPtr returns a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}
