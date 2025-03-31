package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/uopensail/kongming-sdk-go/sdkcore"

	"github.com/uopensail/ulib/commonconfig"
)

type EngineDataConfig struct {
	URL string `toml:"url" yaml:"url"`
}
type GrowthBookSDKConfig struct {
	APIHost   string `toml:"api_host" yaml:"api_host"`
	ClientKey string `toml:"client_key" yaml:"client_key"`
}

type ABConfig struct {
	Type                      string `toml:"type"  yaml:"type"`
	sdkcore.KongMingSDKConfig `toml:"kongming"  yaml:"kongming"`
	GrowthBookSDKConfig       `toml:"growthbook" yaml:"growthbook"`
}

type AppConfig struct {
	EngineDataConfig          `toml:"engine_data" yaml:"engine_data"`
	commonconfig.ServerConfig `toml:"server" yaml:"server"`
	EnvConfig                 `toml:"env" yaml:"env"`
	ABConfig                  `toml:"ab" yaml:"ab"`
	ReportConfig              `toml:"report" yaml:"report"`
}
type SegmentConfig struct {
	Endpoint string `json:"endpoint" toml:"endpoint"  yaml:"endpoint"`
	WriteKey string `json:"write_key" toml:"write_key" yaml:"write_key"`
}

type SLSLogConfig struct {
	Endpoint string `json:"endpoint" toml:"endpoint"  yaml:"endpoint"`
	AK       string `json:"ak" toml:"ak" yaml:"ak"`
	SK       string `json:"sk" toml:"sk" yaml:"sk"`
	RAM      string `json:"ram" toml:"ram" yaml:"ram"`
	Project  string `json:"project" toml:"project" yaml:"project"`
	LogStore string `json:"logstore" toml:"logstore" yaml:"logstore"`
	Region   string `json:"region" toml:"region" yaml:"region"`
}
type ReportConfig struct {
	Type          string `json:"type" toml:"type" yaml:"type"`
	SegmentConfig `json:"segment" toml:"segment" yaml:"segment"`
	SLSLogConfig  `json:"slslog" toml:"slslog" yaml:"slslog"`
}

type EnvConfig struct {
	Finder  commonconfig.FinderConfig `json:"finder" toml:"finder" yaml:"finder"`
	WorkDir string                    `json:"work_dir" toml:"work_dir" yaml:"work_dir"`
}

var AppConfigInstance AppConfig

func (conf *AppConfig) Init(filePath string) {
	fData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("ioutil.ReadFile error: %s", err)
		panic(err)
	}
	_, err = toml.Decode(string(fData), conf)
	if err != nil {
		fmt.Errorf("Unmarshal error: %s", err)
	} else {
		return
	}

	err = yaml.NewDecoder(bytes.NewReader(fData)).Decode(conf)
	if err != nil {
		fmt.Errorf("Unmarshal error: %s", err)
	} else {
		return
	}

	err = json.Unmarshal(fData, conf)
	if err != nil {
		fmt.Errorf("Unmarshal error: %s", err)
	} else {
		return
	}
	panic(err)
	fmt.Printf("InitAppConfig:%v yaml:%s\n", conf, string(fData))
}
