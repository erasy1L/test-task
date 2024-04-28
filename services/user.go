package services

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	db "github.com/erazr/test-task/db"
	models "github.com/erazr/test-task/models"
	"github.com/erazr/test-task/pkg"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo            *db.UserRepository
	refreshTokenTTL time.Duration
	tokenManager    *pkg.Manager
}

func NewUserService(repo *db.UserRepository, refreshTokenTTL string, tm *pkg.Manager) *UserService {
	parsed, err := time.ParseDuration(refreshTokenTTL)
	if err != nil {
		parsed = time.Hour * 24 * 7
	}

	return &UserService{
		repo:            repo,
		tokenManager:    tm,
		refreshTokenTTL: parsed,
	}
}

func (u *UserService) Authenticate(ctx context.Context, guid string) (models.Tokens, error) {
	user, _ := u.repo.GetByGUID(ctx, guid)
	if user.Session.RefreshToken == "" {
		return u.createSession(ctx, guid)
	}

	return models.Tokens{}, errors.New("user not found or session already exists for this user")
}

func (u *UserService) Refresh(refreshToken string) (models.Tokens, error) {
	decodedRefreshToken, err := base64.StdEncoding.DecodeString(refreshToken)
	if err != nil {
		return models.Tokens{}, err
	}

	userGuid, err := u.tokenManager.VerifyRefreshToken(string(decodedRefreshToken))
	if err != nil {
		return models.Tokens{}, err
	}

	user, err := u.repo.GetByGUID(context.Background(), userGuid)
	if err != nil {
		return models.Tokens{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Session.RefreshToken), []byte(decodedRefreshToken)); err != nil {
		return models.Tokens{}, err
	}

	return u.createSession(context.Background(), userGuid)
}

func (u *UserService) createSession(ctx context.Context, guid string) (models.Tokens, error) {
	accessToken, err := u.tokenManager.NewJWT(guid)
	if err != nil {
		return models.Tokens{}, err
	}

	refreshToken, err := u.tokenManager.NewRefreshToken(guid)
	if err != nil {
		return models.Tokens{}, err
	}

	refreshTokenHashed, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return models.Tokens{}, err
	}

	session := models.Session{
		RefreshToken: string(refreshTokenHashed),
		ExpiresAt:    time.Now().Add(u.refreshTokenTTL),
	}

	err = u.repo.SetSession(ctx, guid, session)
	if err != nil {
		return models.Tokens{}, err
	}

	refreshTokenEncoded := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	return models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenEncoded,
	}, nil
}

func (u *UserService) Register(ctx context.Context, user models.User) error {
	return u.repo.Create(ctx, user)
}
