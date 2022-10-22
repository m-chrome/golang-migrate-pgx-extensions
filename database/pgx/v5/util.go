package v5

import (
	nurl "net/url"
	"strings"
)

// Copied from lib/pq implementation: https://github.com/lib/pq/blob/v1.9.0/conn.go#L1611
func quoteIdentifier(name string) string {
	end := strings.IndexRune(name, 0)
	if end > -1 {
		name = name[:end]
	}
	return `"` + strings.Replace(name, `"`, `""`, -1) + `"`
}

const newLine = '\n'

func runesCount(input []rune, target rune) int {
	var count int
	for _, r := range input {
		if r == target {
			count++
		}
	}
	return count
}

func runesLastIndex(input []rune, target rune) int {
	for i := len(input) - 1; i >= 0; i-- {
		if input[i] == target {
			return i
		}
	}
	return -1
}

func computeLineFromPos(s string, pos int) (line uint, col uint, ok bool) {
	// replace crlf with lf
	s = strings.Replace(s, "\r\n", "\n", -1)
	// pg docs: pos uses index 1 for the first character, and positions are measured in characters not bytes
	runes := []rune(s)
	if pos > len(runes) {
		return 0, 0, false
	}
	sel := runes[:pos]
	line = uint(runesCount(sel, newLine) + 1)
	col = uint(pos - 1 - runesLastIndex(sel, newLine))
	return line, col, true
}

// filterCustomQuery filter following query:
// - migrate `x-` options
// - pgx pool options - https://github.com/jackc/pgx/blob/v5.0.3/pgxpool/pool.go
func filterCustomQuery(u *nurl.URL) *nurl.URL {
	ux := *u
	vx := make(nurl.Values)
	for k, v := range ux.Query() {
		if len(k) <= 1 || k[0:2] != "x-" || strings.HasPrefix(k, "pool_") {
			vx[k] = v
		}
	}
	ux.RawQuery = vx.Encode()
	return &ux
}
