package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/peterintech/psocial/internal/env"
)

func main() {
	godotenv.Load(".env")

	cfg := config{
		addr: fmt.Sprintf(":%s", env.GetEnv("PORT", "8080")),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
