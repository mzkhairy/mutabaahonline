package usecase

import (
	"context"
	"errors"
	institutionsEntity "mutabaahapi/internal/institutions/entity"
	"mutabaahapi/internal/institutions/repository"
	"mutabaahapi/internal/users/entity"
	userRepo "mutabaahapi/internal/users/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	InstitutionCode string
	Username        string
	Password        string
	Name            string
	Email           string
	Role            string
	SerialNumber    string
}

type RegisterInstitutionInput struct {
	InstitutionName string
	InstitutionCode string // User yang tentukan kodenya
	Username        string
	Password        string
	Name            string // Nama Admin
	Email           string
}

type Register struct {
	UserRepo        userRepo.UserRepository
	InstitutionRepo repository.InstitutionRepository
}

func NewRegister(u userRepo.UserRepository, i repository.InstitutionRepository) *Register {
	return &Register{UserRepo: u, InstitutionRepo: i}
}

func (uc *Register) Execute(ctx context.Context, in RegisterInput) (*entity.User, error) {
	// 1. Cari Institution
	institution, err := uc.InstitutionRepo.GetByCode(ctx, in.InstitutionCode)
	if err != nil {
		return nil, err
	}
	if institution == nil {
		return nil, ErrInstitutionNotFound
	}

	// 2. Cek Username
	existing, err := uc.UserRepo.FindByUsername(ctx, institution.ID, in.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	// 3. Create
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var sn *string
	if in.SerialNumber != "" {
		sn = &in.SerialNumber
	}

	u := &entity.User{
		InstitutionID: &institution.ID,
		Role:          in.Role,
		Username:      in.Username,
		Name:          in.Name,
		Email:         in.Email,
		SerialNumber:  sn,
		PasswordHash:  string(hash),
	}

	if err := uc.UserRepo.Create(ctx, u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (uc *Register) RegisterInstitution(ctx context.Context, in RegisterInstitutionInput) (*entity.User, error) {
	// 1. Cek apakah Kode Institusi sudah dipakai?
	existingInst, err := uc.InstitutionRepo.GetByCode(ctx, in.InstitutionCode)
	if err != nil {
		return nil, err
	}
	if existingInst != nil {
		return nil, errors.New("institution code already exists")
	}

	// 2. Cek Username Admin (biar gak duplikat global/lokal)
	// Note: Karena belum ada ID institusi, kita cek username nanti setelah insert institusi,
	// atau asumsikan username harus unik per institusi nanti.
	// Untuk safety, kita lanjut dulu.

	// 3. Create Institution
	institution := &institutionsEntity.Institution{ // Pastikan import alias sesuai
		Code: in.InstitutionCode,
		Name: in.InstitutionName,
	}
	// Perlu import entity institutions, pastikan di import bagian atas:
	// "mutabaahapi/internal/institutions/entity"
	if err := uc.InstitutionRepo.Create(ctx, institution); err != nil {
		return nil, err
	}

	// 4. Create Admin User
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	admin := &entity.User{
		InstitutionID: &institution.ID, // Link ke institusi baru
		Role:          "ADMIN",         // Paksa jadi ADMIN
		Username:      in.Username,
		Name:          in.Name,
		Email:         in.Email,
		PasswordHash:  string(hash),
		Status:        "ACTIVE",
	}

	if err := uc.UserRepo.Create(ctx, admin); err != nil {
		return nil, err
	}

	admin.PasswordHash = ""
	return admin, nil
}
