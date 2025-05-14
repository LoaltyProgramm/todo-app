package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LoaltyProgramm/to-do-app/internal/config"
)

func StartServer(cfg *config.Config) error {
	port := cfg.Port
	staticPath := cfg.StaticPath

	http.Handle("/", http.FileServer(http.Dir(staticPath)))

	log.Printf("server start on port: %v", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		return fmt.Errorf("server is not work: %w", err)
	}

	return nil
}