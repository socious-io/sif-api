package views

import (
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"

	"github.com/gin-gonic/gin"
)

func identitiesGroup(router *gin.Engine) {
	g := router.Group("identities")
	g.Use(auth.LoginRequired())

	g.GET("", func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		currentIdentity := c.MustGet("identity").(*models.Identity)
		orgs, _ := models.GetUserOrganizations(user.ID)
		var ids = []interface{}{user.ID}
		for _, o := range orgs {
			ids = append(ids, o.ID)
		}
		identities, err := models.GetIdentities(ids)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		for i := range identities {
			if identities[i].ID == currentIdentity.ID {
				identities[i].Current = true
			}
		}
		c.JSON(http.StatusOK, gin.H{"identities": identities})
	})
}
