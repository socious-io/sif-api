package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

type Vote struct {
	ID uuid.UUID `db:"id" json:"id"`

	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	ProjectID uuid.UUID `db:"project_id" json:"project_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func (Vote) TableName() string {
	return "votes"
}

func (Vote) FetchQuery() string {
	return "votes/fetch"
}

func (v *Vote) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"votes/create",
		v.UserID,
		v.ProjectID,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(v); err != nil {
			return err
		}
	}
	return database.Fetch(v, v.ID)
}

func GetVoteByUserAndProject(userID, projectID uuid.UUID) (*Vote, error) {
	v := new(Vote)
	if err := database.Get(v, "votes/get_by_user_and_project", userID, projectID); err != nil {
		return nil, err
	}
	return v, nil
}
