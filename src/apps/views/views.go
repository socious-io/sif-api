package views

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	authGroup(r)
	projectsGroup(r)
	identitiesGroup(r)
	usersGroup(r)
	mediaGroup(r)
	kybVerificationGroup(r)
	roundsGroup(r)
	organizationsGroup(r)
	syncGroup(r)
}
