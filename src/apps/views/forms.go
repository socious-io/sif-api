package views

import (
	"sif/src/apps/models"

	"github.com/google/uuid"
)

type ProjectForm struct {
	Title                 string                `json:"title" validate:"required"`
	Description           string                `json:"description" validate:"required"`
	Status                *models.ProjectStatus `json:"status"`
	City                  string                `json:"city"`
	Country               string                `json:"country"`
	SocialCause           string                `json:"social_cause" validate:"required"`
	CoverID               *uuid.UUID            `json:"cover_id"`
	WalletAddress         string                `json:"wallet_address" validate:"required"`
	WalletEnv             string                `json:"wallet_env"`
	Website               string                `json:"website"`
	LinkedIn              *string               `json:"linkedin,omitempty"`
	Video                 *string               `json:"video,omitempty"`
	ProblemStatement      *string               `json:"problem_statement,omitempty"`
	Solution              *string               `json:"solution,omitempty"`
	Goals                 *string               `json:"goals,omitempty"`
	TotalRequestedAmount  *int                  `json:"total_requested_amount,omitempty"`
	CostBreakdown         *string               `json:"cost_breakdown,omitempty"`
	ImpactAssessment      *string               `json:"impact_assessment,omitempty"`
	ImpactAssessmentType  *string               `json:"impact_assessment_type"`
	VoluntaryContribution *string               `json:"voluntary_contribution,omitempty"`
	Feasibility           *string               `json:"feasibility,omitempty"`
	Category              *string               `json:"category,omitempty"`
	Email                 string                `json:"email" validate:"required,email"`
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
	PaymentType   models.PaymentType `json:"payment_type"`
	WalletAddress string             `json:"wallet_address" validate:"required"`
	Currency      string             `json:"currency" validate:"required"`
	Rate          float64            `json:"rate"`
	Description   string             `json:"description"`
	Amount        float64            `json:"amount" validate:"required"`
	TxID          string             `json:"txid" validate:"required"`
	Meta          interface{}        `json:"meta" validate:"required"`
	CardToken     *string            `json:"card_token"`
	Anonymous     bool               `json:"anonymous"`
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
