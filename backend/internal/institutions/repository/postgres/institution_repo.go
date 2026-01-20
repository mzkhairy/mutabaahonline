package postgres

import (
	"context"
	"database/sql"
	"errors"
	"mutabaahapi/internal/institutions/entity"
	domrepo "mutabaahapi/internal/institutions/repository"

	"github.com/jmoiron/sqlx"
)

type InstitutionRepository struct {
	db *sqlx.DB
}

func NewInstitutionRepository(db *sqlx.DB) *InstitutionRepository {
	return &InstitutionRepository{db: db}
}

// Pastikan struct ini memenuhi interface
var _ domrepo.InstitutionRepository = (*InstitutionRepository)(nil)

func (r *InstitutionRepository) GetByCode(ctx context.Context, code string) (*entity.Institution, error) {
	const q = `SELECT id, code, name, created_at, updated_at FROM institutions WHERE code = $1`
	var i entity.Institution
	if err := r.db.GetContext(ctx, &i, q, code); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil jika tidak ditemukan (bukan error system)
		}
		return nil, err
	}
	return &i, nil
}

func (r *InstitutionRepository) Create(ctx context.Context, i *entity.Institution) error {
	const q = `INSERT INTO institutions (code, name) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, q, i.Code, i.Name).Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
}
