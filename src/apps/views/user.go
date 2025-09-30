package views

import (
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")
	g.Use(auth.LoginRequired())

	g.GET("/profile", func(c *gin.Context) {
		c.JSON(http.StatusOK, c.MustGet("user"))
	})

	g.GET("", paginate(), func(c *gin.Context) {
		pagination := c.MustGet("paginate").(database.Paginate)

		users, total, err := models.GetAllUsers(pagination)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results": users,
			"total":   total,
			"page":    c.MustGet("page"),
			"limit":   c.MustGet("limit"),
		})
	})

	g.PATCH("/:id/role", auth.AdminOnly(), func(c *gin.Context) {		
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		
		var body map[string]string
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		role, ok := body["role"]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "role is required"})
			return
		}
		
		updatedUser, err := models.UpdateUserRole(userID, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedUser)
	})

}
