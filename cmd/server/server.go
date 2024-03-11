package main

import (
	"pow/pkg/config"
	"pow/pkg/log"
	"pow/pkg/server"
)

func main() {
	cfg := config.ParseConfig()
	s, err := server.NewServer(cfg)
	if err != nil {
		log.Errorf("error starting server: %s", err)
		return
	}
	defer s.Close()

	log.Infof("server listening on port %d", cfg.Port)
	s.Serve()
}
