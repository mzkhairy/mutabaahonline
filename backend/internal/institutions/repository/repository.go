package repository

import (
	"context"
	"mutabaahapi/internal/institutions/entity"
)

type InstitutionRepository interface {
	// Kita butuh ini untuk validasi saat Login/Register
	GetByCode(ctx context.Context, code string) (*entity.Institution, error)
	// Create dipakai nanti saat Admin Super mendaftarkan sekolah baru
	Create(ctx context.Context, i *entity.Institution) error
}
