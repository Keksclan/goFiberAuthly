package app

import (
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/keksclan/goAuthly/authly"

	"goFiberAuthly/internal/config"
	apphttp "goFiberAuthly/internal/http"
)

// Application holds the Fiber app, goAuthly engine, and config.
type Application struct {
	Fiber  *fiber.App
	Engine *authly.Engine
	Config *config.Config
	Ready  bool
}

// New creates and configures the Application.
// It initializes the goAuthly Engine based on the provided config and sets up routes.
func New(cfg *config.Config) (*Application, error) {
	engine, err := buildAuthEngine(cfg)
	if err != nil {
		return nil, fmt.Errorf("init goAuthly engine: %w", err)
	}

	app := &Application{
		Config: cfg,
		Engine: engine,
	}

	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	})

	apphttp.SetupRoutes(fiberApp, engine, cfg, &app.Ready)
	app.Fiber = fiberApp
	app.Ready = true

	slog.Info("application initialized",
		"port", cfg.Server.Port,
		"auth_issuer", cfg.Auth.Issuer,
		"auth_jwks_url", cfg.Auth.JWKSURL,
	)

	return app, nil
}

// buildAuthEngine creates the goAuthly Engine from the application config.
// This is the central integration point – adjust goAuthly Config here.
func buildAuthEngine(cfg *config.Config) (*authly.Engine, error) {
	// Determine OAuth2 mode based on available config.
	oauth2Mode := authly.OAuth2JWTAndOpaque
	if cfg.Auth.HasJWKS() && !cfg.Auth.HasIntrospection() {
		oauth2Mode = authly.OAuth2JWTOnly
	} else if !cfg.Auth.HasJWKS() && cfg.Auth.HasIntrospection() {
		oauth2Mode = authly.OAuth2OpaqueOnly
	}

	// Build audience string for goAuthly.
	audience := cfg.Auth.Audience

	// Build introspection auth if client credentials are provided.
	var introAuth authly.ClientAuth
	if cfg.Auth.ClientID != "" {
		introAuth = authly.ClientAuth{
			Kind:         authly.ClientAuthBasic,
			ClientID:     cfg.Auth.ClientID,
			ClientSecret: cfg.Auth.ClientSecret,
		}
	}

	authlyCfg := authly.Config{
		Mode: authly.AuthModeOAuth2,
		OAuth2: authly.OAuth2Config{
			Mode:     oauth2Mode,
			Issuer:   cfg.Auth.Issuer,
			Audience: audience,
			JWKSURL:  cfg.Auth.JWKSURL,
			Introspection: authly.IntrospectionConfig{
				Endpoint: cfg.Auth.IntrospectionURL,
				Auth:     introAuth,
			},
		},
	}

	return authly.New(authlyCfg)
}
