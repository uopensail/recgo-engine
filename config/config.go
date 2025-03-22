package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/uopensail/kongming-sdk-go/sdkcore"
	"github.com/uopensail/ulib/commonconfig"
)

type DBModelConfig struct {
	URL string `toml:"url"`
}

type AppConfig struct {
	DBModelConfig             `toml:"dbmodel"`
	commonconfig.ServerConfig `toml:"server"`
	EnvConfig                 `toml:"env"`

	sdkcore.KongMingSDKConfig `toml:"absdk"`
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
