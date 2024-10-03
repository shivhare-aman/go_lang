package dao

import (
	"golang/models"

	"gorm.io/gorm"
)

type IUserDao interface {
	Create(user *models.User) error
	GetByID(id uint64) (*models.User, error)
	GetAll(page int, pageSize int, searchQuery string) ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint64) error
	FindByEmail(email string) (*models.User, error)
}

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (u *UserDao) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *UserDao) GetByID(id uint64) (*models.User, error) {
	var user models.User
	err := u.db.Preload("Notes").Preload("CreditCard").First(&user, id).Error
	return &user, err
}

func (u *UserDao) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.db.First(&user, "userName = ?", email).Error
	return &user, err
}

// func (u *UserDao) GetAll() ([]models.User, error) {
// 	var users []models.User
// 	err := u.db.Preload("Notes").Preload("CreditCard").Find(&users).Error
// 	return users, err
// }

// func (u *UserDao) GetAll(offset int, pageSize int) ([]models.User, error) {
// 	var users []models.User
// 	err := u.db.Preload("Notes").Preload("CreditCard").Offset(offset).Limit(pageSize).Find(&users).Error
// 	return users, err
// }

func (u *UserDao) GetAll(offset int, pageSize int, searchQuery string) ([]models.User, error) {
	var users []models.User
	query := u.db.Preload("Notes").Preload("CreditCard")

	if searchQuery != "" {
		query = query.Where("Username ILIKE ?", "%"+searchQuery+"%")
	}

	err := query.Offset(offset).Limit(pageSize).Find(&users).Error

	return users, err
}

func (u *UserDao) Update(user *models.User) error {
	return u.db.Save(user).Error
}

func (u *UserDao) Delete(id uint64) error {
	return u.db.Delete(&models.User{}, id).Error
}
