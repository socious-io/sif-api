package views

import (
	"context"
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/socious-io/goaccount"
)

func authGroup(router *gin.Engine) {
	g := router.Group("auth")

	g.POST("", func(c *gin.Context) {
		form := new(AuthForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		session, authURL, err := goaccount.StartSession(form.RedirectURL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"session":  session,
			"auth_url": authURL,
		})
	})

	g.POST("/session", func(c *gin.Context) {
		form := new(SessionForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		token, err := goaccount.GetSessionToken(form.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var (
			connect *models.OauthConnect
			user    = new(models.User)
		)

		if err := token.GetUserProfile(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if connect, err = models.GetOauthConnectByMUI(user.ID.String(), models.OauthConnectedProvidersSociousID); err != nil {
			connect = &models.OauthConnect{
				Provider:       models.OauthConnectedProvidersSociousID,
				AccessToken:    token.AccessToken,
				RefreshToken:   &token.RefreshToken,
				MatrixUniqueID: user.ID.String(),
			}
		}
		u, err := models.GetUserByUsername(user.Username)
		if err != nil {
			if err := user.Create(c.MustGet("ctx").(context.Context)); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		} else {
			user = u
		}

		var orgs = []models.Organization{}
		token.GetMyOrganizations(orgs)
		for _, o := range orgs {
			o.UpsertAndMember(user.ID)
		}

		connect.IdentityId = user.ID
		if err := connect.Upsert(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		jwt, err := auth.GenerateFullTokens(user.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, jwt)
	})

	g.POST("/refresh", func(c *gin.Context) {
		form := new(RefreshForm)
		if err := c.Bind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		claims, err := auth.VerifyToken(form.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		jwt, err := auth.GenerateFullTokens(claims.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, jwt)
	})
}
