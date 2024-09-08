package http

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/go-chi/chi/v5"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/leekchan/accounting"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

func (s *Server) uiViewAccountBillText(w http.ResponseWriter, r *http.Request) {
	lec, err := s.d.LECGet(&types.LEC{ID: s.strToUint(r.URL.Query().Get("lec"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
	}

	if err := s.bp.Preload(); err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	account, err := s.d.AccountGet(&types.Account{ID: s.strToUint(chi.URLParam(r, "id"))})
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	bill, err := s.bp.BillAccount(account)
	if err != nil {
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err.Error()})
		return
	}

	width := 80

	t := table.NewWriter()
	t.SetColumnConfigs([]table.ColumnConfig{{
		Name:        "Telephone Bill",
		WidthMin:    width - 3,
		Align:       text.AlignCenter,
		AlignHeader: text.AlignCenter,
	}})
	t.Style().Options.DrawBorder = false
	t.SetOutputMirror(w)
	t.SetAllowedRowLength(width)
	t.AppendHeader(table.Row{"Telephone Bill"})
	t.AppendRow(table.Row{lec.BillMsg})
	t.AppendRow(table.Row{lec.Name + " - " + lec.Byline})
	t.AppendRow(table.Row{lec.Website})
	t.Render()
	fmt.Fprintln(w, "")

	t = table.NewWriter()
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "Name",
			WidthMin: (width / 3) - 3,
		},
		{
			Name:     "DBA",
			WidthMin: (width / 3) - 3,
		},
		{
			Name:     "Contact",
			WidthMin: (width / 3) - 3,
		},
	})
	t.Style().Options.DrawBorder = false
	t.SetOutputMirror(w)
	t.SetAllowedRowLength(width)
	t.AppendHeader(table.Row{"Name", "DBA", "Contact"})
	t.AppendRow(table.Row{account.Name, account.Alias, account.Contact})
	t.Render()
	fmt.Fprintln(w, "")

	t = table.NewWriter()
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:        "Fee",
			Align:       text.AlignLeft,
			AlignHeader: text.AlignLeft,
			WidthMin:    width / 4,
		},
		{
			Name:        "Item",
			Align:       text.AlignCenter,
			AlignHeader: text.AlignCenter,
			WidthMin:    width / 2,
		},
		{
			Name:        "Cost",
			Align:       text.AlignRight,
			AlignHeader: text.AlignRight,
			WidthMin:    10,
			WidthMax:    10,
		},
	})
	t.Style().Options.DrawBorder = false
	t.SetOutputMirror(w)
	t.SetAllowedRowLength(width)
	t.AppendHeader(table.Row{"Fee", "Item", "Cost"})
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	for _, item := range bill.Lines {
		t.AppendRow(table.Row{item.Fee, item.Item, ac.FormatMoney(float64(item.Cost) / 100)})
	}
	t.SortBy([]table.SortBy{{Name: "Fee", Mode: table.Asc}})
	t.Render()
	fmt.Fprintln(w, "")

	t = table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.SetOutputMirror(w)
	t.SetAllowedRowLength(width)
	t.AppendRow(table.Row{"Grand Total: " + ac.FormatMoney(float64(bill.Cost())/100)})
	t.Render()
}
