package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ClientConfig struct {
	Types  map[string]string `json:"types"`
	Output string
}

type ServerConfig struct {
	Types  map[string]string `json:"types"`
	Output string
}

type RpcGenConfig struct {
	RpcDefinitionFile string                     `json:"definition"`
	Clients           map[string]ClientConfig    `json:"clients"`
	Servers           map[string]json.RawMessage `json:"servers"`
}

func ReadConfig(file string) (*RpcGenConfig, error) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	f, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	config := new(RpcGenConfig)
	err = json.Unmarshal(f, config)
	if err != nil {
		return nil, err
	}

	return config, err
}
