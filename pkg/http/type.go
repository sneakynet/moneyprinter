package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

// Server handles the HTTP frontend of the application.
type Server struct {
	r chi.Router
	n *http.Server
	d DB

	tpl *pongo2.TemplateSet
}

// Option configures the server in a composeable way.
type Option func(*Server)

// DB encapsulates all the logic that the webserver expects to be able
// to do.
type DB interface {
	AccountCreate(*types.Account) (uint, error)
	AccountList() ([]types.Account, error)
	AccountGet(uint) (types.Account, error)
	AccountGetByName(string) (types.Account, error)

	LineCreate(*types.Line) (uint, error)
	LineListByAccountID(uint) ([]types.Line, error)
	LineGet(uint) (types.Line, error)

	DNCreate(*types.DN) (uint, error)
	DNListByAccountID(uint) ([]types.DN, error)
	DNGet(uint) (types.DN, error)
}
