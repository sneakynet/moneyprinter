package http

import (
	"log/slog"
	"net/http"

	"github.com/flosch/pongo2/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

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
			wcID, err := s.d.WirecenterCreate(&types.Wirecenter{Name: record["WIRECENTER"]})
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
			swID, err := s.d.SwitchCreate(&types.Switch{
				Name:         record["SWITCH"],
				CLLI:         record["CLLI"],
				Type:         record["SWITCHTYPE"],
				WirecenterID: wc.ID,
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

		eqpmnt, err := s.d.EquipmentGet(&types.Equipment{Name: record["EQUIPMENT"]})
		if err != nil {
			eqpmntID, err := s.d.EquipmentCreate(&types.Equipment{
				Name:         record["EQUIPMENT"],
				SwitchID:     sw.ID,
				WirecenterID: wc.ID,
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
				lineID, err := s.d.LineCreate(&types.Line{
					AccountID:    acctID,
					Type:         record["LINETYPE"],
					WirecenterID: wc.ID,
					SwitchID:     sw.ID,
					EquipmentID:  eqpmnt.ID,
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
