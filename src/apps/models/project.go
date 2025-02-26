package models

import (
	"context"
	"time"

	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
)

type WorkSampleDocuments struct {
	Id       string `db:"id" json:"id"`
	Url      string `db:"url" json:"url"`
	Filename string `db:"filename" json:"filename"`
}

type WorkSampleType struct {
	ServiceID uuid.UUID `db:"service_id" json:"service_id"`
	Document  uuid.UUID `db:"document" json:"document"`
}

type Project struct {
	ID uuid.UUID `db:"id" json:"id"`

	Title       *string `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`

	Status ProjectStatus `db:"status" json:"status"`

	City    *string `db:"city" json:"city"`
	Country *string `db:"country" json:"country"`

	SocialCause string `db:"social_cause" json:"social_cause"`

	IdentityID   uuid.UUID      `db:"identity_id" json:"-"`
	Identity     *Identity      `db:"-" json:"identity"`
	IdentityJson types.JSONText `db:"identity" json:"-"`

	CoverID   uuid.UUID      `db:"cover_id" json:"-"`
	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`

	WalletAddress string    `db:"wallet_address" json:"wallet_address"`
	WalletEnv     WalletENV `db:"wallet_env" json:"wallet_env"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	ExpiresAt *time.Time `db:"expires_at" json:"expires_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
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
		// TODO: .....
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
		// TODO: .....
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

func GetProjects(identityId uuid.UUID, p database.Paginate) ([]Project, int, error) {
	var (
		projects  = []Project{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("projects/get", &fetchList, identityId, p.Limit, p.Offet); err != nil {
		return nil, 0, err
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
