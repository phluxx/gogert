package main

import (
	"fmt"
	"net/http"

	"github.com/phluxx/gogert/internal/service/config"
	"github.com/phluxx/gogert/internal/service/v1handler"
)

func main() {
	// build config and launch handlers
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	h := v1handler.NewHttpHandler(cfg)
	h.RegisterHandler()
	http.ListenAndServe(fmt.Sprintf(":%s", cfg.HttpPort), h)
}
