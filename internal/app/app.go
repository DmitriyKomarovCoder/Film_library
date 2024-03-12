package app

import (
	"context"
	"github.com/DmitriyKomarovCoder/Film_library/config"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"log"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	l, err := logger.NewLogger(cfg.Log.Path)
	if err != nil {
		log.Fatalf("Logger initialisation error %s", err)
	}

	<-ctx.Done()
	l.Info("shutting down server gracefully")
}
