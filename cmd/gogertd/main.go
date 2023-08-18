package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/phluxx/gogert/internal/service/config"
	"github.com/phluxx/gogert/internal/service/v1handler"
)

func main() {
	slog.Info("main: building config")
	// build config and launch handlers
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("main: failed to load config", "error", err)
		os.Exit(1)
	}
	h := v1handler.NewHttpHandler(cfg)
	h.RegisterHandler()

	slog.Info("main: starting http server", "port", cfg.HttpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpPort), h)
	if err != nil {
		slog.Error("main: failed to start http server", "error", err)
		os.Exit(1)
	}
}
