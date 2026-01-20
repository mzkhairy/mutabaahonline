package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mutabaahapi/internal/sessions/entity"

	"github.com/jmoiron/sqlx"
)

type SessionRepository struct{ db *sqlx.DB }

func NewSessionRepository(db *sqlx.DB) *SessionRepository { return &SessionRepository{db: db} }

func (r *SessionRepository) Create(ctx context.Context, s *entity.Session) error {
	const q = `
		INSERT INTO sessions (class_id, template_id, date, name, class_activity_data, is_student_input_allowed)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, q,
		s.ClassID, s.TemplateID, s.Date, s.Name, s.ClassActivityData, s.IsStudentInputAllowed,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SessionRepository) IsClassTeacher(ctx context.Context, classID, teacherID string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM classes WHERE id = $1 AND teacher_id = $2)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, q, classID, teacherID)
	return exists, err
}

func (r *SessionRepository) GetByID(ctx context.Context, id string) (*entity.Session, error) {
	const q = `
		SELECT 
			s.id, s.class_id, s.template_id, s.date, s.name, s.class_activity_data, 
			s.is_student_input_allowed, s.created_at, s.updated_at,
			c.teacher_id,
			t.name as template_name,
			t.structure as template_structure
		FROM sessions s
		JOIN classes c ON s.class_id = c.id
		LEFT JOIN mutabaah_templates t ON s.template_id = t.id
		WHERE s.id = $1
	`
	var s entity.Session
	if err := r.db.GetContext(ctx, &s, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *SessionRepository) ListByClass(ctx context.Context, classID string) ([]entity.Session, error) {
	const q = `
		SELECT 
			s.id, s.class_id, s.template_id, s.date, s.name, s.class_activity_data, 
			s.is_student_input_allowed, s.created_at, s.updated_at,
			t.name as template_name
		FROM sessions s
		LEFT JOIN mutabaah_templates t ON s.template_id = t.id
		WHERE s.class_id = $1 
		ORDER BY s.date DESC, s.created_at DESC
	`
	var list []entity.Session
	if err := r.db.SelectContext(ctx, &list, q, classID); err != nil {
		return nil, err
	}
	return list, nil
}

// [FIXED] Query ini sekarang sangat spesifik kolomnya agar tidak error 500
func (r *SessionRepository) ListByStudent(ctx context.Context, studentID string, isReportMode bool, classID string) ([]entity.Session, error) {
	q := `
		SELECT 
			s.id, s.class_id, s.template_id, s.name, s.date, s.is_student_input_allowed, 
			s.created_at, s.updated_at, s.class_activity_data,
			c.name as class_name,
			t.name as template_name, 
			t.structure as template_structure
		FROM sessions s
		JOIN classes c ON s.class_id = c.id
		LEFT JOIN mutabaah_templates t ON s.template_id = t.id 
		LEFT JOIN academic_years ay ON c.academic_year_id = ay.id
		WHERE (
            s.class_id IN (SELECT class_id FROM class_students WHERE student_id = $1)
		    OR 
            s.class_id IN (SELECT id FROM classes WHERE teacher_id = $1)
        )
	`
	args := []interface{}{studentID}
	argIdx := 2

	if classID != "" {
		q += fmt.Sprintf(" AND s.class_id = $%d", argIdx)
		args = append(args, classID)
		argIdx++
	}

	if !isReportMode {
		q += " AND ay.is_active = true"
	}

	q += ` ORDER BY s.date DESC`

	var sessions []entity.Session
	if err := r.db.SelectContext(ctx, &sessions, q, args...); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) Update(ctx context.Context, s *entity.Session) error {
	const q = `
		UPDATE sessions 
		SET 
			class_activity_data = $1,
			name = $2,
			date = $3,
			is_student_input_allowed = $4,
			updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.db.ExecContext(ctx, q,
		s.ClassActivityData,
		s.Name,
		s.Date,
		s.IsStudentInputAllowed,
		s.ID,
	)
	return err
}
