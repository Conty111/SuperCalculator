package initializers

import (
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/app/auth"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/config"
	"github.com/rs/zerolog/log"
)

func InitializeAuth(appCfg *config.App) *auth.Auth {
	authManager, err := auth.NewAuth(appCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while initializing auth")
	}
	return authManager
}