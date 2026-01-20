package postgres

import (
	"context"
	"mutabaahapi/internal/mutabaah/entity"

	"github.com/jmoiron/sqlx"
)

type MutabaahRepository struct{ db *sqlx.DB }

func NewMutabaahRepository(db *sqlx.DB) *MutabaahRepository { return &MutabaahRepository{db: db} }

func (r *MutabaahRepository) Upsert(ctx context.Context, e *entity.Entry) error {
	const q = `
		INSERT INTO mutabaah_entries (
			session_id, student_id, attendance, student_activity_data, status, note, last_updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		ON CONFLICT (session_id, student_id) 
		DO UPDATE SET
			attendance = EXCLUDED.attendance,
			student_activity_data = EXCLUDED.student_activity_data,
			status = EXCLUDED.status,
			note = EXCLUDED.note,
			last_updated_by = EXCLUDED.last_updated_by,
			updated_at = NOW()
		RETURNING id, created_at, updated_at, is_locked
	`
	// Perhatikan urutan variabelnya ($1 sampai $7)
	return r.db.QueryRowxContext(ctx, q,
		e.SessionID,
		e.StudentID,
		e.Attendance,
		e.StudentActivityData, // Kirim JSONMap
		e.Status,
		e.Note,
		e.LastUpdatedBy,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt, &e.IsLocked)
}

func (r *MutabaahRepository) FindBySessionID(ctx context.Context, sessionID string) ([]entity.Entry, error) {
	// Kita JOIN ke tabel users (untuk nama murid) dan sessions (untuk nama sesi)
	const q = `
        SELECT 
            m.*, 
            u.name as student_name,
            s.name as session_name
        FROM mutabaah_entries m
        JOIN users u ON m.student_id = u.id
        JOIN sessions s ON m.session_id = s.id
        WHERE m.session_id = $1 
        ORDER BY m.created_at DESC
    `
	var entries []entity.Entry
	if err := r.db.SelectContext(ctx, &entries, q, sessionID); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *MutabaahRepository) FindByStudentID(ctx context.Context, studentID string) ([]entity.Entry, error) {
	const q = `
        SELECT 
            m.*, 
            u.name as student_name,
            s.name as session_name
        FROM mutabaah_entries m
        JOIN users u ON m.student_id = u.id
        JOIN sessions s ON m.session_id = s.id
        WHERE m.student_id = $1 
        ORDER BY m.created_at DESC
    `
	var entries []entity.Entry
	if err := r.db.SelectContext(ctx, &entries, q, studentID); err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *MutabaahRepository) GetSummaryByStudentID(ctx context.Context, studentID string) (map[string]int, error) {
	const q = `
		SELECT attendance, COUNT(*) as count
		FROM mutabaah_entries
		WHERE student_id = $1
		GROUP BY attendance
	`
	rows, err := r.db.QueryxContext(ctx, q, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := make(map[string]int)
	for rows.Next() {
		var attendance string
		var count int
		if err := rows.Scan(&attendance, &count); err != nil {
			continue
		}
		summary[attendance] = count
	}
	return summary, nil
}
