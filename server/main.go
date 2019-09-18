package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	mw "github.com/cadelaney3/delaneySite/middleware"
	"github.com/cadelaney3/delaneySite/pkg/websocket"
	"github.com/cadelaney3/delaneySite/store/mongodb"
	"github.com/cadelaney3/delaneySite/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var keys = make(map[string]map[string]string)
var client *mongo.Client

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key         = []byte(os.Getenv("SESSION_KEY"))
	cookieStore = sessions.NewCookieStore(key)
)

//var validPath = regexp.MustCompile("^/(ws|edit|save|view)/([a-zA-Z0-9]+)$")
var validPath = regexp.MustCompile("^/(ws|view|home|signin)")

type response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type fact struct {
	Fact string `json:"fact"`
}

type Credentials struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Token    string             `json:"token"`
}

type Token struct {
	UserId string
	jwt.StandardClaims
}

type ArticleStore struct {
	articleStore mongodb.Article
}

type Client struct {
	cli *mongodb.Client
}

func (c *Client) FetchAllArticles(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	articles, err := c.cli.FetchAllArticles(ctx, false)
	if err != nil {
		log.Printf("Error fetching articles: %s", err)
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)
		return
	}
	message := utils.Message(http.StatusOK, "Successfully retrieved articles")
	message["items"] = articles
	w.WriteHeader(http.StatusOK)
	utils.Response(w, message)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	credentials := Credentials{}
	session, _ := cookieStore.Get(r, "cookie-name")

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

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStore.Get(r, "cookie-name")

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

	c := &Client{}
	c.cli, err = mongodb.InitMongodb(keys["MONGODB"]["USER"], keys["MONGODB"]["PASSWORD"])
	fmt.Printf("%T\n", c.cli)
	if err != nil {
		log.Printf("Error connection to mongdb: %s", err)
	}

	defer c.cli.Client.Disconnect(context.TODO())

	pool := websocket.NewPool()
	go pool.Start()

	mux := mux.NewRouter()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
	// mux.HandleFunc("/about", about)
	//mux.HandleFunc("/login", mw.Chain(Authenticate, mw.MethodHandler("POST")))
	//mux.HandleFunc("/addFact", mw.Chain(addFact, mw.MethodHandler("POST")))
	mux.HandleFunc("/articles", mw.Chain(c.FetchAllArticles, mw.MethodHandler("GET", "PUT")))
	//mux.HandleFunc("/articles/{category}", mw.Chain(c.SearchArticlesByCategory, mw.MethodHandler("GET")))
	//mux.HandleFunc("/articles/{category}/{id}", mw.Chain(c.SearchArticlesById, mw.MethodHandler("GET")))
	//mux.HandleFunc("/articles/drafts/", mw.Chain(GetArticles))
	//mux.HandleFunc("/addArticle", mw.Chain(addArticle, mw.MethodHandler("POST")))
	//mux.HandleFunc("/addArticleDraft", mw.Chain(addArticleDraft, mw.MethodHandler("POST")))
	// mux.HandleFunc("/test/", mw.Chain(Articles, mw.MethodHandler("GET")))
	// mux.HandleFunc("/test/{category}/{id}", mw.Chain(Articles, mw.MethodHandler("GET")))

	log.Println("Now server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
