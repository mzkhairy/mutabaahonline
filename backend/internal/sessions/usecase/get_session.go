package usecase

import (
	"context"
	"errors"
	"mutabaahapi/internal/sessions/entity"
	sessRepo "mutabaahapi/internal/sessions/repository"
)

type GetSession struct {
	Repo sessRepo.SessionRepository
}

func NewGetSession(r sessRepo.SessionRepository) *GetSession {
	return &GetSession{Repo: r}
}

func (uc *GetSession) Execute(ctx context.Context, sessionID string) (*entity.Session, error) {
	s, err := uc.Repo.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, errors.New("resource not found") // Pesan standard 404
	}
	return s, nil
}
