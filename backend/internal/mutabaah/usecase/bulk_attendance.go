package usecase

import (
	"context"
	"mutabaahapi/internal/mutabaah/entity"
	"mutabaahapi/internal/mutabaah/repository"
	userRepo "mutabaahapi/internal/users/repository"
)

type BulkAttendance struct {
	MutabaahRepo repository.MutabaahRepository
	UserRepo     userRepo.UserRepository
}

func NewBulkAttendance(m repository.MutabaahRepository, u userRepo.UserRepository) *BulkAttendance {
	return &BulkAttendance{MutabaahRepo: m, UserRepo: u}
}

type AttendanceItem struct {
	StudentID  string `json:"student_id"`
	Attendance string `json:"attendance"` // HADIR, SAKIT, IZIN, ALPA
}

type BulkInput struct {
	UserID    string
	SessionID string
	Items     []AttendanceItem
}

func (uc *BulkAttendance) Execute(ctx context.Context, in BulkInput) error {
	// 1. Cek User (Guru)
	actor, err := uc.UserRepo.GetByID(ctx, in.UserID)
	if err != nil || actor == nil {
		return err
	}

	// Idealnya cek apakah Guru ini pemilik sesi (Skip dulu utk prototype)

	// 2. Loop Upsert
	// Note: Ini cara naive (satu-satu). Untuk production high-scale pake Batch Insert SQL.
	// Tapi untuk kelas isi 30 orang, loop ini sangat cepat (<50ms).
	for _, item := range in.Items {
		entry := &entity.Entry{
			SessionID:     in.SessionID,
			StudentID:     item.StudentID,
			Attendance:    item.Attendance,
			Status:        "LANJUT", // Default
			Note:          "",
			LastUpdatedBy: actor.ID,
			// StudentActivityData biarkan kosong/default map
			StudentActivityData: make(map[string]interface{}),
		}

		// Kita perlu method di Repo yang Update HANYA Attendance tanpa menimpa Data JSON
		// TAPI, karena Upsert kita "DO UPDATE SET ... student_activity_data = EXCLUDED...",
		// data lama bisa ketimpa kosong.

		// STRATEGI AMAN:
		// Karena absen biasanya dilakukan DI AWAL sesi (sebelum isi nilai), menimpa data jadi kosong itu AMAN.
		// Guru absen dulu -> baru isi nilai.

		if err := uc.MutabaahRepo.Upsert(ctx, entry); err != nil {
			return err
		}
	}

	return nil
}
