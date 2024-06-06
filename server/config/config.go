package config

import (
	"fmt"
	"os"

	"gitee.com/openeuler/PilotGo/sdk/logger"
	"gopkg.in/yaml.v2"
)

type PluginPrometheus struct {
	URL string `yaml:"url"`
}
type PrometheusServer struct {
	Addr string `yaml:"addr"`
}
type HttpServer struct {
	Addr string `yaml:"addr"`
}
type PilotGoServer struct {
	Addr string `yaml:"addr"`
}

type MysqlDBInfo struct {
	HostName string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"user"`
	Password string `yaml:"password"`
	DataBase string `yaml:"database"`
}

type ServerConfig struct {
	PluginPrometheus *PluginPrometheus `yaml:"plugin_prometheus"`
	PrometheusServer *PrometheusServer `yaml:"prometheus_server"`
	HttpServer       *HttpServer       `yaml:"http_server"`
	Logopts          *logger.LogOpts   `yaml:"log"`
	Mysql            *MysqlDBInfo      `yaml:"mysql"`
}

const config_file = "./config.yml"

var global_config ServerConfig

func Init() {
	err := readConfig(config_file, &global_config)
	if err != nil {
		fmt.Printf("%v", err.Error())
		os.Exit(-1)
	}
}

func Config() *ServerConfig {
	return &global_config
}

func readConfig(file string, config interface{}) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("open %s failed! err = %s\n", file, err.Error())
		return err
	}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		fmt.Printf("yaml Unmarshal %s failed!\n", string(bytes))
		return err
	}
	return nil
}
