package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	TokenTTL time.Duration  `yaml:"token_ttl" env-default:"1h"`
	GRPC     GRPCConfig     `yaml:"grpc"`
	PSQL     PostgresConfig `yaml:"psql"`
}
type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
type PostgresConfig struct {
	DbHost string `yaml:"DbHost"`
	DbPort int    `yaml:"DbPort"`
	DbUser string `yaml:"DbUser"`
	DbPass string `yaml:"DbPass"`
	DbName string `yaml:"DbName"`
}

func LoadConfig() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config yaml file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
