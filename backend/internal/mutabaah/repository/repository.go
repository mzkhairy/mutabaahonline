package repository

import (
	"context"
	"mutabaahapi/internal/mutabaah/entity"
)

type MutabaahRepository interface {
	Upsert(ctx context.Context, e *entity.Entry) error
	FindBySessionID(ctx context.Context, sessionID string) ([]entity.Entry, error)
	FindByStudentID(ctx context.Context, studentID string) ([]entity.Entry, error)
	GetSummaryByStudentID(ctx context.Context, studentID string) (map[string]int, error)
}
