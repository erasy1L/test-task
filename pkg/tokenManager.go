package pkg

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/erazr/test-task/config"
	"github.com/golang-jwt/jwt"
)

type TokenManager interface {
	NewJWT(ctx context.Context, guid string) (string, error)
	NewRefreshToken(ctx context.Context, token string) (string, error)
	ValidateRefreshToken(ctx context.Context, tokenString string) error
}

type Manager struct {
	secret          string
	aesKey          string
	accesTokenTTL   time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(cfg config.Config) *Manager {
	parsedAccesTokenTTL, err := time.ParseDuration(cfg.AccessTokenTTL)
	if err != nil {
		parsedAccesTokenTTL = time.Minute * 30
	}
	parsedRefreshTokenTTL, err := time.ParseDuration(cfg.RefreshTokenTTL)
	if err != nil {
		parsedRefreshTokenTTL = time.Hour * 24 * 7
	}

	return &Manager{
		secret:          cfg.Secret,
		aesKey:          cfg.AesKey,
		accesTokenTTL:   parsedAccesTokenTTL,
		refreshTokenTTL: parsedRefreshTokenTTL,
	}
}

func (m *Manager) NewJWT(guid string) (string, error) {
	exp := time.Now().Add(m.accesTokenTTL).Unix()

	claims := jwt.MapClaims{}
	claims["guid"] = guid
	claims["exp"] = exp

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(m.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) NewRefreshToken(guid string) (string, error) {
	block, err := aes.NewCipher([]byte(m.aesKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(guid), nil)
	token := fmt.Sprintf("%x", ciphertext)

	return token, nil
}

func (m *Manager) VerifyRefreshToken(refreshToken string) (string, error) {
	block, err := aes.NewCipher([]byte(m.aesKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext, _ := hex.DecodeString(refreshToken)

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	userGuid, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("failed to decrypt")
	}

	return string(userGuid), nil
}
