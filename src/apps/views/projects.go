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
	"github.com/google/martian/v3/log"
	"github.com/google/uuid"
	"github.com/socious-io/goaccount"
	"github.com/socious-io/gopay"
	database "github.com/socious-io/pkg_database"
)

func projectsGroup(router *gin.Engine) {
	g := router.Group("projects")

	g.GET("", auth.LoginOptional(), paginate(), func(c *gin.Context) {
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

	g.GET("/:id", auth.LoginOptional(), func(c *gin.Context) {

		p, err := models.GetProject(uuid.MustParse(c.Param("id")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if u, ok := c.Get("user"); ok {
			v, err := models.GetVoteByUserAndProject(u.(*models.User).ID, p.ID)
			if err == nil && v != nil {
				p.UserVoted = true
			}
		}

		c.JSON(http.StatusOK, p)
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
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

	g.PATCH("/:id", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		form := new(ProjectForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p, err := models.GetProject(uuid.MustParse(id))
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

	g.DELETE("/:id", auth.LoginRequired(), OrganizationRequired(), func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		p, err := models.GetProject(uuid.MustParse(id))
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

	g.POST("/:id/votes", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		if user.IdentityVerifiedAt == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You must verify your identity before voting"})
			return
		}

		project, err := models.GetProject(uuid.MustParse(c.Param("id")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		now := time.Now()
		if !config.Config.Debug && (now.Before(project.Round.VotingStartAt) || now.After(project.Round.VotingEndAt)) {
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

	g.POST("/:id/donates", auth.LoginRequired(), func(c *gin.Context) {
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
		project, err := models.GetProject(uuid.MustParse(c.Param("id")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		now := time.Now()
		if !config.Config.Debug && (now.Before(project.Round.VotingStartAt) || now.After(project.Round.VotingEndAt)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "voting period is closed"})
			return
		}

		donation := &models.Donation{
			UserID:    user.ID,
			ProjectID: project.ID,
			Currency:  form.Currency,
			Amount:    form.Amount,
			Status:    models.DonationStatusPending,
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
			TotalAmount: donation.Amount,
		})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		payment.SetToCryptoMode(form.Currency, 1)
		if _, err := payment.AddIdentity(gopay.IdentityParams{
			ID:      user.ID,
			Account: form.WalletAddress,
		}); err != nil {
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

		vote := &models.Vote{
			UserID:    user.ID,
			ProjectID: project.ID,
		}
		if err := vote.Create(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "already voted"})
			return
		}

		// ADD impact point
		go func() {
			ip := goaccount.ImpactPoint{
				UserID:              user.ID,
				SocialCause:         project.SocialCause,
				SocialCauseCategory: string(utils.GetSDG(project.SocialCause)),
				TotalPoints:         int(donation.Amount),
				Type:                "DONATION",
				Meta: map[string]any{
					"donation": donation,
				},
			}
			if err := ip.AddImpactPoint(); err != nil {
				log.Errorf("Failed to add impact point: %v", err)
			}
		}()

		c.JSON(http.StatusCreated, gin.H{"donation": donation})
	})

	g.GET("/:id/comments", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		projectID := uuid.MustParse(c.Param("id"))
		pagination := c.MustGet("paginate").(database.Paginate)

		comments, total, err := models.GetComments(projectID, identity.ID, pagination)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"results": []gin.H{},
				"total":   0,
				"page":    c.MustGet("page"),
				"limit":   c.MustGet("limit"),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": comments,
			"total":   total,
			"page":    c.MustGet("page"),
			"limit":   c.MustGet("limit"),
		})
	})

	g.GET("/comments/:id", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		id := uuid.MustParse(c.Param("id"))
		comment, err := models.GetComment(id, identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, comment)
	})

	g.POST("/:id/comments", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		form := new(CommentForm)
		if err := c.BindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		comment := &models.Comment{
			ProjectID:  uuid.MustParse(c.Param("id")),
			IdentityID: identity.ID,
			Content:    form.Content,
			MediaID:    form.MediaID,
			ParentID:   form.ParentID,
		}

		if err := comment.Create(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, comment)
	})

	g.PUT("/comments/:id", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		id := uuid.MustParse(c.Param("id"))
		form := new(CommentForm)
		if err := c.BindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
		comment, err := models.GetComment(id, identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if comment.IdentityID != identity.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "not comment owner"})
			return
		}
		comment.Content = form.Content
		comment.MediaID = form.MediaID

		if err := comment.Update(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, comment)
	})

	g.DELETE("/comments/:id", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		commentID := uuid.MustParse(c.Param("id"))
		comment, err := models.GetComment(commentID, identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if comment.IdentityID != identity.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "not comment owner"})
			return
		}
		if err := comment.Delete(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
	})

	g.POST("/comments/:id/likes", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)

		like := &models.Like{
			CommentID:  uuid.MustParse(c.Param("id")),
			IdentityID: identity.ID,
		}
		if err := like.Create(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, like)
	})

	g.DELETE("/comments/:id/likes", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		like, err := models.GetLike(uuid.MustParse(c.Param("id")), identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := like.Delete(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to unlike"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "unliked"})
	})

	g.POST("/comments/:id/reactions", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		form := new(ReactionForm)
		if err := c.BindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		reaction := &models.Reaction{
			CommentID:  uuid.MustParse(c.Param("id")),
			IdentityID: identity.ID,
			Reaction:   form.Reaction,
		}
		if err := reaction.Create(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to add reaction"})
			return
		}
		c.JSON(http.StatusCreated, reaction)
	})

	g.DELETE("/comments/:id/reactions", auth.LoginRequired(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		reaction, err := models.GetReaction(uuid.MustParse(c.Param("id")), identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := reaction.Delete(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "reaction removed"})
	})

}
