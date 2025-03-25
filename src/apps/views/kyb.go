package views

import (
	"context"
	"fmt"
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"
	"sif/src/apps/utils"
	"sif/src/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func kybVerificationGroup(router *gin.Engine) {
	g := router.Group("kybs")

	g.POST("", auth.LoginRequired(), OrganizationRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		identity := c.MustGet("identity").(*models.Identity)
		ctx := c.MustGet("ctx")

		form := new(KYBVerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, _ := models.GetOrganization(identity.ID)

		kyb := &models.KYBVerification{
			UserID: user.ID,
			OrgID:  org.ID,
		}

		if err := kyb.Create(ctx.(context.Context), form.Documents); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org.Status = models.OrganizationStatusPending
		if err := org.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		utils.DiscordSendTextMessage(
			config.Config.Discord.Channel,
			createDiscordReviewMessage(kyb, user, org),
		)

		c.JSON(http.StatusOK, kyb)
	})

	g.GET("", auth.LoginRequired(), OrganizationRequired(), paginate(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)

		kyb, err := models.GetKybByOrganization(identity.ID)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "no kyb found for this identity"})
			return
		}
		c.JSON(http.StatusOK, kyb)
	})

	g.GET("/:id/approve", adminAccessRequired(), func(c *gin.Context) {

		ctx := c.MustGet("ctx").(context.Context)
		verificationId := c.Param("id")

		verification, err := models.GetKyb(uuid.MustParse(verificationId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := verification.ChangeStatus(ctx, models.KYBStatusApproved); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, _ := models.GetOrganization(verification.OrgID)

		org.Verified = true

		if err := org.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	g.GET("/:id/reject", adminAccessRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		verificationId := c.Param("id")

		verification, err := models.GetKyb(uuid.MustParse(verificationId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := verification.ChangeStatus(ctx, models.KYBStatusRejected); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

}

func createDiscordReviewMessage(kyb *models.KYBVerification, u *models.User, org *models.Organization) string {

	documents := ""
	for i, document := range kyb.Documents {
		documents = fmt.Sprintf("%s\n%v. %s", documents, i, document.Url)
	}

	message := fmt.Sprintf("ID: %s\n", kyb.ID)
	message += "\nUser--------------------------------\n"
	message += fmt.Sprintf("ID: %s\n", u.ID)

	if u.FirstName != "" {
		message += fmt.Sprintf("Firstname: %s\n", u.FirstName)
	} else {
		message += "Firstname: N/A\n"
	}

	if u.LastName != "" {
		message += fmt.Sprintf("Lastname: %s\n", u.LastName)
	} else {
		message += "Lastname: N/A\n"
	}

	message += fmt.Sprintf("Email: %s\n", u.Email)
	message += "\nOrganization------------------------\n"
	message += fmt.Sprintf("ID: %s\n", org.ID)
	message += fmt.Sprintf("Name: %v\n", org.Name)
	message += fmt.Sprintf("Description: %v\n", org.Description)
	message += fmt.Sprintf("\nDocuments---------------------------%s\n\n", documents)
	message += "\nReviewing----------------------------\n"
	message += fmt.Sprintf("Approve: <%s/kyb/%s/approve?admin_access_token=%s>\n", config.Config.Host, kyb.ID, config.Config.AdminToken)
	message += fmt.Sprintf("Reject: <%s/kyb/%s/reject?admin_access_token=%s>\n", config.Config.Host, kyb.ID, config.Config.AdminToken)

	return message

}
