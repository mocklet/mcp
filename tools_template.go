package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"mocklet-mcp/client"
)

func createTemplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")
	description := request.GetString("description", "")
	harPath, err := request.RequireString("har_file_path")
	if err != nil {
		return mcp.NewToolResultError("har_file_path must be a string"), nil
	}

	fields := map[string]string{}
	if name != "" {
		fields["name"] = name
	}
	if description != "" {
		fields["description"] = description
	}
	ttlVal := request.GetFloat("default_ttl_seconds", 0)
	if ttlVal > 0 {
		fields["default_ttl_seconds"] = strconv.Itoa(int(ttlVal))
	}

	contentType, body, err := createMultipartFileWithFields(harPath, fields)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to prepare request: %v", err)), nil
	}

	resp, err := apiClient.CreateTemplateFileWithBodyWithResponse(ctx, contentType, body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 201 && resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON201, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Template created successfully:\n```json\n%s\n```", string(b))), nil
}

func listTemplatesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	limit := int(request.GetFloat("limit", 5))
	
	params := &client.ListTemplatesParams{
		Limit: &limit,
	}

	resp, err := apiClient.ListTemplatesWithResponse(ctx, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON200, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Templates:\n```json\n%s\n```", string(b))), nil
}

func spawnMockHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	templateId, err := request.RequireString("template_public_id")
	if err != nil {
		return mcp.NewToolResultError("template_public_id must be a string"), nil
	}

	resp, err := apiClient.SpawnMockFromTemplateWithResponse(ctx, templateId, client.SpawnMockFromTemplateJSONRequestBody{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 201 && resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON201, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Mock spawned successfully:\n```json\n%s\n```", string(b))), nil
}

func uploadTemplateRevisionHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	templateId, err := request.RequireString("template_public_id")
	if err != nil {
		return mcp.NewToolResultError("template_public_id must be a string"), nil
	}
	harPath, err := request.RequireString("har_file_path")
	if err != nil {
		return mcp.NewToolResultError("har_file_path must be a string"), nil
	}

	contentType, body, err := createMultipartFileWithFields(harPath, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to prepare request: %v", err)), nil
	}

	resp, err := apiClient.UploadTemplateRevisionWithBodyWithResponse(ctx, templateId, contentType, body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON201, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Template updated successfully:\n```json\n%s\n```", string(b))), nil
}
