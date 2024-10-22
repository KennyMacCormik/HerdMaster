package routes

import "github.com/gin-gonic/gin"

func OwnersHandlers(router *gin.Engine) {
	router.GET("/owner/:id", func(c *gin.Context) {
		_ = c.Param("key")
	})
	router.PUT("/owner/:id", func(c *gin.Context) {
		_ = c.Param("key")
	})
	router.DELETE("/owner/:id", func(c *gin.Context) {
		_ = c.Param("key")
	})
}
