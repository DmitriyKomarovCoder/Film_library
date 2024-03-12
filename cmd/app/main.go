package main

import (
	"github.com/DmitriyKomarovCoder/Film_library/config"
	"github.com/DmitriyKomarovCoder/Film_library/internal/app"
	"log"
)

const path = "./config/config.yaml"

func main() {
	cfg, err := config.NewConfig(path)
	if err != nil {
		log.Fatalf("Config error %s", err)
	}

	app.Run(cfg)
}
