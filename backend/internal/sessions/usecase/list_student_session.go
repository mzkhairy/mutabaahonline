package usecase

import (
	"context"
	"mutabaahapi/internal/sessions/entity"
	"mutabaahapi/internal/sessions/repository"
)

type ListStudentSession struct {
	SessionRepo repository.SessionRepository
}

func NewListStudentSession(r repository.SessionRepository) *ListStudentSession {
	return &ListStudentSession{SessionRepo: r}
}

// Execute: Default (FILTERED) - Hanya field Student Activity untuk Dashboard Input
func (uc *ListStudentSession) Execute(ctx context.Context, studentID string, classID string) ([]entity.Session, error) {
	// Mode biasa: isReportMode = false
	return uc.SessionRepo.ListByStudent(ctx, studentID, false, classID)
}

// [FIX] Tambahkan parameter classID
func (uc *ListStudentSession) ExecuteReport(ctx context.Context, studentID string, classID string) ([]entity.Session, error) {
	// Mode laporan: isReportMode = true (ambil semua history)
	return uc.SessionRepo.ListByStudent(ctx, studentID, true, classID)
}
