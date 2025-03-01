package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type Identity struct {
	ID        uuid.UUID              `db:"id" json:"id"`
	Type      IdentityType           `db:"type" json:"type"`
	MetaMap   map[string]interface{} `db:"-" json:"meta"`
	Meta      types.JSONText         `db:"meta" json:"-"`
	CreatedAt time.Time              `db:"created_at" json:"created_at"`
}

func (Identity) TableName() string {
	return "identities"
}

func (Identity) FetchQuery() string {
	return "identities/fetch"
}

func GetIdentity(id uuid.UUID) (*Identity, error) {
	i := new(Identity)
	if err := database.Fetch(i, id); err != nil {
		return nil, err
	}
	return i, nil
}
