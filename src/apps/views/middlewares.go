package views

import (
	"fmt"
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"
	"sif/src/config"
	"strconv"
	"strings"

	database "github.com/socious-io/pkg_database"

	"github.com/gin-gonic/gin"
)

func paginate() gin.HandlerFunc {
	return func(c *gin.Context) {

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			page = 1
		}

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		}
		if page < 1 {
			page = 1
		}
		if limit > 100 || limit < 1 {
			limit = 10
		}
		filters := make([]database.Filter, 0)
		for key, values := range c.Request.URL.Query() {
			if strings.Contains(key, "filter.") && len(values) > 0 {
				filters = append(filters, database.Filter{
					Key:   strings.Replace(key, "filter.", "", -1),
					Value: values[0],
				})
			}
		}

		c.Set("paginate", database.Paginate{
			Limit:   limit,
			Offet:   (page - 1) * limit,
			Filters: filters,
		})
		c.Set("limit", limit)
		c.Set("page", page)
		c.Next()
	}
}

// Administration
func adminAccessRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		access_token := c.Query("admin_access_token")
		isAdmin := access_token == config.Config.Admin.AccessToken

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func OrganizationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		if identity.Type != models.IdentityTypeOrganizations {
			c.JSON(http.StatusForbidden, gin.H{"error": "Organization identity required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AccountCenterRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := fmt.Sprintf("%s:%s", config.Config.GoAccounts.ID, config.Config.GoAccounts.Secret)
		hash, _ := auth.HashPassword(raw)
		if hash != c.Request.Header.Get("x-account-center") {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account center required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
