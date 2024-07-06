package http

import (
	"encoding/json"
	"net/http"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

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
