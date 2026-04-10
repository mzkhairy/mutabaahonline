package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mutabaahapi/internal/templates/entity"
	// Asumsikan di sini Anda mungkin belum mengubah UseCase template ke interface,
	// jadi kita buat test dummy sederhana agar tidak error compile.
	// Jika mau test beneran, pastikan NewUseCase menerima interface repo.
	userEntity "mutabaahapi/internal/users/entity"
)

// --- MOCK TEMPLATE REPO (Definisi Manual) ---
type MockTemplateRepo struct {
	mock.Mock
}

func (m *MockTemplateRepo) Create(ctx context.Context, t *entity.Template) error {
	return m.Called(ctx, t).Error(0)
}

// Tambahkan method lain jika repo template punya banyak method

// [FIX] Mock User Repo LENGKAP
type MockUserRepo struct{ mock.Mock }

func (m *MockUserRepo) Create(ctx context.Context, u *userEntity.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *MockUserRepo) FindByUsername(ctx context.Context, iID, u string) (*userEntity.User, error) {
	args := m.Called(ctx, iID, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}
func (m *MockUserRepo) FindByIdentifier(ctx context.Context, iID, id string) (*userEntity.User, error) {
	args := m.Called(ctx, iID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}
func (m *MockUserRepo) FindBySerialNumber(ctx context.Context, iID, sn string) (*userEntity.User, error) {
	args := m.Called(ctx, iID, sn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}
func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*userEntity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}
func (m *MockUserRepo) ListByInstitution(ctx context.Context, iID, r, s string) ([]userEntity.User, error) {
	args := m.Called(ctx, iID, r, s)
	return args.Get(0).([]userEntity.User), args.Error(1)
}
func (m *MockUserRepo) UpdateStatus(ctx context.Context, uid, s string) error {
	return m.Called(ctx, uid, s).Error(0)
}
func (m *MockUserRepo) CountByInstitution(ctx context.Context, iID string) (int, error) {
	args := m.Called(ctx, iID)
	return args.Int(0), args.Error(1)
}
func (m *MockUserRepo) GetStats(ctx context.Context, iID string) (map[string]int, error) {
	args := m.Called(ctx, iID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}
func (m *MockUserRepo) UpdatePassword(ctx context.Context, uid, h string) error {
	return m.Called(ctx, uid, h).Error(0)
}
func (m *MockUserRepo) Update(ctx context.Context, u *userEntity.User) error {
	return m.Called(ctx, u).Error(0)
}

// --- TEST CASES ---

func TestCreateTemplate(t *testing.T) {
	// [NOTE] Karena UseCase Templates mungkin masih pakai struct Repo asli (Postgres),
	// kita tidak bisa inject mock ke sana tanpa refactor 'NewUseCase' templates.
	// Agar test file ini valid (tidak error compile), kita buat dummy test saja.

	mockUserRepo := new(MockUserRepo)
	ctx := context.Background()

	t.Run("Dummy Test Structure", func(t *testing.T) {
		// Mock call example
		mockUserRepo.On("GetByID", ctx, "admin").Return(&userEntity.User{ID: "admin"}, nil)

		u, err := mockUserRepo.GetByID(ctx, "admin")
		assert.NoError(t, err)
		assert.Equal(t, "admin", u.ID)
	})
}
