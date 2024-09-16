package middlewares

import (
	"errors"
	"github.com/mathehluiz/plant-care-tracker/pkg/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

func ValidateRoles() Middleware {
	return func(c *gin.Context) *result {
		token := c.Request.Header.Get("Authorization")

		token = strings.ReplaceAll(token, "Bearer ", "")
		if token == "" {
			return &result{
				Status: 401,
				Error:  errors.New("no authorization token provided"),
			}
		}

		id, verified, roles, tokenErr := jwt.ValidateToken(token)
		if tokenErr != nil {
			if errors.Is(jwt.ErrExpiredToken, tokenErr) {
				return &result{
					Status: 401,
					Error:  errors.New("token has expired"),
				}
			}

			return &result{
				Status: 400,
				Error:  errors.New("invalid token format"),
			}
		}

		c.Set("auth:type", "token")
		c.Set("auth:bearer:id", id)
		c.Set("auth:bearer:verified", verified)
		c.Set("auth:bearer:roles", roles)

		return nil
	}
}
