package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewAccountCreateForm(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "p2/views/account_create.p2", nil)
}

func (s *Server) uiViewAccountBulkForm(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "p2/views/account_bulk.p2", nil)
}

func (s *Server) uiViewAccountsList(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.d.AccountList()
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "p2/views/account_list.p2", pongo2.Context{"accounts": accounts})
}

func (s *Server) uiViewAccount(w http.ResponseWriter, r *http.Request) {
	account, err := s.d.AccountGet(s.strToUint(chi.URLParam(r, "id")))
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "p2/views/account.p2", pongo2.Context{"account": account})
}

func (s *Server) uiHandleAccountCreateSingle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	aName := r.FormValue("account_name")
	aContact := r.FormValue("account_contact")

	slog.Debug("Want to create account", "name", aName, "contact", aContact)

	id, err := s.d.AccountCreate(&types.Account{Name: aName, Contact: aContact})

	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ui/account/%d", id), http.StatusSeeOther)
}

func (s *Server) uiHandleAccountCreateBulk(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	f, _, err := r.FormFile("accounts_file")
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	defer f.Close()
	accounts := s.csvToMap(f)

	// Ideally we'd submit all these in a single batch to the
	// database for efficiency, but this is easier to filter for
	// dupes or CSV issues, and in reality the scale that
	// moneyprinter is expected to see is fairly small, so the
	// performance "gain" isn't worth the added complexity.
	names := make(map[string]struct{})
	for _, account := range accounts {
		if _, seen := names[account["Name"]]; len(account["Name"]) == 0 || seen {
			continue
		}

		_, err := s.d.AccountCreate(&types.Account{Name: account["Name"], Contact: account["Contact"]})
		if err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}
		names[account["Name"]] = struct{}{}
	}

	http.Redirect(w, r, "/ui/accounts", http.StatusSeeOther)
}
