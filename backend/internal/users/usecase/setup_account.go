package usecase

import (
	"context"
	"errors"
	"mutabaahapi/internal/users/repository"

	"golang.org/x/crypto/bcrypt"
)

type SetupAccount struct {
	UserRepo repository.UserRepository
}

func NewSetupAccount(r repository.UserRepository) *SetupAccount {
	return &SetupAccount{UserRepo: r}
}

type SetupInput struct {
	NewUsername string
	NewPassword string
	Email       string
}

func (uc *SetupAccount) Execute(ctx context.Context, userID string, in SetupInput) error {
	// 1. Ambil User
	user, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 2. Validasi Username Baru (Cek Unik)
	if user.InstitutionID != nil {
		existing, _ := uc.UserRepo.FindByUsername(ctx, *user.InstitutionID, in.NewUsername)
		if existing != nil && existing.ID != user.ID {
			return errors.New("username already taken")
		}
	}

	// 3. Hash Password
	hash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 4. Update Entity
	user.Username = in.NewUsername
	user.Email = in.Email
	user.PasswordHash = string(hash)
	user.MustChangePassword = false
	if err := uc.UserRepo.Update(ctx, user); err != nil {
		return err
	}
	// B. Update Password & Flag
	return uc.UserRepo.UpdatePassword(ctx, user.ID, string(hash))
}
