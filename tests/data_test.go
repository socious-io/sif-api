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
		{
			Username:  "test2",
			Email:     "test2@test.com",
			FirstName: "test2",
			LastName:  "test2",
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

	commentsData = []gin.H{
		{"content": "test"},
	}
)
