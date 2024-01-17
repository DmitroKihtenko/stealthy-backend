package services

import (
	"SharingBackend/api"
	"SharingBackend/base"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type JWTClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthorizationService struct {
	JwtConfig *base.JwtConfig
}

func (service AuthorizationService) GenerateToken(user *api.User) (string, error) {
	tokenLifespan := service.JwtConfig.DaysLifespan

	claims := &JWTClaim{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * time.Duration(tokenLifespan)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(service.JwtConfig.Secret))
	if err != nil {
		return "", base.ServiceError{
			Summary: "Token generation error",
			Detail:  err.Error(),
		}
	}
	return tokenString, nil
}

func (service AuthorizationService) ParseToken(tokenString string) (*api.User, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(service.JwtConfig.Secret), nil
		},
	)
	if err != nil {
		return nil, base.ServiceError{
			Summary: "Invalid token",
			Status:  http.StatusForbidden,
		}
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, base.ServiceError{
			Summary: "Invalid token",
			Status:  http.StatusForbidden,
		}
	}
	expiration := time.Unix(claims.ExpiresAt, 0)
	if time.Now().Compare(expiration) > 0 {
		return nil, base.ServiceError{
			Summary: "Token expired. Login to your account again",
			Detail:  fmt.Sprintf("Expration date: %v", expiration.Format(time.RFC3339)),
			Status:  http.StatusUnauthorized,
		}
	}
	return &api.User{
		Username: claims.Username,
	}, nil
}
