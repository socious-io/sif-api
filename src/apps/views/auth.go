package views

import (
	"github.com/gin-gonic/gin"
)

func authGroup(router *gin.Engine) {
	router.Group("auth")

}
