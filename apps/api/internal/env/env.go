package env

import (
	"fmt"

	troncenv "github.com/FacileStudio/tronc/env"
)

// Config is the process configuration: the tronc core plus the Capsule-specific
// maximum paste size in bytes.
type Config struct {
	troncenv.Core
	MaxPasteSize int
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

	return Config{Core: core, MaxPasteSize: maxPasteSize}, nil
}
