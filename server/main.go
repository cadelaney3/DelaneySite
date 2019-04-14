package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/cadelaney3/delaneySite/pkg/websocket"
)

//var validPath = regexp.MustCompile("^/(ws|edit|save|view)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(ws|view)")


// define our WebSocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}
	go websocket.Writer(ws)
	websocket.Reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
	// map our `/ws` endpoint to the `serveWs` function
	http.HandleFunc("/ws", serveWs)
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

func main() {
	http.HandleFunc("/ws", makeHandler(serveWs))
	//setupRoutes()
	http.HandleFunc("/view/", makeHandler(handler))
	log.Println("Now server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
