package views

import (
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"

	"github.com/gin-gonic/gin"
	database "github.com/socious-io/pkg_database"
)

func roundsGroup(router *gin.Engine) {
	g := router.Group("rounds")

	g.GET("", func(c *gin.Context) {
		round, err := models.GetRoundLatestRound()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, round)
	})

	g.GET("/all", auth.LoginOptional(), paginate(), func(c *gin.Context) {
		p := c.MustGet("paginate").(database.Paginate)

		rounds, total, err := models.GetAllRounds(p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  rounds,
			"total": total,
		})
	})
}
