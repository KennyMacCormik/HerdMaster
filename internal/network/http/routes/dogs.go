package routes

import (
	"github.com/gin-gonic/gin"
)

func DogsHandlers(router *gin.Engine) {
	router.GET("/dog/:id", func(c *gin.Context) {
		_ = c.Param("key")
	})
	router.PUT("/dog/:id", func(c *gin.Context) {
		_ = c.Param("key")
	})
	router.DELETE("/dog/:id", func(c *gin.Context) {
		_ = c.Param("key")
	})
}
