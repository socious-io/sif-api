package views

import (
	"context"
	"net/http"
	"socious/src/apps/models"
	"socious/src/apps/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

func projectsGroup(router *gin.Engine) {
	g := router.Group("projects")

	g.GET("", paginate(), func(c *gin.Context) {
		identity, _ := c.Get("identity")
		page, _ := c.Get("paginate")

		projects, total, err := models.GetProjects(identity.(*models.Identity).ID, page.(database.Paginate))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"results": projects,
			"total":   total,
		})
	})

	g.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")

		p, err := models.GetProject(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	g.POST("", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		identity, _ := c.Get("identity")

		form := new(ProjectForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		p := new(models.Project)
		utils.Copy(form, p)
		p.IdentityID = identity.(*models.Identity).ID
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
		p := new(models.Project)
		utils.Copy(form, p)
		p.ID = uuid.MustParse(id)
		if err := p.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	g.DELETE("/:id", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		p, err := models.GetProject(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
}
