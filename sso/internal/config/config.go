package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" end-required:"true"`
	StoragePath string        `yaml:"storage_path" end-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" end-required:"true"`
	GRPC        GRPConfig     `yaml:"grpc"`
}

type GRPConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file does not exist: " + path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Println(err)
		panic("failed to read config file: " + path)

	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or enviroment variable.
// Prioryty: flag > env > default.
// Default value is empty string
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
