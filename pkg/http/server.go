package http

import (
	"context"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New retruns a ready to serve instance of the HTTP server.
func New(opts ...Option) (*Server, error) {
	var tplRoot fs.FS

	if tpath := os.Getenv("MONEYD_TEMPLATE_PATH"); tpath != "" {
		slog.Warn("Loading templates from debug path", "path", tpath)
		tplRoot = os.DirFS(tpath)
	} else {
		tplRoot, _ = fs.Sub(efs, "ui")
	}
	p2Root, _ := fs.Sub(tplRoot, "p2")
	ldr := pongo2.NewFSLoader(p2Root)

	s := new(Server)
	s.r = chi.NewRouter()
	s.n = &http.Server{}
	s.tpl = pongo2.NewSet("html", ldr)

	for _, o := range opts {
		o(s)
	}

	pongo2.RegisterFilter("money", s.filterFormatMoney)
	pongo2.RegisterFilter("bytesAsString", s.filterBytesAsString)
	pongo2.RegisterFilter("decodeBase64", s.filterDecodeBase64)

	s.r.Use(middleware.Heartbeat("/ping"))
	s.r.Use(middleware.Logger)

	s.r.Handle("/static/*", http.FileServer(http.FS(tplRoot)))

	s.r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusSeeOther)
	})
	s.r.Get("/ui", s.uiViewLanding)

	s.r.Route("/ui/bulk", func(r chi.Router) {
		r.Get("/", s.uiViewBulkLanding)
		r.Get("/omni", s.uiViewBulkOmniForm)
		r.Post("/omni", s.uiViewBulkOmniCreate)
		r.Get("/line-card", s.uiViewBulkLinecardForm)
		r.Post("/line-card", s.uiViewBulkLinecardCreate)
		r.Get("/accounts", s.uiViewBulkAccountsForm)
		r.Post("/accounts", s.uiViewBulkAccountsCreate)
	})

	s.r.Route("/ui/accounts", func(r chi.Router) {
		r.Get("/", s.uiViewAccountsList)
		r.Post("/", s.uiHandleAccountCreateSingle)
		r.Get("/new", s.uiViewAccountCreateForm)
		r.Get("/{id}", s.uiViewAccount)
		r.Get("/{id}/bill", s.uiViewAccountBill)
		r.Get("/{id}/provision-line", s.uiViewAccountProvisionLineForm)
		r.Post("/{id}/provision-line", s.uiViewAccountProvisionLine)
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
			er.Get("/", s.uiViewSwitchEquipment)
			er.Get("/{eid}", s.uiViewSwitchEquipmentDetail)
			er.Get("/new", s.uiViewSwitchEquipmentFormCreate)
			er.Post("/new", s.uiViewSwitchEquipmentUpsert)
			er.Get("/{eid}/edit", s.uiViewSwitchEquipmentFormEdit)
			er.Post("/{eid}/edit", s.uiViewSwitchEquipmentUpsert)
			er.Post("/{eid}/delete", s.uiViewSwitchEquipmentDelete)
			er.Get("/filter/{eName}", s.uiViewSwitchEquipment)
		})

		r.Route("/{id}/lines", func(lr chi.Router) {
			lr.Get("/", s.uiViewSwitchLineList)
			lr.Get("/{lid}", s.uiViewSwitchLineDetail)
			lr.Get("/new", s.uiViewSwitchLineFormCreate)
			lr.Post("/new", s.uiViewSwitchLineUpsert)
			lr.Get("/{lid}/edit", s.uiViewSwitchLineFormEdit)
			lr.Post("/{lid}/edit", s.uiViewSwitchLineUpsert)
			lr.Post("/{lid}/delete", s.uiViewSwitchLineDelete)
		})
	})

	s.r.Route("/ui/lines", func(r chi.Router) {
		r.Get("/", s.uiViewSwitchLineListAll)
		r.Get("/{lid}", s.uiViewSwitchLineDetail)
	})

	s.r.Route("/ui/dn", func(r chi.Router) {
		r.Get("/", s.uiViewDNList)
		r.Get("/{id}", s.uiViewDNDetail)
		r.Get("/new", s.uiViewDNFormCreate)
		r.Post("/new", s.uiViewDNUpsert)
		r.Get("/{id}/edit", s.uiViewDNFormEdit)
		r.Post("/{id}/edit", s.uiViewDNUpsert)
		r.Post("/{id}/delete", s.uiViewDNDelete)
	})

	s.r.Route("/ui/wirecenters", func(r chi.Router) {
		r.Get("/", s.uiViewWirecenterList)
		r.Get("/new", s.uiViewWirecenterFormCreate)
		r.Post("/new", s.uiViewWirecenterUpsert)
		r.Get("/{id}", s.uiViewWirecenterDetail)
		r.Get("/{id}/edit", s.uiViewWirecenterFormEdit)
		r.Post("/{id}/edit", s.uiViewWirecenterUpsert)
		r.Get("/{id}/delete", s.uiViewWirecenterDelete)
	})

	s.r.Route("/ui/lecs", func(r chi.Router) {
		r.Get("/", s.uiViewLECList)
		r.Get("/new", s.uiViewLECFormCreate)
		r.Post("/new", s.uiViewLECUpsert)
		r.Get("/{id}", s.uiViewLECDetail)
		r.Get("/{id}/edit", s.uiViewLECFormEdit)
		r.Post("/{id}/edit", s.uiViewLECUpsert)
		r.Get("/{id}/set-logo", s.uiViewLogoForm)
		r.Post("/{id}/set-logo", s.uiViewLogoSet)
		r.Post("/{id}/delete", s.uiViewLECDelete)
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

func (s *Server) uiViewLanding(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "views/landing.p2", nil)
}
