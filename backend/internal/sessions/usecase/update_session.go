package usecase

import (
	"context"
	"errors"
	"mutabaahapi/internal/sessions/entity"
	"mutabaahapi/internal/sessions/repository"
)

type UpdateSession struct {
	SessionRepo repository.SessionRepository
}

func NewUpdateSession(r repository.SessionRepository) *UpdateSession {
	return &UpdateSession{SessionRepo: r}
}

type UpdateSessionInput struct {
	Name                  *string // Pointer agar bisa detect null (optional update)
	Date                  *string
	IsStudentInputAllowed *bool
	ClassActivityData     map[string]interface{}
}

func (uc *UpdateSession) Execute(ctx context.Context, sessionID string, input UpdateSessionInput) error {
	// 1. Ambil Data Lama
	sess, err := uc.SessionRepo.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if sess == nil {
		return errors.New("session not found")
	}

	// 2. Update Field (Hanya jika input tidak nil)
	if input.Name != nil {
		sess.Name = *input.Name
	}
	if input.Date != nil {
		sess.Date = *input.Date
	}
	if input.IsStudentInputAllowed != nil {
		sess.IsStudentInputAllowed = *input.IsStudentInputAllowed
	}
	if input.ClassActivityData != nil {
		sess.ClassActivityData = entity.JSONMap(input.ClassActivityData)
	}

	// 3. Simpan
	return uc.SessionRepo.Update(ctx, sess)
}
