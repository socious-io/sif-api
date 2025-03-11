package views

import (
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"

	"github.com/gin-gonic/gin"
)

func roundsGroup(router *gin.Engine) {
	g := router.Group("rounds")
	g.Use(auth.LoginRequired())

	g.GET("", func(c *gin.Context) {
		round, err := models.GetRoundLatestRound()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, round)
	})
}
