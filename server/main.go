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
	"time"
	"html"

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

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Credentials struct {
	Username string `json:"username", db:"username"`
	Password string `json:"password", db:"password"`
	Email string `json:"email", db:"email"`
}

type aboutResponse struct {
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

type fact struct {
	Fact string `json:"fact"`
}

type article struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Category string `json:"category"`
	Topic string `json:"topic"`
	Description string `json:"description"`
	Content string `json:"content"`
	Date string `json:"date"`
}

type articleResponse struct {
	Articles []article `json:"article"`
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

func about(w http.ResponseWriter, r *http.Request) {

	var factList []string

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

	rows, err := azureDB.Query("select fact from facts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var f string
		err = rows.Scan(&f)
		if err != nil {
			panic(err)
		}
		factList = append(factList, f)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	
	res := aboutResponse{
		Body: body,
		Facts: factList,
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

func articles(w http.ResponseWriter, r *http.Request) {
	var articleList []article

	u := html.EscapeString(r.URL.Path)
	log.Println(u)
	
	cat := html.EscapeString(r.URL.Query().Get("cat"))
	id := html.EscapeString(r.URL.Query().Get("id"))
	title := html.EscapeString(r.URL.Query().Get("title"))

	sqlStmt := "select id, title, author, category, topic, description, article_content, date_created from articles"
	
	if cat != "" {
		sqlStmt += " where category = " + "'" + cat + "'"
	}
	if cat == "drafts" {
		sqlStmt = "select id, title, author, category, topic, description, article_content, date_created from article_drafts"
	}
	if id != "" {
		sqlStmt += " where id=" + "'" + id + "'"
	}
	if title != "" {
		sqlStmt += " where title=" + "'" + title + "'"
	}
	rows, err := azureDB.Query(sqlStmt)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var tempDate time.Time
	for rows.Next() {
		var art article
		err = rows.Scan(&art.Id, &art.Title, &art.Author, &art.Category, &art.Topic, &art.Description, &art.Content, &tempDate)
		if err != nil {
			panic(err)
		}
		layoutUS  := "January 2, 2006"
		art.Date = tempDate.Format(layoutUS)
		articleList = append(articleList, art)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	res := articleResponse{
		Articles: articleList,
	}
	resB, err := json.Marshal(res)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	}

	log.Println(articleList)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Write(resB)
}

func addArticle(w http.ResponseWriter, r *http.Request) {
	incomingArticle := article{}

	err := json.NewDecoder(r.Body).Decode(&incomingArticle)
	if err != nil {
		log.Println(err)
	}

	log.Println(incomingArticle)

	if incomingArticle != (article{}) {
		log.Println(incomingArticle)
		statement := `insert into articles(title, author, category, topic, description, article_content)
				      values (@title, @author, @category, @topic, @description, @content)`
		_, err = azureDB.Exec(statement, sql.Named("title", incomingArticle.Title), sql.Named("author", incomingArticle.Author),
								sql.Named("category", incomingArticle.Category), sql.Named("topic", incomingArticle.Topic),
								sql.Named("description", incomingArticle.Description), sql.Named("content", incomingArticle.Content))
		if err != nil {
			respo := response{
				Status: 500,
				Message: "Internal server error",
			}
			resp, _ := json.Marshal(respo)
			w.Write(resp)
			panic(err)
			return
		}
		respo := response{
			Status: 200,
			Message: "Successful entry",
		}
		resp, _ := json.Marshal(respo)
		w.Write(resp)
		return
	}
	respo := response{
		Status: 401,
		Message: "Invalid entry",
	}
	resp, _ := json.Marshal(respo)
	w.Write(resp)
}

func addArticleDraft(w http.ResponseWriter, r *http.Request) {
	incomingArticle := article{}

	err := json.NewDecoder(r.Body).Decode(&incomingArticle)
	if err != nil {
		log.Println(err)
	}

	log.Println(incomingArticle)

	if incomingArticle != (article{}) {
		log.Println(incomingArticle)
		statement := `insert into article_drafts(title, author, category, topic, description, article_content)
				      values (@title, @author, @category, @topic, @description, @content)`
		_, err = azureDB.Exec(statement, sql.Named("title", incomingArticle.Title), sql.Named("author", incomingArticle.Author),
								sql.Named("category", incomingArticle.Category), sql.Named("topic", incomingArticle.Topic),
								sql.Named("description", incomingArticle.Description), sql.Named("content", incomingArticle.Content))
		if err != nil {
			respo := response{
				Status: 500,
				Message: "Internal server error",
			}
			resp, _ := json.Marshal(respo)
			w.Write(resp)
			panic(err)
			return
		}
		respo := response{
			Status: 200,
			Message: "Successful entry",
		}
		resp, _ := json.Marshal(respo)
		w.Write(resp)
		return
	}
	respo := response{
		Status: 401,
		Message: "Invalid entry",
	}
	resp, _ := json.Marshal(respo)
	w.Write(resp)
}

func addFact(w http.ResponseWriter, r *http.Request) {

	incomingFact := fact{}

	err := json.NewDecoder(r.Body).Decode(&incomingFact)
	if err != nil {
		log.Println(err)		
	}

	if incomingFact != (fact{}) {
		log.Println(incomingFact.Fact)
		statement := "insert into facts(fact) values (@fact)"
		_, err = azureDB.Exec(statement, sql.Named("fact", incomingFact.Fact))
		if err != nil {
			respo := response{
				Status: 500,
				Message: "Internal server error",
			}
			resp, _ := json.Marshal(respo)
			w.Write(resp)
			panic(err)
			return
		}
		respo := response{
			Status: 200,
			Message: "Successful entry",
		}
		resp, _ := json.Marshal(respo)
		w.Write(resp)
		return
	}
	respo := response{
		Status: 401,
		Message: "Invalid entry",
	}
	resp, _ := json.Marshal(respo)
	w.Write(resp)
}

func signIn(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie-name")

	credents := creds{}

	err := json.NewDecoder(r.Body).Decode(&credents)
	if err != nil {
		log.Println(err)		
	}
	if credents == (creds{}) {
		return
	}

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

func methodHandler(method string) Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
			// if r.Method != method {
			// 	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			// 	return
			// }
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST")
			fn(w, r)
		}
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
	http.HandleFunc("/about", about)
	http.HandleFunc("/signin", Chain(signIn, methodHandler("POST")))
	http.HandleFunc("/addFact", Chain(addFact, methodHandler("POST")))
	http.HandleFunc("/articles", Chain(articles, methodHandler("GET")))
	http.HandleFunc("/addArticle", Chain(addArticle, methodHandler("POST")))
	http.HandleFunc("/addArticleDraft", Chain(addArticleDraft, methodHandler("POST")))
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

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func main() {

	f, err := ioutil.ReadFile("./keys.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &keys)

	setupRoutes()
	azureDB = db.InitAzureDB(keys)

	log.Println("Now server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
