package initializers

import (
	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog/log"
)

// InitializeEnvs intializes envy
func InitializeEnvs() {
	if err := envy.Load("enviroments/sys.env", "enviroments/db.sys.env", "enviroments/orkestrator.sys.env", "enviroments/kafka.sys.env"); err != nil {
		log.Error().Err(err).Msg("can not load sys.env files")
	}
}
