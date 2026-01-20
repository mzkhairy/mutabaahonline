package usecase

import (
	"context"
	"errors"
	"mutabaahapi/internal/templates/entity"
	"mutabaahapi/internal/templates/repository/postgres"
	userRepo "mutabaahapi/internal/users/repository"
)

type UseCase struct {
	Repo     *postgres.Repository
	UserRepo userRepo.UserRepository
}

func NewUseCase(r *postgres.Repository, u userRepo.UserRepository) *UseCase {
	return &UseCase{Repo: r, UserRepo: u}
}

func (uc *UseCase) Create(ctx context.Context, adminID string, input *entity.Template) error {
	// 1. Ambil Institution ID Admin
	user, err := uc.UserRepo.GetByID(ctx, adminID)
	if err != nil {
		return err
	}
	if user == nil || user.InstitutionID == nil {
		return errors.New("invalid admin")
	}

	// 2. Validasi Structure (Minimal harus ada isinya)
	if len(input.Structure) == 0 {
		return errors.New("template structure cannot be empty")
	}
	for _, item := range input.Structure {
		if item.Type != "class_activity" && item.Type != "student_activity" {
			return errors.New("invalid item type: " + item.Type)
		}
	}

	input.InstitutionID = *user.InstitutionID
	return uc.Repo.Create(ctx, input)
}

func (uc *UseCase) List(ctx context.Context, adminID string) ([]entity.Template, error) {
	user, err := uc.UserRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	return uc.Repo.List(ctx, *user.InstitutionID)
}

func (uc *UseCase) Update(ctx context.Context, id string, input *entity.Template) error {
	input.ID = id
	// Validasi structure lagi jika perlu
	return uc.Repo.Update(ctx, input)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.Repo.Delete(ctx, id)
}
