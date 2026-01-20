package repository

import (
	"context"
	"mutabaahapi/internal/academicyears/entity"
)

type Repository interface {
	Create(ctx context.Context, e *entity.AcademicYear) error
	Update(ctx context.Context, e *entity.AcademicYear) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.AcademicYear, error)
	List(ctx context.Context, institutionID string) ([]entity.AcademicYear, error)
	ResetStatusByInstitution(ctx context.Context, institutionID string) error
}
