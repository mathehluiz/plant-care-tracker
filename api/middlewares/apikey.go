package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ValidateAPIKey(keys []string) Middleware {
	return func(c *gin.Context) *result {
		key := c.Request.Header.Get("x-api-key")

		for _, k := range keys {
			if key == k {
				return nil
			}
		}

		return &result{
			Status: http.StatusNotFound,
		}
	}
}
