package services

import (
	"errors"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/clierrs"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/enums"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/interfaces"
	"github.com/cristalhq/jwt/v5"
	"gorm.io/gorm"
	"slices"
)

var notAllowedFieldsToUpdate = []string{
	"Password",
	"Email",
}

type UserService struct {
	UserRepo interfaces.UserManager
	Auth     interfaces.AuthManager
}

func NewUserService(userRepo interfaces.UserManager, authManager interfaces.AuthManager) *UserService {
	return &UserService{
		UserRepo: userRepo,
		Auth:     authManager,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *models.User) (*jwt.Token, error) {
	exist, err := s.UserRepo.UserEmailExists(user.Email)
	if exist || errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, clierrs.ErrUserAlreadyExist
	}

	// every self-registered user gets a "common" role by default
	user.Role = enums.Common

	err = s.setPassword(user)
	if err != nil {
		return nil, err
	}

	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return s.createToken(user.ID)
}

// setPassword hashing user password
func (s *UserService) setPassword(user *models.User) error {
	hashed, err := s.toHash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return nil
}

// GetUserByID returns user by the provided id or an error if the record is not found
func (s *UserService) GetUserByID(userID, callerID uint) (*models.User, error) {
	if userID != callerID {
		err := s.isCallerAdmin(callerID)
		if err != nil {
			return nil, err
		}
	}

	userExists, err := s.UserRepo.UserExists(userID)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, clierrs.ErrUserNotFound
	}

	return s.UserRepo.GetUserByID(userID)
}

// GetUserByEmail returns user by the provided email or an error if the record is not found
func (s *UserService) GetUserByEmail(userEmail string) (*models.User, error) {
	userExists, err := s.UserRepo.UserEmailExists(userEmail)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, clierrs.ErrUserNotFound
	}
	return s.UserRepo.GetUserByEmail(userEmail)
}

// UpdateUserParamByID updates param of the user with provided userID with value or returns an error if the record is not found
func (s *UserService) UpdateUserParamByID(userID uint, param string, value interface{}, callerID uint) error {
	err := s.isCallerAdmin(callerID)
	if err != nil {
		return err
	}

	userExists, err := s.UserRepo.UserExists(userID)
	if err != nil {
		return err
	}
	if !userExists {
		return clierrs.ErrUserNotFound
	}

	if slices.Contains(notAllowedFieldsToUpdate, param) {
		return clierrs.ErrUpdateForbidden
	}

	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	return s.UserRepo.UpdateUser(user, param, value)
}

// DeleteUserByID deletes user with provided userID or returns an error if the record is not found
func (s *UserService) DeleteUserByID(userID, callerID uint) error {
	err := s.isCallerAdmin(callerID)
	if err != nil {
		return err
	}

	userExists, err := s.UserRepo.UserExists(userID)
	if err != nil {
		return err
	}
	if !userExists {
		return clierrs.ErrUserNotFound
	}

	user, err := s.UserRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	return s.UserRepo.DeleteUser(user)
}

// Login checks if the user with provided email and password exists and returns new auth token for him containing his ID
func (s *UserService) Login(email, password string) (*models.User, *jwt.Token, error) {
	userExists, err := s.UserRepo.UserEmailExists(email)
	if err != nil {
		return nil, nil, err
	}
	if !userExists {
		return nil, nil, clierrs.ErrUserNotFound
	}

	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	hashString, err := s.toHash(password)
	if err != nil {
		return nil, nil, err
	}

	if user.Password != hashString {
		return nil, nil, clierrs.ErrInvalidCredentials
	}

	token, err := s.createToken(user.ID)
	return user, token, err
}

// toHash converts string to hash
func (s *UserService) toHash(text string) (string, error) {
	return s.Auth.HashString(text)
}

// createToken creates JWT token with data as models.Token
func (s *UserService) createToken(userID uint) (*jwt.Token, error) {
	return s.Auth.BuildToken(userID)
}

// GetAllUsers gets all available users from database
func (s *UserService) GetAllUsers(callerID uint) ([]*models.User, error) {
	err := s.isCallerAdmin(callerID)
	if err != nil {
		return nil, err
	}

	users, err := s.UserRepo.GetAllUsers(callerID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetAllUsers check if the caller is admin
func (s *UserService) isCallerAdmin(callerID uint) error {
	caller, err := s.UserRepo.GetUserByID(callerID)
	if err != nil {
		return err
	}

	if caller.Role != enums.Admin {
		return clierrs.ErrPermissionAdmin
	}

	return nil
}
