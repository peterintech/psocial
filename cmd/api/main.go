package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/peterintech/psocial/internal/auth"
	"github.com/peterintech/psocial/internal/db"
	"github.com/peterintech/psocial/internal/env"
	"github.com/peterintech/psocial/internal/mailer"
	"github.com/peterintech/psocial/internal/store"
	"go.uber.org/zap"
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

//	@BasePath					/v1
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	godotenv.Load(".env")

	cfg := config{
		addr:        fmt.Sprintf(":%s", env.GetEnv("PORT", "8080")),
		apiURL:      env.GetEnv("API_URL", "localhost:8080"),
		frontendURL: env.GetEnv("FRONTEND_URL", "http://localhost:5173/"),
		db: dbConfig{
			addr:         env.GetEnv("DB_ADDR", "postgres://postgres:postgres@localhost:5432/psocial?sslmode=disable"),
			maxOpenConns: env.GetEnvAsInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetEnvAsInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetEnv("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetEnv("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetEnv("MAIL_FROM_EMAIL", "peter@cloverkrafts.com"),
			sendGrid: sendGridConfig{
				apiKey: env.GetEnv("SENDGRID_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicAuthConfig{
				user: env.GetEnv("BASIC_AUTH_USER", "admin"),
				pass: env.GetEnv("BASIC_AUTH_PASS", "password"),
			},
			token: tokenAuthConfig{
				secret: env.GetEnv("JWT_SECRET", "mysecretkey"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "psocial",
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//Database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)

	if err != nil {
		logger.Fatal("Error connecting to the database: %v", err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewStorage(db)

	mailer := mailer.NewSendGridMailer(cfg.mail.fromEmail, cfg.mail.sendGrid.apiKey)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss) // 24 hours

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
