package rest

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/induzo/crud"
	"github.com/rs/xid"
	"github.com/induzo/gohttperror"
)

var (
	// ErrNotFound is used for Get and GetList when no entity is found
	ErrNotFound = errors.New("entity not found")
	// ErrBadRequest is used when the request is malformed or wrong
	ErrBadRequest = errors.New("bad request")
	// ErrForbidden is used when the request is malformed or wrong
	ErrForbidden = errors.New("forbidden")
)

// mgrMock is a mock for the mgr interface
type mgrMock struct {
	wantCreateError        bool
	wantDeleteError        bool
	wantGetError           bool
	wantGetListError       bool
	wantUpdateError        bool
	wantPartialUpdateError bool
	EntityList             map[xid.ID]*entityMock
}

// entityMock is respecting the CRUD interface
type entityMock struct {
	ID       xid.ID `json:"id"`
	StatusID int    `json:"status_id"`
}

func newMgrMock() *mgrMock {
	return &mgrMock{
		EntityList: make(map[xid.ID]*entityMock),
	}
}

func (m *mgrMock) NewEmptyEntity() interface{} {
	return &entityMock{}
}

func (m *mgrMock) Create(
	ctx context.Context,
	e interface{},
	pl io.Reader,
) (interface{}, error) {
	if m.wantCreateError {
		return nil, fmt.Errorf("Error create")
	}
	ec, ok := e.(*entityMock)
	if !ok {
		return nil,
			fmt.Errorf("mgrMock Create: impossible to cast e to entityMock")
	}
	ec.ID = xid.New()
	m.EntityList[ec.ID] = ec
	return ec, nil
}

func (m *mgrMock) Delete(ctx context.Context, id xid.ID) error {
	if m.wantDeleteError {
		return fmt.Errorf("Error delete")
	}
	if _, ok := m.EntityList[id]; !ok {
		return ErrNotFound
	}

	delete(m.EntityList, id)
	return nil
}

func (m *mgrMock) Get(ctx context.Context, id xid.ID) (interface{}, error) {
	if m.wantGetError {
		return nil, fmt.Errorf("Error get")
	}
	if ent, ok := m.EntityList[id]; ok {
		return ent, nil
	}
	return nil, ErrNotFound
}

func (m *mgrMock) GetList(
	context.Context,
	crud.ListModifiers,
) (interface{}, error) {
	if m.wantGetListError {
		return nil, fmt.Errorf("Error getlist")
	}
	if len(m.EntityList) == 0 {
		return nil, ErrNotFound
	}

	v := make([]*entityMock, 0, len(m.EntityList))
	for _, value := range m.EntityList {
		v = append(v, value)
	}

	return v, nil
}

func (m *mgrMock) Update(
	ctx context.Context,
	id xid.ID,
	newE interface{},
	pl io.Reader,
) (interface{}, error) {
	if m.wantUpdateError {
		return nil, fmt.Errorf("Error update")
	}
	_, ok := m.EntityList[id]
	if !ok {
		return nil, ErrNotFound
	}

	newEC, okC := newE.(*entityMock)
	if !okC {
		return ErrBadRequest, nil
	}
	m.EntityList[id] = newEC
	return newEC, nil
}

func (m *mgrMock) PartialUpdate(
	ctx context.Context,
	id xid.ID,
	pud crud.PartialUpdateData,
	pl io.Reader,
) error {
	if m.wantPartialUpdateError {
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

func (m *mgrMock) MapErrorToHTTPError(e error) *gohttperror.ErrResponse {
	switch e {
	case ErrNotFound:
		return gohttperror.ErrNotFound
	case ErrForbidden:
		return gohttperror.ErrForbidden(e)
	default:
		return gohttperror.ErrInternal(e)
	}
}
