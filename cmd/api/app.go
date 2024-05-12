package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	commandDHTTP "github.com/janaridev/postgres_pro/internal/command/dhttp"
	commandRepo "github.com/janaridev/postgres_pro/internal/command/repo"
	commandService "github.com/janaridev/postgres_pro/internal/command/service"
	"github.com/janaridev/postgres_pro/internal/config"
	"github.com/janaridev/postgres_pro/pkg/db"
	"github.com/janaridev/postgres_pro/pkg/logger"
	"github.com/janaridev/postgres_pro/pkg/pool"
)

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logger := logger.New(cfg.Env)

	logger.Info("init app context")
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		os.Interrupt,
	)
	defer cancel()

	logger.Info("init db")
	conn := db.GenerateConnString(cfg.DB.User, cfg.DB.Password, cfg.DB.Name,
		cfg.DB.Host, cfg.DB.Port, cfg.DB.UseSSL)

	db, err := db.New(conn)
	if err != nil {
		return err
	}

	logger.Info("init pool")
	pool := pool.New(&pool.PoolOptions{Max: 10})

	logger.Info("init command repo")
	commandRepo := commandRepo.New(db)

	logger.Info("init command service")
	commandService := commandService.New(logger, commandRepo, pool)

	logger.Info("init command routes")
	commandDHTTP := commandDHTTP.New(logger, commandService)

	logger.Info("init mux router")
	mux := http.NewServeMux()

	mux.HandleFunc("/api/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK\n")
	})

	mux.HandleFunc("GET /api/command", commandDHTTP.List)
	mux.HandleFunc("GET /api/command/{id}", commandDHTTP.Get)
	mux.HandleFunc("POST /api/command", commandDHTTP.Create)
	mux.HandleFunc("DELETE /api/command/{id}", commandDHTTP.Stop)

	go func() {
		url := net.JoinHostPort(cfg.HTTPServer.Host, strconv.Itoa(cfg.HTTPServer.Port))
		logger.Info("server running on " + url)
		if err := http.ListenAndServe(url, mux); err != nil {
			logger.Error("", err)
		}
	}()
	<-ctx.Done()

	return nil
}
