package usecase

import (
	"context"
	"errors"
	"mutabaahapi/internal/academicyears/entity"
	"mutabaahapi/internal/academicyears/repository"
	userRepo "mutabaahapi/internal/users/repository"
	"time"
)

type CreateInput struct {
	Name      string
	StartDate string // Input String
	EndDate   string // Input String
	IsActive  bool
}

type UpdateInput struct {
	Name      string
	StartDate string
	EndDate   string
	IsActive  bool
}

type UseCase struct {
	Repo     repository.Repository
	UserRepo userRepo.UserRepository
}

func NewUseCase(r repository.Repository, u userRepo.UserRepository) *UseCase {
	return &UseCase{Repo: r, UserRepo: u}
}

func (uc *UseCase) List(ctx context.Context, actorID string) ([]entity.AcademicYear, error) {
	actor, err := uc.UserRepo.GetByID(ctx, actorID)
	if err != nil || actor == nil || actor.InstitutionID == nil {
		return nil, errors.New("unauthorized")
	}
	return uc.Repo.List(ctx, *actor.InstitutionID)
}

func (uc *UseCase) Create(ctx context.Context, actorID string, in CreateInput) error {
	actor, err := uc.UserRepo.GetByID(ctx, actorID)
	if err != nil || actor == nil || actor.InstitutionID == nil {
		return errors.New("unauthorized")
	}

	// [FIX] Validasi format tanggal saja (biar ga error DB), tapi datanya tetap String
	if _, err := time.Parse("2006-01-02", in.StartDate); err != nil {
		return errors.New("invalid start_date format (YYYY-MM-DD)")
	}
	if _, err := time.Parse("2006-01-02", in.EndDate); err != nil {
		return errors.New("invalid end_date format (YYYY-MM-DD)")
	}

	// [FIX LOGIC] Jika ini di-set aktif, matikan yang lain dulu
	if in.IsActive {
		if err := uc.Repo.ResetStatusByInstitution(ctx, *actor.InstitutionID); err != nil {
			return err
		}
	}

	e := &entity.AcademicYear{
		InstitutionID: *actor.InstitutionID,
		Name:          in.Name,
		StartDate:     in.StartDate, // Tetap string (sesuai entity Anda)
		EndDate:       in.EndDate,   // Tetap string (sesuai entity Anda)
		IsActive:      in.IsActive,
	}
	return uc.Repo.Create(ctx, e)
}

func (uc *UseCase) Update(ctx context.Context, actorID string, id string, in UpdateInput) error {
	actor, err := uc.UserRepo.GetByID(ctx, actorID)
	if err != nil || actor == nil || actor.InstitutionID == nil {
		return errors.New("unauthorized")
	}

	existing, err := uc.Repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("academic year not found")
	}

	if existing.InstitutionID != *actor.InstitutionID {
		return errors.New("forbidden")
	}

	// [FIX] Validasi format saja
	if _, err := time.Parse("2006-01-02", in.StartDate); err != nil {
		return errors.New("invalid start_date format")
	}
	if _, err := time.Parse("2006-01-02", in.EndDate); err != nil {
		return errors.New("invalid end_date format")
	}

	// [FIX LOGIC] Reset jika diubah jadi aktif
	if in.IsActive {
		if err := uc.Repo.ResetStatusByInstitution(ctx, *actor.InstitutionID); err != nil {
			return err
		}
	}

	existing.Name = in.Name
	existing.StartDate = in.StartDate // Tetap string
	existing.EndDate = in.EndDate     // Tetap string
	existing.IsActive = in.IsActive

	return uc.Repo.Update(ctx, existing)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.Repo.Delete(ctx, id)
}
