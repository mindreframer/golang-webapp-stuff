package stats

import (
	"sync"
	"net/http"
	"encoding/json"
)

type ToNotify interface {
	SendJSON(value interface{}) error
}

var Destination ToNotify

var trigger chan struct{}

var stats struct {
	data map[string]int
	sync.Mutex
}

func ChangeStat(statName string, change int) {
	stats.Lock()
	defer stats.Unlock()
	stats.data[statName] = stats.data[statName] + change
	go func() {
		if Destination != nil {
			Destination.SendJSON(stats.data)
		}
	}()
}

func init() {
	trigger = make(chan struct{})
	stats.data = make(map[string]int)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.Encode(stats.data)
}
