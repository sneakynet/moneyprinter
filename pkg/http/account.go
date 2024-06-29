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

	lines, err := s.d.LineListByAccountID(account.ID)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"account": account,
		"lines": lines,
	}

	s.doTemplate(w, r, "p2/views/account.p2", ctx)
}

func (s *Server) uiHandleAccountCreateSingle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	aName := r.FormValue("account_name")
	aContact := r.FormValue("account_contact")
	aAlias := r.FormValue("account_alias")

	slog.Debug("Want to create account", "name", aName, "contact", aContact)

	id, err := s.d.AccountCreate(&types.Account{
		Name:    aName,
		Contact: aContact,
		Alias:   aAlias,
	})

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
	records := s.csvToMap(f)

	// Ideally we'd submit all these in a single batch to the
	// database for efficiency, but this is easier to filter for
	// dupes or CSV issues, and in reality the scale that
	// moneyprinter is expected to see is fairly small, so the
	// performance "gain" isn't worth the added complexity.
	for _, record := range records {
		if len(record["Name"]) == 0 {
			continue
		}

		acct, err := s.d.AccountGetByName(record["Name"])
		acctID := acct.ID
		if err != nil {
			slog.Warn("Error fetching account by name", "error", err)
			acctID, err = s.d.AccountCreate(&types.Account{
				Name:    record["Name"],
				Contact: record["Contact"],
				Alias:   record["Alias"],
			})
			if err != nil {
				s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
				return
			}
		}

		if _, err := s.d.DNGetByNumber(s.strToUint(record["DN"])); err != nil {
			slog.Warn("Error fetching DN by number", "error", err)
			slog.Debug("Want to create a line", "linetype", record["LINETYPE"])
			if record["LINETYPE"] == "FXS-LOOP-START" {
				lineID, err := s.d.LineCreate(&types.Line{
					AccountID:  acctID,
					Type:       record["LINETYPE"],
					Switch:     record["SWITCH"],
					Equipment:  record["EQUIPMENT"],
					Wirecenter: record["WIRECENTER"],
				})
				if err != nil {
					s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
					return
				}

				_, err = s.d.DNCreate(&types.DN{
					Number:    s.strToUint(record["DN"]),
					Display:   record["CNAM"],
					AccountID: acctID,
					LineID:    lineID,
				})
				if err != nil {
					s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
					return
				}
			}
		}
	}

	http.Redirect(w, r, "/ui/accounts", http.StatusSeeOther)
}
