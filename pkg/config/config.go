package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type ServerConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type GRPCServerConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Config struct {
	Env        string           `yaml:"env" env-default:"local"`
	APIServer  ServerConfig     `yaml:"api_server"`
	GRPCServer GRPCServerConfig `yaml:"grpc_server"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic(any("Файл конфигурации по указанному пути отсутствует"))
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(any(fmt.Sprintf("Файл конфигурации не найден: %s", configPath)))
	}

	cfg := new(Config)

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		panic(any(fmt.Sprintf("Ошибка чтения файла конфигурации: %v", err)))
	}

	return cfg
}

// fetchConfigPath - парсинг пути к конфигурации из флага или переменной окружения
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	if res == "" {
		path := "local.yaml"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			path = "./config/local.yaml"
		}
		res = path
	}

	return res
}
