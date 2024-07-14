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

	s.r.Route("/ui/accounts", func(r chi.Router) {
		r.Get("/", s.uiViewAccountsList)
		r.Post("/", s.uiHandleAccountCreateSingle)
		r.Get("/new", s.uiViewAccountCreateForm)
		r.Get("/bulk", s.uiViewAccountBulkForm)
		r.Post("/bulk", s.uiHandleAccountCreateBulk)
		r.Get("/{id}", s.uiViewAccount)
		r.Get("/{id}/bill", s.uiViewAccountBill)
	})

	s.r.Get("/ui/cdrs", s.uiViewCDRs)

	s.r.Route("/ui/fees", func(r chi.Router) {
		r.Get("/", s.uiViewFeeList)
		r.Get("/new", s.uiViewFeeCreate)
		r.Post("/new", s.uiViewFeeUpsertPost)
		r.Get("/{id}/edit", s.uiViewFeeEditForm)
		r.Post("/{id}/edit", s.uiViewFeeUpsertPost)
		r.Post("/{id}/delete", s.uiViewFeeDelete)
	})

	s.r.Route("/ui/switches", func(r chi.Router) {
		r.Get("/", s.uiViewSwitchList)
		r.Get("/new", s.uiViewSwitchFormCreate)
		r.Post("/new", s.uiViewSwitchUpsert)
		r.Get("/{id}", s.uiViewSwitchDetail)
		r.Get("/{id}/edit", s.uiViewSwitchFormEdit)
		r.Post("/{id}/edit", s.uiViewSwitchUpsert)
		r.Post("/{id}/delete", s.uiViewSwitchDelete)

		r.Route("/{id}/equipment", func(er chi.Router) {
			r.Get("/", s.uiViewSwitchEquipment)
			r.Get("/{eid}", s.uiViewSwitchEquipmentDetail)
			r.Get("/new", s.uiViewSwitchEquipmentFormCreate)
			r.Post("/new", s.uiViewSwitchEquipmentUpsert)
			r.Get("/{eid}/edit", s.uiViewSwitchEquipmentFormEdit)
			r.Post("/{eid}/edit", s.uiViewSwitchEquipmentUpsert)
			r.Post("/{eid}/delete", s.uiViewSwitchEquipmentDelete)
			r.Get("/filter/{eName}", s.uiViewSwitchEquipment)
		})
	})

	s.r.Route("/ui/wirecenters", func(r chi.Router) {
		r.Get("/", s.uiViewWirecenterList)
		r.Get("/new", s.uiViewWirecenterFormCreate)
		r.Post("/new", s.uiViewWirecenterUpsert)
		r.Get("/{id}/edit", s.uiViewWirecenterFormEdit)
		r.Post("/{id}/edit", s.uiViewWirecenterUpsert)
		r.Get("/{id}/delete", s.uiViewWirecenterDelete)
	})

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
