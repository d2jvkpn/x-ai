package chatgpt

import (
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

const (
	English = "English"
	Chinese = "Chinese"
)

var (
	//go:embed chatgpt.yaml
	_Yaml         string
	_Data         *viper.Viper
	_RE_ImageSize *regexp.Regexp
)

func init() {
	_Data = viper.New()
	_Data.SetConfigType("yaml")
	_ = _Data.ReadConfig(strings.NewReader(_Yaml))

	_RE_ImageSize = regexp.MustCompile(`^\d+x\d+$`)
}

func Version() string {
	return _Data.GetString("version")
}

func default_model() string {
	return _Data.GetString("defaults.model")
}

func default_temperature() float32 {
	return float32(_Data.GetFloat64("defaults.temperature"))
}

func chat_url() string {
	return fmt.Sprintf("%s%s", _Data.GetString("url"), _Data.GetString("apis.completions"))
}

func img_gen_url() string {
	return fmt.Sprintf("%s%s", _Data.GetString("url"), _Data.GetString("apis.images_generations"))
}

func img_edits_url() string {
	return fmt.Sprintf("%s%s", _Data.GetString("url"), _Data.GetString("apis.images_edits"))
}
