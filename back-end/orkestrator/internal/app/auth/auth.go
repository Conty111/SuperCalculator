package auth

import (
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/config"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/helpers"
	webHelpers "github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/transport/web/helpers"
	"github.com/cristalhq/jwt/v5"
	"github.com/gofrs/uuid"
	"time"
)

type Auth struct {
	TokenBuilder *jwt.Builder
	HashHelper   *helpers.HashHelper
	TokenTTL     time.Duration
}

func NewAuth(appCfg *config.App) (*Auth, error) {
	tokenBuilder, err := webHelpers.NewTokenBuilder(appCfg.AuthPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	hashHelper, err := helpers.NewHasher("sha256")
	if err != nil {
		return nil, err
	}
	return &Auth{TokenBuilder: tokenBuilder, HashHelper: hashHelper, TokenTTL: appCfg.TokenTTL}, nil
}

func (a *Auth) BuildToken(userID uuid.UUID) (*jwt.Token, error) {
	tokenData := models.Token{
		UserID:  userID.String(),
		Expires: time.Now().Add(a.TokenTTL),
	}
	return a.TokenBuilder.Build(tokenData)
}

func (a *Auth) HashString(text string) (string, error) {
	return a.HashHelper.HashString(text)
}

func (a *Auth) GetTokenTTL() time.Duration {
	return a.TokenTTL
}