package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	database "github.com/socious-io/pkg_database"
)

type Donation struct {
	ID uuid.UUID `json:"id" db:"id"`

	UserID   uuid.UUID      `json:"user_id" db:"user_id"`
	User     User           `json:"user" db:"-"`
	UserJson types.JSONText `json:"-" db:"user"`

	ProjectID uuid.UUID `json:"project_id" db:"project_id"`

	Currency  string  `json:"currency" db:"currency"`
	Amount    float64 `json:"amount" db:"amount"`
	Rate      float64 `json:"rate" db:"rate"`
	Anonymous bool    `json:"anonymous" db:"anonymous"`

	Status               DonationStatus `json:"status" db:"status"`
	TransactionID        string         `json:"transaction_id" db:"transaction_id"`
	ReleaseTransactionID *string        `json:"release_transaction_id" db:"release_transaction_id"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (Donation) TableName() string {
	return "donations"
}

func (Donation) FetchQuery() string {
	return "donations/fetch"
}

func (d *Donation) Create(ctx context.Context) error {
	rows, err := database.Query(
		ctx,
		"donations/create",
		d.UserID,
		d.ProjectID,
		d.Amount,
		d.Currency,
		d.Status,
		d.Anonymous,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(d); err != nil {
			return err
		}
	}
	return database.Fetch(d, d.ID)

}

func (d *Donation) Update(ctx context.Context) error {

	rows, err := database.Query(
		ctx,
		"donations/update",
		d.ID,
		d.Status,
		d.TransactionID,
		d.ReleaseTransactionID,
	)

	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.StructScan(d); err != nil {
			return err
		}
	}
	return database.Fetch(d, d.ID)
}

func GetDonation(id interface{}) (*Donation, error) {
	d := new(Donation)
	if err := database.Fetch(d, id); err != nil {
		return nil, err
	}
	return d, nil
}

func GetDonations(projectID interface{}, p database.Paginate) ([]Donation, int, error) {
	var (
		donations = []Donation{}
		fetchList []database.FetchList
		ids       []interface{}
	)
	if err := database.QuerySelect("donations/get", &fetchList, projectID, p.Limit, p.Offet); err != nil {
		return nil, 0, err
	}

	if len(fetchList) < 1 {
		return donations, 0, nil
	}

	for _, f := range fetchList {
		ids = append(ids, f.ID)
	}

	if err := database.Fetch(&donations, ids...); err != nil {
		return nil, 0, err
	}

	return donations, fetchList[0].TotalCount, nil
}
