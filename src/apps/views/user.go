package views

import (
	"net/http"
	"sif/src/apps/auth"

	"github.com/gin-gonic/gin"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(auth.LoginRequired())

	g.GET("/profile", func(c *gin.Context) {
		c.JSON(http.StatusOK, c.MustGet("user"))
	})

	g.GET("", func(c *gin.Context) {

	})
}
