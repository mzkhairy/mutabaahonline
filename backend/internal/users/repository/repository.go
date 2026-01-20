package repository

import (
	"context"
	"mutabaahapi/internal/users/entity"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	FindByUsername(ctx context.Context, institutionID string, username string) (*entity.User, error)
	FindByIdentifier(ctx context.Context, institutionID string, identifier string) (*entity.User, error)
	FindBySerialNumber(ctx context.Context, institutionID string, serialNumber string) (*entity.User, error)

	GetByID(ctx context.Context, id string) (*entity.User, error)
	ListByInstitution(ctx context.Context, institutionID string, roleFilter string, search string) ([]entity.User, error)
	UpdateStatus(ctx context.Context, userID string, status string) error
	CountByInstitution(ctx context.Context, institutionID string) (int, error)
	GetStats(ctx context.Context, institutionID string) (map[string]int, error)
	UpdatePassword(ctx context.Context, userID, hash string) error
	Update(ctx context.Context, u *entity.User) error
}
