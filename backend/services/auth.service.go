package services

import (
	"context"
	"errors"
	"time"

	"github.com/BalkanID-University/vit-2026-capstone-internship-hiring-task-joel2607/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService provides methods for user authentication, including registration,
// login, and token validation. It interacts with the database to manage user records.
type AuthService struct {
	DB *gorm.DB
}

// NewAuthService creates and returns a new instance of AuthService.
// It requires a GORM database connection.
//
// Inputs:
// - db: A pointer to a gorm.DB instance.
//
// Outputs:
// - A pointer to the newly created AuthService.
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

// Register handles the creation of a new user account.
// It hashes the user's password for security before storing it in the database.
//
// Inputs:
// - input: A models.RegisterInput struct containing the new user's details
//          (username, email, password).
//
// Outputs:
// - A pointer to the newly created models.User object.
// - An error if password hashing or database creation fails.
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

// Login authenticates a user based on their email and password.
// If the credentials are valid, it generates and returns a JWT.
//
// Inputs:
// - email: The user's email address.
// - password: The user's plain-text password.
//
// Outputs:
// - A string containing the signed JWT.
// - A pointer to the authenticated models.User object.
// - An error if the user is not found, the password is invalid, or JWT signing fails.
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

// GetUserFromToken validates a JWT string and retrieves the corresponding user from the database.
// It checks for token expiration and validity before fetching the user.
//
// Inputs:
// - tokenString: The JWT string to be validated.
//
// Outputs:
// - A pointer to the models.User object associated with the token.
// - An error if the token is expired, invalid, or the user is not found.
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

// SearchUsers performs a case-insensitive search for users with usernames or emails that contain the query string.
func (s *AuthService) SearchUsers(ctx context.Context, query string) ([]*models.User, error) {
	var users []*models.User
	searchQuery := "%" + query + "%"
	err := s.DB.Where("username ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).Find(&users).Error
	return users, err
}
