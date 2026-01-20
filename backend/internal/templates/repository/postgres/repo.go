package postgres

import (
	"context"
	"mutabaahapi/internal/templates/entity"

	"github.com/jmoiron/sqlx"
)

type Repository struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(ctx context.Context, t *entity.Template) error {
	const q = `
		INSERT INTO mutabaah_templates (institution_id, name, structure)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, q,
		t.InstitutionID, t.Name, t.Structure,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *Repository) List(ctx context.Context, institutionID string) ([]entity.Template, error) {
	const q = `SELECT * FROM mutabaah_templates WHERE institution_id = $1 ORDER BY created_at DESC`
	var list []entity.Template
	if err := r.db.SelectContext(ctx, &list, q, institutionID); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*entity.Template, error) {
	const q = `SELECT * FROM mutabaah_templates WHERE id = $1`
	var t entity.Template
	if err := r.db.GetContext(ctx, &t, q, id); err != nil {
		return nil, err
	}
	return &t, nil
}

// Update Template (Rename / Change Structure)
func (r *Repository) Update(ctx context.Context, t *entity.Template) error {
	const q = `
        UPDATE mutabaah_templates 
        SET name = $1, structure = $2, updated_at = NOW()
        WHERE id = $3
    `
	_, err := r.db.ExecContext(ctx, q, t.Name, t.Structure, t.ID)
	return err
}

// Delete Template
func (r *Repository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM mutabaah_templates WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}
