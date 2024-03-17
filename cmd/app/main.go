package main

import (
	"log"

	"github.com/DmitriyKomarovCoder/Film_library/config"
	"github.com/DmitriyKomarovCoder/Film_library/internal/app"
)

const path = "./config/config.yaml"

// @title		FilmLibrary API
// @version		1.0.0
// @description	Server API for FilmLibrary Application

// @contact.name   FilmLibrary API Support
// @contact.email  dimka.komarov@bk.ru
// @contact.url    https://t.me/Kosmatoff

// @host		127.0.0.1:8080
// @BasePath	/
func main() {
	cfg, err := config.NewConfig(path)
	if err != nil {
		log.Fatalf("Config error %s", err)
	}

	app.Run(cfg)
}
