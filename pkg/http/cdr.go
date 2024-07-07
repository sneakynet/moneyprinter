package http

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/flosch/pongo2/v5"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewCDRs(w http.ResponseWriter, r *http.Request) {
	cdrs := []types.CDR{}

	if dn := r.URL.Query().Get("dn"); dn != "" {
		s1, _ := s.d.CDRList(&types.CDR{CLID: s.strToUint(dn)})
		cdrs = append(cdrs, s1...)

		s2, _ := s.d.CDRList(&types.CDR{DNIS: s.strToUint(dn)})
		cdrs = append(cdrs, s2...)
	}

	sort.Slice(cdrs, func(i, j int) bool {
		return cdrs[i].Start.Before(cdrs[j].Start)
	})

	s.doTemplate(w, r, "p2/views/cdr_list.p2", pongo2.Context{"cdrs": cdrs})
}

func (s *Server) apiCreateCDR(w http.ResponseWriter, r *http.Request) {
	cdr := new(types.CDR)

	if err := json.NewDecoder(r.Body).Decode(cdr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	id, err := s.d.CDRCreate(cdr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]uint{"ID": id})
}
