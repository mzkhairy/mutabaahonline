package usecase

import (
	"context"
	"errors"

	classRepo "mutabaahapi/internal/classes/repository"
	"mutabaahapi/internal/mutabaah/entity"
	mutabaahRepo "mutabaahapi/internal/mutabaah/repository"
	sessRepo "mutabaahapi/internal/sessions/repository"
	userRepo "mutabaahapi/internal/users/repository"
)

type InputMutabaah struct {
	MutabaahRepo mutabaahRepo.MutabaahRepository
	SessionRepo  sessRepo.SessionRepository
	UserRepo     userRepo.UserRepository
	ClassRepo    classRepo.ClassRepository // [BARU] Tambah ini
}

func NewInputMutabaah(
	m mutabaahRepo.MutabaahRepository,
	s sessRepo.SessionRepository,
	u userRepo.UserRepository,
	c classRepo.ClassRepository, // [BARU] Tambah parameter
) *InputMutabaah {
	return &InputMutabaah{MutabaahRepo: m, SessionRepo: s, UserRepo: u, ClassRepo: c}
}

type InputData struct {
	UserID              string
	SessionID           string
	StudentID           string
	Attendance          string
	StudentActivityData map[string]interface{}
	Status              string
	Note                string
}

func (uc *InputMutabaah) Execute(ctx context.Context, in InputData) (*entity.Entry, error) {
	// 1. Cek User
	actor, err := uc.UserRepo.GetByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	if actor == nil {
		return nil, errors.New("user not found")
	}

	// 2. Cek Session
	session, err := uc.SessionRepo.GetByID(ctx, in.SessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found")
	}

	// 3. Authorization Logic
	targetStudentID := in.StudentID

	if actor.Role == "MURID" {
		// Rule A: Murid hanya boleh input untuk dirinya sendiri
		if targetStudentID != "" && targetStudentID != actor.ID {
			return nil, errors.New("forbidden: students can only input for themselves")
		}
		targetStudentID = actor.ID

		// Rule B: Sesi harus mengizinkan input murid
		if !session.IsStudentInputAllowed {
			return nil, errors.New("student input is currently closed/locked for this session")
		}

		// Rule C: Murid harus anggota kelas tersebut
		isMember, err := uc.ClassRepo.IsStudentInClass(ctx, session.ClassID, actor.ID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New("forbidden: you are not a member of this class")
		}

		// Rule D: Cek apakah data sudah dikunci (Approved/Locked)
		existingEntries, err := uc.MutabaahRepo.FindByStudentID(ctx, targetStudentID)
		if err == nil {
			for _, e := range existingEntries {
				if e.SessionID == in.SessionID && e.IsLocked {
					return nil, errors.New("data has been locked by teacher and cannot be edited")
				}
			}
		}

	} else if actor.Role == "GURU" {
		// Validasi Guru Pemilik Kelas
		// if session.TeacherID == nil || *session.TeacherID != actor.ID {
		// 	return nil, errors.New("you are not the teacher of this class session")
		// }
		// (Opsional: dinonaktifkan dulu agar Admin/Guru lain bisa bantu input jika perlu)

		if targetStudentID == "" {
			return nil, errors.New("student_id is required for teacher input")
		}
	}

	// 4. Save
	entry := &entity.Entry{
		SessionID:           in.SessionID,
		StudentID:           targetStudentID,
		Attendance:          in.Attendance,
		StudentActivityData: in.StudentActivityData,
		Status:              in.Status,
		Note:                in.Note,
		LastUpdatedBy:       actor.ID,
	}

	// Status default jika kosong
	if entry.Status == "" {
		entry.Status = "LANJUT"
	}

	if err := uc.MutabaahRepo.Upsert(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}
