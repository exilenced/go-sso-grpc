package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"time"
)

type Auth struct {
	log         *slog.Logger
	userStorage UserStorage
	appProvider AppProvider
	tokenTTL    time.Duration
}
type UserStorage interface {
	SaveUser(ctx context.Context, username string, passHash []byte) (uid int64, err error)
	User(ctx context.Context, username string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

// New returns a new instance of the Auth service.
func New(log *slog.Logger, userStorage UserStorage, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		userStorage: userStorage,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, username string, password string, appID int) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op))

	user, err := a.userStorage.User(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", err.Error())
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid password", err.Error())
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user logged")
	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", err.Error())
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil

}
func (a *Auth) RegisterNewUser(ctx context.Context, username string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("registering user")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err.Error())
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	id, err := a.userStorage.SaveUser(ctx, username, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", err.Error())

			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	log.Info("user registered")
	return id, nil
}
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	log.Info("checking if user is admin")
	isAdmin, err := a.userStorage.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", err.Error())
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checking if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
