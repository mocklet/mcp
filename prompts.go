package main

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func registerPrompts() {
	mcpServer.AddPrompt(mcp.Prompt{
		Name:        "spawn_dependency_mock",
		Description: "Spawn a mock for a dependency",
	}, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: "Find a Mocklet template for the `PaymentGateway` service, spawn an ephemeral mock, and configure my local `.env` to point to the new mock URL.",
					},
				},
			},
		}, nil
	})

	mcpServer.AddPrompt(mcp.Prompt{
		Name:        "debug_mock_usage",
		Description: "Debug 404s and missing routes for a mock",
	}, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: "Fetch the statistics for mock `mock_123` and analyze why my recent frontend requests resulted in 404s (check the missed routes and ignore rules).",
					},
				},
			},
		}, nil
	})

	mcpServer.AddPrompt(mcp.Prompt{
		Name:        "generate_frontend_with_mock",
		Description: "Build a frontend prototype using a mock backend",
	}, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: "I need a React dashboard for user analytics. Please generate a Mocklet template containing realistic user analytics data, spawn a mock, and build the React components that fetch and display this data from the mock URL.",
					},
				},
			},
		}, nil
	})

	mcpServer.AddPrompt(mcp.Prompt{
		Name:        "integration_testing_setup",
		Description: "Setup isolated mocks for integration testing",
	}, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: "We are writing integration tests for the checkout flow. Save the attached HAR data to a local file, upload it as a new Mocklet template named 'Checkout Flow', spawn a mock, and write a Cypress test script that uses this mock as the backend API.",
					},
				},
			},
		}, nil
	})

	mcpServer.AddPrompt(mcp.Prompt{
		Name:        "har_validation_cleanup",
		Description: "Validate and clean up HAR payload before mocking",
	}, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Messages: []mcp.PromptMessage{
				{
					Role: "user",
					Content: mcp.TextContent{
						Type: "text",
						Text: "Save the HAR data in my clipboard to a local file and validate it using Mocklet. If it's valid, create a template from it; if there are errors, explain what I need to fix.",
					},
				},
			},
		}, nil
	})
}
