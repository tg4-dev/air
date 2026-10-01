package discovery

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Endpoint struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Addrs []net.IP          `json:"addrs"`
	Port  int               `json:"port"`
	Meta  map[string]string `json:"meta"`
}

type endpointConfig struct {
	ID string `json:"id"`
}

func NewEndpoint(name string, addrs []net.IP, port int, meta map[string]string) (*Endpoint, error) {
	id, err := loadOrCreateID()
	if err != nil {
		return nil, err
	}
	return &Endpoint{
		ID:    id,
		Name:  name,
		Addrs: addrs,
		Port:  port,
		Meta:  meta,
	}, nil
}

func loadOrCreateID() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(home, ".air")
	path := filepath.Join(dir, "config.json")

	data, err := os.ReadFile(path)
	if err == nil {
		var cfg endpointConfig
		if json.Unmarshal(data, &cfg) == nil {
			if _, err := uuid.Parse(cfg.ID); err == nil {
				return cfg.ID, nil
			}
		}
	}

	id := uuid.New().String()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}

	cfg := endpointConfig{ID: id}
	data, err = json.Marshal(cfg)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return "", err
	}
	return id, nil
}
