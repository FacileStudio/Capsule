package env

import (
	"fmt"

	troncenv "github.com/FacileStudio/tronc/env"
)

type Config struct {
	troncenv.Core
	MaxPasteSize int
}

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
