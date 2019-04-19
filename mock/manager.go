package mock

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/induzo/crud"
	"github.com/induzo/gohttperror"
	"github.com/rs/xid"
)

var (
	// ErrNotFound is used for Get and GetList when no entity is found
	ErrNotFound = errors.New("entity not found")
	// ErrBadRequest is used when the request is malformed or wrong
	ErrBadRequest = errors.New("bad request")
	// ErrForbidden is used when the request is malformed or wrong
	ErrForbidden = errors.New("forbidden")
)

// Mgr is a mock for the mgr interface
type Mgr struct {
	WantCreateError        bool
	WantDeleteError        bool
	WantGetError           bool
	WantGetListError       bool
	WantUpdateError        bool
	WantPartialUpdateError bool
	EntityList             map[xid.ID]*Entity
}

func NewMgr() *Mgr {
	return &Mgr{
		EntityList: make(map[xid.ID]*Entity),
	}
}

func (m *Mgr) NewEmptyEntity() interface{} {
	return &Entity{}
}

func (m *Mgr) Create(
	ctx context.Context,
	e interface{},
	pl io.Reader,
) (interface{}, error) {
	if m.WantCreateError {
		return nil, fmt.Errorf("Error create")
	}
	ec, ok := e.(*Entity)
	if !ok {
		return nil,
			fmt.Errorf("Mgr Create: impossible to cast e to Entity")
	}
	ec.ID = xid.New()
	m.EntityList[ec.ID] = ec
	return ec, nil
}

func (m *Mgr) Delete(ctx context.Context, id xid.ID) error {
	if m.WantDeleteError {
		return fmt.Errorf("Error delete")
	}
	if _, ok := m.EntityList[id]; !ok {
		return ErrNotFound
	}

	delete(m.EntityList, id)
	return nil
}

func (m *Mgr) Get(ctx context.Context, id xid.ID) (interface{}, error) {
	if m.WantGetError {
		return nil, fmt.Errorf("Error get")
	}
	if ent, ok := m.EntityList[id]; ok {
		return ent, nil
	}
	return nil, ErrNotFound
}

func (m *Mgr) GetList(
	context.Context,
	crud.ListModifiers,
) (interface{}, error) {
	if m.WantGetListError {
		return nil, fmt.Errorf("Error getlist")
	}
	if len(m.EntityList) == 0 {
		return nil, ErrNotFound
	}

	v := make([]*Entity, 0, len(m.EntityList))
	for _, value := range m.EntityList {
		v = append(v, value)
	}

	return v, nil
}

func (m *Mgr) Update(
	ctx context.Context,
	id xid.ID,
	newE interface{},
	pl io.Reader,
) (interface{}, error) {
	if m.WantUpdateError {
		return nil, fmt.Errorf("Error update")
	}
	_, ok := m.EntityList[id]
	if !ok {
		return nil, ErrNotFound
	}

	newEC, okC := newE.(*Entity)
	if !okC {
		return ErrBadRequest, nil
	}
	m.EntityList[id] = newEC
	return newEC, nil
}

func (m *Mgr) PartialUpdate(
	ctx context.Context,
	id xid.ID,
	pud crud.PartialUpdateData,
	pl io.Reader,
) error {
	if m.WantPartialUpdateError {
		return fmt.Errorf("Error partial update")
	}
	if _, ok := pud["status_id"]; !ok {
		return ErrBadRequest
	}
	if _, ok := m.EntityList[id]; !ok {
		return ErrNotFound
	}

	m.EntityList[id].StatusID = int(pud["status_id"].(float64))

	return nil
}

func (m *Mgr) MapErrorToHTTPError(e error) *gohttperror.ErrResponse {
	switch e {
	case ErrNotFound:
		return gohttperror.ErrNotFound
	case ErrForbidden:
		return gohttperror.ErrForbidden(e)
	default:
		return gohttperror.ErrInternal(e)
	}
}
