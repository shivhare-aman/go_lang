package services

import (
	"golang/dao"
	"golang/models"
)

type IUserService interface {
	Create(user *models.User) error
	GetByID(id uint64) (*models.User, error)
	GetAll(page int, pageSize int, searchQuery string) ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint64) error
}

type UserService struct {
	userDao dao.IUserDao
}

func NewUserService(userDao dao.IUserDao) *UserService {
	return &UserService{userDao: userDao}
}

func (u *UserService) Create(user *models.User) error {
	return u.userDao.Create(user)
}

func (u *UserService) GetByID(id uint64) (*models.User, error) {
	return u.userDao.GetByID(id)
}

// func (u *UserService) GetAll() ([]models.User, error) {
// 	users, err := u.userDao.GetAll()
// 	return users, err
// }

// func (u *UserService) GetAll(page int, pageSize int) ([]models.User, error) {
// 	offset := (page - 1) * pageSize

// 	users, err := u.userDao.GetAll(offset, pageSize)
// 	return users, err
// }

func (u *UserService) GetAll(page int, pageSize int, searchQuery string) ([]models.User, error) {
	offset := (page - 1) * pageSize

	users, err := u.userDao.GetAll(offset, pageSize, searchQuery)
	return users, err
}

func (u *UserService) Update(user *models.User) error {
	return u.userDao.Update(user)
}

func (u *UserService) Delete(id uint64) error {
	return u.userDao.Delete(id)
}
