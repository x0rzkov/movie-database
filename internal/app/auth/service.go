package auth

import (
	"context"
	"time"

	"github.com/dqkcode/movie-database/internal/app/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"

	"github.com/sirupsen/logrus"
)

type (
	config struct {
		key string `envconfig:"JWT_KEY" default:"cold water"`
	}
	Service struct {
		conf   config
		usrSrv UserService
	}
	UserService interface {
		FindUserByEmail(ctx context.Context, email string) (*types.UserInfo, error)
	}
)

func NewService(userSvc UserService) *Service {
	return &Service{
		usrSrv: userSvc,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {
	if err := validator.New().Struct(req); err != nil {
		logrus.Errorf("failed to validation, err: %v", err)
		return "", err
	}
	user, err := s.usrSrv.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}
	if user.Locked {
		return "", ErrUserIsLocked
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logrus.Errorf("Compare hash and password failed")
		return "", err
	}
	expirationTime := time.Now().Add(500 * time.Minute)
	claims := &Claims{
		Email: req.Email,
		Id:    user.ID,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(s.conf.key))
	if err != nil {
		logrus.Errorf("Signing string fail")
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(token string) (*Claims, error) {
	var tokenValidType string
	if token[:7] == "Bearer " {
		tokenValidType = token[7:]
	} else {
		tokenValidType = token
	}
	claims := &Claims{}
	c := config{}
	t, err := jwt.ParseWithClaims(tokenValidType, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.key), nil
	})
	if err != nil {
		return nil, ErrCompareToken
	}
	if !t.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}
