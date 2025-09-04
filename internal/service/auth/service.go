package auth

import (
	"errors"
	"time"

	"github.com/ferdy-adr/elibrary-backend/internal/configs"
	"github.com/ferdy-adr/elibrary-backend/internal/model"
	userRepo "github.com/ferdy-adr/elibrary-backend/internal/repository/users"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepository *userRepo.Repository
}

func NewService(userRepository *userRepo.Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) Register(req model.RegisterRequest) (*model.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepository.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
	}

	err = s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepository.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *Service) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.Get().JWT.SecretKey))
}
