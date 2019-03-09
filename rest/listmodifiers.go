package rest

import (
	"net/url"

	"github.com/induzo/crud"
)

// ListModifiersFromURL will parse the url into CRUD ListModifiers
func ListModifiersFromURL(u *url.URL) crud.ListModifiers {
	rf := make(crud.ListModifiers, len(u.Query()))
	for k, v := range u.Query() {
		rf[k] = v
	}
	return rf
}
