package nodeinfo

import (
	"context"
	"strings"
	"time"

	"github.com/tg4-dev/air/src/internal/engines"
)

type StatusProvider interface {
	CurrentStatus(ctx context.Context) NodeStatus
}

type statusProvider struct {
	NodeID    string
	Hostname  string
	StartedAt time.Time
}

func NewStatusProvider(nodeID string, hostname string, startedAt time.Time) StatusProvider {
	return &statusProvider{
		NodeID:    nodeID,
		Hostname:  hostname,
		StartedAt: startedAt,
	}
}

func (p *statusProvider) CurrentStatus(ctx context.Context) NodeStatus {
	gpuName := strings.TrimSpace(DetectGPUName())
	if gpuName == "" {
		gpuName = "unknown"
	}

	var engineList []engines.EngineStatus

	ollamaCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	engineList = append(engineList, engines.ProbeOllama(ollamaCtx))

	// llamacpp probe and append

	// other engine probe and append...

	return NodeStatus{
		ProtocolVersion: ProtocolVersion,
		NodeID:          p.NodeID,
		Hostname:        p.Hostname,
		UptimeSeconds:   int64(time.Since(p.StartedAt).Seconds()),
		GPU:             GPUInfo{Name: gpuName},
		Engines:         engineList,
	}
}
