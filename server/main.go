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
	"strings"
	"context"
	// "net/url"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/cadelaney3/delaneySite/pkg/websocket"
	"github.com/cadelaney3/delaneySite/pkg/db"
	"github.com/cadelaney3/delaneySite/utils"
	mw "github.com/cadelaney3/delaneySite/middleware"
	"github.com/gorilla/sessions"
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var keys = make(map[string]map[string]string)
var azureDB *sql.DB
var client *mongo.Client

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key = []byte(os.Getenv("SESSION_KEY"))
	store = sessions.NewCookieStore(key)

)

//var validPath = regexp.MustCompile("^/(ws|edit|save|view)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(ws|view|home|signin)")


type aboutResponse struct {
	Body string `json:"body"`
	Facts []string `json:"facts"`
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

type Credentials struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token string `json:"token"`
}

type Token struct {
	UserId string
	jwt.StandardClaims
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	credentials := Credentials{}
	session, _ := store.Get(r, "cookie-name")

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		log.Println(err)		
	}
	if credentials == (Credentials{}) {
		return
	}

	storedCreds := Credentials{}
	collection := client.Database("delaney-db").Collection("users")
	filter := bson.M{"username": credentials.Username}
	documentReturned := collection.FindOne(context.TODO(), filter)
	err = documentReturned.Decode(&storedCreds)
	if err != nil {
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		utils.Response(w, message)
		log.Println(err)
		return		
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(credentials.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		message := utils.Message(http.StatusUnauthorized, "Invalid username or password")
		utils.Response(w, message)
		return
	}

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)

	//Create JWT token
	tk := &Token{UserId: storedCreds.ID.String()}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	storedCreds.Token = tokenString //Store the token in the response
	storedCreds.Password = ""

	message := utils.Message(http.StatusOK, "Success")
	message["account"] = storedCreds

	utils.Response(w, message)
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
	// log.Println("title: ", strings.Split(title, "-"))
	title = strings.Join(strings.Split(title, "-"), " ")

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
		sqlStmt += " where LOWER(title)=" + "'" + title + "'"
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

	f, err := ioutil.ReadFile("./keys.json")
	if err != nil {
		panic(err)
	}
	
	err = json.Unmarshal(f, &keys)
	if err != nil {
		log.Println(err)
	}

	client = db.InitMongodb(keys["MONGODB"]["USER"], keys["MONGODB"]["PASSWORD"])
	defer client.Disconnect(context.TODO())
	err = client.Ping(context.Background(), readpref.Primary())
	
    if err != nil {
        log.Fatal("Couldn't connect to the database", err)
    } else {
        log.Println("Connected!")
	}
	
	pool := websocket.NewPool()
	go pool.Start()

	mux := mux.NewRouter()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
	mux.HandleFunc("/about", about)
	mux.HandleFunc("/login", mw.Chain(Authenticate, mw.MethodHandler("POST")))
	mux.HandleFunc("/addFact", mw.Chain(addFact, mw.MethodHandler("POST")))
	mux.HandleFunc("/articles", mw.Chain(Articles, mw.MethodHandler("GET", "PUT")))
	mux.HandleFunc("/articles/{category}", mw.Chain(GetArticlesOfCategory, mw.MethodHandler("GET")))
	mux.HandleFunc("/articles/{category}/{id}", mw.Chain(GetArticle, mw.MethodHandler("GET")))
	mux.HandleFunc("/addArticle", mw.Chain(addArticle, mw.MethodHandler("POST")))
	mux.HandleFunc("/addArticleDraft", mw.Chain(addArticleDraft, mw.MethodHandler("POST")))
	// mux.HandleFunc("/test/", mw.Chain(Articles, mw.MethodHandler("GET")))
	// mux.HandleFunc("/test/{category}/{id}", mw.Chain(Articles, mw.MethodHandler("GET")))

	// setupRoutes()
	// azureDB = db.InitAzureDB(keys)

	log.Println("Now server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
