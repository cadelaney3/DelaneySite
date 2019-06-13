package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"encoding/json"
	"database/sql"
	"io/ioutil"

	_ "github.com/lib/pq"
	"github.com/cadelaney3/delaneySite/pkg/websocket"
)

var keys = make(map[string]map[string]string)

type postgresConn struct {
	host string
	port int 
	user string
	password string 
	dbname string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "cdswaggy"
	password = "Theucanes3"
	dbname   = "delaneysite"
)

//var validPath = regexp.MustCompile("^/(ws|edit|save|view)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(ws|view|home)")

type homeResponse struct {
	Body string `json:"body"`
	Facts []string `json:"facts"`
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
	body := `I am Chris Delaney. I was born and raised in Frenchtown, Montana. I do not a ride my horse to school or say
			 hello by tipping my cowboy hat, but I do have Montanan characteristics. I enjoy the outdoors
			 and outdoor activities, like hiking, fishing, hunting. My main outdoor activity is golfing,
			 however. I have been golfing for the last 8 years or so. Lately, golf serves more as a
			 hobby, but I still enjoy competing and practicing to get better.

			 My other main interest anymore is programming and computer science, which is kind of why I am
			 here. I didn't know anything about computer science going in to freshman year in college, but 
			 chose it as my major because I didn't know what else I would study and because I heard it was
			 a good field to get into. I would describe my experience with computer science as I would my
			 experience with coffee. At first, it was bitter, difficult to swallow, and not all that enjoyable.
			 But life comes at you fast and classes get harder, so I started taking in more as a means of 
			 performing better. Over time, my palette refined and I started to actually enjoy it. Nowadays,
			 I wake up excited for it, and a day doesn't seem quite right without a healthy dosage.

			 This is my basic description. But if you are itching for more, I have much more interesting things about me
			 listed below.`
	
	facts := []string{`I made up the name "CDSwaggy" back when "swag" was starting to become overused and kinda dumb.
					   I used it joking for usernames, and it got good reactions from people I knew.
					   Because of that and because it was easy to remember, I started using it for basically everything.
					   I am not particularly swaggy at all, which kind of adds to the joke. So that's that story.`,
					  `I've been to a rodeo where I caught a chicken and got to keep it 
					  (I named it Curly, since I packed it home in a curly fry box)`,
					  "I have helped my grandpa give shots to cows.",
					  `My lowest score in golf was a 65 (-5) at King Ranch Golf Course, the course I grew up on.
					   The score I most proud of was a 69 (-3) at Stock Farm Golf Club, where I used to work. I
					   should have shot a 75 or so, but I putted out of mind that round.`,
					   `I did karate for about 7 years and was a belt and a stripe from black belt. I recently watched
					   a few kung-fu movies with Bruce Lee and Jackie Chan, which has really got the kung-fu juices flowing.`,
					   `I love love love spicy food. It has to be tasty spicy though. So Tapatio is my go-to choice, but if
					   I'm sweatin' form some extra spicy Tai food, those are tears of joy that I am crying.`,
					   `More facts and tidbits to come as I think of them.`}
	
	res := homeResponse{
		Body: body,
		Facts: facts,
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
	f, err := ioutil.ReadFile("../keys.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &keys)

	postgres := postgresConn{
		host: "localhost",
		port: 5432,
		user: keys["POSTGRES"]["USER"],
		password: keys["POSTGRES"]["PASSWORD"],
		dbname: "delaneysite",
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
	postgres.host, postgres.port, postgres.user, postgres.password, postgres.dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
	  panic(err)
	} 
	fmt.Println("Successfully connected!")

	setupRoutes()
	log.Println("Now server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
