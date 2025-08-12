package views

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sif/src/apps/auth"
	"sif/src/apps/models"
	"sif/src/apps/utils"
	"sif/src/config"
	"strconv"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	database "github.com/socious-io/pkg_database"
	"github.com/unrolled/secure"

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
		isAdmin := access_token == config.Config.AdminToken

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
		id := c.Request.Header.Get("x-account-center-id")
		secret := c.Request.Header.Get("x-account-center-secret")
		hash, _ := auth.HashPassword(secret)

		if id != config.Config.GoAccounts.ID || auth.CheckPasswordHash(secret, hash) != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Account center required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func SecureHeaders(env string) gin.HandlerFunc {

	IsDevelopment := env != "production"
	options := secure.Options{
		FrameDeny:          true, // X-Frame-Options: DENY
		ContentTypeNosniff: true, // X-Content-Type-Options: nosniff
		BrowserXssFilter:   true, // X-XSS-Protection: 1; mode=block (legacy)
		// ReferrerPolicy:        "no-referrer",
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' $NONCE; img-src 'self' https: http:;", // Very important for XSS
		// HSTS:
		SSLRedirect:          true,
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           31536000,
		STSIncludeSubdomains: true,
		STSPreload:           true,
		IsDevelopment:        IsDevelopment,
	}

	return func(c *gin.Context) {
		s := secure.New(options)
		nonce, err := s.ProcessAndReturnNonce(c.Writer, c.Request)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		c.Set("nonce", nonce)

		c.Next()
	}
}

func SecureRequest(p *bluemonday.Policy) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check content type
		isUrlEncodedContent := strings.Contains(c.GetHeader("Content-Type"), "application/x-www-form-urlencoded")
		isMultipartContent := strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data")
		isJsonContent := strings.Contains(c.GetHeader("Content-Type"), "application/json")

		// --- 1. Sanitize Query Parameters ---
		q := c.Request.URL.Query()
		utils.SanitizeURLValues(q, p)
		c.Request.URL.RawQuery = q.Encode()

		// --- 2. Sanitize Form Data (application/x-www-form-urlencoded or multipart) ---
		if isUrlEncodedContent || isMultipartContent {
			if err := c.Request.ParseForm(); err != nil {
				c.AbortWithStatusJSON(400, gin.H{
					"error": fmt.Sprintf("Invalid body payload, err: %v", err),
				})
				return
			}
			utils.SanitizeURLValues(c.Request.PostForm, p)
		} else if isJsonContent {
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
			}

			if len(bodyBytes) > 0 {
				var data map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &data); err != nil {
					c.AbortWithStatusJSON(400, gin.H{
						"error": fmt.Sprintf("Invalid body payload, err: %v", err),
					})
					return
				}

				utils.SanitizeMap(data, p)
				safeBody, _ := json.Marshal(data)
				c.Request.Body = io.NopCloser(bytes.NewReader(safeBody))
			}
		}
		c.Next()
	}
}
