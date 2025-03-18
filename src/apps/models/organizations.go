package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Shortname   string    `db:"shortname" json:"shortname"`
	Name        *string   `db:"name" json:"name"`
	Bio         *string   `db:"bio" json:"bio"`
	Description *string   `db:"description" json:"description"`
	Email       *string   `db:"email" json:"email"`
	Phone       *string   `db:"phone" json:"phone"`

	City    *string `db:"city" json:"city"`
	Country *string `db:"country" json:"country"`
	Address *string `db:"address" json:"address"`
	Website *string `db:"website" json:"website"`

	Mission *string `db:"mission" json:"mission"`
	Culture *string `db:"culture" json:"culture"`

	Logo     *Media         `db:"-" json:"logo"`
	LogoJson types.JSONText `db:"logo" json:"-"`

	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`

	Status OrganizationStatus `db:"status" json:"status"`

	VerifiedImpact bool `db:"verified_impact" json:"verified_impact"`
	Verified       bool `db:"verified" json:"verified"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type OrganizationMember struct {
	ID             uuid.UUID `db:"id" json:"id"`
	OrganizationID uuid.UUID `db:"organization_id" json:"organization_id"`
	UserID         uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (Organization) FetchQuery() string {
	return "organizations/fetch"
}

func (om *OrganizationMember) Create(ctx context.Context) error {
	_, err := database.Query(ctx, "organizations/create_member", om.OrganizationID, om.UserID)
	return err
}

func (o *Organization) Create(ctx context.Context, userID uuid.UUID) error {

	tx, err := database.GetDB().Beginx()
	if err != nil {
		return err
	}

	if o.Logo != nil {
		b, _ := json.Marshal(o.Logo)
		o.LogoJson.Scan(b)
	}
	if o.Cover != nil {
		b, _ := json.Marshal(o.Cover)
		o.CoverJson.Scan(b)
	}
	if o.ID == uuid.Nil {
		newID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		o.ID = newID
	}

	if o.Verified || o.VerifiedImpact {
		o.Status = OrganizationStatusActive
	} else {
		o.Status = OrganizationStatusNotActive
	}
	rows, err := database.TxQuery(
		ctx,
		tx,
		"organizations/create",
		o.ID,
		o.Shortname,
		o.Name,
		o.Bio,
		o.Description,
		o.Email,
		o.Phone,
		o.City,
		o.Country,
		o.Address,
		o.Website,
		o.Mission,
		o.Culture,
		o.Status,
		o.VerifiedImpact,
		o.Verified,
		o.LogoJson,
		o.CoverJson,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(o); err != nil {
			tx.Rollback()
			return err
		}
	}

	if _, err := database.TxQuery(
		ctx,
		tx,
		"organizations/add_member",
		o.ID,
		userID,
	); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return database.Fetch(o, o.ID)
}

func (o *Organization) Update(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"organizations/update",
		o.ID,
		o.Shortname,
		o.Name,
		o.Bio,
		o.Description,
		o.Email,
		o.Phone,
		o.City,
		o.Country,
		o.Address,
		o.Website,
		o.Mission,
		o.Culture,
		o.Logo,
		o.Cover,
		o.Status,
		o.VerifiedImpact,
		o.Verified,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(o); err != nil {
			return err
		}
	}
	return database.Fetch(o, o.ID)
}

func (*Organization) Remove() error {
	return nil
}

func (*Organization) UpdateDID() error {
	return nil
}

func (*Organization) ToggleHiring() error {
	return nil
}

func GetOrganization(id uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, id.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func Member(orgID, userID uuid.UUID) (*OrganizationMember, error) {
	om := new(OrganizationMember)
	if err := database.Get(om, "organizations/get_member", orgID, userID); err != nil {
		return nil, err
	}
	return om, nil
}

func GetUserOrganizations(userId uuid.UUID) ([]Organization, error) {
	var (
		orgs      = []Organization{}
		fetchList []database.FetchList
		ids       []interface{}
	)

	if err := database.QuerySelect("organizations/get_by_member", &fetchList, userId); err != nil {
		return orgs, err
	}

	if len(fetchList) < 1 {
		return orgs, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&orgs, ids...); err != nil {
		return orgs, err
	}
	return orgs, nil
}
