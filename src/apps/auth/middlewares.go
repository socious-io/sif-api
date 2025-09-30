package auth

import (
	"net/http"
	"sif/src/apps/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		splited := strings.Split(tokenStr, " ")
		if len(splited) > 1 {
			tokenStr = splited[1]
		} else {
			tokenStr = splited[0]
		}
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		claims, err := VerifyToken(tokenStr)

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		u, err := models.GetUser(uuid.MustParse(claims.ID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", u)

		var identity *models.Identity

		identityStr := c.GetHeader(http.CanonicalHeaderKey("current-identity"))
		if identityUUID, err := uuid.Parse(identityStr); err == nil {
			identity, _ = models.GetIdentity(identityUUID)
		}
		if identity == nil {
			identity, _ = models.GetIdentity(u.ID)
		}

		if identity.Type == models.IdentityTypeOrganizations {
			if _, err := models.Member(identity.ID, u.ID); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Identitiy not allowed"})
				c.Abort()
				return
			}
		}

		c.Set("identity", identity)

		c.Next()
	}
}

func LoginOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		splited := strings.Split(tokenStr, " ")
		if len(splited) > 1 {
			tokenStr = splited[1]
		} else {
			tokenStr = splited[0]
		}
		if tokenStr == "" {
			c.Next()
			return
		}

		claims, err := VerifyToken(tokenStr)

		if err != nil {
			c.Next()
			return
		}

		u, err := models.GetUser(uuid.MustParse(claims.ID))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", u)

		var identity *models.Identity

		identityStr := c.GetHeader(http.CanonicalHeaderKey("current-identity"))
		if identityUUID, err := uuid.Parse(identityStr); err == nil {
			identity, _ = models.GetIdentity(identityUUID)
		}
		if identity == nil {
			identity, _ = models.GetIdentity(u.ID)
		}

		if identity.Type == models.IdentityTypeOrganizations {
			if _, err := models.Member(identity.ID, u.ID); err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Identitiy not allowed"})
				c.Abort()
				return
			}
		}

		c.Set("identity", identity)

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userI, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			c.Abort()
			return
		}

		user := userI.(*models.User)

		if user.Role != "ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}

