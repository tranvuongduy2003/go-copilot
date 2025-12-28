package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/tranvuongduy2003/go-copilot/pkg/logger"
)

type DocsHandler struct {
	logger      logger.Logger
	openAPIPath string
}

func NewDocsHandler(logger logger.Logger) *DocsHandler {
	workingDirectory, _ := os.Getwd()
	openAPIPath := filepath.Join(workingDirectory, "docs", "openapi.yaml")

	return &DocsHandler{
		logger:      logger,
		openAPIPath: openAPIPath,
	}
}

func (handler *DocsHandler) SwaggerUI(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(swaggerUIHTML))
}

func (handler *DocsHandler) OpenAPISpec(writer http.ResponseWriter, request *http.Request) {
	specContent, err := os.ReadFile(handler.openAPIPath)
	if err != nil {
		handler.logger.Error("failed to read openapi spec", logger.Err(err), logger.String("path", handler.openAPIPath))
		http.Error(writer, "OpenAPI spec not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/yaml")
	writer.Header().Set("Cache-Control", "public, max-age=3600")
	writer.WriteHeader(http.StatusOK)
	writer.Write(specContent)
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation - Go Copilot</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
        .swagger-ui .topbar { display: none; }
        .swagger-ui .info { margin: 20px 0; }
        .swagger-ui .info .title { font-size: 2em; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            SwaggerUIBundle({
                url: "/docs/openapi.yaml",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                validatorUrl: null,
                supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
                defaultModelsExpandDepth: 1,
                defaultModelExpandDepth: 1,
                displayRequestDuration: true,
                filter: true,
                showExtensions: true,
                showCommonExtensions: true
            });
        };
    </script>
</body>
</html>`
