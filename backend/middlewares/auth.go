package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/WillieBam/support_copilot/backend/app/config"
)

type authContextKey string

const (
	authTokenContextKey authContextKey = "authToken"
	firebaseUIDContextKey authContextKey = "firebaseUID"
)

var (
	initFirebaseClientOnce sync.Once
	cachedFirebaseClient   *auth.Client
	cachedFirebaseErr      error
)

func getFirebaseClient(ctx context.Context) (*auth.Client, error) {
	initFirebaseClientOnce.Do(func() {
		cfg := config.Get()
		firebaseCfg := &firebase.Config{}
		if cfg.Firebase.ProjectID != "" {
			firebaseCfg.ProjectID = cfg.Firebase.ProjectID
		}

		app, err := firebase.NewApp(ctx, firebaseCfg)
		if err != nil {
			cachedFirebaseErr = fmt.Errorf("failed to init firebase app: %w", err)
			return
		}

		cachedFirebaseClient, cachedFirebaseErr = app.Auth(ctx)
		if cachedFirebaseErr != nil {
			cachedFirebaseErr = fmt.Errorf("failed to init firebase auth client: %w", cachedFirebaseErr)
		}
	})

	if cachedFirebaseErr != nil {
		return nil, cachedFirebaseErr
	}

	return cachedFirebaseClient, nil
}

func getBearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}

	if len(header) >= 7 && strings.EqualFold(header[:7], "Bearer ") {
		return strings.TrimSpace(header[7:])
	}

	return ""
}

func nestedMapValue(m map[string]any, key string) map[string]any {
	v, ok := m[key]
	if !ok {
		return nil
	}

	obj, ok := v.(map[string]any)
	if !ok {
		return nil
	}

	return obj
}

func tokenHasTOTPClaim(claims map[string]any) bool {
	firebaseClaims := nestedMapValue(claims, "firebase")
	if firebaseClaims == nil {
		return false
	}

	methodRaw, ok := firebaseClaims["sign_in_second_factor"]
	if !ok {
		return false
	}

	method, ok := methodRaw.(string)
	if !ok {
		return false
	}

	return strings.EqualFold(method, "totp")
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authClient, err := getFirebaseClient(r.Context())
		if err != nil {
			http.Error(w, "Auth service unavailable", http.StatusInternalServerError)
			return
		}

		verified, err := authClient.VerifyIDToken(r.Context(), token)
		if err != nil {
			if auth.IsIDTokenExpired(err) {
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		cfg := config.Get()
		if cfg.Auth.TOTPRequired && !tokenHasTOTPClaim(verified.Claims) {
			http.Error(w, "TOTP second factor required", http.StatusUnauthorized)
			return
		}

		// Preserve the auth token in context for downstream verification/authorization layers.
		ctx := context.WithValue(r.Context(), authTokenContextKey, token)
		ctx = context.WithValue(ctx, firebaseUIDContextKey, verified.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
