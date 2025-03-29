package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/uopensail/kongming-sdk-go/sdkcore"
	"github.com/uopensail/ulib/commonconfig"
)

type EngineDataConfig struct {
	URL string `toml:"url"`
}

type GrowthBookSDKConfig struct {
	APIHost   string `toml:"api_host"`
	ClientKey string `toml:"client_key"`
}

type ABConfig struct {
	Type                      string `toml:"type"`
	sdkcore.KongMingSDKConfig `toml:"kongming"`
	GrowthBookSDKConfig       `toml:"growthbook"`
}

type AppConfig struct {
	EngineDataConfig          `toml:"engine_data"`
	commonconfig.ServerConfig `toml:"server"`
	EnvConfig                 `toml:"env"`
	ABConfig                  `toml:"ab"`
}
type EnvConfig struct {
	Finder  commonconfig.FinderConfig `json:"finder" toml:"finder"`
	WorkDir string                    `json:"work_dir" toml:"work_dir"`
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
		panic(err)
	}
	fmt.Printf("InitAppConfig:%v yaml:%s\n", conf, string(fData))
}
