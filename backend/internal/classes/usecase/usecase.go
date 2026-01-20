package usecase

import (
	"context"
	"errors"
	"mutabaahapi/internal/classes/entity"
	"mutabaahapi/internal/classes/repository"
	userRepo "mutabaahapi/internal/users/repository"
)

type UseCase struct {
	Repo     repository.ClassRepository
	UserRepo userRepo.UserRepository
}

func NewUseCase(r repository.ClassRepository, u userRepo.UserRepository) *UseCase {
	return &UseCase{Repo: r, UserRepo: u}
}

// [ADMIN] Create
func (uc *UseCase) Create(ctx context.Context, userID string, input *entity.Class) error {
	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || user.InstitutionID == nil {
		return errors.New("user institution not found")
	}

	input.InstitutionID = *user.InstitutionID

	// Validasi TeacherID (Pointer)
	if input.TeacherID != nil && *input.TeacherID != "" {
		teacher, err := uc.UserRepo.GetByID(ctx, *input.TeacherID)
		if err != nil {
			return err
		}
		if teacher == nil || teacher.Role != "GURU" {
			return errors.New("invalid teacher id")
		}
	}

	return uc.Repo.Create(ctx, input)
}

// [ADMIN] List
func (uc *UseCase) ListByInstitution(ctx context.Context, actorID string, academicYearID string) ([]entity.Class, error) {
	actor, err := uc.UserRepo.GetByID(ctx, actorID)
	if err != nil {
		return nil, err
	}
	if actor.InstitutionID == nil {
		return nil, errors.New("unauthorized")
	}
	return uc.Repo.ListByInstitution(ctx, *actor.InstitutionID, academicYearID)
}

// [GURU] List
func (uc *UseCase) ListByTeacher(ctx context.Context, teacherID string, academicYearID string) ([]entity.Class, error) {
	return uc.Repo.ListByTeacher(ctx, teacherID, academicYearID)
}

// [MURID] List
func (uc *UseCase) ListByStudent(ctx context.Context, studentID string) ([]entity.Class, error) {
	return uc.Repo.ListByStudent(ctx, studentID)
}

// Get Detail
func (uc *UseCase) GetByID(ctx context.Context, id string) (*entity.Class, error) {
	return uc.Repo.GetByID(ctx, id)
}

// Update
func (uc *UseCase) Update(ctx context.Context, id string, input *entity.Class) error {
	if input.TeacherID != nil && *input.TeacherID != "" {
		teacher, err := uc.UserRepo.GetByID(ctx, *input.TeacherID)
		if err != nil {
			return err
		}
		if teacher == nil || teacher.Role != "GURU" {
			return errors.New("invalid teacher id")
		}
	}
	input.ID = id
	return uc.Repo.Update(ctx, input)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.Repo.Delete(ctx, id)
}

// Student Management
func (uc *UseCase) AddStudent(ctx context.Context, classID, studentID string) error {
	return uc.Repo.AddStudent(ctx, classID, studentID)
}

func (uc *UseCase) RemoveStudent(ctx context.Context, classID, studentID string) error {
	return uc.Repo.RemoveStudent(ctx, classID, studentID)
}

func (uc *UseCase) ListStudents(ctx context.Context, classID string) ([]entity.StudentInClass, error) {
	return uc.Repo.GetStudents(ctx, classID)
}
