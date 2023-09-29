package service

import (
	"errors"
	"os"
	"time"

	"github.com/gin-api/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) userService {
	return userService{userRepo: userRepo}
}

func (s userService) Login(auth AuthInput) (*string, error) {
	user, err := s.userRepo.GetUser(auth.UserName)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}

	checkPass, err := checkHashPassword(auth.Password, user.Password)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}

	if !checkPass {
		return nil, errors.New("Unauthorized")
	}

	claims := &JwtCustomClaims{
		Id:       user.ID.Hex(),
		UserName: user.UserName,
		Admin:    true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, err := goDotEnvVariable("SECRET_KEY")
	if err != nil {
		return nil, err
	}
	t, err := token.SignedString([]byte(*key))
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s userService) RegisterUser(register RegisterUser) (*string, error) {
	if register.Password != register.ConfirmPassword {
		return nil, errors.New("passwords don't match")
	}

	password, err := hashPassword(register.Password)
	if err != nil {
		return nil, err
	}

	user := repository.CreateUser{
		UserName: register.UserName,
		Password: password,
		Email:    register.Email,
	}

	id, err := s.userRepo.CreateUser(user)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("data with the same unique value for username already exists")
		}
		return nil, err
	}

	objIDString := id.(primitive.ObjectID).Hex()

	message := "User id " + objIDString + " created successfully"

	return &message, nil
}

func checkHashPassword(password string, passwordDB string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(password))
	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else {
		return false, err
	}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func goDotEnvVariable(key string) (*string, error) {

	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	value := os.Getenv(key)

	return &value, nil
}
