package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"mutabaahapi/internal/classes/entity"
	"mutabaahapi/internal/classes/usecase"
	userEntity "mutabaahapi/internal/users/entity"
)

// --- MOCKS ---

type MockClassRepo struct {
	mock.Mock
}

func (m *MockClassRepo) Create(ctx context.Context, c *entity.Class) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockClassRepo) AddStudent(ctx context.Context, classID, studentID string) error {
	return m.Called(ctx, classID, studentID).Error(0)
}
func (m *MockClassRepo) GetByID(ctx context.Context, id string) (*entity.Class, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Class), args.Error(1)
}
func (m *MockClassRepo) Update(ctx context.Context, c *entity.Class) error {
	return m.Called(ctx, c).Error(0)
}
func (m *MockClassRepo) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockClassRepo) ListByInstitution(ctx context.Context, iID, yID string) ([]entity.Class, error) {
	args := m.Called(ctx, iID, yID)
	return args.Get(0).([]entity.Class), args.Error(1)
}
func (m *MockClassRepo) ListByTeacher(ctx context.Context, tID, yID string) ([]entity.Class, error) {
	args := m.Called(ctx, tID, yID)
	return args.Get(0).([]entity.Class), args.Error(1)
}
func (m *MockClassRepo) ListByStudent(ctx context.Context, sID string) ([]entity.Class, error) {
	args := m.Called(ctx, sID)
	return args.Get(0).([]entity.Class), args.Error(1)
}
func (m *MockClassRepo) RemoveStudent(ctx context.Context, cID, sID string) error {
	return m.Called(ctx, cID, sID).Error(0)
}
func (m *MockClassRepo) GetStudents(ctx context.Context, cID string) ([]entity.StudentInClass, error) {
	args := m.Called(ctx, cID)
	return args.Get(0).([]entity.StudentInClass), args.Error(1)
}
func (m *MockClassRepo) IsStudentInClass(ctx context.Context, cID, sID string) (bool, error) {
	args := m.Called(ctx, cID, sID)
	return args.Bool(0), args.Error(1)
}

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

func TestCreateClass(t *testing.T) {
	mockRepo := new(MockClassRepo)
	mockUserRepo := new(MockUserRepo)
	uCase := usecase.NewUseCase(mockRepo, mockUserRepo)

	ctx := context.Background()
	instID := "inst-1"
	adminUser := &userEntity.User{ID: "admin", InstitutionID: &instID}

	teacherID := "guru-1"
	teacherUser := &userEntity.User{ID: teacherID, Role: "GURU"}

	t.Run("Success Create Class with Teacher", func(t *testing.T) {
		input := &entity.Class{
			Name:      "Kelas 7A",
			Level:     "7",
			TeacherID: &teacherID,
		}

		mockUserRepo.On("GetByID", ctx, "admin").Return(adminUser, nil).Once()
		mockUserRepo.On("GetByID", ctx, teacherID).Return(teacherUser, nil).Once()

		mockRepo.On("Create", ctx, mock.MatchedBy(func(c *entity.Class) bool {
			return c.Name == "Kelas 7A" && c.InstitutionID == instID
		})).Return(nil).Once()

		err := uCase.Create(ctx, "admin", input)
		assert.NoError(t, err)
	})

	t.Run("Fail Invalid Teacher", func(t *testing.T) {
		invalidID := "murid-1"
		muridUser := &userEntity.User{ID: invalidID, Role: "MURID"}

		input := &entity.Class{Name: "X", TeacherID: &invalidID}

		mockUserRepo.On("GetByID", ctx, "admin").Return(adminUser, nil).Once()
		mockUserRepo.On("GetByID", ctx, invalidID).Return(muridUser, nil).Once()

		err := uCase.Create(ctx, "admin", input)
		assert.Error(t, err)
		assert.Equal(t, "invalid teacher id", err.Error())
	})
}

func TestAssignStudent(t *testing.T) {
	mockRepo := new(MockClassRepo)
	mockUserRepo := new(MockUserRepo)
	uCase := usecase.NewUseCase(mockRepo, mockUserRepo)

	t.Run("Success Add Student", func(t *testing.T) {
		mockRepo.On("AddStudent", context.Background(), "class-1", "student-1").Return(nil)

		err := uCase.AddStudent(context.Background(), "class-1", "student-1")
		assert.NoError(t, err)
	})
}
