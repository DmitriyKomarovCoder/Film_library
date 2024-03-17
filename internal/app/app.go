package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/DmitriyKomarovCoder/Film_library/config"
	deliveryActor "github.com/DmitriyKomarovCoder/Film_library/internal/actor/delivery/http"
	repoActor "github.com/DmitriyKomarovCoder/Film_library/internal/actor/repository"
	usecaseActor "github.com/DmitriyKomarovCoder/Film_library/internal/actor/usecase"
	routerInit "github.com/DmitriyKomarovCoder/Film_library/internal/app/router"
	deliveryFilm "github.com/DmitriyKomarovCoder/Film_library/internal/movie/delivery/http"
	repoFilm "github.com/DmitriyKomarovCoder/Film_library/internal/movie/repository"
	usecaseFilm "github.com/DmitriyKomarovCoder/Film_library/internal/movie/usecase"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/closer"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/logger"
	"github.com/DmitriyKomarovCoder/Film_library/pkg/postgres"
)

func Run(cfg *config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	l, err := logger.NewLogger(cfg.Log.Path)
	if err != nil {
		log.Fatalf("Logger initialisation error %s", err)
	}

	pg, err := postgres.New(cfg.PG.Host, cfg.PG.User, cfg.PG.Password, cfg.PG.Name, cfg.PG.Port, cfg.PG.PoolMax)
	if err != nil {
		l.Fatal(fmt.Errorf("error: postgres.New: %w", err))
	}

	repMovie := repoFilm.NewRepository(pg.Pool)
	repActor := repoActor.NewRepository(pg.Pool)

	useActor := usecaseActor.NewUsecase(repActor)
	useMovie := usecaseFilm.NewUsecase(repMovie, useActor)

	handlerFilm := deliveryFilm.NewHandler(useMovie, *l)
	handlerActor := deliveryActor.NewHandler(useActor, *l)

	router := *routerInit.NewRouter(handlerFilm, handlerActor, l)

	httpServer := &http.Server{
		Addr:         cfg.Http.Host + ":" + cfg.Http.Port,
		Handler:      router,
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,
	}

	c := &closer.Closer{}
	c.Add(httpServer.Shutdown)
	c.Add(pg.Close)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatalf("Erorr starting server: %v", err)
		}
	}()
	l.Infof("server start in port: %v", cfg.Http.Port)

	<-ctx.Done()
	l.Info("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := c.Close(shutdownCtx); err != nil {
		l.Fatalf("closer: %v", err)
	}

	l.Info("Service close without error")
}
