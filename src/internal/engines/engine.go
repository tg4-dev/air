package engines

type Engine interface {
	TestRequest()
}

func WorkingEngine(engine_name string) Engine {
	switch engine_name {
	case "ollama":
		return &ollamaEngine{} // Fill with config-like json TODO
	case "llamacpp":
		return &llamacppEngine{}
	}
	return &ollamaEngine{} // TODO
}
