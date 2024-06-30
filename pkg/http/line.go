package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewLineCreateForm(w http.ResponseWriter, r *http.Request) {
	accountID := s.strToUint(r.URL.Query().Get("account"))

	var accounts = []types.Account{{ID: accountID}}
	if accountID == 0 {
		var err error
		accounts, err = s.d.AccountList(nil)
		if err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}
	} else {
		acct, err := s.d.AccountGet(&types.Account{ID: accountID})
		if err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}
		accounts = []types.Account{acct}
	}

	ctx := pongo2.Context{
		"accounts": accounts,
	}

	s.doTemplate(w, r, "p2/views/line_create.p2", ctx)
}
