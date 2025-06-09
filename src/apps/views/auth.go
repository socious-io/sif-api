package views

import (
	"context"
	"log"
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"
	"sif/src/apps/utils"

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
			ctx     = c.MustGet("ctx").(context.Context)
		)

		u, err := token.GetUserProfile()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.Copy(u, user)

		if connect, err = models.GetOauthConnectByMUI(user.ID.String(), models.OauthConnectedProvidersSociousID); err != nil {
			connect = &models.OauthConnect{
				Provider:       models.OauthConnectedProvidersSociousID,
				AccessToken:    token.AccessToken,
				RefreshToken:   &token.RefreshToken,
				MatrixUniqueID: user.ID.String(),
				IdentityId:     user.ID,
			}
		}

		if err := user.Upsert(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orgs, err := token.GetMyOrganizations()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, org := range orgs {
			var o = new(models.Organization)
			utils.Copy(org, o)
			if err := o.Create(ctx, user.ID); err != nil {
				log.Println(err.Error(), o)
			}
		}

		if err := connect.Upsert(ctx); err != nil {
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
