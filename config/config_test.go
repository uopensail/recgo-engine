package config

import (
	"fmt"
	"testing"
)

func TestAppConfig_Init(t *testing.T) {
	conf := AppConfig{}
	conf.Init("./../conf/local/config.toml")
	fmt.Println(conf)
}
