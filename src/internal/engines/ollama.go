package engines

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// TODO fetch port from config
const ollamaBaseURL = "http://127.0.0.1:11434"

type ollamaTagsResponse struct {
	Models []ollamaModel `json:"models"`
}

type ollamaModel struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func ProbeOllama(ctx context.Context) EngineStatus {
	status := EngineStatus{
		Name:      "ollama",
		Port:      "11434", // TODO
		Available: false,
		Models:    []ModelInfo{},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ollamaBaseURL+"/api/tags", nil)
	if err != nil {
		return status
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return status
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return status
	}

	var body ollamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return status
	}

	status.Available = true
	for _, m := range body.Models {
		status.Models = append(status.Models, ModelInfo{
			Name:      m.Name,
			SizeBytes: m.Size,
		})
	}
	if status.Models == nil {
		status.Models = []ModelInfo{}
	}
	return status
}

type ollamaEngine struct {
	port string
}

func (e *ollamaEngine) TestRequest() {}
