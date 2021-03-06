package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

var cmdDynos = &Command{
	Run:   runDynos,
	Usage: "dynos [-l] [name...]",
	Short: "list dynos",
	Long: `
Lists dynos.

Options:

    -l       long listing

Long listing for shows the name, state, age, and command.

Examples:

    $ hk dynos
    run.3794
    web.1
    web.2

    $ hk dynos web
    web.1
    web.2

    $ hk dynos -l
    run.3794  up   1m  bash
    web.1     up  15h  "blog /app /tmp/dst"
    web.2     up   8h  "blog /app /tmp/dst"
`,
}

func init() {
	cmdDynos.Flag.BoolVar(&flagLong, "l", false, "long listing")
}

func runDynos(cmd *Command, names []string) {
	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	defer w.Flush()
	listDynos(w, names)
}

func listDynos(w io.Writer, names []string) {
	var dynos []*Dyno
	must(Get(&dynos, "/apps/"+mustApp()+"/dynos"))
	sort.Sort(DynosByName(dynos))

	if len(names) == 0 {
		for _, d := range dynos {
			listDyno(w, d)
		}
		return
	}

	for _, name := range names {
		for _, d := range dynos {
			if !strings.Contains(name, ".") {
				if strings.HasPrefix(d.Name, name+".") {
					listDyno(w, d)
				}
			} else {
				if d.Name == name {
					listDyno(w, d)
				}
			}
		}
	}
}

func listDyno(w io.Writer, d *Dyno) {
	if flagLong {
		listRec(w,
			d.Name,
			d.State,
			prettyDuration{d.Age()},
			maybeQuote(d.Command),
		)
	} else {
		fmt.Fprintln(w, d.Name)
	}
}

// quotes s as a json string if it contains any weird chars
// currently weird is anything other than [alnum]_-
func maybeQuote(s string) string {
	for _, r := range s {
		if !('0' <= r && r <= '9' || 'a' <= r && r <= 'z' ||
			'A' <= r && r <= 'Z' || r == '-' || r == '_') {
			return quote(s)
		}
	}
	return s
}

// quotes s as a json string
func quote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

type DynosByName []*Dyno

func (p DynosByName) Len() int      { return len(p) }
func (p DynosByName) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p DynosByName) Less(i, j int) bool {
	return p[i].Type < p[j].Type || (p[i].Type == p[j].Type && p[i].Seq() < p[j].Seq())
}

type prettyDuration struct {
	time.Duration
}

func (a prettyDuration) String() string {
	switch d := a.Duration; {
	case d > 2*24*time.Hour:
		return a.Unit(24*time.Hour, "d")
	case d > 2*time.Hour:
		return a.Unit(time.Hour, "h")
	case d > 2*time.Minute:
		return a.Unit(time.Minute, "m")
	}
	return a.Unit(time.Second, "s")
}

func (a prettyDuration) Unit(u time.Duration, s string) string {
	return fmt.Sprintf("%2d", roundDur(a.Duration, u)) + s
}

func roundDur(d, k time.Duration) int {
	return int((d + k/2 - 1) / k)
}

func abbrev(s string, n int) string {
	if len(s) > n {
		return s[:n-1] + "…"
	}
	return s
}
