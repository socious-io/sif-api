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
	Website       string                `json:"website"`
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

type KYBVerificationForm struct {
	Documents []string `json:"documents"`
}

type ApikeyForm struct {
	Name string `json:"name"`
}

type DnateDepositForm struct {
	Currency    string      `json:"currency" validate:"required"`
	Description string      `json:"description"`
	Amount      float64     `json:"amount" validate:"required"`
	TxID        string      `json:"txid" validate:"required"`
	Meta        interface{} `json:"meta" validate:"required"`
}

type SyncForm struct {
	Organizations []models.Organization `json:"organization"`
	User          models.User           `json:"users" validate:"required"`
}
