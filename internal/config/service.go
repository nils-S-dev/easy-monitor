package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Cron     string    `json:"cron"`
	Monitors []Monitor `json:"monitors"`
	Notify   []string  `json:"notify"`
}

type Monitor struct {
	Name     string          `json:"name"`
	Endpoint string          `json:"endpoint"`
	Method   string          `json:"method"`
	Body     string          `json:"body"`
	Cron     string          `json:"cron"`
	Notify   []string        `json:"notify"`
	Expect   MonitorExpected `json:"expect"`
}

type MonitorExpected struct {
	Status  int               `json:"status,omitempty"`
	Body    string            `json:"body,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

func GetConfig() *Config {
	// Open the JSON file
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode into struct
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}

	return &config
}
