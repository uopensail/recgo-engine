package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/uopensail/recgo-engine/model"
	"github.com/uopensail/ulib/commonconfig"
)

type ResourceConfig struct {
	Name string `json:"name"`
	Dir  string `json:"dir"`
}

type AppConfig struct {
	commonconfig.ServerConfig `json:"server"`
	ReportConfig              `json:"report"`
	Feeds                     []model.PipelineConfigure `json:"feeds"`
	Related                   []model.PipelineConfigure `json:"related"`
	Indexes                   []ResourceConfig          `json:"indexes"`
	Items                     ResourceConfig            `json:"items"`
}

type SegmentConfig struct {
	Endpoint string `json:"endpoint" toml:"endpoint"  yaml:"endpoint"`
	WriteKey string `json:"write_key" toml:"write_key" yaml:"write_key"`
}

type SLSLogConfig struct {
	Endpoint string `json:"endpoint"`
	AK       string `json:"ak"`
	SK       string `json:"sk"`
	RAM      string `json:"ram"`
	Project  string `json:"project"`
	LogStore string `json:"logstore"`
	Region   string `json:"region"`
}

type ReportConfig struct {
	Type          string `json:"type"`
	SegmentConfig `json:"segment"`
	SLSLogConfig  `json:"slslog"`
}

var AppConfigInstance AppConfig

func (conf *AppConfig) Init(filePath string) {
	fData, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("os.ReadFile error: %s", err))
	}

	err = json.Unmarshal(fData, conf)
	if err != nil {
		panic(fmt.Errorf("Unmarshal error: %s", err))
	}
	fmt.Printf("InitAppConfig:%v yaml:%s\n", conf, string(fData))

}
