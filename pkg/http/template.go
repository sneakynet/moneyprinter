package http

import (
	"embed"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/leekchan/accounting"
)

//go:embed ui
var efs embed.FS

func (s *Server) templateErrorHandler(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, "Error while rendering template: %s\n", err)
}

func (s *Server) doTemplate(w http.ResponseWriter, r *http.Request, tmpl string, ctx pongo2.Context) {
	if ctx == nil {
		ctx = pongo2.Context{}
	}
	t, err := s.tpl.FromCache(tmpl)
	if err != nil {
		s.templateErrorHandler(w, err)
		return
	}
	if err := t.ExecuteWriter(ctx, w); err != nil {
		s.templateErrorHandler(w, err)
	}
}

func (s *Server) filterFormatMoney(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	cents, ok := in.Interface().(int)
	if !ok {
		slog.Warn("Got something that wasn't a number in formatMoney", "something", in)
		return pongo2.AsValue(""), nil
	}
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	return pongo2.AsValue(ac.FormatMoney(float64(cents) / 100)), nil
}

func (s *Server) filterBytesAsString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	bytes, ok := in.Interface().([]byte)
	if !ok {
		slog.Warn("Got something that wasn't a byte array in bytesAsString", "something", in)
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(string(bytes)), nil
}

func (s *Server) filterDecodeBase64(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	bytes, ok := in.Interface().([]byte)
	if !ok {
		slog.Warn("Got something that wasn't a byte array in decodeBase64", "something", in)
		return pongo2.AsValue(""), nil
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(bytes)))
	n, err := base64.StdEncoding.Decode(dst, bytes)
	if err != nil {
		slog.Warn("Decode error in fitlerDecodeBase64", "error", err)
		return pongo2.AsValue(""), nil
	}
	dst = dst[:n]
	return pongo2.AsSafeValue(string(dst)), nil
}
