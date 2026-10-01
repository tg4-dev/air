package nodeinfo

import (
	"context"
	"strings"
	"time"
)

type StatusProvider interface {
	CurrentStatus(ctx context.Context) NodeStatus
}

type defaultProvider struct {
	nodeID    string
	hostname  string
	startedAt time.Time
}

func NewDefaultProvider(nodeID, hostname string, startedAt time.Time) StatusProvider {
	return &defaultProvider{
		nodeID:    nodeID,
		hostname:  hostname,
		startedAt: startedAt,
	}
}

func (p *defaultProvider) CurrentStatus(ctx context.Context) NodeStatus {
	gpuName := strings.TrimSpace(DetectGPUName())
	if gpuName == "" {
		gpuName = "unknown"
	}

	ollamaCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	engine := probeOllama(ollamaCtx)

	return NodeStatus{
		ProtocolVersion: ProtocolVersion,
		NodeID:          p.nodeID,
		Hostname:        p.hostname,
		UptimeSeconds:   int64(time.Since(p.startedAt).Seconds()),
		GPU:             GPUInfo{Name: gpuName},
		Engines:         []EngineStatus{engine},
	}
}
