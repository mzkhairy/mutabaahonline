package postgres

import (
	"context"
	"database/sql"
	"errors"
	"mutabaahapi/internal/academicyears/entity"
	domrepo "mutabaahapi/internal/academicyears/repository"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

var _ domrepo.Repository = (*Repository)(nil)

func (r *Repository) Create(ctx context.Context, e *entity.AcademicYear) error {
	const q = `
		INSERT INTO academic_years (institution_id, name, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, q, e.InstitutionID, e.Name, e.StartDate, e.EndDate, e.IsActive).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func (r *Repository) Update(ctx context.Context, e *entity.AcademicYear) error {
	const q = `
		UPDATE academic_years
		SET name=$1, start_date=$2, end_date=$3, is_active=$4, updated_at=NOW()
		WHERE id=$5
	`
	_, err := r.db.ExecContext(ctx, q, e.Name, e.StartDate, e.EndDate, e.IsActive, e.ID)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM academic_years WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*entity.AcademicYear, error) {
	const q = `SELECT id, institution_id, name, start_date, end_date, is_active, created_at, updated_at FROM academic_years WHERE id=$1`
	var e entity.AcademicYear
	if err := r.db.GetContext(ctx, &e, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *Repository) List(ctx context.Context, institutionID string) ([]entity.AcademicYear, error) {
	const q = `
		SELECT id, institution_id, name, start_date, end_date, is_active, created_at, updated_at 
		FROM academic_years 
		WHERE institution_id = $1
		ORDER BY start_date DESC
	`
	var list []entity.AcademicYear
	if err := r.db.SelectContext(ctx, &list, q, institutionID); err != nil {
		return nil, err
	}
	return list, nil
}

// [BARU] Implementasi Reset
func (r *Repository) ResetStatusByInstitution(ctx context.Context, institutionID string) error {
	const q = `UPDATE academic_years SET is_active = false WHERE institution_id = $1`
	_, err := r.db.ExecContext(ctx, q, institutionID)
	return err
}
