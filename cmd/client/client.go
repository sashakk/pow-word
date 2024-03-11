package main

import (
	"pow/pkg/client"
	"pow/pkg/config"
	"pow/pkg/log"
)

func main() {
	cfg := config.ParseConfig()
	c, err := client.NewClient(cfg.Host, cfg.Port, cfg.ConnTimeout)
	if err != nil {
		log.Errorf("error creating a new client: %s", err)
		return
	}
	defer c.Close()

	quote, err := c.Authorize()
	if err != nil {
		log.Errorf("error to authorize: %s", err)
	}
	log.Infof("successfully authorized: %s", quote)
}
