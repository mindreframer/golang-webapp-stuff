// package sockgroup implements a convienient way to multiplex messages to a group
// of websocket clients
package sockgroup

import (
	"encoding/json"
	"github.com/igm/sockjs-go/sockjs"
	"net/http"
)

type sockgroup struct {
	members    map[sockjs.Conn]bool
	messages   chan []byte
	register   chan sockjs.Conn
	unregister chan sockjs.Conn
}

func NewGroup() *sockgroup {
	return &sockgroup{
		members:    make(map[sockjs.Conn]bool),
		messages:   make(chan []byte),
		register:   make(chan sockjs.Conn),
		unregister: make(chan sockjs.Conn),
	}
}

func (sg *sockgroup) Start() {
	go func() {
		for {
			select {
			case member := <-sg.register:
				sg.members[member] = true
			case member := <-sg.unregister:
				member.Close()
				delete(sg.members, member)
			case message := <-sg.messages:
				for member, _ := range sg.members {
					go func(member sockjs.Conn) {
						defer func() {
							if r := recover(); r != nil {
								sg.unregister <- member
							}
						}()
						_, err := member.WriteMessage(message)
						if err != nil {
							sg.unregister <- member
						}
					}(member)
				}
			}
		}
	}()
}

func (sg *sockgroup) SendJSON(value interface{}) error {
	//if bytes, err := json.MarshalIndent(value, "", "  "); err != nil {
	if bytes, err := json.Marshal(value); err != nil {
		return err
	} else {
		sg.messages <- bytes
	}
	return nil
}

//func (sg *sockgroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
func (sg *sockgroup) Handler(baseUrl string) http.Handler {
	return sockjs.NewRouter(baseUrl, func(sock sockjs.Conn) {
		sg.register <- sock

		// incoming message pump
		go func() {
			var err error
			for err == nil {
				_, err = sock.ReadMessage()
			}
			sg.unregister <- sock
		}()
	}, sockjs.DefaultConfig)
}
