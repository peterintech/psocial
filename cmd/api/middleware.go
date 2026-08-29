package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/peterintech/psocial/internal/store"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("missing authorization header"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("invalid authorization header format"))
			return
		}

		tokenString := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(tokenString)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID := fmt.Sprintf("%v", claims["sub"])

		user, err := app.getUser(r.Context(), userID)
		if err != nil {
			app.unauthorizedError(w, r, fmt.Errorf("user not found"))
			return
		}

		// Add the user
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, http.ErrNoCookie)
				return
			}
			//parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("auth header is not in the correct format"))
				return
			}
			//decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			//check the credentials
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})

	}
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getUserFromContext(r)
		post, err := getPostFromContext(r.Context())
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		//check if the user is the owner of the post
		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}

		//role precedence
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		if !allowed {
			app.forbiddenError(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil
}

func (app *application) getUser(ctx context.Context, userID string) (*store.User, error) {
	if !app.config.redisCfg.enabled {
		fmt.Println("-- redis is disabled --")
		return app.store.Users.GetByID(ctx, userID)
	}

	user, err := app.cacheStorage.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = app.store.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	} else {
		fmt.Println("-- redis has the user --")

	}

	return user, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			// RemoteAddr is "host:port" — extract just the host
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			if allowed, retryAfter := app.rateLimiter.Allow(ip); !allowed {
				app.tooManyRequestsError(w, r, retryAfter.String())
				return
			}

		}

		next.ServeHTTP(w, r)
	})
}
