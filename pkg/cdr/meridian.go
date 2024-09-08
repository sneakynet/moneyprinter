package cdr

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/sneakynet/moneyprinter/pkg/types"
)

// Meridian implements the Parser interface and resolves the Meridian
// "New" format.
type Meridian struct{}

// Parse will read records from the provided reader and will return a
// slice of strongly typed CDRs.
func (m *Meridian) Parse(r io.Reader, clli string) ([]types.CDR, error) {
	scanner := bufio.NewScanner(bufio.NewReader(r))
	dateRegexp := regexp.MustCompile(`([\d]{2}\/[\d]{2} [\d]{2}:[\d]{2}:[\d]{2})`)
	durRegexp := regexp.MustCompile(`([\d]{2}:[\d]{2}:[\d]{2})`)
	numRegexp := regexp.MustCompile(`\w([\d]+)\w`)

	out := []types.CDR{}
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "L") && !strings.HasPrefix(line, "D") {
			continue
		}

		// MO/DD HH:MI:SS
		tsIdx := dateRegexp.FindStringIndex(line)
		timestamp, err := time.Parse("2006 01/02 15:04:05", fmt.Sprintf("%d %s", time.Now().Year(), dateRegexp.FindString(line)))
		if err != nil {
			slog.Warn("Error parsing timestamp in record", "record", line, "error", err)
			continue
		}

		// This is dumb and will break calls that are up more
		// than 24 hours, but it mostly works.
		durStr := durRegexp.FindString(line[tsIdx[1]:])
		durIdx := durRegexp.FindStringIndex(line[tsIdx[1]:])
		durTS, err := time.Parse("15:04:05", durStr)
		if err != nil {
			slog.Warn("Could not parse duration string", "record", line, "error", err)
			continue
		}

		h, m, s := durTS.Clock()
		dur, err := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))
		if err != nil {
			slog.Warn("Error parsing duration string", "string", durStr, "error", err)
			continue
		}

		dnis := strToUint(line[17:21])
		if line[:1] == "D" {
			dnisStr := numRegexp.FindString(line[durIdx[1]+tsIdx[1]:])
			dnis = strToUint(dnisStr)
		}

		t := types.CDR{
			CLLI:    clli,
			OrigID:  strToUint(line[4:6]),
			LogTime: timestamp,
			CLID:    strToUint(line[9:13]),
			DNIS:    dnis,
			Start:   timestamp.Add(dur * -1),
			End:     timestamp,
		}

		if t.CLID == t.DNIS {
			continue
		}

		out = append(out, t)
	}

	return out, nil
}
