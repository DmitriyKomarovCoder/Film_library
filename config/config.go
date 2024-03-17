package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Http            http          `yaml:"http"`
	Log             logCustom     `yaml:"log_file"`
	PG              postgres      `yaml:"postgres"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type http struct {
	Host         string        `yaml:"host"`
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type logCustom struct {
	Path string `yaml:"path"`
}

type postgres struct {
	Name     string `env:"DB_NAME"`
	User     string `env:"DB_USER"`
	Port     int    `env:"DB_PORT"`
	Password string `env:"DB_PASSWORD"`
	Host     string `env:"DB_HOST"`
	PoolMax  int32  `yaml:"pool_max"`
}

func NewConfig(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return &cfg, err
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		fmt.Println(err)
		return &cfg, err
	}

	log.Println("Parsed Configuration")
	log.Println(cfg)
	return &cfg, nil
}
