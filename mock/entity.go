package mock

import (
	"github.com/rs/xid"
)

// Entity is respecting the CRUD interface
type Entity struct {
	ID       xid.ID `json:"id"`
	StatusID int    `json:"status_id"`
}
