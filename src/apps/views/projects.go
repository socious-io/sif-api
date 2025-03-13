package views

import (
	"context"
	"fmt"
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"
	"sif/src/apps/utils"
	"sif/src/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socious-io/gopay"
	database "github.com/socious-io/pkg_database"
)

func projectsGroup(router *gin.Engine) {
	g := router.Group("projects")
	g.Use(auth.LoginRequired())

	g.GET("", paginate(), func(c *gin.Context) {
		pagination := c.MustGet("paginate").(database.Paginate)

		projects, total, err := models.GetProjects(pagination)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": projects,
			"total":   total,
			"page":    c.MustGet("page"),
			"limit":   c.MustGet("limit"),
		})
	})

	g.GET("/:id", func(c *gin.Context) {

		p, err := models.GetProject(uuid.MustParse(c.Param("id")), c.MustGet("user").(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	g.POST("", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		identity, _ := c.Get("identity")

		r, err := models.GetRoundLatestRound()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		if !config.Config.Debug && (now.Before(r.SubmissionStartAt) || now.After(r.SubmissionEndAt)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Submission period is closed"})
			return
		}

		form := new(ProjectForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p := new(models.Project)
		utils.Copy(form, p)
		p.IdentityID = identity.(*models.Identity).ID
		if form.Status == nil {
			p.Status = models.ProjectStatusActive
		}
		if err := p.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, p)
	})

	g.PATCH("/:id", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		form := new(ProjectForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p, err := models.GetProject(uuid.MustParse(id), c.MustGet("user").(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if p.IdentityID != c.MustGet("identity").(*models.Identity).ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		r, err := models.GetRound(p.RoundID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		if !config.Config.Debug && (now.Before(r.SubmissionStartAt) || now.After(r.SubmissionEndAt)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Submission period is closed"})
			return
		}

		utils.Copy(form, p)
		p.ID = uuid.MustParse(id)
		if err := p.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	g.DELETE("/:id", OrganizationRequired(), func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		p, err := models.GetProject(uuid.MustParse(id), c.MustGet("user").(*models.User).ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if p.IdentityID != c.MustGet("identity").(*models.Identity).ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		r, err := models.GetRound(p.RoundID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		if !config.Config.Debug && (now.Before(r.SubmissionStartAt) || now.After(r.SubmissionEndAt)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Submission period is closed"})
			return
		}

		if err := p.Delete(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	g.POST("/:id/votes", func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		if user.IdentityVerifiedAt == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You must verify your identity before voting"})
			return
		}

		project, err := models.GetProject(uuid.MustParse(c.Param("id")), user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		r, err := models.GetRound(project.RoundID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now()
		if !config.Config.Debug && (now.Before(r.VotingStartAt) || now.After(r.VotingEndAt)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "voting period is closed"})
			return
		}

		vote := &models.Vote{
			UserID:    user.ID,
			ProjectID: project.ID,
		}
		if err := vote.Create(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already voted"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"vote": vote})
	})

	g.POST("/:id/donates", func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		if user.IdentityVerifiedAt == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You must verify your identity before donating"})
			return
		}

		form := new(DnateDepositForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		project, err := models.GetProject(uuid.MustParse(c.Param("id")), user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		donation := &models.Donation{
			UserID:      user.ID,
			ProjectID:   project.ID,
			Currency:    form.Currency,
			TotalAmount: form.TotalAmount,
			Status:      models.DonationStatusPending,
		}

		if err := donation.Create(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Start a payment session
		payment, err := gopay.New(gopay.PaymentParams{
			Tag:         fmt.Sprintf("Donation-%s-%s", *project.Title, user.Username),
			Description: form.Description,
			Ref:         donation.ID.String(),
			Type:        gopay.CRYPTO,
			Currency:    gopay.USD,
			TotalAmount: donation.TotalAmount,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := payment.ConfirmDeposit(form.TxID, form.Meta); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		donation.Status = models.DonationStatusApproved
		if err := donation.Update(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"donation": donation})
	})
}
