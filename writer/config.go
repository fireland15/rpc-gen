package writer

import (
	"encoding/json"
	"os"
)

type ClientConfig struct {
	Types map[string]string `json:"types"`
}

type WriterConfig struct {
	Clients map[string]ClientConfig `json:"clients"`
}

func ReadConfig(file string) (*WriterConfig, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := new(WriterConfig)
	err = json.Unmarshal(f, config)
	if err != nil {
		return nil, err
	}

	return config, err
}
