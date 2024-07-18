package http

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewSwitchList(w http.ResponseWriter, r *http.Request) {
	switches, err := s.d.SwitchList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/switch/list.p2", pongo2.Context{"switches": switches})
}

func (s *Server) uiViewSwitchDetail(w http.ResponseWriter, r *http.Request) {
	sw, err := s.d.SwitchGet(&types.Switch{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	s.doTemplate(w, r, "views/switch/detail.p2", pongo2.Context{"switch": sw})
}

func (s *Server) uiViewSwitchUpsert(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	sw := types.Switch{
		ID:   s.strToUint(chi.URLParam(r, "id")),
		Name: r.FormValue("switch_name"),
		CLLI: r.FormValue("switch_clli"),
		Type: r.FormValue("switch_type"),
	}

	_, err := s.d.SwitchSave(&sw)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, "/ui/switches", http.StatusSeeOther)
}

func (s *Server) uiViewSwitchFormCreate(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "views/switch/create.p2", nil)
}

func (s *Server) uiViewSwitchFormEdit(w http.ResponseWriter, r *http.Request) {
	sw, err := s.d.SwitchGet(&types.Switch{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "views/switch/create.p2", pongo2.Context{"switch": sw})
}

func (s *Server) uiViewSwitchDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.d.SwitchDelete(&types.Switch{ID: s.strToUint(chi.URLParam(r, "id"))}); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/ui/switches", http.StatusSeeOther)
}

func (s *Server) uiViewSwitchEquipment(w http.ResponseWriter, r *http.Request) {
	swID := s.strToUint(chi.URLParam(r, "id"))
	eName := chi.URLParam(r, "eName")

	sw, err := s.d.SwitchGet(&types.Switch{ID: swID})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	equipment, err := s.d.EquipmentList(&types.Equipment{SwitchID: sw.ID, Name: eName})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	slices.SortFunc(equipment, func(a, b types.Equipment) int {
		return cmp.Or(
			cmp.Compare(a.Name, b.Name),
			cmp.Compare(a.Port, b.Port),
		)
	})

	ctx := pongo2.Context{
		"switch":    sw,
		"equipment": equipment,
	}

	s.doTemplate(w, r, "views/equipment/list.p2", ctx)
}

func (s *Server) uiViewSwitchEquipmentDetail(w http.ResponseWriter, r *http.Request) {
	eqID := s.strToUint(chi.URLParam(r, "eid"))

	eq, err := s.d.EquipmentGet(&types.Equipment{ID: eqID})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "views/equipment/detail.p2", pongo2.Context{"equipment": eq})
}

func (s *Server) uiViewSwitchEquipmentFormCreate(w http.ResponseWriter, r *http.Request) {
	wirecenters, err := s.d.WirecenterList(nil)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "views/equipment/create.p2", pongo2.Context{"wirecenters": wirecenters})
}

func (s *Server) uiViewSwitchEquipmentFormEdit(w http.ResponseWriter, r *http.Request) {
	id := s.strToUint(chi.URLParam(r, "eid"))

	eq, err := s.d.EquipmentGet(&types.Equipment{ID: id})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	s.doTemplate(w, r, "views/equipment/create.p2", pongo2.Context{"equipment": eq})
}

func (s *Server) uiViewSwitchEquipmentUpsert(w http.ResponseWriter, r *http.Request) {
	swID := s.strToUint(chi.URLParam(r, "id"))
	eqID := s.strToUint(chi.URLParam(r, "eid"))

	if err := r.ParseForm(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	eq := types.Equipment{
		ID:           eqID,
		SwitchID:     swID,
		WirecenterID: s.strToUint(r.FormValue("equipment_wirecenter_id")),
		Name:         r.FormValue("equipment_name"),
		Port:         r.FormValue("equipment_port"),
		Description:  r.FormValue("equipment_desc"),
		Type:         r.FormValue("equipment_type"),
	}

	if _, err := s.d.EquipmentSave(&eq); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/ui/switches/%d/equipment/%d", swID, eqID), http.StatusSeeOther)
}

func (s *Server) uiViewSwitchEquipmentDelete(w http.ResponseWriter, r *http.Request) {
	if err := s.d.EquipmentDelete(&types.Equipment{ID: s.strToUint(chi.URLParam(r, "eid"))}); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}
	http.Redirect(w, r, "/ui/switches/"+chi.URLParam(r, "id")+"/equipment", http.StatusSeeOther)
}
