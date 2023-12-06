package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/HironixRotifer/grpc-authorization-service/internal/domain/models"
	"github.com/HironixRotifer/grpc-authorization-service/internal/lib/jwt"
	"github.com/HironixRotifer/grpc-authorization-service/internal/lib/logger/sl"
	"github.com/HironixRotifer/grpc-authorization-service/internal/storage"

	"golang.org/x/crypto/bcrypt"
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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app ID")
)

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

// Login cheks if user with given credentials exist in the system.
//
// If user doesn't exist, returns error.
// If user exist, but password is incorrect, returns error.
func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("attempting to login user")

	user, err := a.UserProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("user not found: %s, err: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("Error to get user", sl.Err(err))

		return "", fmt.Errorf("user not found: %s, err: %w", op, err)
	}

	if err := VerefiPassword(user.HashPassword, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(ErrInvalidCredentials))

		return "", fmt.Errorf("invalid credentials: %s, err: %w", op, err)
	}

	app, err := a.AppProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("app not found: %s, err: %w", op, err)
	}

	log.Info("user logged is successful")

	token, err := jwt.NewToken(user, app, a.TokenTTL)
	if err != nil {
		a.log.Error("Error to generate token", sl.Err(err))

		return "", fmt.Errorf("error generating token: %s, err: %w", op, err)
	}

	return token, nil
}

// RegisterUser registers new user in the system and returns user ID.
//
// If user already exists, returns error.
// If user exists, but password is incorrect, returns error.
func (a *Auth) RegisterUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterUser"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("registering user")

	hashPassword, err := HashPassword(password)
	if err != nil {
		log.Error("Error hashing password in: %s, err: %w", op, err)

		return 0, err
	}

	id, err := a.UserSaver.SaveUser(ctx, email, hashPassword)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, err)
		}
		log.Error("Error saving user, %w, err: %w", op, err)

		return 0, err
	}

	log.Info("user registered")

	return id, nil
}

// IsAdmin checks if user is admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.UserProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("user is not found", sl.Err(err))

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		return false, fmt.Errorf("error checking if user is admin: %s, err: %w", op, err)
	}

	log.Info("checking if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

// HashPassword take password string and return hash
func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// VerefiPassword function validates the received password with the database password
func VerefiPassword(userPassword []byte, givenPassword []byte) error {
	if err := bcrypt.CompareHashAndPassword(userPassword, givenPassword); err != nil {
		return err
	}
	return nil
}
