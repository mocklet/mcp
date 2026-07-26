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

func getTemplateOpenApiHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	templateID, err := request.RequireString("template_public_id")
	if err != nil {
		return mcp.NewToolResultError("template_public_id must be a string"), nil
	}
	format := request.GetString("format", "yaml")

	params := &client.GetTemplateOpenAPIParams{
		Format: &format,
	}

	resp, err := apiClient.GetTemplateOpenAPIWithResponse(ctx, templateID, params)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	contentType := resp.ContentType()
	if contentType == "" {
		contentType = "text/yaml"
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("OpenAPI Spec:\n```%s\n%s\n```", format, string(resp.Body))), nil
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

	resp, err := apiClient.UploadTemplateRevisionFileWithBodyWithResponse(ctx, templateId, contentType, body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 201 && resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON201, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Template updated successfully:\n```json\n%s\n```", string(b))), nil
}

func updateTemplateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	templateID, err := request.RequireString("template_public_id")
	if err != nil {
		return mcp.NewToolResultError("template_public_id must be a string"), nil
	}

	reqBody := client.TemplateUpdateRequest{}

	name := request.GetString("name", "")
	if name != "" {
		reqBody.Name = &name
	}
	desc := request.GetString("description", "")
	if desc != "" {
		reqBody.Description = &desc
	}

	resp, err := apiClient.UpdateTemplateWithResponse(ctx, templateID, reqBody)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON200, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Template metadata updated:\n```json\n%s\n```", string(b))), nil
}

func downloadTemplateHarHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	templateID, err := request.RequireString("template_public_id")
	if err != nil {
		return mcp.NewToolResultError("template_public_id must be a string"), nil
	}

	resp, err := apiClient.DownloadTemplateHarWithResponse(ctx, templateID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Request failed: %v", err)), nil
	}

	if resp.StatusCode() != 200 {
		return formatError(resp.StatusCode(), resp.Body), nil
	}

	b, _ := json.MarshalIndent(resp.JSON200, "", "  ")
	return mcp.NewToolResultText(fmt.Sprintf("Template HAR:\n```json\n%s\n```", string(b))), nil
}
