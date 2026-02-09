package service

import (
	"errors"
	"go-login/models"
	"go-login/repository"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		log.Printf("Login failed: user not found - %s", email)
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Login failed: invalid password - %s", email)
		return "", errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user.ID, email)
	if err != nil {
		log.Printf("JWT generation failed for user: %s", email)
		return "", errors.New("failed to generate token")
	}

	log.Printf("Login successful: %s", email)
	return token, nil
}

func (s *AuthService) generateJWT(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) Register(username, email, password string) error {
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to create user")
	}

	user := &models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return errors.New("failed to create user")
	}

	return nil
}
