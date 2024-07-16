package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/flosch/pongo2/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewBulkLanding(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "p2/views/bulk_landing.p2", nil)
}

func (s *Server) uiViewBulkOmniForm(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "p2/views/bulk_omni.p2", nil)
}

func (s *Server) uiViewBulkOmniCreate(w http.ResponseWriter, r *http.Request) {
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

		acct, err := s.d.AccountGet(&types.Account{Name: record["Name"]})
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

		wc, err := s.d.WirecenterGet(&types.Wirecenter{Name: record["WIRECENTER"]})
		if err != nil {
			wcID, err := s.d.WirecenterSave(&types.Wirecenter{Name: record["WIRECENTER"]})
			if err != nil {
				s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
				return
			}

			wc, err = s.d.WirecenterGet(&types.Wirecenter{ID: wcID})
			if err != nil {
				s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
				return
			}
		}

		sw, err := s.d.SwitchGet(&types.Switch{Name: record["SWITCH"]})
		if err != nil {
			swID, err := s.d.SwitchSave(&types.Switch{
				Name: record["SWITCH"],
				CLLI: record["CLLI"],
				Type: record["SWITCHTYPE"],
			})
			if err != nil {
				s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
				return
			}

			sw, err = s.d.SwitchGet(&types.Switch{ID: swID})
			if err != nil {
				s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
				return
			}
		}

		eqpmnt, err := s.d.EquipmentGet(&types.Equipment{Name: record["EQUIPMENT"], Port: record["PORT"]})
		if err != nil {
			eqpmntID, err := s.d.EquipmentSave(&types.Equipment{
				Name:         record["EQUIPMENT"],
				Port:         record["PORT"],
				SwitchID:     sw.ID,
				WirecenterID: wc.ID,
				Type:         record["LINETYPE"],
			})
			if err != nil {
				s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
				return
			}

			eqpmnt, err = s.d.EquipmentGet(&types.Equipment{ID: eqpmntID})
			if err != nil {
				s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
				return
			}
		}

		if _, err := s.d.DNGet(&types.DN{Number: s.strToUint(record["DN"])}); err != nil {
			slog.Warn("Error fetching DN by number", "error", err)
			slog.Debug("Want to create a line", "linetype", record["LINETYPE"])
			if record["LINETYPE"] == "FXS-LOOP-START" {
				lineID, err := s.d.LineSave(&types.Line{
					AccountID:   acctID,
					SwitchID:    sw.ID,
					EquipmentID: eqpmnt.ID,
				})
				if err != nil {
					s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
					return
				}

				_, err = s.d.DNSave(&types.DN{
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

func (s *Server) uiViewBulkLinecardForm(w http.ResponseWriter, r *http.Request) {
	switches, err := s.d.SwitchList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	wirecenters, err := s.d.WirecenterList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	ctx := pongo2.Context{
		"switches":    switches,
		"wirecenters": wirecenters,
	}

	s.doTemplate(w, r, "p2/views/bulk_linecard.p2", ctx)
}

func (s *Server) uiViewBulkLinecardCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	tpl, err := pongo2.FromString(r.FormValue("port_tmpl"))
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	swID := s.strToUint(r.FormValue("switch_id"))
	wcID := s.strToUint(r.FormValue("wirecenter_id"))
	eqName := r.FormValue("card_name")
	eqType := r.FormValue("equipment_type")

	for id := range s.strToUint(r.FormValue("port_count")) {
		port, err := tpl.Execute(pongo2.Context{"id": id})
		if err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}

		eq := types.Equipment{
			SwitchID:     swID,
			WirecenterID: wcID,
			Name:         eqName,
			Type:         eqType,
			Port:         port,
		}

		if _, err := s.d.EquipmentSave(&eq); err != nil {
			s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
			return
		}

		if eq.Type == "FXS-LOOP-START" && r.FormValue("allocate_lines") != "" {
			l := types.Line{SwitchID: swID, EquipmentID: eq.ID}
			if _, err := s.d.LineSave(&l); err != nil {
				s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
				return
			}
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/ui/switches/%d/equipment/filter/%s", swID, eqName), http.StatusSeeOther)
}
