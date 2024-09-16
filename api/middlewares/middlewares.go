package middlewares

import "github.com/gin-gonic/gin"

type Middleware func(c *gin.Context) *result

type result struct {
	Status int
	Error  error
}

func AddMiddlewares(middlewares ...Middleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errs []*result

		for _, m := range middlewares {
			if res := m(c); res != nil {
				errs = append(errs, res)
			}
		}

		if len(middlewares) > len(errs) {
			c.Next()
			return
		}

		for _, res := range errs {
			if res.Error == nil {
				c.AbortWithStatus(res.Status)
				return
			}

			c.AbortWithStatusJSON(res.Status, gin.H{"error": res.Error.Error()})
			return
		}
	}
}
