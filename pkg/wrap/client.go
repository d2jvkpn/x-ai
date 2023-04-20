package wrap

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"golang.org/x/net/proxy"
)

type Config struct {
	Url     string `mapstructure:"url"`
	API_Key string `mapstructure:"api_key"`
	ORG_ID  string `mapstructure:"org_id"`

	Proxy         string `mapstructure:"proxy"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	TlsSkipVerify bool   `mapstructure:"tls_skip_verify"`
}

func NewOpenAiClient(fp, key string) (client *openai.Client, err error) {
	var (
		cfg    *Config
		dialer proxy.Dialer
		config openai.ClientConfig

		proxyURL  *url.URL
		transport *http.Transport
		auth      *proxy.Auth
	)

	vp := viper.New()
	vp.SetConfigType("yaml")

	vp.SetConfigFile(fp)
	if err = vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ReadInConfig(): %q, %w", fp, err)
	}

	cfg = new(Config)
	if err = vp.UnmarshalKey(key, cfg); err != nil {
		return nil, err
	}

	if cfg.API_Key == "" {
		return nil, fmt.Errorf("empty api_key")
	}

	config = openai.DefaultConfig(cfg.API_Key)
	if cfg.Url != "" {
		config.BaseURL = cfg.Url
	}
	if cfg.ORG_ID != "" {
		config.OrgID = cfg.ORG_ID
	}

	if cfg.Proxy == "" {
		return openai.NewClientWithConfig(config), nil
	}

	if cfg.User != "" && cfg.Password != "" {
		auth = &proxy.Auth{User: cfg.User, Password: cfg.Password}
	}

	if proxyURL, err = url.Parse(cfg.Proxy); err != nil {
		return nil, err
	}

	switch proxyURL.Scheme {
	case "socks5":
		if dialer, err = proxy.SOCKS5("tcp", proxyURL.Host, auth, nil); err != nil {
			return nil, err
		}
		transport = &http.Transport{Dial: dialer.Dial}
	case "http", "https":
		transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		if auth != nil {
			bts := []byte(auth.User + ":" + auth.Password)
			basicAuth := "Basic " + base64.StdEncoding.EncodeToString(bts)
			transport.ProxyConnectHeader.Add("Proxy-Authorization", basicAuth)
		}
	default:
		return nil, fmt.Errorf("unknow proxy schema")
	}

	if strings.HasPrefix(cfg.Url, "https") && cfg.TlsSkipVerify {
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS12, InsecureSkipVerify: true,
		}
	}
	config.HTTPClient.Transport = transport

	return openai.NewClientWithConfig(config), nil
}
