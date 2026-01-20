package usecase

import (
	"context"
	"mutabaahapi/internal/mutabaah/repository"
	userRepo "mutabaahapi/internal/users/repository"
)

type MutabaahSummary struct {
	MutabaahRepo repository.MutabaahRepository
	UserRepo     userRepo.UserRepository
}

func NewMutabaahSummary(m repository.MutabaahRepository, u userRepo.UserRepository) *MutabaahSummary {
	return &MutabaahSummary{MutabaahRepo: m, UserRepo: u}
}

type SummaryOutput struct {
	TotalAttendance   int            `json:"total_attendance"`
	AttendanceDetails map[string]int `json:"attendance_details"`
}

func (uc *MutabaahSummary) Execute(ctx context.Context, userID string) (*SummaryOutput, error) {
	// 1. Cek User
	u, err := uc.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Jika Guru request summary, idealnya dia harus kirim student_id target (skip dulu untuk MVP)
	// Kita asumsikan endpoint ini khusus "My Summary" (Context User)

	// 2. Ambil Data
	stats, err := uc.MutabaahRepo.GetSummaryByStudentID(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	// 3. Hitung Total
	total := 0
	for _, v := range stats {
		total += v
	}

	return &SummaryOutput{
		TotalAttendance:   total,
		AttendanceDetails: stats,
	}, nil
}
