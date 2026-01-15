package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the downhub configuration structure
type Config struct {
	Defaults    Defaults     `yaml:"defaults"`
	Repositories []Repository `yaml:"repositories"`
	FileFilters FileFilters  `yaml:"file_filters"`
	Download    Download     `yaml:"download"`
	Logging     Logging      `yaml:"logging"`
	Advanced    Advanced     `yaml:"advanced"`
}

// Defaults contains default configuration values
type Defaults struct {
	BaseDataDir string `yaml:"base_data_dir"`
	DocsDir     string `yaml:"docs_dir"`
	SourceDir   string `yaml:"source_dir"`
	DocsPath    string `yaml:"docs_path"`
	MaxConcurrentDownloads int    `yaml:"max_concurrent_downloads"`
	Proxy                 string `yaml:"proxy"`
}

// Repository represents a GitHub repository configuration
type Repository struct {
	Name           string `yaml:"name"`
	URL            string `yaml:"url"`
	DownloadDocs   bool   `yaml:"download_docs"`
	DownloadSource bool   `yaml:"download_source"`
	OutputDir      string `yaml:"output_dir"`
	DocsPath       string `yaml:"docs_path"`
}

// FileFilters contains file inclusion and exclusion patterns
type FileFilters struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

// Download contains download-related settings
type Download struct {
	Timeout     int `yaml:"timeout"`
	Retries     int `yaml:"retries"`
	RetryDelay  int `yaml:"retry_delay"`
	UserAgent   string `yaml:"user_agent"`
}

// Logging contains logging configuration
type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// Advanced contains advanced configuration options
type Advanced struct {
	PreserveStructure  bool `yaml:"preserve_structure"`
	CreateReadme       bool `yaml:"create_readme"`
	ValidateChecksums  bool `yaml:"validate_checksums"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to a YAML file
func SaveConfig(config *Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Defaults: Defaults{
			BaseDataDir: "data",
			DocsDir:     "docs",
			SourceDir:   "source",
			DocsPath:    "docs",
			MaxConcurrentDownloads: 5,
			Proxy:                 "http://localhost:7890",
		},
		FileFilters: FileFilters{
			Include: []string{"*.md", "*.txt", "*.yaml", "*.yml"},
			Exclude: []string{"node_modules/*", ".git/*", "vendor/*"},
		},
		Download: Download{
			Timeout:    300,
			Retries:    3,
			RetryDelay: 5,
			UserAgent:  "Downhub/1.0",
		},
		Logging: Logging{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Advanced: Advanced{
			PreserveStructure: true,
			CreateReadme:      true,
			ValidateChecksums: false,
		},
	}
}
