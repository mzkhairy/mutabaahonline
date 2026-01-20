package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mutabaahapi/internal/classes/entity"
	domrepo "mutabaahapi/internal/classes/repository"

	"github.com/jmoiron/sqlx"
)

type Repository struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) *Repository { return &Repository{db: db} }

// Verifikasi Interface
var _ domrepo.ClassRepository = (*Repository)(nil)

func (r *Repository) Create(ctx context.Context, c *entity.Class) error {
	const q = `
		INSERT INTO classes (institution_id, academic_year_id, teacher_id, name, level)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, q,
		c.InstitutionID, c.AcademicYearID, c.TeacherID, c.Name, c.Level,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*entity.Class, error) {
	const q = `
		SELECT 
			c.*,
			u.name as teacher_name,
			ay.name as academic_year_name,
			(SELECT COUNT(*) FROM class_students cs WHERE cs.class_id = c.id) as student_count
		FROM classes c
		LEFT JOIN users u ON c.teacher_id = u.id
		LEFT JOIN academic_years ay ON c.academic_year_id = ay.id
		WHERE c.id = $1
	`
	var c entity.Class
	if err := r.db.GetContext(ctx, &c, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) Update(ctx context.Context, c *entity.Class) error {
	const q = `
        UPDATE classes 
        SET name = $1, level = $2, teacher_id = $3, updated_at = NOW()
        WHERE id = $4
    `
	_, err := r.db.ExecContext(ctx, q, c.Name, c.Level, c.TeacherID, c.ID)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM classes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

// [ADMIN] List
func (r *Repository) ListByInstitution(ctx context.Context, institutionID string, academicYearID string) ([]entity.Class, error) {
	q := `
		SELECT 
			c.id, c.institution_id, c.academic_year_id, c.teacher_id, c.name, c.level, c.created_at, c.updated_at,
			u.name as teacher_name,
			ay.name as academic_year_name,
			(SELECT COUNT(*) FROM class_students cs WHERE cs.class_id = c.id) as student_count
		FROM classes c
		LEFT JOIN users u ON c.teacher_id = u.id
		LEFT JOIN academic_years ay ON c.academic_year_id = ay.id
		WHERE c.institution_id = $1
	`
	args := []interface{}{institutionID}
	argIdx := 2

	if academicYearID != "" {
		q += fmt.Sprintf(" AND c.academic_year_id = $%d", argIdx)
		args = append(args, academicYearID)
	}

	q += " ORDER BY ay.start_date DESC, c.name ASC"

	var classes []entity.Class
	if err := r.db.SelectContext(ctx, &classes, q, args...); err != nil {
		return nil, err
	}
	return classes, nil
}

// [GURU] List
func (r *Repository) ListByTeacher(ctx context.Context, teacherID string, academicYearID string) ([]entity.Class, error) {
	q := `
		SELECT 
			c.id, c.institution_id, c.academic_year_id, c.teacher_id, c.name, c.level, c.created_at, c.updated_at,
			u.name as teacher_name,
			ay.name as academic_year_name,
			(SELECT COUNT(*) FROM class_students cs WHERE cs.class_id = c.id) as student_count
		FROM classes c
		LEFT JOIN users u ON c.teacher_id = u.id
		LEFT JOIN academic_years ay ON c.academic_year_id = ay.id
		WHERE c.teacher_id = $1
	`
	args := []interface{}{teacherID}
	argIdx := 2

	if academicYearID != "" {
		q += fmt.Sprintf(" AND c.academic_year_id = $%d", argIdx)
		args = append(args, academicYearID)
	} else {
		// Default: Hanya yang aktif
		q += " AND ay.is_active = true"
	}

	q += " ORDER BY c.name ASC"

	var classes []entity.Class
	if err := r.db.SelectContext(ctx, &classes, q, args...); err != nil {
		return nil, err
	}
	return classes, nil
}

// [MURID] List (Dropdown)
func (r *Repository) ListByStudent(ctx context.Context, studentID string) ([]entity.Class, error) {
	const q = `
		SELECT 
			c.id, c.institution_id, c.academic_year_id, c.teacher_id, c.name, c.level, c.created_at, c.updated_at,
			u.name as teacher_name,
			ay.name as academic_year_name,
			(SELECT COUNT(*) FROM class_students cs2 WHERE cs2.class_id = c.id) as student_count
		FROM classes c
		JOIN class_students cs ON c.id = cs.class_id
		LEFT JOIN users u ON c.teacher_id = u.id
		LEFT JOIN academic_years ay ON c.academic_year_id = ay.id
		WHERE cs.student_id = $1
		ORDER BY ay.start_date DESC, c.name ASC
	`
	var classes []entity.Class
	if err := r.db.SelectContext(ctx, &classes, q, studentID); err != nil {
		return nil, err
	}
	return classes, nil
}

// --- STUDENT MANAGEMENT ---

func (r *Repository) AddStudent(ctx context.Context, classID, studentID string) error {
	const q = `
		INSERT INTO class_students (class_id, student_id)
		VALUES ($1, $2)
		ON CONFLICT (class_id, student_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, q, classID, studentID)
	return err
}

func (r *Repository) RemoveStudent(ctx context.Context, classID, studentID string) error {
	const q = `DELETE FROM class_students WHERE class_id = $1 AND student_id = $2`
	_, err := r.db.ExecContext(ctx, q, classID, studentID)
	return err
}

func (r *Repository) GetStudents(ctx context.Context, classID string) ([]entity.StudentInClass, error) {
	const q = `
		SELECT u.id, u.name, COALESCE(u.serial_number, '') as serial_number, u.status
		FROM class_students cs
		JOIN users u ON cs.student_id = u.id
		WHERE cs.class_id = $1
		ORDER BY u.name ASC
	`
	var list []entity.StudentInClass
	if err := r.db.SelectContext(ctx, &list, q, classID); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) IsStudentInClass(ctx context.Context, classID, studentID string) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM class_students WHERE class_id = $1 AND student_id = $2)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, q, classID, studentID)
	return exists, err
}
