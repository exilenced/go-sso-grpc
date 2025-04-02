package postgres

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sso/internal/config"
	"sso/internal/domain/models"
	"sso/internal/storage"
)

type Storage struct {
	db *gorm.DB
}

func New(config *config.Config) (*Storage, error) {
	const op = "storage.postgres.New"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.PSQL.DbHost,
		config.PSQL.DbUser,
		config.PSQL.DbPass,
		config.PSQL.DbName,
		config.PSQL.DbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = db.AutoMigrate(&models.User{}, &models.App{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, username string, passHash []byte) (uid int64, err error) {
	const op = "storage.postgres.SaveUser"
	user := &models.User{
		Username: username,
		PassHash: passHash,
	}
	result := s.db.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return 0, storage.ErrUserExists
		}
		return 0, fmt.Errorf("%s: %w", op, result.Error)
	}
	return user.ID, nil
}

func (s *Storage) User(ctx context.Context, username string) (models.User, error) {
	const op = "storage.postgres.User"
	var user models.User
	result := s.db.Model(&user).Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"
	user := &models.User{
		ID: userID,
	}
	var isAdmin bool
	result := s.db.Model(&user).Select("is_admin").Where("id = ?", userID).Scan(&isAdmin)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, storage.ErrUserNotFound
		}
		return false, fmt.Errorf("%s: %w", op, result.Error)
	}
	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.App"
	var app models.App
	result := s.db.Model(&app).Where("app_id = ?", appID).First(&app)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
	}
	return app, nil
}
