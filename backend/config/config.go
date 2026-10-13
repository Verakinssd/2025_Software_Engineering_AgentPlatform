package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type InterfaceConfig struct {
	Database DatabaseConfig
	Jwt      JwtConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Charset  string
}

type JwtConfig struct {
	Secret string
	Token  TokenConfig
}

type ServerConfig struct {
	Port        string
	Mode        string
	Upload      string
	Environment string
}

type TokenConfig struct {
	Name                    string
	Expiration              int64
	RefreshBeforeExpiration int64
}

var Config InterfaceConfig = InterfaceConfig{}

func init() {
	dataBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(dataBytes, &Config)
	if err != nil {
		panic(err)
	}
}
