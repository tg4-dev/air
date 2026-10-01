package nodeinfo

const ProtocolVersion = 1

type NodeStatus struct {
	ProtocolVersion int            `json:"protocol_version"`
	NodeID          string         `json:"node_id"`
	Hostname        string         `json:"hostname"`
	UptimeSeconds   int64          `json:"uptime_seconds"`
	GPU             GPUInfo        `json:"gpu"`
	Engines         []EngineStatus `json:"engines"`
}

type EngineStatus struct {
	Name      string      `json:"name"`
	Available bool        `json:"available"`
	Models    []ModelInfo `json:"models"`
}

type ModelInfo struct {
	Name      string `json:"name"`
	SizeBytes int64  `json:"size_bytes"`
}
