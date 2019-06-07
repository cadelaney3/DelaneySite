package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"encoding/json"

	"github.com/cadelaney3/delaneySite/pkg/websocket"
)

//var validPath = regexp.MustCompile("^/(ws|edit|save|view)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(ws|view|home)")

type response struct {
	Items []map[string]string `json:"Items"`
	Type int `json:"Type"`
}

// define our WebSocket endpoint
func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client {
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	res := response{
		Items: []map[string]string{{"body": "I am Chris Delaney"}},
		Type: 1,
	}
	resB, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(resB)
	log.Println("I am Chris Delaney")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
	http.HandleFunc("/view/", makeHandler(handler))
	http.HandleFunc("/home", makeHandler(homeHandler))
}

func main() {
	setupRoutes()
	log.Println("Now server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
