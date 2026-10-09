package nodeinfo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tg4-dev/air/src/internal/engines"
)

type fakeProvider struct {
	status NodeStatus
}

func (f *fakeProvider) CurrentStatus(_ context.Context) NodeStatus {
	return f.status
}

func TestHandleHealth(t *testing.T) {
	s := NewServer(&fakeProvider{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	s.handleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want 200", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status = %q, want ok", body["status"])
	}
}

func TestHandleStatus(t *testing.T) {
	want := NodeStatus{
		ProtocolVersion: 1,
		NodeID:          "test-id",
		Hostname:        "test-host",
		UptimeSeconds:   42,
		GPU:             GPUInfo{Name: "unknown"},
		Engines: []engines.EngineStatus{
			{Name: "ollama", Port: "8080", Available: false, IsDefault: true, Models: []engines.ModelInfo{}},
		},
	}
	s := NewServer(&fakeProvider{status: want})
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()

	s.handleStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want 200", rec.Code)
	}
	var got NodeStatus
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.ProtocolVersion != want.ProtocolVersion {
		t.Errorf("protocol_version = %d, want %d", got.ProtocolVersion, want.ProtocolVersion)
	}
	if got.NodeID != want.NodeID {
		t.Errorf("node_id = %q, want %q", got.NodeID, want.NodeID)
	}
	if got.GPU.Name != "unknown" {
		t.Errorf("gpu.name = %q, want unknown", got.GPU.Name)
	}
	if len(got.Engines) != 1 || got.Engines[0].Available {
		t.Errorf("engines = %+v, want ollama unavailable", got.Engines)
	}
}
