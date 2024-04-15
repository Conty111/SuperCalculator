package initializers

import (
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
)

// InitializeEnvs intializes envy
func InitializeEnvs() {
	if err := envy.Load(); err != nil {
		log.Error().Err(err).Msg("can not load env files")

		envy.Reload()
	}
}
