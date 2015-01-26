package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type tracesWatchers struct {
	sync.Mutex
	conns map[string]*websocket.Conn
	ch    chan trace
}

func (trs *tracesWatchers) add(conn *websocket.Conn) {
	trs.Lock()
	trs.conns[conn.RemoteAddr().String()] = conn
	trs.Unlock()
}

func (trs *tracesWatchers) sendTraces() {
	for trace := range trs.ch {
		trs.Lock()
		for key, conn := range trs.conns {
			if err := conn.WriteJSON(trace); err != nil {
				conn.Close()
				delete(trs.conns, key)
			}
		}
		trs.Unlock()
	}
}

func newTracesWatchers() *tracesWatchers {
	return &tracesWatchers{
		conns: make(map[string]*websocket.Conn),
		ch:    make(chan trace),
	}
}
