package env

import (
	"fmt"

	"github.com/FacileStudio/porte"
	troncenv "github.com/FacileStudio/tronc/env"
)

// OIDCConfig is the single sign-on client configuration, present only when
// OIDC_ISSUER is set.
type OIDCConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	SuccessURL   string
}

// Config is the process configuration: the tronc core, the Capsule-specific
// maximum paste size in bytes, and the optional single sign-on client.
type Config struct {
	troncenv.Core
	MaxPasteSize int
	OIDC         *OIDCConfig
}

// Load reads the environment and validates it, defaulting MAX_PASTE_SIZE to
// 1 MiB when unset.
func Load() (Config, error) {
	core, err := troncenv.LoadCore()
	if err != nil {
		return Config{}, err
	}
	if core.Port < 1 || core.Port > 65535 {
		return Config{}, fmt.Errorf("PORT must be a valid TCP port")
	}

	maxPasteSize, err := troncenv.Int("MAX_PASTE_SIZE", 1048576)
	if err != nil {
		return Config{}, err
	}
	if maxPasteSize < 1 {
		return Config{}, fmt.Errorf("MAX_PASTE_SIZE must be a positive integer")
	}

	cfg := Config{Core: core, MaxPasteSize: maxPasteSize}
	if err := cfg.loadOIDC(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) loadOIDC() error {
	issuer := troncenv.String("OIDC_ISSUER", "")
	if issuer == "" {
		return nil
	}

	clientID := troncenv.String("OIDC_CLIENT_ID", "")
	clientSecret := troncenv.String("OIDC_CLIENT_SECRET", "")
	redirectURL := troncenv.String("OIDC_REDIRECT_URL", "")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return fmt.Errorf("OIDC_CLIENT_ID, OIDC_CLIENT_SECRET, and OIDC_REDIRECT_URL are required when OIDC_ISSUER is set")
	}

	successURL := troncenv.String("OIDC_SUCCESS_URL", "")
	if successURL == "" && len(c.CORSAllowedOrigins) > 0 {
		successURL = c.CORSAllowedOrigins[0]
	}

	c.OIDC = &OIDCConfig{
		Issuer:       issuer,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		SuccessURL:   successURL,
	}
	return nil
}

// Porte translates this config into porte's. Sign-in is single sign-on only:
// Capsule has no password path. FailureURL is the home page, since Capsule has
// no /login route for porte's default to land on.
func (c Config) Porte() porte.Config {
	cfg := porte.Config{SSOOnly: true}
	if c.OIDC == nil {
		return cfg
	}
	cfg.Issuer = c.OIDC.Issuer
	cfg.ClientID = c.OIDC.ClientID
	cfg.ClientSecret = c.OIDC.ClientSecret
	cfg.RedirectURL = c.OIDC.RedirectURL
	cfg.SuccessURL = c.OIDC.SuccessURL
	cfg.FailureURL = c.OIDC.SuccessURL
	return cfg
}
