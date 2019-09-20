package api

import (
	"context"
	"fmt"
	"time"

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

type ArticleStore interface {
	FetchAllArticles(context.Context) (*[]Article, error)
	SaveArticle(context.Context, *Article) (string, error)
	DeleteArticleById(context.Context, string) error
	GetArticleById(context.Context, string) (*Article, error)
	GetArticlesByCategory(context.Context, string) ([]*Article, error)
}

func (a *Article) String() string {
	return fmt.Sprintf(`
		{
			"id": %s,
			"title": %s,
			"category": %s,
			"tags": %v,
			"description": %s,
			"content": %s,
			"date": %v
		}
	`,
		a.ID, a.Title, a.Category, a.Tags,
		a.Description, a.Content, a.Date,
	)
}
