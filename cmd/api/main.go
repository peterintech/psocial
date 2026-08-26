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

//	@title			Psocial API
//	@version		1.0
//	@description	API documentation for Psocial API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//	@securityDefinitions.basic	ApiKeyAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	godotenv.Load(".env")

	cfg := config{
		addr:   fmt.Sprintf(":%s", env.GetEnv("PORT", "8080")),
		apiURL: env.GetEnv("API_URL", "localhost:8080"),
		env:    env.GetEnv("ENV", "development"),
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
