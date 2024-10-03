// controllers/user_controller_test.go
package controllers

import (
	"errors"
	"golang/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the IUserService interface
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetByID(id uint64) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetAll(page int, pageSize int, searchQuery string) ([]models.User, error) {
	args := m.Called(page, pageSize, searchQuery)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Delete(id uint64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserController_CreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	r.POST("/users", controller.CreateUser)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{Username: "john", Password: "password"}
		mockService.On("Create", user).Return(nil)

		jsonStr := `{"Username":"john","Password":"password"}`
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonStr))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request", func(t *testing.T) {
		jsonStr := `{"Username":"john"}`
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonStr))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		user := &models.User{Username: "john", Password: "password"}
		mockService.On("Create", user).Return(errors.New("error creating user"))

		jsonStr := `{"Username":"john","Password":"password"}`
		req, _ := http.NewRequest("POST", "/users", strings.NewReader(jsonStr))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserController_GetUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	r.GET("/users/:id", controller.GetUserById)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{ID: 1, Username: "john", Password: "password"}
		mockService.On("GetByID", uint64(1)).Return(user, nil)

		req, _ := http.NewRequest("GET", "/users/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockService.On("GetByID", uint64(999)).Return(nil, errors.New("user not found"))

		req, _ := http.NewRequest("GET", "/users/999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/abc", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_GetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	r.GET("/users", controller.GetAllUsers)

	t.Run("Success", func(t *testing.T) {
		users := []models.User{
			{ID: 1, Username: "john", Password: "password"},
			{ID: 2, Username: "jane", Password: "password"},
		}
		mockService.On("GetAll", 0, 10, "").Return(users, nil)

		req, _ := http.NewRequest("GET", "/users?page=1&pageSize=10", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockService.On("GetAll", 0, 10, "").Return(nil, errors.New("error fetching users"))

		req, _ := http.NewRequest("GET", "/users?page=1&pageSize=10", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserController_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	r.PUT("/users/:id", controller.UpdateUser)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{ID: 1, Username: "john", Password: "newpassword"}
		mockService.On("Update", user).Return(nil)

		jsonStr := `{"ID":1,"Username":"john","Password":"newpassword"}`
		req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(jsonStr))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Bad Request", func(t *testing.T) {
		jsonStr := `{"ID":1,"Username":"john"}`
		req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(jsonStr))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		user := &models.User{ID: 1, Username: "john", Password: "newpassword"}
		mockService.On("Update", user).Return(errors.New("error updating user"))

		jsonStr := `{"ID":1,"Username":"john","Password":"newpassword"}`
		req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(jsonStr))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserController_DeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockService := new(MockUserService)
	controller := NewUserController(mockService)

	r.DELETE("/users/:id", controller.DeleteUser)

	t.Run("Success", func(t *testing.T) {
		mockService.On("Delete", uint64(1)).Return(nil)

		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		mockService.On("Delete", uint64(999)).Return(errors.New("error deleting user"))

		req, _ := http.NewRequest("DELETE", "/users/999", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/abc", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
