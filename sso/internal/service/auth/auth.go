package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/HironixRotifer/grpc-authorization-service/internal/domain/models"
)

type Auth struct {
	log          *slog.Logger
	UserSaver    UserSaver
	UserProvider UserProvider
	AppProvider  AppProvider
	TokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, emain string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// NewAuth creates a new instance of the Auth service.
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appPrivider AppProvider,
	tokenTLL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		UserSaver:    userSaver,
		UserProvider: userProvider,
		AppProvider:  appPrivider,
		TokenTTL:     tokenTLL,
	}
}
