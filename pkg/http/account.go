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
	s.doTemplate(w, r, "views/account/create.p2", nil)
}

func (s *Server) uiViewAccountsList(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.d.AccountList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "views/account/list.p2", pongo2.Context{"accounts": accounts})
}

func (s *Server) uiViewAccount(w http.ResponseWriter, r *http.Request) {
	account, err := s.d.AccountGet(&types.Account{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	lines, err := s.d.LineList(&types.Line{AccountID: account.ID})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"account": account,
		"lines":   lines,
	}

	s.doTemplate(w, r, "views/account/detail.p2", ctx)
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

	http.Redirect(w, r, fmt.Sprintf("/ui/accounts/%d", id), http.StatusSeeOther)
}

func (s *Server) uiViewAccountBill(w http.ResponseWriter, r *http.Request) {
	if err := s.bp.Preload(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	account, err := s.d.AccountGet(&types.Account{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	bill, err := s.bp.BillAccount(account)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/account/bill.p2", pongo2.Context{"account": account, "bill": bill})
}

func (s *Server) uiViewAccountProvisionLineForm(w http.ResponseWriter, r *http.Request) {
	var lines []types.Line
	err := s.d.Raw().Where(map[string]interface{}{"account_id": 0}).Preload("Equipment").Find(&lines).Error
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{"lines": lines}

	s.doTemplate(w, r, "views/account/provision_line.p2", ctx)
}

func (s *Server) uiViewAccountProvisionLine(w http.ResponseWriter, r *http.Request) {
	acctID := s.strToUint(chi.URLParam(r, "id"))
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	line, err := s.d.LineGet(&types.Line{ID: s.strToUint(r.FormValue("line_id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	line.AccountID = acctID

	if _, err := s.d.LineSave(&line); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	dn := types.DN{
		Number:    s.strToUint(r.FormValue("dn_number")),
		Display:   r.FormValue("dn_display"),
		LineID:    line.ID,
		AccountID: acctID,
	}

	if _, err := s.d.DNSave(&dn); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ui/accounts/%d", acctID), http.StatusSeeOther)
}
