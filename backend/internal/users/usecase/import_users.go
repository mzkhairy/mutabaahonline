package usecase

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mutabaahapi/internal/users/entity"
	"mutabaahapi/internal/users/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type ImportUsers struct {
	UserRepo repository.UserRepository
}

func NewImportUsers(u repository.UserRepository) *ImportUsers {
	return &ImportUsers{UserRepo: u}
}

type ImportResult struct {
	TotalProcessed int      `json:"total_processed"`
	SuccessCount   int      `json:"success_count"`
	SkippedCount   int      `json:"skipped_count"`
	SkippedDetails []string `json:"skipped_details"`
}

// [PENTING] Method Wrapper ini wajib ada
func (uc *ImportUsers) ExecuteWrapper(ctx context.Context, adminID string, fileReader io.Reader) (*ImportResult, error) {
	admin, err := uc.UserRepo.GetByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if admin == nil || admin.InstitutionID == nil || admin.Role != "ADMIN" {
		return nil, errors.New("unauthorized import")
	}
	return uc.Execute(ctx, *admin.InstitutionID, fileReader)
}

func (uc *ImportUsers) Execute(ctx context.Context, institutionID string, fileReader io.Reader) (*ImportResult, error) {
	reader := csv.NewReader(fileReader)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("failed to parse CSV file")
	}

	result := &ImportResult{
		SkippedDetails: []string{},
	}

	defaultPassHash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	for i, row := range records {
		if i == 0 {
			continue
		} // Skip Header

		result.TotalProcessed++

		if len(row) < 3 {
			result.SkippedCount++
			result.SkippedDetails = append(result.SkippedDetails, fmt.Sprintf("Baris %d: Format kolom tidak lengkap", i+1))
			continue
		}

		name := strings.TrimSpace(row[0])
		serialNumber := strings.TrimSpace(row[1])
		role := strings.ToUpper(strings.TrimSpace(row[2]))

		if name == "" || serialNumber == "" || role == "" {
			result.SkippedCount++
			result.SkippedDetails = append(result.SkippedDetails, fmt.Sprintf("Baris %d: Data wajib ada yang kosong", i+1))
			continue
		}

		if role != "GURU" && role != "MURID" {
			result.SkippedCount++
			result.SkippedDetails = append(result.SkippedDetails, fmt.Sprintf("Baris %d: Role '%s' tidak valid", i+1, role))
			continue
		}

		existing, err := uc.UserRepo.FindBySerialNumber(ctx, institutionID, serialNumber)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			result.SkippedCount++
			result.SkippedDetails = append(result.SkippedDetails, fmt.Sprintf("%s (%s) - Nomor Induk sudah terdaftar", serialNumber, name))
			continue
		}

		username := serialNumber
		existingUser, _ := uc.UserRepo.FindByUsername(ctx, institutionID, username)
		if existingUser != nil {
			result.SkippedCount++
			result.SkippedDetails = append(result.SkippedDetails, fmt.Sprintf("%s (%s) - Username konflik", serialNumber, name))
			continue
		}

		newUser := &entity.User{
			InstitutionID:      &institutionID,
			Name:               name,
			Username:           username,
			SerialNumber:       &serialNumber,
			Role:               role,
			Email:              "",
			PasswordHash:       string(defaultPassHash),
			Status:             "ACTIVE",
			MustChangePassword: true,
		}

		if err := uc.UserRepo.Create(ctx, newUser); err != nil {
			result.SkippedCount++
			result.SkippedDetails = append(result.SkippedDetails, fmt.Sprintf("%s (%s) - Gagal simpan DB", serialNumber, name))
			continue
		}

		result.SuccessCount++
	}

	return result, nil
}
