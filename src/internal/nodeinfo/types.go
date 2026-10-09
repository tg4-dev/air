package nodeinfo

import "github.com/tg4-dev/air/src/internal/engines"

const ProtocolVersion = 1

type NodeStatus struct {
	ProtocolVersion int                    `json:"protocol_version"`
	NodeID          string                 `json:"node_id"`
	Hostname        string                 `json:"hostname"`
	UptimeSeconds   int64                  `json:"uptime_seconds"`
	GPU             GPUInfo                `json:"gpu"`
	Engines         []engines.EngineStatus `json:"engines"`
}
