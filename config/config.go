package config

import (
	"github.com/go-errors/errors"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Configuration struct {
	TimeOut int `yaml:"time_out"`
	Services struct{
		Source interface{} `yaml:"source"`
		Destination interface{} `yaml:"destination"`
		Crash interface{} `yaml:"crash"`
	} `yaml:"services"`
	Log struct{
		Path string `yaml:"path"`
		Format string `yaml:"format"`
	} `yaml:"log"`
}

var conf Configuration

func Init(f []byte) error {
	err := yaml.Unmarshal(f, &conf)
	if err!=nil{
		return errors.Wrap(err, -1)
	}
	file, err := os.OpenFile(conf.Log.Path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(file)
	return nil
}

func GetConfig() *Configuration {
	return &conf
}
