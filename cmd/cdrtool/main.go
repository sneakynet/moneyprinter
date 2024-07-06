package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/sneakynet/moneyprinter/pkg/cdr"
)

var (
	cdrType = flag.String("type", "", "Type of CDR")
	cdrFile = flag.String("file", "", "Path to the CDR")
	cdrCLLI = flag.String("clli", "", "Source switch CLLI")
	upload  = flag.Bool("upload", false, "Upload CDRs to MoneyPrinter")
	mpAddr  = flag.String("addr", "localhost:8080", "MoneyPrinter address")
)

func main() {
	flag.Parse()

	f, err := os.Open(*cdrFile)
	if err != nil {
		slog.Error("Error loading CDR", "error", err)
		return
	}
	defer f.Close()

	var parser cdr.Parser

	switch *cdrType {
	case "cisco":
		parser = new(cdr.Cisco)
	default:
		slog.Error("Invalid CDR type; valid options are 'cisco'")
		return
	}

	records, err := parser.Parse(f, *cdrCLLI)
	if err != nil {
		slog.Error("Error parsing CDRs")
		return
	}

	for i, r := range records {
		slog.Info("record", "number", i, "from", r.CLID, "to", r.DNIS, "duration", r.End.Sub(r.Start))

		if *upload {
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(r)
			resp, err := http.Post("http://"+*mpAddr+"/api/cdr", "application/json", buf)
			if err != nil {
				slog.Error("Error uploading CDR", "error", err)
				continue
			}
			res := make(map[string]uint)
			json.NewDecoder(resp.Body).Decode(&res)
			slog.Info("Created new CDR", "ID", res["ID"])
		}
	}
}
