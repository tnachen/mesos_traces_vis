package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	port      = flag.Int("p", 3000, "Port to listen on")
	redisAddr = flag.String("r", "127.0.0.1:6379", "Redis where to connect")
	upgrader  = websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
)

type trace struct {
	Name           string `redis:"name" json:"name,omitempty"`
	Span_id        string `redis:"span_id" json:"span_id,omitempty"`
	Span_parent_id string `redis:"span_parent_id" json:"span_parent_id,omitempty"`
	From_upid      string `redis:"from_upid" json:"from_upid,omitempty"`
	To_upid        string `redis:"to_upid" json:"to_upid,omitempty"`
	Time_nanos     string `redis:"time_nanos" json:"time_nanos,omitempty"`
	Stage          string `redis:"stage" json:"stage,omitempty"`
	Id             string `redis:"id" json:"id,omitempty"`
}

func main() {
	flag.Parse()
	pool := newPool(*redisAddr)
	tws := newTracesWatchers()
	go watchTraces(pool, tws.ch)
	go tws.sendTraces()
	r := mux.NewRouter()
	r.Path("/traces/ws").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			return
		}
		tws.add(conn)
	})
	r.Path("/traces").Methods("GET").HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		traces, err := getTraces(pool.Get())
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(traces)
	})
	r.Path("/trace/{id}").Methods("GET").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		trcs, err := getTrace(pool.Get(), mux.Vars(req)["id"])
		if err != nil {
			http.Error(res, err.Error(), 500)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(trcs)

	})
	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))
	log.Printf("Example app listening at http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
