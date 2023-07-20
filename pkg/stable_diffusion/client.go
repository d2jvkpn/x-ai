package stable_diffusion

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

type Config struct {
	Url        string `mapstructure:"url"`
	Controlnet struct {
		Module string   `mapstructure:"module"`
		Models []string `mapstructure:"models"`
	} `mapstructure:"controlnet"`
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

	if config.Controlnet.Module == "" {
		return nil, fmt.Errorf("controlnet.module is unset")
	}

	if len(config.Controlnet.Models) == 0 {
		return nil, fmt.Errorf("controlnet.models is empty")
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

func (client *Client) ValidateControlnet(item *Controlnet) (err error) {
	cn := &client.config.Controlnet
	if len(item.InputImage) == 0 {
		return fmt.Errorf("has no input_image")
	}

	if item.Module == "" {
		item.Module = cn.Module
	}

	if item.Model == "" {
		return fmt.Errorf("model is unset")
	}

	match := false
	for i := range cn.Models {
		if cn.Models[i] == item.Model {
			match = true
			break
		}
	}
	if !match {
		return fmt.Errorf("model is invalid")
	}

	return nil
}
