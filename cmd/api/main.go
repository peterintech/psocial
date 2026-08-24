package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/peterintech/psocial/internal/db"
	"github.com/peterintech/psocial/internal/env"
	"github.com/peterintech/psocial/internal/store"
)

const version = "1.0.0"

func main() {
	godotenv.Load(".env")

	cfg := config{
		addr: fmt.Sprintf(":%s", env.GetEnv("PORT", "8080")),
		env:  env.GetEnv("ENV", "development"),
		db: dbConfig{
			addr:         env.GetEnv("DB_ADDR", "postgres://postgres:postgres@localhost:5432/psocial?sslmode=disable"),
			maxOpenConns: env.GetEnvAsInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetEnvAsInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetEnv("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Panicf("Error connecting to the database: %v", err)
	}

	defer db.Close()

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
