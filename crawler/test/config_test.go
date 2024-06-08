package test

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

type RSS struct {
	URL         string `yaml:"url"`
	LastUpdated int    `yaml:"lastUpdated"`
}

type Crawler struct {
	ID          int    `yaml:"id"`
	Name        string `yaml:"name"`
	Initialized bool   `yaml:"initialized"`
	RSS         RSS    `yaml:"rss"`
}

type Config struct {
	Crawlers []Crawler `yaml:"crawlers"`
}

// time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC).Unix(),
// // lastUpdated: time.Now().Unix()

func GetConfig(filePath string) Config {
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

func UpdateConfig(config *Config) {
	// Update the lastUpdated value for the crawler with name "lineCrawler"
	for i := range config.Crawlers {
		if config.Crawlers[i].Name == "lineCrawler" {
			config.Crawlers[i].RSS.LastUpdated = int(time.Now().Unix())
		}
	}
}
func WriteConfig(filePath string, config Config) {
	// Marshal the updated config back to YAML
	updatedData, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error marshalling YAML: %v", err)
	}

	// Open the YAML file for writing
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("error opening file for writing: %v", err)
	}
	defer file.Close()

	// Write the updated YAML back to the file
	_, err = file.Write(updatedData)
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}

func TestReadUpdateWriteConfig(t *testing.T) {
	filePath := "../config-crawler.yaml"
	config := GetConfig(filePath)
	UpdateConfig(&config)
	WriteConfig(filePath, config)
}
