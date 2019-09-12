package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	//"html"
	//"net/url"
	"time"
	"strings"
	"github.com/cadelaney3/delaneySite/utils"
	//"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title string `json:"title"`
	Category string `json:"category"`
	Tags []string `json:"tags"`
	Description string `json:"description"`
	Content string `json:"content"`
	Date time.Time `json:"creation_time" bson:"creation_date"`
}

func Articles(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		GetAllArticles(w, r)
	} else if r.Method == "PUT" {
		PutArticle(w, r)
	}
	return
}

func PutArticle(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := client.Database("delaney-db").Collection("articles")
	incomingArticle := Article{}

	err := json.NewDecoder(r.Body).Decode(&incomingArticle)
	if err != nil {
		log.Println(err)		
	}
	incomingArticle.Date = time.Now()
	res, err := collection.InsertOne(ctx, incomingArticle)
	if err != nil {
		log.Println("Error inserting document into MongoDB: ", err)
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		utils.Response(w, message)
		return
	}
	fmt.Println("res: ", res)
	message := utils.Message(http.StatusOK, "Successfully inserted document")
	w.WriteHeader(http.StatusOK)
	utils.Response(w, message)
}

func GetAllArticles(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := client.Database("delaney-db").Collection("articles")
	documentsReturned, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Error getting documents: %v\n", err)
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)
		return
	}
	defer documentsReturned.Close(ctx)
	articles := make([]*Article, 0)
	for documentsReturned.Next(ctx) {
		article := &Article{}
		err = documentsReturned.Decode(article)
		if err != nil {
			log.Printf("Error decoding article: %v\n", err)
			message := utils.Message(http.StatusInternalServerError, "Internal server error")
			w.WriteHeader(http.StatusNoContent)
			utils.Response(w, message)
			return
		}
		articles = append(articles, article)
	}
	if err = documentsReturned.Err(); err != nil {
		log.Println("Error return documents: %v", err)
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)
		return
	}
	if len(articles) == 0 {
		message := utils.Message(http.StatusNoContent, "No articles found")
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)
		return				
	}
	message := utils.Message(http.StatusOK, "Successfully retrieved articles")
	message["items"] = articles
	w.WriteHeader(http.StatusOK)
	utils.Response(w, message)		
	w.WriteHeader(http.StatusOK)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("delaney-db").Collection("articles")
	path := strings.Split(r.URL.Path, "/")
	objId, err := primitive.ObjectIDFromHex(path[len(path)-1])
	if err != nil {
		log.Println("Error converting to type ObjectId: %v", err)
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusInternalServerError)
		utils.Response(w, message)
		return
	}
	article := &Article{}
	filter := bson.M{"_id": objId}
	documentReturned := collection.FindOne(context.TODO(), filter)
	err = documentReturned.Decode(article)
	if err != nil {
		log.Println("Error decoding article: %v", err)
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)
		return
	}
	if article.Title != "" {
		message := utils.Message(http.StatusNoContent, "No article with id: " + path[len(path)-1])
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)	
		return	
	}
	message := utils.Message(http.StatusOK, "Successfully retrieved article")
	message["item"] = article
	w.WriteHeader(http.StatusOK)
	utils.Response(w, message)			
}

func GetArticlesOfCategory(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	collection := client.Database("delaney-db").Collection("articles")
	path := strings.Split(r.URL.Path, "/")
	filter := bson.M{"category": path[len(path)-1]}
	documentsReturned, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error getting category documents: %v\n", err)
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)
		return
	}
	defer documentsReturned.Close(ctx)
	articles := make([]*Article, 0)
	for documentsReturned.Next(ctx) {
		article := &Article{}
		err = documentsReturned.Decode(article)
		if err != nil {
			log.Printf("Error decoding article: %v\n", err)
			message := utils.Message(http.StatusInternalServerError, "Internal server error")
			w.WriteHeader(http.StatusNoContent)
			utils.Response(w, message)
			return
		}
		articles = append(articles, article)
	}
	if err = documentsReturned.Err(); err != nil {
		log.Println("Error return documents: %v", err)
		message := utils.Message(http.StatusInternalServerError, "Internal server error")
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)
		return
	}
	if len(articles) == 0 {
		message := utils.Message(http.StatusNoContent, "No article of category: " + path[len(path)-1])
		w.WriteHeader(http.StatusNoContent)
		utils.Response(w, message)
		return				
	}
	message := utils.Message(http.StatusOK, "Successfully retrieved articles")
	message["items"] = articles
	w.WriteHeader(http.StatusOK)
	utils.Response(w, message)		
}