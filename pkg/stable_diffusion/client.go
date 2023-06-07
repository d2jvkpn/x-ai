package stable_diffusion

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type Config struct {
	Url        string     `mapstructure:"url"`
	Controlnet Controlnet `mapstructure:"controlnet"`
}

type Client struct {
	config *Config
	cli    *http.Client
}

func NewConfg(fp, key string) (config *Config, err error) {
	vp := viper.New()
	vp.SetConfigType("yaml")

	vp.SetConfigFile(fp)
	if err = vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ReadInConfig(): %q, %w", fp, err)
	}

	config = new(Config)
	if err = vp.UnmarshalKey(key, config); err != nil {
		return nil, err
	}

	return config, nil
}

func NewClient(fp, key string) (client *Client, err error) {
	var config *Config

	if config, err = NewConfg(fp, key); err != nil {
		return nil, err
	}

	client = &Client{config: config, cli: new(http.Client)}

	return client, nil
}

func ClientFromViper(vp *viper.Viper, field string) (client *Client, err error) {
	var config Config

	if err = vp.UnmarshalKey(field, &config); err != nil {
		return nil, err
	}

	return &Client{config: &config, cli: new(http.Client)}, nil
}
