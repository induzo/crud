package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rs/xid"
)

func parseIDFromRequest(r *http.Request) (xid.ID, error) {
	idStr := chi.URLParam(r, "ID")
	ID, errConv := xid.FromString(idStr)
	if ID.IsNil() || errConv != nil {
		return xid.NilID(),
			fmt.Errorf("parseIDFromRequest(%s): %v", idStr, errConv)
	}

	return ID, nil
}

// GetTestContextWithID will return a context with an xid as chi URL Params
func GetTestContextWithID(
	ctx context.Context,
	ID xid.ID,
) context.Context {

	// Set the URL param
	ctxR := chi.NewRouteContext()
	ctxR.URLParams.Add("ID", ID.String())
	newCtx := context.WithValue(ctx, chi.RouteCtxKey, ctxR)

	return newCtx
}
