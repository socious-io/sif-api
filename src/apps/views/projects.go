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

	g.GET("/preview", auth.LoginOptional(), paginate(), func(c *gin.Context) {
		pagination := c.MustGet("paginate").(database.Paginate)

		projects, total, err := models.GetProjects(pagination)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		projectsPreview := new([]models.ProjectPreview)
		utils.Copy(projects, projectsPreview)

		c.JSON(http.StatusOK, gin.H{
			"results": projectsPreview,
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

		c.JSON(http.StatusOK, p)
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		identity, _ := c.MustGet("identity").(*models.Identity)

		// r, err := models.GetRoundLatestRound()
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }
		// now := time.Now()
		// if !config.Config.Debug && (now.Before(r.SubmissionStartAt) || now.After(r.SubmissionEndAt)) {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Submission period is closed"})
		// 	return
		// }

		form := new(ProjectForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p := new(models.Project)
		utils.Copy(form, p)
		p.IdentityID = identity.ID
		if form.Status == nil {
			p.Status = models.ProjectStatusActive
		}
		if err := p.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, p)
	})

	g.PATCH("/:id", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		identity := c.MustGet("identity").(*models.Identity)
		id := uuid.MustParse(c.Param("id"))

		form := new(ProjectForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p, err := models.GetProject(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if p.IdentityID != identity.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		utils.Copy(form, p)
		p.ID = id
		if err := p.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	g.DELETE("/:id", auth.LoginRequired(), OrganizationRequired(), func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		identity := c.MustGet("identity").(*models.Identity)
		id := uuid.MustParse(c.Param("id"))

		p, err := models.GetProject(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if p.IdentityID != identity.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		if err := p.Delete(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	// g.POST("/:id/votes", auth.LoginRequired(), func(c *gin.Context) {
	// 	user := c.MustGet("user").(*models.User)
	// 	// Allow people vote and donate without verify
	// 	/* if user.IdentityVerifiedAt == nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "You must verify your identity before voting"})
	// 		return
	// 	} */

	// 	project, err := models.GetProject(uuid.MustParse(c.Param("id")))
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	now := time.Now()
	// 	if !config.Config.Debug && (now.Before(project.Round.VotingStartAt) || now.After(project.Round.VotingEndAt)) {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "voting period is closed"})
	// 		return
	// 	}

	// 	vote := &models.Vote{
	// 		UserID:    user.ID,
	// 		ProjectID: project.ID,
	// 	}
	// 	if err := vote.Create(c.MustGet("ctx").(context.Context)); err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "already voted"})
	// 		return
	// 	}

	// 	go func() {
	// 		ip := goaccount.ImpactPoint{
	// 			UserID:              user.ID,
	// 			SocialCause:         project.SocialCause,
	// 			SocialCauseCategory: string(utils.GetSDG(project.SocialCause)),
	// 			TotalPoints:         1,
	// 			Type:                "VOTING",
	// 			UniqueTag:           vote.ID.String(),
	// 			Value:               float64(0),
	// 			Meta: map[string]any{
	// 				"vote": vote,
	// 			},
	// 		}
	// 		if err := ip.AddImpactPoint(); err != nil {
	// 			log.Errorf("Failed to add impact point: %v", err)
	// 		}

	// 		ra := goaccount.ReferralAchievement{
	// 			RefereeID:       user.ID,
	// 			AchievementType: "VOTE",
	// 			Meta: map[string]any{
	// 				"vote": vote,
	// 			},
	// 		}
	// 		if err := ra.AddReferralAchievement(); err != nil {
	// 			log.Errorf("Failed to add achievement: %v", err)
	// 		}
	// 	}()

	// 	c.JSON(http.StatusCreated, gin.H{"vote": vote})
	// })

	g.POST("/:id/donates", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		ctx, _ := c.MustGet("ctx").(context.Context)
		id := uuid.MustParse(c.Param("id"))
		// Allow people vote and donate without verify
		/* if user.IdentityVerifiedAt == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You must verify your identity before donating"})
			return
		} */

		form := new(DnateDepositForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		project, err := models.GetProject(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		rate := form.Rate
		if rate <= 0 || rate >= 2 {
			rate = 1
		}

		donation := &models.Donation{
			UserID:    user.ID,
			ProjectID: project.ID,
			Currency:  form.Currency,
			Amount:    form.Amount,
			Status:    models.DonationStatusPending,
			Rate:      rate,
			Anonymous: form.Anonymous,
			PaidAs:    form.PaidAs,
		}

		if err := donation.Create(ctx); err != nil {
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
		pID := payment.ID.String()
		donation.TransactionID = &pID

		// impactPoints := int(donation.Amount * rate)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if form.PaymentType == models.Fiat {
			fiatService := config.Config.Payment.Fiats[0]

			payment.Currency = gopay.Currency(form.Currency)
			payment.SetToFiatMode(fiatService.Name)
			if form.CardToken == nil && user.StripeCustomerID == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "payment source card could not be found"})
				return
			}
			if user.StripeCustomerID == nil {
				cus, err := fiatService.AddCustomer(user.Email)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				user.StripeCustomerID = &cus.ID
				if err := user.Upsert(c.MustGet("ctx").(context.Context)); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
			}

			if form.CardToken != nil {
				if _, err := fiatService.AttachPaymentMethod(*user.StripeCustomerID, *form.CardToken); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
			}

			if _, err := payment.AddIdentity(gopay.IdentityParams{
				ID:      user.ID,
				Account: *user.StripeCustomerID,
			}); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if err := payment.Deposit(); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}

		} else {
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
		}

		if payment.Status == gopay.ON_HOLD || *payment.TransactionStatus == gopay.ACTION_REQUIRED {
			c.JSON(http.StatusAccepted, gin.H{
				"donation":        donation,
				"message":         "payment is on hold",
				"action_required": true,
				"client_secret":   payment.ClientSecret,
			})
			return
		}
		donation.Status = models.DonationStatusApproved
		if err := donation.Update(c.MustGet("ctx").(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// vote := &models.Vote{
		// 	UserID:    user.ID,
		// 	ProjectID: project.ID,
		// }
		// if err := vote.Create(c.MustGet("ctx").(context.Context)); err != nil {
		// 	log.Infof("Failed to create vote: %v", err)
		// } else {
		// 	impactPoints += 1
		// }

		// go func() {
		// 	ip := goaccount.ImpactPoint{
		// 		UserID:              user.ID,
		// 		SocialCause:         project.SocialCause,
		// 		SocialCauseCategory: string(utils.GetSDG(project.SocialCause)),
		// 		TotalPoints:         impactPoints,
		// 		Type:                "DONATION",
		// 		UniqueTag:           donation.ID.String(),
		// 		Value:               float64(impactPoints),
		// 		Meta: map[string]any{
		// 			"donation": donation,
		// 		},
		// 	}
		// 	if err := ip.AddImpactPoint(); err != nil {
		// 		log.Errorf("Failed to add impact point: %v", err)
		// 	}
		// }()

		c.JSON(http.StatusCreated, gin.H{"donation": donation})
	})

	g.GET("/:id/donates", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		pagination := c.MustGet("paginate").(database.Paginate)
		donations, total, err := models.GetDonations(c.Param("id"), pagination)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": donations,
			"total":   total,
			"page":    c.MustGet("page"),
			"limit":   c.MustGet("limit"),
		})
	})

	g.PUT("/donates/:id/confirm", auth.LoginRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)

		form := new(DonateDepositConfirmForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		donation, err := models.GetDonation(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		payment, err := gopay.FetchByUniqueRef(donation.ID.String())
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := payment.ConfirmPayment(form.PaymentIntentID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		donation.Status = models.DonationStatusApproved
		pID := payment.ID.String()
		donation.TransactionID = &pID
		if err := donation.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"donation": donation,
		})
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
