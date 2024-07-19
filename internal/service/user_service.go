package service

import (
	"errors"
	"time"

	"github.com/andrMaulana/Go-Task-Management-System/internal/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

var jwtKey = []byte("your_secret_key")

func (s *UserService) Register(user *models.User) error {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user with this email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	result := s.db.Create(user)
	return result.Error
}

func (s *UserService) Login(email, password string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
