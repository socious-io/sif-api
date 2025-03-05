package views

import (
	"sif/src/apps/models"

	"github.com/google/uuid"
)

type ProjectForm struct {
	Title         string                `json:"title" validate:"required"`
	Description   string                `json:"description" validate:"required"`
	Status        *models.ProjectStatus `json:"status"`
	City          string                `json:"city"`
	Country       string                `json:"country"`
	SocialCause   string                `json:"social_cause" validate:"required"`
	CoverID       *uuid.UUID            `json:"cover_id"`
	WalletAddress string                `json:"wallet_address" validate:"required"`
	WalletEnv     string                `json:"wallet_env"`
}

type AuthForm struct {
	RedirectURL string `json:"redirect_url" validate:"required"`
}

type SessionForm struct {
	Code string `json:"code" validate:"required"`
}

type RefreshForm struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
