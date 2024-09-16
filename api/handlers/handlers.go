package handlers

import "github.com/gin-gonic/gin"

func DefaultError(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}
