package models

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`

	Email string `db:"email" json:"email"`

	City    *string `db:"city" json:"city"`
	Country *string `db:"country" json:"country"`
	Address *string `db:"address" json:"address"`

	Avatar     *Media         `db:"-" json:"avatar"`
	AvatarJson types.JSONText `db:"avatar" json:"-"`

	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`

	Language *string `db:"language" json:"language"`

	ImpactPoints      float32 `db:"impact_points" json:"impact_points"`
	Donates           float32 `db:"donates" json:"donates"`
	ProjectsSupported int     `db:"project_supported" json:"project_supported"`
	StripeCustomerID  *string `db:"stripe_customer_id" json:"stripe_customer_id"`

	ReferredBy *uuid.UUID `db:"referred_by" json:"referred_by"`

	IdentityVerifiedAt *time.Time `db:"identity_verified_at" json:"identity_verified_at"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (User) FetchQuery() string {
	return "users/fetch"
}

func (u *User) Upsert(ctx context.Context) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Avatar != nil {
		b, _ := json.Marshal(u.Avatar)
		u.AvatarJson.Scan(b)
	}
	if u.Cover != nil {
		b, _ := json.Marshal(u.Cover)
		u.CoverJson.Scan(b)
	}
	rows, err := database.Query(
		ctx,
		"users/upsert",
		u.ID,
		u.FirstName, u.LastName, u.Username, u.Email,
		u.City, u.Country, u.AvatarJson, u.CoverJson,
		u.Language, u.ImpactPoints, u.IdentityVerifiedAt,
		u.StripeCustomerID, u.ReferredBy,
	)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(u); err != nil {
			return err
		}
	}
	return nil
}

func GetUser(id uuid.UUID) (*User, error) {
	u := new(User)
	if err := database.Fetch(u, id.String()); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByOrg(orgId uuid.UUID) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_org", orgId); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByEmail(email string) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_email", email); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByUsername(username string) (*User, error) {
	u := new(User)
	if err := database.Get(u, "users/fetch_by_username", username); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) Delete(ctx context.Context) error {
	if _, err := database.Query(
		ctx,
		"users/delete",
		u.ID,
	); err != nil {
		return err
	}

	return nil
}
