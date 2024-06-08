package main

import (
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type RSS struct {
	URL         string `yaml:"url"`
	LastUpdated int64  `yaml:"lastUpdated"`
}

type Crawler struct {
	ID          uint64 `yaml:"id"`
	Name        string `yaml:"name"`
	Initialized bool   `yaml:"initialized"`
	RSS         RSS    `yaml:"rss"`
}

type Config struct {
	Crawlers []Crawler `yaml:"crawlers"`
}

func getConfig(filePath string) Config {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	var data []byte
	data, err = io.ReadAll(file)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}
	return config
}

func getCrawlerIdx(config Config, id uint64) int {
	for i := range config.Crawlers {
		if config.Crawlers[i].ID == id {
			return i
		}
	}
	return -1
}

func updateConfig(config *Config, id uint64) {
	i := getCrawlerIdx(*config, id)
	if i == -1 {
		log.Fatalf("ID [%d] not exists", id)
	}
	config.Crawlers[i].RSS.LastUpdated = time.Now().Unix()
}
func writeConfig(filePath string, config Config) {
	updatedData, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error marshalling YAML: %v", err)
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("error opening file for writing: %v", err)
	}
	defer file.Close()

	_, err = file.Write(updatedData)
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}
