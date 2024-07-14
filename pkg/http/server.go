package http

import (
	"context"
	"io/fs"
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New retruns a ready to serve instance of the HTTP server.
func New(opts ...Option) (*Server, error) {
	sub, _ := fs.Sub(efs, "ui")
	ldr := pongo2.NewFSLoader(sub)

	s := new(Server)
	s.r = chi.NewRouter()
	s.n = &http.Server{}
	s.tpl = pongo2.NewSet("html", ldr)

	for _, o := range opts {
		o(s)
	}

	pongo2.RegisterFilter("money", s.filterFormatMoney)

	s.r.Use(middleware.Heartbeat("/ping"))
	s.r.Use(middleware.Logger)

	s.r.Handle("/static/*", http.FileServer(http.FS(sub)))

	s.r.Get("/ui/accounts", s.uiViewAccountsList)
	s.r.Post("/ui/accounts", s.uiHandleAccountCreateSingle)
	s.r.Get("/ui/accounts/new", s.uiViewAccountCreateForm)
	s.r.Get("/ui/accounts/bulk", s.uiViewAccountBulkForm)
	s.r.Post("/ui/accounts/bulk", s.uiHandleAccountCreateBulk)
	s.r.Get("/ui/account/{id}", s.uiViewAccount)
	s.r.Get("/ui/account/{id}/bill", s.uiViewAccountBill)

	s.r.Get("/ui/cdrs", s.uiViewCDRs)

	s.r.Get("/ui/fees", s.uiViewFeeList)
	s.r.Get("/ui/fees/new", s.uiViewFeeCreate)
	s.r.Post("/ui/fees/new", s.uiViewFeeUpsertPost)
	s.r.Get("/ui/fees/{id}/edit", s.uiViewFeeEditForm)
	s.r.Post("/ui/fees/{id}/edit", s.uiViewFeeUpsertPost)
	s.r.Post("/ui/fees/{id}/delete", s.uiViewFeeDelete)

	s.r.Get("/ui/switches", s.uiViewSwitchList)
	s.r.Get("/ui/switches/new", s.uiViewSwitchFormCreate)
	s.r.Post("/ui/switches/new", s.uiViewSwitchUpsert)
	s.r.Get("/ui/switches/{id}", s.uiViewSwitchDetail)
	s.r.Get("/ui/switches/{id}/edit", s.uiViewSwitchFormEdit)
	s.r.Post("/ui/switches/{id}/edit", s.uiViewSwitchUpsert)
	s.r.Post("/ui/switches/{id}/delete", s.uiViewSwitchDelete)
	s.r.Get("/ui/switches/{id}/equipment", s.uiViewSwitchEquipment)
	s.r.Get("/ui/switches/{id}/equipment/{eid}", s.uiViewSwitchEquipmentDetail)
	s.r.Get("/ui/switches/{id}/equipment/new", s.uiViewSwitchEquipmentFormCreate)
	s.r.Post("/ui/switches/{id}/equipment/new", s.uiViewSwitchEquipmentUpsert)
	s.r.Get("/ui/switches/{id}/equipment/{eid}/edit", s.uiViewSwitchEquipmentFormEdit)
	s.r.Post("/ui/switches/{id}/equipment/{eid}/edit", s.uiViewSwitchEquipmentUpsert)
	s.r.Post("/ui/switches/{id}/equipment/{eid}/delete", s.uiViewSwitchEquipmentDelete)
	s.r.Get("/ui/switches/{id}/equipment/filter/{eName}", s.uiViewSwitchEquipment)

	s.r.Get("/ui/wirecenters", s.uiViewWirecenterList)
	s.r.Get("/ui/wirecenters/new", s.uiViewWirecenterFormCreate)
	s.r.Post("/ui/wirecenters/new", s.uiViewWirecenterUpsert)
	s.r.Get("/ui/wirecenters/{id}/edit", s.uiViewWirecenterFormEdit)
	s.r.Post("/ui/wirecenters/{id}/edit", s.uiViewWirecenterUpsert)
	s.r.Get("/ui/wirecenters/{id}/delete", s.uiViewWirecenterDelete)

	s.r.Post("/api/cdr", s.apiCreateCDR)

	return s, nil
}

// Serve binds and serves http on the bound socket.  An error will be
// returned if the server cannot initialize.
func (s *Server) Serve(bind string) error {
	s.n.Addr = bind
	s.n.Handler = s.r
	return s.n.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.n.Shutdown(ctx)
}
