package views

import (
	"sif/src/apps/models"
	"time"

	"github.com/google/uuid"
)

type ProjectForm struct {
	Title                string                `json:"title" validate:"required"`
	Description          string                `json:"description" validate:"required"`
	Status               *models.ProjectStatus `json:"status"`
	CoverID              *uuid.UUID            `json:"cover_id"`
	TotalRequestedAmount float64               `json:"total_requested_amount"`
	SchoolName           string                `json:"school_name"`
	SchoolSize           int                   `json:"school_size"`
	Kpw                  float64               `json:"kpw"` // kWp capacity
	KwhPerYear           float64               `json:"kwh_per_year"`
	Co2PerYear           float64               `json:"co2_per_year"`
	ExpiresAt            time.Time             `json:"expires_at,omitempty"`
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
	PaymentType   models.PaymentType        `json:"payment_type"`
	WalletAddress string                    `json:"wallet_address" validate:"required"`
	Currency      string                    `json:"currency" validate:"required"`
	Rate          float64                   `json:"rate"`
	Description   string                    `json:"description"`
	Amount        float64                   `json:"amount" validate:"required"`
	TxID          string                    `json:"txid" validate:"required"`
	Meta          interface{}               `json:"meta" validate:"required"`
	CardToken     *string                   `json:"card_token"`
	Anonymous     bool                      `json:"anonymous"`
	PaidAs        models.DonationPaidAsType `json:"paid_as"`
}

type DonateDepositConfirmForm struct {
	PaymentIntentID string `json:"payment_intent_id"`
}

type SyncForm struct {
	Organizations []models.Organization `json:"organizations"`
	User          models.User           `json:"user" validate:"required"`
}

type CommentForm struct {
	Content  string     `json:"content" validate:"required"`
	ParentID *uuid.UUID `json:"parent_id"`
	MediaID  *uuid.UUID `json:"media_id"`
}

type ReactionForm struct {
	Reaction string `json:"reaction" validate:"required"`
}
