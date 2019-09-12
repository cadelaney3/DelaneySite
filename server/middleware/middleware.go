package middleware

import (
	"context"
	"net/http"
	"strings"
	"os"
	"fmt"
	"log"
	"github.com/cadelaney3/delaneySite/utils"
	jwt "github.com/dgrijalva/jwt-go"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type Token struct {
	UserId string
	jwt.StandardClaims
}

// Auth is a middleware to check authentication when trying to login to the admin page
// based off: https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
func Auth() Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			auth := "/swagmaster"
			requestPath :=  r.URL.Path

			if auth != requestPath {
				return
			}

			response := make(map[string] interface{})
			tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
	
			if tokenHeader == "" { //Token is missing, returns with error code 403 Forbidden
				response = utils.Message(http.StatusForbidden, "Missing auth token")
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				utils.Response(w, response)
				return
			}

			splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
			if len(splitted) != 2 {
				response = utils.Message(http.StatusForbidden, "Invalid/Malformed auth token")
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				utils.Response(w, response)
				return
			}
	
			tokenPart := splitted[1] //Grab the token part, what we are truly interested in
			tk := &Token{}
	
			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("token_password")), nil
			})
	
			if err != nil { //Malformed token, returns with http code 403 as usual
				response = utils.Message(http.StatusForbidden, "Malformed authentication token")
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				utils.Response(w, response)
				return
			}
	
			if !token.Valid { //Token is invalid, maybe not signed on this server
				response = utils.Message(http.StatusForbidden, "Token is not valid.")
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				utils.Response(w, response)
				return
			}
	
			//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
			fmt.Sprintf("User %", tk.UserId) //Useful for monitoring
			ctx := context.WithValue(r.Context(), "user", tk.UserId)
			r = r.WithContext(ctx)
			fn(w, r)
		}
	}
}

func MethodHandler(methods ...string) Middleware {
	return func(fn http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method)
			for _, method := range methods {
				if r.Method == method || r.Method == "OPTIONS" {
					w.Header().Set("Content-Type", "application/json")
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Headers", "*")
					w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT")
					fn(w, r)
					return
				}
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
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