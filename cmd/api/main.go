package main

import (
	"log"
	"ssb/internal/app"
)

func main() {
	cfg := app.LoadConfig()
	a := app.NewApp(cfg)

	if err := a.Run(); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
