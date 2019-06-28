package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"encoding/json"
	"database/sql"
	"io/ioutil"
	"golang.org/x/crypto/bcrypt"
	"os"

	_ "github.com/lib/pq"
	"github.com/cadelaney3/delaneySite/pkg/websocket"
	"github.com/cadelaney3/delaneySite/pkg/db"
	"github.com/gorilla/sessions"
)

var keys = make(map[string]map[string]string)
var azureDB *sql.DB
var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)

)

//var validPath = regexp.MustCompile("^/(ws|edit|save|view)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(ws|view|home|signin)")

type Credentials struct {
	Username string `json:"username", db:"username"`
	Password string `json:"password", db:"password"`
	Email string `json:"email", db:"email"`
}

type homeResponse struct {
	Body string `json:"body"`
	Facts []string `json:"facts"`
}

type creds struct {
	Username string `json:"username"`
	Password string `password:"password"`
}

type response struct {
	Status int `json:"status"`
	Message string `json:"message"`
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
	body := `I am Chris Delaney. I was born and raised in Frenchtown, Montana. I had an active childhood,
			 playing soccer, baseball, basketball, and football. I practiced karated for about seven years
			 before high school. I did some hunting, fishing, and boating. But when I went into 8th grade, 
			 I became obsessed with golf. Over the next few years, I spent countless hours practicing and
			 playing. It didn't much matter if it was hot as the dickens, pouring rain, or even snowing.
			 There was golf to be played. I dreamed of becoming a professional golfer, and was willing to
			 put in the work. It always felt like I wasn't getting out as much as I was putting in though.
			 Nonetheless, I still have a few nice accomplishments under my belt, like taking 2nd in state
			 in high school and getting my handicap down to 0. I did not play golf in college, though I 
			 wanted to and did try. Now, I no longer have the burning desire to practice for hours every
			 day to become pro, but I still love the game and look forward to competing in the future.

			 Nowadays I also enjoy working on coding projects, this website being one. It's hard to say I
			 have a favorite area of computer science since I find so many areas fascinating. Machine
			 learning, parallel computing, and math-heavy computing are particularly interesting to me, but 
			 I am happy to learn anything new. My favorite programming language at the moment is Go. I like 
			 learning lots of different languages though since they all seem to have unique properties
			 that make them especially handy for making certain types of programs. 

			 I recently graduated from Gonzaga University with a B.S. in computer science and am excited to
			 start life as a software engineer. I hope to work on all sorts of projects so that I can have
			 the chance to experience many different roles and learn a wide variety of skills. I think that
			 would keep things interesting and allow me to discover what I truly enjoy the most.
			 
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
}

func signIn(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	credents := creds{}

	err := json.NewDecoder(r.Body).Decode(&credents)
	if err != nil {
		log.Println(err)		
	}
	if credents == (creds{}) {
		return
	}

	//result := db.QueryRow("select password from account where username=$1", credents.Username)
	result := azureDB.QueryRow("select password from users where username=@username", sql.Named("username", credents.Username))
	if err != nil {
		// If there is an issue with the database, return a 500 error
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		respo := response{
			Status: 500,
			Message: "Internal server error",
		}
		resp, _ := json.Marshal(respo)
		w.Write(resp)
		return
	}

	storedCreds := &creds{}

	// Store the obtained password in `storedCreds`
	err = result.Scan(&storedCreds.Password)

	if err != nil {
		log.Println(err)
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			respo := response{
				Status: 401,
				Message: "Invalid username",
			}
			resp, _ := json.Marshal(respo)
			w.Write(resp)
			return
		}
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		respo := response{
			Status: 500,
			Message: "Internal server error",
		}
		resp, _ := json.Marshal(respo)
		w.Write(resp)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(credents.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		respo := response{
			Status: 401,
			Message: "Incorrect password",
		}
		resp, _ := json.Marshal(respo)
		w.Write(resp)
		return
	}

	respo := response{
		Status: 200,
		Message: "All good",
	}
	resp, err := json.Marshal(respo)
	if err != nil {
		log.Println(err)
	}
	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)

	w.Write(resp)
}

func signInHandler(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		fn(w, r)
	}
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
	http.HandleFunc("/signin", signInHandler(signIn))
}

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func main() {

	f, err := ioutil.ReadFile("../keys.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &keys)

	setupRoutes()
	// connDB(keys)
	azureDB = db.InitAzureDB(keys)

	log.Println("Now server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
