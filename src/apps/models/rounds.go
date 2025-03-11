package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

// Round represent a round of event
type Round struct {
	ID uuid.UUID `db:"id" json:"id"`

	Name       string `db:"name" json:"name"`
	PoolAmount int    `db:"pool_amount" json:"pool_amount"`

	CoverID   *uuid.UUID     `db:"cover_id" json:"cover_id"`
	Cover     *Media         `db:"-" json:"cover"`
	CoverJson types.JSONText `db:"cover" json:"-"`

	VotingStartAt time.Time `db:"voting_start_at" json:"voting_start_at"`
	VotingEndAt   time.Time `db:"voting_end_at" json:"voting_end_at"`

	SubmissionStartAt time.Time `db:"submission_start_at" json:"submission_start_at"`
	SubmissionEndAt   time.Time `db:"submission_end_at" json:"submission_end_at"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// TableName returns the table name for the Round model
func (Round) TableName() string {
	return "rounds"
}

func (Round) Fetch() string {
	return "rounds/fetch"
}

func GetRoundLatestRound() (*Round, error) {
	r := new(Round)
	if err := database.Get(r, "rounds/get_latest_round"); err != nil {
		return nil, err
	}
	return r, nil
}
