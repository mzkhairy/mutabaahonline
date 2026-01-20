package usecase

import (
	"context"
	"errors"
	"mutabaahapi/internal/users/entity"
	"mutabaahapi/internal/users/repository"

	"golang.org/x/crypto/bcrypt"
)

type ManageUsers struct {
	UserRepo repository.UserRepository
}

func NewManageUsers(u repository.UserRepository) *ManageUsers {
	return &ManageUsers{UserRepo: u}
}

type UserListOutput struct {
	Stats struct {
		TotalRegistered int `json:"total_registered"`
		GuruActive      int `json:"guru_active"`
		MuridActive     int `json:"murid_active"`
		MaxAllowed      int `json:"max_allowed"`
	} `json:"stats"`
	Users []entity.User `json:"users"`
}

type CreateUserManualInput struct {
	Name         string
	Username     string
	SerialNumber string
	Role         string
	Password     string // Bisa kosong, kalau kosong pakai default
}

type UpdateUserInput struct {
	Name         string
	Username     string
	SerialNumber string
}

func (uc *ManageUsers) ExecuteList(ctx context.Context, actorID string, roleFilter string, search string) (*UserListOutput, error) {
	actor, err := uc.UserRepo.GetByID(ctx, actorID)
	if err != nil {
		return nil, err
	}
	if actor == nil || actor.InstitutionID == nil {
		return nil, errors.New("unauthorized")
	}

	users, err := uc.UserRepo.ListByInstitution(ctx, *actor.InstitutionID, roleFilter, search)
	if err != nil {
		return nil, err
	}

	stats, err := uc.UserRepo.GetStats(ctx, *actor.InstitutionID)

	guruActive := 0
	muridActive := 0
	totalUsers := 0

	if err == nil && stats != nil {
		guruActive = stats["guru_active"]
		muridActive = stats["murid_active"]
		totalUsers = stats["total_users"]
	}

	return &UserListOutput{
		Stats: struct {
			TotalRegistered int `json:"total_registered"`
			GuruActive      int `json:"guru_active"`
			MuridActive     int `json:"murid_active"`
			MaxAllowed      int `json:"max_allowed"`
		}{
			TotalRegistered: totalUsers,
			GuruActive:      guruActive,
			MuridActive:     muridActive,
			MaxAllowed:      100,
		},
		Users: users,
	}, nil
}

func (uc *ManageUsers) ResetPassword(ctx context.Context, adminID, targetUserID, newPassword string) error {
	admin, err := uc.UserRepo.GetByID(ctx, adminID)
	if err != nil || admin.Role != "ADMIN" {
		return errors.New("unauthorized")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return uc.UserRepo.UpdatePassword(ctx, targetUserID, string(hash))
}

func (uc *ManageUsers) ChangeStatus(ctx context.Context, userID, newStatus string) error {
	if newStatus != "ACTIVE" && newStatus != "INACTIVE" && newStatus != "CUTI" {
		return errors.New("invalid status value")
	}
	return uc.UserRepo.UpdateStatus(ctx, userID, newStatus)
}

// [PENTING] Method ini wajib ada
func (uc *ManageUsers) UpdateUser(ctx context.Context, adminID, targetUserID string, in UpdateUserInput) error {
	admin, err := uc.UserRepo.GetByID(ctx, adminID)
	if err != nil || admin.Role != "ADMIN" {
		return errors.New("unauthorized")
	}

	target, err := uc.UserRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("user not found")
	}

	target.Name = in.Name
	target.Username = in.Username

	if in.SerialNumber != "" {
		target.SerialNumber = &in.SerialNumber
	} else {
		target.SerialNumber = nil
	}

	return uc.UserRepo.Update(ctx, target)
}

// Method CreateUser (Manual by Admin)
func (uc *ManageUsers) CreateUser(ctx context.Context, adminID string, in CreateUserManualInput) error {
	// 1. Cek Admin
	admin, err := uc.UserRepo.GetByID(ctx, adminID)
	if err != nil || admin.Role != "ADMIN" || admin.InstitutionID == nil {
		return errors.New("unauthorized")
	}

	// 2. Validasi Username Unik
	exists, err := uc.UserRepo.FindByUsername(ctx, *admin.InstitutionID, in.Username)
	if exists != nil {
		return errors.New("username already exists")
	}

	// 3. Validasi NIP Unik (Jika ada)
	if in.SerialNumber != "" {
		existsNIP, err := uc.UserRepo.FindBySerialNumber(ctx, *admin.InstitutionID, in.SerialNumber)
		if err != nil {
			return err
		}
		if existsNIP != nil {
			return errors.New("serial number already exists")
		}
	}

	// 4. Set Password (Default: 123456 jika kosong)
	rawPass := in.Password
	if rawPass == "" {
		rawPass = "123456"
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(rawPass), bcrypt.DefaultCost)

	// 5. Create Entity
	var sn *string
	if in.SerialNumber != "" {
		sn = &in.SerialNumber
	}

	newUser := &entity.User{
		InstitutionID:      admin.InstitutionID,
		Name:               in.Name,
		Username:           in.Username,
		SerialNumber:       sn,
		Role:               in.Role,
		PasswordHash:       string(hash),
		Status:             "ACTIVE",
		MustChangePassword: true,
	}

	return uc.UserRepo.Create(ctx, newUser)
}

func (uc *ManageUsers) GetUser(ctx context.Context, actorID, targetUserID string) (*entity.User, error) {
	// 1. Cek Actor (yang me-request)
	actor, err := uc.UserRepo.GetByID(ctx, actorID)
	if err != nil {
		return nil, err
	}
	if actor == nil || actor.InstitutionID == nil {
		return nil, errors.New("unauthorized")
	}

	// 2. Ambil User Target
	target, err := uc.UserRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, errors.New("user not found")
	}

	//Cek Institusi
	// Pastikan hanya bisa lihat user di institusi yang sama
	if *target.InstitutionID != *actor.InstitutionID {
		return nil, errors.New("forbidden: different institution")
	}

	return target, nil
}
