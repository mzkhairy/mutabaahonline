package usecase

import (
	"context"
	"mutabaahapi/internal/sessions/entity"

	// Gunakan interface repo
	sessRepo "mutabaahapi/internal/sessions/repository"
)

type ListSession struct {
	Repo sessRepo.SessionRepository
}

func NewListSession(r sessRepo.SessionRepository) *ListSession {
	return &ListSession{Repo: r}
}

func (uc *ListSession) Execute(ctx context.Context, classID string) ([]entity.Session, error) {
	// Di masa depan bisa tambah validasi apakah user berhak melihat kelas ini
	return uc.Repo.ListByClass(ctx, classID)
}
