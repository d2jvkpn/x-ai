package stable_diffusion

import (
	_ "embed"
	"strings"

	"github.com/spf13/viper"
)

var (
	//go:embed stable_diffusion.yaml
	_Yaml string
	_Data *viper.Viper
)

func init() {
	_Data = viper.New()
	_Data.SetConfigType("yaml")
	_ = _Data.ReadConfig(strings.NewReader(_Yaml))
}

func Version() string {
	return _Data.GetString("version")
}
