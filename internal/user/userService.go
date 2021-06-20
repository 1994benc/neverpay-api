package user

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type UserService interface {
	CreateUser(user UserModel) (UserModel, error)
	FindUserByEmail(email string) (UserModel, error)
}

type DefaultUserService struct {
	DB *gorm.DB
}

// Returns a new User service
func NewUserService(db *gorm.DB) *DefaultUserService {
	return &DefaultUserService{
		DB: db,
	}
}

func (s *DefaultUserService) CreateUser(u UserModel) (UserModel, error) {
	result := s.DB.Save(&u)
	return u, result.Error
}

func (s *DefaultUserService) FindUserByEmail(email string) (UserModel, error) {
	var user UserModel
	result := s.DB.Where("email = ?", email).First(&user)
	return user, result.Error
}

func (s *DefaultUserService) GetAllUsers() ([]UserModel, error) {
	var users []UserModel
	result := s.DB.Find(&users)
	return users, result.Error
}

func (s *DefaultUserService) GenerateJWT(email string, role string) (string, error) {
	var mySigningKey = []byte(os.Getenv("AUTH_SECRET")) // TODO: use a secure secretkey
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Printf("Error generating JWT: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
