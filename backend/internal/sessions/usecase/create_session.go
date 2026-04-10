package usecase

import (
	"context"
	"errors"
	"mutabaahapi/internal/sessions/entity"

	// [FIX] Gunakan Interface (repository), BUKAN Postgres struct
	"mutabaahapi/internal/sessions/repository"
)

type CreateSession struct {
	// [FIX] Ubah tipe data menjadi Interface
	Repo repository.SessionRepository
}

// [FIX] Ubah argumen konstruktor menjadi Interface
func NewCreateSession(r repository.SessionRepository) *CreateSession {
	return &CreateSession{Repo: r}
}

type CreateSessionInput struct {
	UserID                string // Guru yang request
	ClassID               string
	TemplateID            string
	Date                  string
	Name                  string
	ClassActivityData     map[string]interface{}
	IsStudentInputAllowed bool
}

func (uc *CreateSession) Execute(ctx context.Context, in CreateSessionInput) (*entity.Session, error) {
	// 1. Cek Otoritas: Apakah UserID adalah guru dari ClassID?
	isTeacher, err := uc.Repo.IsClassTeacher(ctx, in.ClassID, in.UserID)
	if err != nil {
		return nil, err
	}
	if !isTeacher {
		return nil, errors.New("you are not the teacher of this class")
	}

	// 2. Siapkan Data
	session := &entity.Session{
		ClassID:               in.ClassID,
		TemplateID:            in.TemplateID,
		Date:                  in.Date,
		Name:                  in.Name,
		ClassActivityData:     in.ClassActivityData,
		IsStudentInputAllowed: in.IsStudentInputAllowed,
	}

	// 3. Simpan
	if err := uc.Repo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}
