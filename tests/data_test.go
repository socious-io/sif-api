package tests_test

import (
	"sif/src/apps/models"

	"github.com/gin-gonic/gin"
)

var (
	users = []*models.User{
		{
			Username:  "test",
			Email:     "test@test.com",
			FirstName: "test",
			LastName:  "test",
		},
	}
	usersAuths []string

	projectsData = []gin.H{
		{
			"title":          "test",
			"description":    "test",
			"social_cause":   "test",
			"wallet_address": "test",
		},
	}
)
