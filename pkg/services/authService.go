package services

import (
	"errors"
	"goticka/pkg/adapters/repositories"
	"goticka/pkg/config"
	"goticka/pkg/dependencies"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	userRepository repositories.UserRepositoryInterface
}

func NewAuthService() AuthService {
	return AuthService{
		userRepository: dependencies.DI().UserRepository,
	}
}

type AuthData struct {
	ID       int64
	Username string
	Roles    []string
}

type Claims struct {
	ID       int64    `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

func (as AuthService) PasswordAuthentication(username, password string) (string, error) {
	user, err := NewUserService().GetByUserNameAndPassword(username, password)
	if err != nil {
		log.Printf("[ERROR] %s\n", err)
		return "", errors.New("wrong username and/or password")
	}

	secrets := config.GetConfig().Secrets

	claims := &Claims{
		ID:       user.ID,
		Username: user.UserName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(secrets.JWTTTL)),
		},
		Roles: make([]string, 0),
	}

	for _, role := range user.Roles {
		claims.Roles = append(claims.Roles, role.Name)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secrets.JWTSecret))
}

func (as AuthService) VerifyJWT(JWT string) (AuthData, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(JWT, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Secrets.JWTSecret), nil
	})
	if err != nil {
		log.Printf("[ERROR] %s\n", err)
		return AuthData{}, errors.New("invalid token")
	}
	if !tkn.Valid {
		return AuthData{}, errors.New("invalid token")
	}
	if time.Until(claims.ExpiresAt.Time) <= 0 {
		return AuthData{}, errors.New("expired token")
	}
	return AuthData{
		ID:       claims.ID,
		Username: claims.Username,
		Roles:    claims.Roles,
	}, nil
}
