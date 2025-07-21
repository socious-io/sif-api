package models

import (
	"context"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type Project struct {
	ID uuid.UUID `db:"id" json:"id"`

	Title       *string `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`

	Status ProjectStatus `db:"status" json:"status"`

	City    *string `db:"city" json:"city"`
	Country *string `db:"country" json:"country"`
	Website *string `db:"website" json:"website"`

	SocialCause string `db:"social_cause" json:"social_cause"`

	IdentityID   uuid.UUID      `db:"identity_id" json:"-"`
	Identity     *Identity      `db:"-" json:"identity"`
	IdentityJson types.JSONText `db:"identity" json:"-"`

	CoverID   *uuid.UUID     `db:"cover_id" json:"cover_id"`
	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`

	RoundID   uuid.UUID      `db:"round_id" json:"-"`
	Round     *Round         `db:"-" json:"round"`
	RoundJson types.JSONText `db:"round" json:"-"`

	TotalVotes     int            `db:"total_votes" json:"total_votes"`
	TotalDonations types.JSONText `db:"total_donations" json:"total_donations"`
	// TotalDonations float64 `db:"total_donations" json:"total_donations"`

	WalletAddress string    `db:"wallet_address" json:"wallet_address"`
	WalletEnv     WalletENV `db:"wallet_env" json:"wallet_env"`

	LinkedIn              *string               `db:"linkedin" json:"linkedin"`
	Video                 *string               `db:"video" json:"video"`
	ProblemStatement      *string               `db:"problem_statement" json:"problem_statement"`
	Solution              *string               `db:"solution" json:"solution"`
	Goals                 *string               `db:"goals" json:"goals"`
	TotalRequestedAmount  float64               `db:"total_requested_amount" json:"total_requested_amount"`
	CostBreakdown         *string               `db:"cost_beakdown" json:"cost_breakdown"`
	ImpactAssessment      *string               `db:"impact_assessment" json:"impact_assessment"`
	ImpactAssessmentType  *ImpactAssessmentType `db:"impact_assessment_type" json:"impact_assessment_type"`
	VoluntaryContribution *string               `db:"voluntery_contribution" json:"voluntery_contribution"`
	Feasibility           *string               `db:"feasibility" json:"feasibility"`
	Email                 *string               `db:"email" json:"email"`

	Category ProjectCategory `db:"category" json:"category"`

	UserVoted bool `db:"-" json:"user_voted"`

	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	ExpiresAt     *time.Time `db:"expires_at" json:"expires_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
	NotEligibleAt *time.Time `db:"not_eligible_at" json:"not_eligible_at"`
}

func (Project) TableName() string {
	return "projects"
}

func (Project) FetchQuery() string {
	return "projects/fetch"
}

func (p *Project) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"projects/create",
		p.Title,
		p.Description,
		p.Status,
		p.City,
		p.Country,
		p.SocialCause,
		p.IdentityID,
		p.CoverID,
		p.WalletAddress,
		p.WalletEnv,
		p.Website,
		p.LinkedIn,
		p.Video,
		p.ProblemStatement,
		p.Solution,
		p.Goals,
		p.TotalRequestedAmount,
		p.CostBreakdown,
		p.ImpactAssessment,
		p.ImpactAssessmentType,
		p.VoluntaryContribution,
		p.Feasibility,
		p.Category,
		p.Email,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(p); err != nil {
			return err
		}
	}
	return database.Fetch(p, p.ID)
}

func (p *Project) Update(ctx context.Context) error {

	rows, err := database.Query(
		ctx,
		"projects/update",
		p.ID,
		p.Title,
		p.Description,
		p.Status,
		p.City,
		p.Country,
		p.SocialCause,
		p.CoverID,
		p.WalletAddress,
		p.WalletEnv,
		p.Website,
		p.LinkedIn,
		p.Video,
		p.ProblemStatement,
		p.Solution,
		p.Goals,
		p.TotalRequestedAmount,
		p.CostBreakdown,
		p.ImpactAssessment,
		p.ImpactAssessmentType,
		p.VoluntaryContribution,
		p.Feasibility,
		p.Category,
		p.Email,
	)

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(p); err != nil {
			return err
		}
	}
	return database.Fetch(p, p.ID)
}

func (p *Project) Delete(ctx context.Context) error {
	rows, err := database.Query(ctx, "projects/delete", p.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func GetProjects(p database.Paginate) ([]Project, int, error) {
	var (
		projects  = []Project{}
		fetchList []database.FetchList
		ids       []interface{}
	)
	if len(p.Filters) > 0 {
		var identityID, roundID, category string
		for _, filter := range p.Filters {
			if filter.Key == "identity_id" || filter.Key == "identity" {
				identityID = filter.Value
			}
			if filter.Key == "round_id" {
				roundID = filter.Value
			}
			if filter.Key == "category" {
				category = filter.Value
			}
		}
		if identityID != "" {
			if err := database.QuerySelect("projects/get_by_identity", &fetchList, identityID, p.Limit, p.Offet); err != nil {
				return nil, 0, err
			}
		} else if roundID != "" {
			if err := database.QuerySelect("projects/get_by_round", &fetchList, roundID, p.Limit, p.Offet); err != nil {
				return nil, 0, err
			}
		} else if category != "" {
			if err := database.QuerySelect("projects/get_by_category", &fetchList, category, p.Limit, p.Offet); err != nil {
				return nil, 0, err
			}
		}
	} else {
		if err := database.QuerySelect("projects/get", &fetchList, p.Limit, p.Offet); err != nil {
			return nil, 0, err
		}
	}

	if len(fetchList) < 1 {
		return projects, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&projects, ids...); err != nil {
		return nil, 0, err
	}

	return projects, fetchList[0].TotalCount, nil
}

func GetProject(id uuid.UUID) (*Project, error) {
	p := new(Project)
	if err := database.Fetch(p, id); err != nil {
		return nil, err
	}
	return p, nil
}
