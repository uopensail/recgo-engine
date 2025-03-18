package resources

import (
	"fmt"
	"testing"

	"github.com/uopensail/recgo-engine/config"
	"github.com/uopensail/ulib/commonconfig"
)

func Test_ResourceLoad(t *testing.T) {
	ress, err := NewResource(config.EnvConfig{
		Finder: commonconfig.FinderConfig{
			Type: "local",
		},
		WorkDir: "/tmp",
	}, "../test/test_data")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(ress.Array))
}
