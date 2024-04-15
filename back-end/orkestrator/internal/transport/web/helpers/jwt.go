package helpers

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/cristalhq/jwt/v5"
)

type JWTVerifier interface {
	Verify(token string) (*models.Token, error)
}

func NewTokenBuilder(privateKeyPath string) (*jwt.Builder, error) {
	privateKey, err := readPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	signer, err := jwt.NewSignerES(jwt.ES256, privateKey)
	if err != nil {
		return nil, err
	}

	return jwt.NewBuilder(signer), nil
}

type verifierImpl struct {
	verifier jwt.Verifier
}

func NewJWTVerifier(publicKeyPath string) (JWTVerifier, error) {
	pub, err := readPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key: %w", err)
	}
	verifier, err := jwt.NewVerifierES(jwt.ES256, pub)
	if err != nil {
		return nil, fmt.Errorf("while creating jwt verififer: %w", err)
	}
	return &verifierImpl{
		verifier: verifier,
	}, nil
}

func (j *verifierImpl) Verify(token string) (*models.Token, error) {
	validated, err := jwt.Parse([]byte(token), j.verifier)
	if err != nil {
		return nil, fmt.Errorf("failed to validate jwt token: %w", err)
	}
	return decodeClaims(validated)
}

func decodeClaims(token *jwt.Token) (*models.Token, error) {
	var claims models.Token
	if err := token.DecodeClaims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}
	return &claims, nil
}

func readPrivateKey(filepath string) (*ecdsa.PrivateKey, error) {
	privateKey, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("decoding private key PEM failed")
	}

	ecdsaPrivate, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return ecdsaPrivate, nil
}

func readPublicKey(filepath string) (*ecdsa.PublicKey, error) {
	publicKey, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("decoding public key PEM failed")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key at %s is not of type *ecdsa.AuthPublicKey", filepath)
	}
	return ecdsaPub, nil
}
