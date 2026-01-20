package usecase

import (
	"context"
	"errors"

	"mutabaahapi/internal/mutabaah/entity"
	mutabaahRepo "mutabaahapi/internal/mutabaah/repository"
	userRepo "mutabaahapi/internal/users/repository"
)

type ListMutabaah struct {
	MutabaahRepo mutabaahRepo.MutabaahRepository
	UserRepo     userRepo.UserRepository
}

func NewListMutabaah(m mutabaahRepo.MutabaahRepository, u userRepo.UserRepository) *ListMutabaah {
	return &ListMutabaah{MutabaahRepo: m, UserRepo: u}
}

type ListFilter struct {
	UserID    string
	SessionID string
	StudentID string
}

func (uc *ListMutabaah) Execute(ctx context.Context, f ListFilter) ([]entity.Entry, error) {
	// 1. Cek User
	actor, err := uc.UserRepo.GetByID(ctx, f.UserID)
	if err != nil {
		return nil, err
	}
	if actor == nil {
		return nil, errors.New("user not found")
	}

	// 2. Logic Authorization & Filtering

	// Jika Actor adalah MURID, paksa filter StudentID ke ID dia sendiri
	if actor.Role == "MURID" {
		f.StudentID = actor.ID
	}

	// Skenario A: Filter by SessionID (Prioritas 1)
	if f.SessionID != "" {
		// Jika murid request by session, pastikan session itu miliknya (ideally check enrollment)
		// Tapi karena repo FindBySessionID mengembalikan semua entry di sesi itu,
		// kita perlu filter manual di code atau query khusus jika user adalah murid.
		// Sederhananya: Murid sebaiknya tidak pakai endpoint ini tanpa filter student_id.
		// Untuk keamanan, kita arahkan ke logic FindByStudentID di bawah jika SessionID kosong,
		// TAPI jika SessionID ada, Murid tetap harus difilter student_id nya.

		// [FIX] Jangan return langsung jika murid.
		if actor.Role != "MURID" {
			return uc.MutabaahRepo.FindBySessionID(ctx, f.SessionID)
		}
	}

	// Skenario B: Filter by StudentID (Prioritas 2 / Default Murid)
	targetStudent := f.StudentID

	if actor.Role == "GURU" {
		if targetStudent == "" && f.SessionID == "" {
			return nil, errors.New("teacher must provide session_id or student_id")
		}
	}

	// [FIX] Panggil FindByStudentID yang sudah ada
	entries, err := uc.MutabaahRepo.FindByStudentID(ctx, targetStudent)
	if err != nil {
		return nil, err
	}

	// Jika ada filter SessionID tambahan (misal murid mau spesifik sesi tertentu)
	if f.SessionID != "" {
		var filtered []entity.Entry
		for _, e := range entries {
			if e.SessionID == f.SessionID {
				filtered = append(filtered, e)
			}
		}
		return filtered, nil
	}

	return entries, nil
}
