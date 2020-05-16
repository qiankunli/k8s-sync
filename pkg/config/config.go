package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	ApiServer `yaml:"apiServer"`
	Clean     `yaml:"clean"`
	Log       `yaml:"log"`
}

func (config *Config) Check() error {
	if err := config.ApiServer.Check(); err != nil {
		return err
	}
	if err := config.Clean.Check(); err != nil {
		return err
	}
	if err := config.Log.Check(); err != nil {
		return err
	}
	return nil
}

type Clean struct {
	IntervalSeconds int32 `yaml:"intervalSeconds"`
}

func (clean *Clean) Check() error {
	if clean.IntervalSeconds == 0 {
		return errors.New("you must set clean.IntervalSeconds")
	}
	return nil
}

type Log struct {
	Debug bool `yaml:"debug"`
}

func (log *Log) Check() error {
	return nil
}

type ApiServer struct {
	Env         string `yaml:"env"`
	Host        string `yaml:"host"`
	BearerToken string `yaml:"bearerToken"`
}

func (apiServer *ApiServer) Check() error {
	if len(apiServer.Host) == 0 {
		return errors.New("apiServer.Host can not be null")
	}
	if len(apiServer.BearerToken) == 0 {
		return errors.New("apiServer.BearerToken can not be null")
	}
	if len(apiServer.Env) == 0 {
		return errors.New("apiServer.Env can not be null")
	}
	return nil
}

func ReadYaml(path string) (*Config, error) {
	conf := &Config{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(f).Decode(conf)
	}
	if err := conf.Check(); err != nil {
		return nil, err
	}
	return conf, nil
}
