package http

import (
	"net/http"

	"github.com/flosch/pongo2/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewLineCreateForm(w http.ResponseWriter, r *http.Request) {
	accountID := s.strToUint(r.URL.Query().Get("account"))
	switchID := s.strToUint(r.URL.Query().Get("switch"))
	equipmentID := s.strToUint(r.URL.Query().Get("equipment"))

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

	switches := []types.Switch{{ID: switchID}}
	if switchID == 0 {
		var err error
		switches, err = s.d.SwitchList(nil)
		if err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}
	} else {
		var err error
		switches, err = s.d.SwitchList(&types.Switch{ID: switchID})
		if err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}
	}

	equipment := []types.Equipment{{ID: equipmentID}}
	if equipmentID == 0 {
		var err error
		equipment, err = s.d.EquipmentList(nil)
		if err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}
	} else {
		var err error
		equipment, err = s.d.EquipmentList(&types.Equipment{ID: equipmentID})
		if err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}
	}
	eNames := map[string]uint{}
	for _, eqpmnt := range equipment {
		eNames[eqpmnt.Name] = eqpmnt.ID
	}

	ctx := pongo2.Context{
		"accounts":  accounts,
		"switches":  switches,
		"linetypes": map[string]string{"FXS-LOOP-START": "FXS Loop Start"},
		"eNames":    eNames,
		"equipment": equipment,
	}

	s.doTemplate(w, r, "p2/views/line_create.p2", ctx)
}
