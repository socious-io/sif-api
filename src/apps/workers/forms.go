package workers

import (
	"sif/src/apps/models"
)

type SyncForm struct {
	Organizations []models.Organization `json:"organizations"`
	User          models.User           `json:"user" validate:"required"`
}

type DeleteUserForm struct {
	User   models.User `json:"user" validate:"required"`
	Reason string      `json:"reason" validate:"required"`
}
