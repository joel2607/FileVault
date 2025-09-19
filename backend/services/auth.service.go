package services

import (
	"errors"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(input models.RegisterInput) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) Login(email string, password string) (string, *models.User, error) {
	var user models.User
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	expirationHours := viper.GetInt("auth.jwt_expiration_hours")
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"role": user.Role,
		"exp": time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
	})

	token, err := claims.SignedString([]byte(viper.GetString("JWT_AUTH_SECRET")))
	if err != nil {
		return "", nil, err
	}

	return token, &user, nil
}

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT_AUTH_SECRET")), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token is expired")
		}
		return nil, errors.New("could not parse token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from claims
		sub, ok := claims["sub"].(float64)
		if !ok {
			return nil, errors.New("invalid token: sub claim is not a number")
		}
		userID := uint(sub)

		var user models.User
		if err := s.DB.First(&user, userID).Error; err != nil {
			return nil, errors.New("user not found")
		}
		return &user, nil
	}

	return nil, errors.New("invalid token")
}
