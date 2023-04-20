package chatgpt

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
	"golang.org/x/net/proxy"
)

type Config struct {
	Url     string `mapstructure:"url"`
	API_Key string `mapstructure:"api_key"`
	ORG_ID  string `mapstructure:"org_id"`

	Proxy         string `mapstructure:"proxy"`
	TlsSkipVerify bool   `mapstructure:"tls_skip_verify"`
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
	if config.Url == "" {
		config.Url = _Data.GetString("url")
	}

	return config, nil
}

func NewClient(fp, key string) (client *Client, err error) {
	var (
		cli      *http.Client
		config   *Config
		dialer   proxy.Dialer
		proxyURL *url.URL
	)

	if config, err = NewConfg(fp, key); err != nil {
		return nil, err
	}

	cli = new(http.Client)
	if config.TlsSkipVerify {
		cli.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true},
		}
	}

	client = &Client{config: config, cli: cli}

	if config.Proxy == "" {
		return client, nil
	}

	if proxyURL, err = url.Parse(config.Proxy); err != nil {
		return nil, err
	}

	switch proxyURL.Scheme {
	case "socks5":
		if dialer, err = proxy.SOCKS5("tcp", proxyURL.Host, nil, nil); err != nil {
			return nil, err
		}
		client.cli.Transport = &http.Transport{Dial: dialer.Dial}
	case "http", "https":
		client.cli.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	default:
		return nil, fmt.Errorf("unknow proxy schema")
	}

	return
}

func (client *Client) setAuth(request *http.Request, isJson bool) {
	if isJson {
		request.Header.Set("Content-Type", "application/json")
	}

	if client.config.API_Key != "" {
		request.Header.Set("Authorization", "Bearer "+client.config.API_Key)
	}

	if client.config.ORG_ID != "" {
		request.Header.Set("OpenAI-Organization", client.config.ORG_ID)
	}
}
