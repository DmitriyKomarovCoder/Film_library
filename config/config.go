package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Http http      `yaml:"http"`
	Log  logCustom `yaml:"log_file"`
	PG   postgres
}

type http struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
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
}

// временное решение, потом буду делать это в compose
func setEnvValues() error {
	err := os.Setenv("DB_NAME", "FilmLibrary")
	if err != nil {
		return fmt.Errorf("Error setting port, err = %v", err)
	}

	err = os.Setenv("DB_USER", "kosmatoff")
	if err != nil {
		return fmt.Errorf("Error setting jwt secret, err = %v", err)
	}

	err = os.Setenv("DB_PASSWORD", "2003")
	if err != nil {
		return fmt.Errorf("Error setting jwt secret, err = %v", err)
	}

	err = os.Setenv("DB_HOST", "postgres")
	if err != nil {
		return fmt.Errorf("Error setting jwt secret, err = %v", err)
	}

	err = os.Setenv("DB_PORT", "5432")
	if err != nil {
		return fmt.Errorf("Error setting jwt secret, err = %v", err)
	}

	return nil
}

func NewConfig(path string) (*Config, error) {
	var cfg Config

	err := setEnvValues()
	if err != nil {
		panic(err)
	}

	err = cleanenv.ReadEnv(&cfg)
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
