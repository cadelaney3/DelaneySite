package mongodb

import (
	"context"
	"fmt"

	//"html"
	//"net/url"

	"time"

	"github.com/cadelaney3/delaneySite/server/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Category    string             `json:"category"`
	Tags        []string           `json:"tags"`
	Description string             `json:"description"`
	Content     string             `json:"content"`
	Date        time.Time          `json:"creation_time" bson:"creation_date"`
}

func (c *Client) FetchAllArticles(ctx context.Context, draft bool) (api.Article, error) {
	var collection *mongo.Collection
	if draft {
		collection = c.db.Collection("articleDrafts")
	} else {
		collection = c.db.Collection("articles")
	}
	documentsReturned, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("Error retrieving articles: %s", err)
	}
	defer documentsReturned.Close(ctx)
	articles := make([]*api.Article, 0)
	for documentsReturned.Next(ctx) {
		article := &api.Article{}
		err = documentsReturned.Decode(article)
		if err != nil {
			return nil, fmt.Errorf("Error decoding article: %s", err)
		}
		articles = append(articles, article)
	}
	if err = documentsReturned.Err(); err != nil {
		return nil, fmt.Errorf("Error return documents: %s", err)
	}
	return articles, nil
}

func (c *Client) SaveArticle(ctx context.Context, article *api.Article, draft bool) (string, error) {
	var collection *mongo.Collection
	if draft {
		collection = c.db.Collection("articleDrafts")
	} else {
		collection = c.db.Collection("articles")
	}
	article.Date = time.Now()
	_, err := collection.InsertOne(ctx, article)
	if err != nil {
		return "Unsuccessful document insert", fmt.Errorf("Error inserting document: %s", err)
	}

	return "Successfully inserted document", nil
}

func (c *Client) DeleteArticleById(ctx context.Context, id string, draft bool) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Error converting id to bson objectId: %s", err)
	}

	var collection *mongo.Collection
	if draft {
		collection = c.db.Collection("articleDrafts")
	} else {
		collection = c.db.Collection("articles")
	}

	err = collection.DeleteOne(bson.M{"_id": objId})
	if err != nil {
		return fmt.Errorf("Error deleting document with id %s: %s", id, err)
	}
	return nil
}

func (c *Client) SearchArticlesById(ctx context.Context, id string) (*api.Article, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("Error converting id to bson objectId: %s", err)
	}
	collection := c.db.Collection("articles")

	article := &api.Article{}
	filter := bson.M{"_id": objId}
	documentReturned := collection.FindOne(context.TODO(), filter)
	err = documentReturned.Decode(article)
	if err != nil {
		return nil, fmt.Errorf("Error decoding article: %s", err)
	}
	return article, nil
}

func (c *Client) SearchArticlesByCategory(ctx context.Context, category string) ([]*api.Article, error) {
	collection := c.db.Collection("articles")
	filter := bson.M{"category": category}
	documentsReturned, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving articles of category '%s': %s", category, err)
	}
	defer documentsReturned.Close(ctx)
	articles := make([]*api.Article, 0)
	for documentsReturned.Next(ctx) {
		article := &api.Article{}
		err = documentsReturned.Decode(article)
		if err != nil {
			return nil, fmt.Errorf("Error decoding article: %s", err)
		}
		articles = append(articles, article)
	}
	if err = documentsReturned.Err(); err != nil {
		return nil, fmt.Errorf("Error returning documents: %s", err)
	}
	return articles, nil
}
