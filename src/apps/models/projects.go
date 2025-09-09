package models

import (
	"context"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

/*
export interface Project {
  id: string;
  title: string;
  description: string;
  status: ProjectStatus;
  identity: Identity;
  cover: Media;
  total_donations: { [currency: string]: number };
  total_investments: { [currency: string]: number };
  school_name: string;
  school_size: number;
  kpw: number;
  kwh_per_year: number;
  co2_per_year: number;
  target_amount: number;
  created_at: Date;
  updated_at: Date;
  expires_at: Date;
  deleted_at: Date;
}
*/

type Project struct {
	ID uuid.UUID `db:"id" json:"id"`

	Title       *string `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`

	Status ProjectStatus `db:"status" json:"status"`

	IdentityID   uuid.UUID      `db:"identity_id" json:"-"`
	Identity     *Identity      `db:"-" json:"identity"`
	IdentityJson types.JSONText `db:"identity" json:"-"`

	CoverID   *uuid.UUID     `db:"cover_id" json:"cover_id"`
	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`

	TotalDonations   types.JSONText `db:"total_donations" json:"total_donations"`
	TotalInvestments types.JSONText `db:"total_investments" json:"total_investments"`

	TotalRequestedAmount float64 `db:"total_requested_amount" json:"total_requested_amount"`

	//Specifics What about others?
	SchoolName string  `db:"school_name" json:"school_name"`
	SchoolSize int     `db:"school_size" json:"school_size"`   //-> isn't this number or enum? (optional?)
	Kpw        float64 `db:"kpw" json:"kpw"`                   //-> isn't this number? (optional?)
	KwhPerYear float64 `db:"kwh_per_year" json:"kwh_per_year"` //-> isn't this number? (optional?)
	Co2PerYear float64 `db:"co2_per_year" json:"co2_per_year"` //-> isn't this number? (optional?)
	// TargetAmount string `db:"target_amount" json:"targetAmount"` //-> isn't this number? (optional?) -> use TotalRequestedAmount

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`

	SearchVector string `db:"search_vector"`
}

func (Project) TableName() string {
	return "projects"
}

func (Project) FetchQuery() string {
	return "projects/fetch"
}

type ProjectPreview struct {
	ID uuid.UUID `db:"id" json:"id"`

	Title       *string `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`

	SocialCause string `db:"social_cause" json:"social_cause"`

	IdentityID   uuid.UUID      `db:"identity_id" json:"-"`
	Identity     *Identity      `db:"-" json:"identity"`
	IdentityJson types.JSONText `db:"identity" json:"-"`

	CoverID   *uuid.UUID     `db:"cover_id" json:"cover_id"`
	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`
}

func (p *Project) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"projects/create",
		p.IdentityID,
		p.Title,
		p.Description,
		p.Status,
		p.CoverID,
		p.TotalRequestedAmount,
		p.SchoolName,
		p.SchoolSize,
		p.Kpw,
		p.KwhPerYear,
		p.Co2PerYear,
		p.ExpiresAt,
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
		p.CoverID,
		p.TotalRequestedAmount,
		p.SchoolName,
		p.SchoolSize,
		p.Kpw,
		p.KwhPerYear,
		p.Co2PerYear,
		p.ExpiresAt,
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
		projects  []Project
		fetchList []database.FetchList
		ids       []interface{}
	)

	search := ""
	for _, filter := range p.Filters {
		switch filter.Key {
		case "q":
			search = filter.Value
		}
	}

	if err := database.QuerySelect("projects/get_filtered", &fetchList, search, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) == 0 {
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
