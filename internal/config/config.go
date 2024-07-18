package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Server struct {
		Host          string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
		Port          int    `yaml:"port" env:"SERVER_PORT" env-default:"8081"`
		Secure        bool   `yaml:"secure" env:"SERVER_SECURE" env-default:"false"`
		Secret        string `yaml:"secret" env:"SERVER_SECRET" env-default:""`
		TokenDuration int    `yaml:"token_duration" env:"TOKEN_DURATION" env-default:"300"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
		Port     int    `yaml:"port" env:"PORT" env-default:"5432"`
		Name     string `yaml:"name" env:"NAME" env-default:"owl"`
		User     string `yaml:"username"  env-default:"postgres"`
		Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	} `yaml:"database"`
	Logger struct {
		Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	} `yaml:"log"`
}

func (s *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", s.Server.Host, s.Server.Port)
}

func (s *ServerConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", s.Database.Host, s.Database.User, s.Database.Password, s.Database.Name)
}

func NewServerConfig() (*ServerConfig, error) {
	config := ServerConfig{}

	if err := cleanenv.ReadConfig("./config/server_cfg.yaml", &config); err != nil {
		return nil, fmt.Errorf("cannot read server config: %s", err)
	}

	return &config, nil
}

func (c *ClientConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

type ClientConfig struct {
	Server struct {
		Host string `yaml:"host" env:"API_HOST" env-default:"localhost"`
		Port int    `yaml:"port" env:"API_PORT" env-default:"8081"`
	} `yaml:"server"`
	Cert struct {
		CA   string `yaml:"ca" env:"CA_PATH" env-default:"cert/CertAuth.crt"`
		Cert string `yaml:"cert" env:"CLIENT_CERT_PATH" env-default:"cert/cli.crt"`
		Key  string `yaml:"key" env:"CLIENT_KEY_PATH" env-default:"cert/cli.key"`
	} `yaml:"cert"`
	Auth struct {
		RefreshPeriod int `yaml:"refresh_period" env:"REFRESH_PERIOD" env-default:"300"`
	} `yaml:"auth"`
}

func NewClientConfig() (*ClientConfig, error) {
	config := ClientConfig{}

	if err := cleanenv.ReadConfig("./config/client_cfg.yaml", &config); err != nil {
		return nil, fmt.Errorf("cannot read cli config: %s", err)
	}

	return &config, nil
}
