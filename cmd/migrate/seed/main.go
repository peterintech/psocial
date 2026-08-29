package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/peterintech/psocial/internal/db"
	"github.com/peterintech/psocial/internal/env"
	"github.com/peterintech/psocial/internal/store"
)

func main() {
	godotenv.Load(".env")

	addr := env.GetEnv("DB_ADDR", "postgres://postgres:postgres@localhost:5432/psocial?sslmode=disable")
	conn, err := db.New(addr, 30, 30, "15ms")

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)

	if err := db.Seed(store, conn); err != nil {
		panic(err)
	}
	log.Println("Database seeding completed successfully.")
}
