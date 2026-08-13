package models

import "github.com/lib/pq"

// Form represents a saved, embeddable public subscription form tied to a
// set of lists.
type Form struct {
	Base

	UUID        string        `db:"uuid" json:"uuid"`
	Name        string        `db:"name" json:"name"`
	ListIDs     pq.Int64Array `db:"list_ids" json:"list_ids"`
	RedirectURL string        `db:"redirect_url" json:"redirect_url"`
}
