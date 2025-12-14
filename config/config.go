package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/ulib/commonconfig"
	"gopkg.in/yaml.v3"
)

type ResourceConfig struct {
	Name string `json:"name" yaml:"name" toml:"name"`
	Dir  string `json:"dir" yaml:"dir" toml:"dir"`
}

type AppConfig struct {
	commonconfig.ServerConfig `json:"server" yaml:"server" toml:"server"`
	ReportConfig              `json:"report" yaml:"report" toml:"report"`
	Feeds                     []model.PipelineConfigure `json:"feeds" yaml:"feeds" toml:"feeds"`
	Related                   []model.PipelineConfigure `json:"related" yaml:"related" toml:"related"`
	Indexes                   []ResourceConfig          `json:"indexes" yaml:"indexes" toml:"indexes"`
	Items                     ResourceConfig            `json:"items" yaml:"items" toml:"items"`
}

type SegmentConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint" toml:"endpoint"`
	WriteKey string `json:"write_key" yaml:"write_key" toml:"write_key"`
}

type SLSLogConfig struct {
	Endpoint string `json:"endpoint" yaml:"endpoint" toml:"endpoint"`
	AK       string `json:"ak" yaml:"ak" toml:"ak"`
	SK       string `json:"sk" yaml:"sk" toml:"sk"`
	RAM      string `json:"ram" yaml:"ram" toml:"ram"`
	Project  string `json:"project" yaml:"project" toml:"project"`
	LogStore string `json:"logstore" yaml:"logstore" toml:"logstore"`
	Region   string `json:"region" yaml:"region" toml:"region"`
}

type ReportConfig struct {
	Type          string `json:"type" yaml:"type" toml:"type"`
	SegmentConfig `json:"segment" yaml:"segment" toml:"segment"`
	SLSLogConfig  `json:"slslog" yaml:"slslog" toml:"slslog"`
}

// Global AppConfig instance
var AppConfigInstance *AppConfig

// Init loads configuration from a JSON, YAML, or TOML file.
// Returns an error instead of panic, allowing caller to decide handling.
func (conf *AppConfig) Init(filePath string) error {
	fData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("os.ReadFile error: %w", err)
	}

	ext := filepath.Ext(filePath)
	switch ext {
	case ".json":
		if err := json.Unmarshal(fData, conf); err != nil {
			return fmt.Errorf("json.Unmarshal error: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(fData, conf); err != nil {
			return fmt.Errorf("yaml.Unmarshal error: %w", err)
		}
	case ".toml":
		if err := toml.Unmarshal(fData, conf); err != nil {
			return fmt.Errorf("toml.Unmarshal error: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	fmt.Printf("InitAppConfig: %+v\n", conf)
	return nil
}
