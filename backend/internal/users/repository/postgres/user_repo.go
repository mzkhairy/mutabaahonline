package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mutabaahapi/internal/users/entity"
	domrepo "mutabaahapi/internal/users/repository"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct{ db *sqlx.DB }

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ domrepo.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, u *entity.User) error {
	// [FIX] Gunakan placeholder $1, $2... karena kita mengirim argumen satu per satu
	const q = `
		INSERT INTO users (
			institution_id, role, username, name, email, 
			password_hash, serial_number, status, must_change_password
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9
		) RETURNING id, created_at, updated_at
	`

	// Default status ACTIVE jika kosong
	if u.Status == "" {
		u.Status = "ACTIVE"
	}

	return r.db.QueryRowxContext(ctx, q,
		u.InstitutionID,
		u.Role,
		u.Username,
		u.Name,
		u.Email,
		u.PasswordHash,
		u.SerialNumber, // Bisa nil
		u.Status,
		u.MustChangePassword, // [FIX] Field ini wajib dikirim sebagai $9
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) FindByUsername(ctx context.Context, institutionID string, username string) (*entity.User, error) {
	const q = `
		SELECT id, institution_id, role, username, name, email, password_hash, created_at, updated_at 
		FROM users 
		WHERE institution_id = $1 AND username = $2;
	`
	var u entity.User
	if err := r.db.GetContext(ctx, &u, q, institutionID, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	const q = `
		SELECT id, institution_id, role, username, name, email, password_hash, 
		       serial_number, status, created_at, updated_at, must_change_password
		FROM users WHERE id = $1;
	`
	var u entity.User
	if err := r.db.GetContext(ctx, &u, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) ListByInstitution(ctx context.Context, institutionID string, roleFilter string, search string) ([]entity.User, error) {
	query := `
		SELECT id, institution_id, role, username, name, email, serial_number, status, created_at, updated_at
		FROM users 
		WHERE institution_id = $1
	`
	args := []interface{}{institutionID}
	argIdx := 2

	if roleFilter != "" {
		query += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, roleFilter)
		argIdx++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR username ILIKE $%d)", argIdx, argIdx)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern)
		argIdx++
	}

	query += " ORDER BY name ASC"

	var users []entity.User
	if err := r.db.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetStats(ctx context.Context, institutionID string) (map[string]int, error) {
	const q = `
        SELECT 
            COUNT(*) FILTER (WHERE role = 'GURU' AND status = 'ACTIVE') as guru_active,
            COUNT(*) FILTER (WHERE role = 'MURID' AND status = 'ACTIVE') as murid_active,
            COUNT(*) as total_users
        FROM users 
        WHERE institution_id = $1
    `
	var stats struct {
		GuruActive  int `db:"guru_active"`
		MuridActive int `db:"murid_active"`
		TotalUsers  int `db:"total_users"`
	}

	if err := r.db.GetContext(ctx, &stats, q, institutionID); err != nil {
		return nil, err
	}

	return map[string]int{
		"guru_active":  stats.GuruActive,
		"murid_active": stats.MuridActive,
		"total_users":  stats.TotalUsers,
	}, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, hash string) error {
	// [FIX] Reset must_change_password jadi FALSE setelah update password
	const q = `
		UPDATE users 
		SET password_hash = $1, must_change_password = FALSE, updated_at = NOW() 
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, q, hash, userID)
	return err
}

func (r *UserRepository) UpdateStatus(ctx context.Context, userID string, status string) error {
	const q = `UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, status, userID)
	return err
}

func (r *UserRepository) CountByInstitution(ctx context.Context, institutionID string) (int, error) {
	const q = `SELECT COUNT(*) FROM users WHERE institution_id = $1`
	var count int
	err := r.db.GetContext(ctx, &count, q, institutionID)
	return count, err
}

func (r *UserRepository) FindBySerialNumber(ctx context.Context, institutionID string, serialNumber string) (*entity.User, error) {
	const q = `
		SELECT id, institution_id, role, username, name, email, password_hash, serial_number, status, created_at, updated_at 
		FROM users 
		WHERE institution_id = $1 AND serial_number = $2;
	`
	var u entity.User
	if err := r.db.GetContext(ctx, &u, q, institutionID, serialNumber); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *entity.User) error {
	const q = `
		UPDATE users 
		SET name = $1, username = $2, serial_number = $3, updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, q, u.Name, u.Username, u.SerialNumber, u.ID)
	return err
}

func (r *UserRepository) FindByIdentifier(ctx context.Context, institutionID string, identifier string) (*entity.User, error) {
	const q = `
		SELECT id, institution_id, role, username, name, email, 
		       password_hash, serial_number, status, must_change_password,
		       created_at, updated_at 
		FROM users 
		WHERE institution_id = $1 
		  AND (username = $2 OR serial_number = $2);
	`
	var u entity.User
	if err := r.db.GetContext(ctx, &u, q, institutionID, identifier); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Tidak ketemu
		}
		return nil, err
	}
	return &u, nil
}
