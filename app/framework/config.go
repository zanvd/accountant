package framework

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BaseUrl string `yaml:"base_url"`
	Cache   struct {
		Database int    `yaml:"database"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
	}
	Database struct {
		Host         string `yaml:"host"`
		Name         string `yaml:"name"`
		Password     string `yaml:"password"`
		Port         string `yaml:"port"`
		RootPassword string `yaml:"root_password"`
		RootUsername string `yaml:"root_username"`
		Username     string `yaml:"username"`
	}
	Development bool `yaml:"development"`
	Mail        struct {
		DefaultSender string `yaml:"default_sender"`
		Host          string `yaml:"host"`
		Password      string `yaml:"password"`
		Port          string `yaml:"port"`
		Username      string `yaml:"username"`
	}
	Session struct {
		CookieName   string `yaml:"cookie_name"`
		KeyPrefix    string `yaml:"key_prefix"`
		SecureCookie bool   `yaml:"secure_cookie"`
	}
}

func NewConfig() (*Config, error) {
	c := &Config{}
	f, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.UnmarshalStrict(f, c)

	return c, err
}
