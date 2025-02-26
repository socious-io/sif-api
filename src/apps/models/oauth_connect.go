package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type OauthConnect struct {
	ID             uuid.UUID               `db:"id" json:"id"`
	IdentityId     uuid.UUID               `db:"identity_id" json:"identity_id"`
	Provider       OauthConnectedProviders `db:"provider" json:"provider"`
	MatrixUniqueId string                  `db:"matrix_unique_id" json:"matrix_unique_id"`
	AccessToken    string                  `db:"access_token" json:"access_token"`
	RefreshToken   *string                 `db:"refresh_token" json:"refresh_token"`
	Meta           *types.JSONText         `db:"meta" json:"meta"`
	ExpiredAt      *time.Time              `db:"expired_at" json:"expired_at"`
	CreatedAt      time.Time               `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time               `db:"updated_at" json:"updated_at"`
}

func (OauthConnect) TableName() string {
	return "oauth_connects"
}

func (OauthConnect) FetchQuery() string {
	return "oauth_connects/fetch"
}

func GetOauthConnectByIdentityId(identityId uuid.UUID, provider OauthConnectedProviders) (*OauthConnect, error) {
	oc := new(OauthConnect)
	if err := database.Get(oc, "oauth_connects/get_by_identityid", identityId, provider); err != nil {
		return nil, err
	}
	return oc, nil
}
