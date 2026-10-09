package engines

type EngineStatus struct {
	Name      string      `json:"name"`
	Port      string      `json:"port"` // TODO should be fetched from config
	Available bool        `json:"available"`
	IsDefault bool        `json:"is_default"`
	Models    []ModelInfo `json:"models"`
}

type ModelInfo struct {
	Name      string `json:"name"`
	SizeBytes int64  `json:"size_bytes"`
}
