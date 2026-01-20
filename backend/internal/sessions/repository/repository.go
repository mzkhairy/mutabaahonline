package repository

import (
	"context"
	"mutabaahapi/internal/sessions/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, s *entity.Session) error
	GetByID(ctx context.Context, id string) (*entity.Session, error)
	ListByClass(ctx context.Context, classID string) ([]entity.Session, error)

	// [FIX] Update signature agar sesuai dengan implementasi Postgres
	ListByStudent(ctx context.Context, studentID string, isReportMode bool, classID string) ([]entity.Session, error)

	IsClassTeacher(ctx context.Context, classID, teacherID string) (bool, error)
	Update(ctx context.Context, s *entity.Session) error
}
