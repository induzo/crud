package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/induzo/crud"

	"github.com/go-chi/render"
	"github.com/induzo/gohttperror"
)

// POSTHandler will handle data from request
// and returns bytes to be written to response
func POSTHandler(
	cmgr crud.MgrI,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("POSTHandler Render: %v", errRender)
			}
		}()

		var payload bytes.Buffer
		pl := io.TeeReader(r.Body, &payload)

		ent := cmgr.NewEmptyEntity()
		if errJSON := json.NewDecoder(pl).Decode(ent); errJSON != nil {
			errRender = render.Render(w, r,
				gohttperror.ErrBadRequest(errJSON),
			)
			return
		}

		e, err := cmgr.Create(r.Context(), ent, &payload)
		if err != nil {
			errRender = render.Render(w, r, gohttperror.ErrInternal(err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.DefaultResponder(w, r, e)
	}
}

// GETListHandler will handle data from request
// and returns bytes to be written to response
func GETListHandler(
	cmgr crud.MgrI,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("GETListHandler Render: %v", errRender)
			}
		}()

		es, errGL := cmgr.GetList(
			r.Context(),
			ListModifiersFromURL(r.URL),
		)
		if errGL != nil {
			errRender = render.Render(w, r, cmgr.MapErrorToHTTPError(errGL))
			return
		}

		render.DefaultResponder(w, r, es)
	}
}

// GETHandler returns a unique entity
func GETHandler(
	cmgr crud.MgrI,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("GETHandler Render: %v", errRender)
			}
		}()

		ID, errParse := parseIDFromRequest(r)
		if errParse != nil {
			errRender = render.Render(w, r,
				gohttperror.ErrBadRequest(errParse),
			)
			return
		}

		e, errG := cmgr.Get(r.Context(), ID)
		if errG != nil {
			errRender = render.Render(w, r, cmgr.MapErrorToHTTPError(errG))
			return
		}

		render.DefaultResponder(w, r, e)
	}
}

// DELETEHandler will delete a specific entity
func DELETEHandler(
	cmgr crud.MgrI,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("DELETEHandler Render: %v", errRender)
			}
		}()

		ID, errParse := parseIDFromRequest(r)
		if errParse != nil {
			errRender = render.Render(w, r,
				gohttperror.ErrBadRequest(errParse),
			)
			return
		}

		if err := cmgr.Delete(
			r.Context(),
			ID,
		); err != nil {
			errRender = render.Render(w, r, cmgr.MapErrorToHTTPError(err))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

// PUTHandler will update all data for a specific entity
func PUTHandler(
	cmgr crud.MgrI,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("PUTHandler Render: %v", errRender)
			}
		}()

		ID, errParse := parseIDFromRequest(r)
		if errParse != nil {
			errRender = render.Render(w, r,
				gohttperror.ErrBadRequest(errParse),
			)
			return
		}

		var payload bytes.Buffer
		pl := io.TeeReader(r.Body, &payload)

		ent := cmgr.NewEmptyEntity()
		if errJSON := json.NewDecoder(pl).Decode(&ent); errJSON != nil {
			errRender = render.Render(w, r,
				gohttperror.ErrBadRequest(errJSON),
			)
			return
		}

		e, errU := cmgr.Update(
			r.Context(),
			ID,
			ent,
			&payload,
		)
		if errU != nil {
			errRender = render.Render(w, r, cmgr.MapErrorToHTTPError(errU))
			return
		}

		w.WriteHeader(http.StatusOK)
		render.DefaultResponder(w, r, e)
	}
}

// PATCHHandler will update specific data for a specific entity
// Following https://tools.ietf.org/html/rfc7386
// Content-Type: application/merge-patch+json
func PATCHHandler(
	cmgr crud.MgrI,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var errRender error
		defer func() {
			if errRender != nil {
				log.Fatalf("PATCHHandler Render: %v", errRender)
			}
		}()

		ID, errParse := parseIDFromRequest(r)
		if errParse != nil {
			errRender = render.Render(w, r,
				gohttperror.ErrBadRequest(errParse),
			)
			return
		}

		var payload bytes.Buffer
		pl := io.TeeReader(r.Body, &payload)

		updates := crud.PartialUpdateData{}
		if errJSON := json.NewDecoder(pl).Decode(&updates); errJSON != nil {
			errRender = render.Render(w, r,
				gohttperror.ErrBadRequest(errJSON),
			)
			return
		}

		if err := cmgr.PartialUpdate(
			r.Context(),
			ID,
			updates,
			&payload,
		); err != nil {
			errRender = render.Render(w, r, cmgr.MapErrorToHTTPError(err))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
