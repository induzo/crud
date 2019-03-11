package crud

import (
	"context"
	"io"

	"github.com/rs/xid"
	"github.com/induzo/gohttperror"
)

// MgrI is the interface to initialize the new entity mgr
type MgrI interface {
	NewEmptyEntity() interface{}
	Create(context.Context, interface{}, io.Reader) (interface{}, error)
	Delete(context.Context, xid.ID) error
	Get(context.Context, xid.ID) (interface{}, error)
	GetList(context.Context, ListModifiers) (interface{}, error)
	Update(context.Context, xid.ID, interface{}, io.Reader) (interface{}, error)
	PartialUpdate(context.Context, xid.ID, PartialUpdateData, io.Reader) error
	MapErrorToHTTPError(error) *gohttperror.ErrResponse
}
