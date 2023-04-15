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

func (client *Client) models_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.models"))
}

func (client *Client) completions_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.completions"))
}

func (client *Client) img_gen_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.images_generations"))
}

func (client *Client) img_edits_url() string {
	return fmt.Sprintf("%s%s", client.config.Url, _Data.GetString("apis.images_edits"))
}
