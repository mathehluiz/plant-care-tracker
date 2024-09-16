package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var secret = "coxinha123"
var ErrExpiredToken = errors.New("the token has expired")

type customClaims struct {
	jwt.StandardClaims
	Roles    []string
	Verified bool
}

func GenerateToken(id interface{}, roles []string, verified bool) (string, error) {
	stdClaims := customClaims{
		jwt.StandardClaims{
			Subject:   fmt.Sprintf("%v", id),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "plant-care-tracker",
			IssuedAt:  time.Now().Unix(),
		},
		roles,
		verified,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
	ss, tokenErr := token.SignedString([]byte(secret))
	if tokenErr != nil {
		return "", tokenErr
	}
	return ss, nil
}

func ValidateToken(token string) (string, bool, []string, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &customClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
	if err != nil {
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return "", false, nil, ErrExpiredToken
		}
		return "", false, nil, err
	}
	claims, ok := parsedToken.Claims.(*customClaims)
	if !ok {
		return "", false, nil, errors.New("error on parsing the claims")
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return "", false, nil, ErrExpiredToken
	}
	validErr := claims.Valid()
	if validErr != nil {
		return "", false, nil, validErr
	}

	return claims.Subject, claims.Verified, claims.Roles, nil
}
