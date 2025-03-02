package models

import (
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

	Status string `db:"status" json:"status"` //type -> org_status DEFAULT 'ACTIVE'

	VerifiedImpact bool `db:"verified_impact" json:"verified_impact"`
	Verified       bool `db:"verified" json:"verified"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (Organization) TableName() string {
	return "organizations"
}

func (Organization) FetchQuery() string {
	return "organizations/fetch"
}

func (*Organization) Create() error {
	return nil
}

func (*Organization) Update() error {
	return nil
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

func GetOrganization(id uuid.UUID, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, id.String()); err != nil {
		return nil, err
	}
	return o, nil
}

func getManyOrganizations(ids []uuid.UUID, identity uuid.UUID) ([]Organization, error) {
	result := []Organization{}
	return result, nil
}

func GetOrganizationByShortname(shortname string, identity uuid.UUID) (*Organization, error) {
	o := new(Organization)
	if err := database.Fetch(o, identity.String()); err != nil {
		return nil, err
	}
	return o, nil
}
