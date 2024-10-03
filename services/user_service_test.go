// services/user_service_test.go
package services

import (
	"golang/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the IUserDao interface
type MockUserDao struct {
	mock.Mock
}

func (m *MockUserDao) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDao) GetByID(id uint64) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserDao) GetAll(page int, pageSize int, searchQuery string) ([]models.User, error) {
	args := m.Called(page, pageSize, searchQuery)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserDao) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDao) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

// FindByEmail implements dao.IUserDao.
func (m *MockUserDao) FindByEmail(userName string) (*models.User, error) {
	args := m.Called(userName)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_Create(t *testing.T) {
	mockDao := new(MockUserDao)
	userService := NewUserService(mockDao)

	user := &models.User{Username: "john", Password: "password"}

	mockDao.On("Create", user).Return(nil)

	err := userService.Create(user)

	assert.NoError(t, err)
	mockDao.AssertExpectations(t)
}

func TestUserService_GetByID(t *testing.T) {
	mockDao := new(MockUserDao)
	userService := NewUserService(mockDao)

	user := &models.User{ID: 1, Username: "john", Password: "password"}

	mockDao.On("GetByID", uint64(1)).Return(user, nil)

	fetchedUser, err := userService.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, user, fetchedUser)
	mockDao.AssertExpectations(t)
}

func TestUserService_GetAll(t *testing.T) {
	mockDao := new(MockUserDao)
	userService := NewUserService(mockDao)

	users := []models.User{
		{ID: 1, Username: "john", Password: "password"},
		{ID: 2, Username: "jane", Password: "password"},
	}

	mockDao.On("GetAll", 0, 10, "").Return(users, nil)

	fetchedUsers, err := userService.GetAll(0, 10, "")

	assert.NoError(t, err)
	assert.Equal(t, users, fetchedUsers)
	mockDao.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	mockDao := new(MockUserDao)
	userService := NewUserService(mockDao)

	user := &models.User{ID: 1, Username: "john", Password: "password"}

	mockDao.On("Update", user).Return(nil)

	err := userService.Update(user)

	assert.NoError(t, err)
	mockDao.AssertExpectations(t)
}

func TestUserService_Delete(t *testing.T) {
	mockDao := new(MockUserDao)
	userService := NewUserService(mockDao)

	mockDao.On("Delete", uint64(1)).Return(nil)

	err := userService.Delete(1)

	assert.NoError(t, err)
	mockDao.AssertExpectations(t)
}
