package repository

import (
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	Database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Database: db}
}

// CreateUser adds new user to the database
func (ur *UserRepository) CreateUser(user *models.User) error {
	tx := ur.Database.Create(user)
	return tx.Error
}

// GetUserByID returns user record by the provided id
func (ur *UserRepository) GetUserByID(userID uint) (*models.User, error) {
	var user models.User

	err := ur.Database.
		Where("id = ?", userID).
		Preload(clause.Associations).
		Find(&user).
		Error
	return &user, err
}

// GetUserByEmail returns user record by the provided email
func (ur *UserRepository) GetUserByEmail(userEmail string) (*models.User, error) {
	var user models.User

	err := ur.Database.Where("email = ?", userEmail).Preload(clause.Associations).Find(&user).Error
	return &user, err
}

// UserEmailExists returns whether the user exists or not
func (ur *UserRepository) UserEmailExists(userEmail string) (bool, error) {
	var (
		user models.User
	)
	user.Email = userEmail

	r := ur.Database.Model(&models.User{}).
		Where("email = ?", userEmail).
		Find(&user).
		Limit(1)
	if r.Error != nil {
		return false, r.Error
	}
	return r.RowsAffected > 0, nil
}

// UserExists returns whether the user exists or not
func (ur *UserRepository) UserExists(userID uint) (bool, error) {
	var (
		user models.User
	)
	r := ur.Database.Model(&models.User{}).
		Where("id = ?", userID).
		Find(&user).
		Limit(1)

	if r.Error != nil {
		return false, r.Error
	}
	return r.RowsAffected > 0, nil
}

// UpdateUser sets param of the user with provided userID to value
func (ur *UserRepository) UpdateUser(user *models.User, param, value string) error {
	err := ur.Database.Model(user).Update(param, value).Error
	return err
}

// GetLastID returns the ID of the last created user
func (ur *UserRepository) GetLastID() (uint, error) {
	var user models.User
	err := ur.Database.Last(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return uint(0), nil
	}
	if err != nil {
		return uint(0), err
	}

	return user.ID, nil
}

// DeleteUser deletes user with provided id
func (ur *UserRepository) DeleteUser(user *models.User) error {
	return ur.Database.Delete(user).Error
}

func (ur *UserRepository) GetAllUsers(callerID uint) ([]*models.User, error) {
	var users []*models.User
	r := ur.Database.Model(models.User{}).
		Where("users.id != ?", callerID)

	err := r.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
