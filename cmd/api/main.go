package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/peterintech/psocial/internal/env"
	"github.com/peterintech/psocial/internal/store"
)

func main() {
	godotenv.Load(".env")

	cfg := config{
		addr: fmt.Sprintf(":%s", env.GetEnv("PORT", "8080")),
	}

	store := store.NewStorage(nil)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
