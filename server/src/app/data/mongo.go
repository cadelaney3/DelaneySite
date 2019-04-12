package mongo

import (
	"context"
	//"main"

	"github.com/mongodb/mongo-go-driver/mongo"
)

//var password = main.MongoPassword
var password = "Theucanes3!"
var Client, err = mongo.Connect(context.Background(), "mongodb+srv://cadelaney3:"+password+"@delaneycluster-i5f6m.gcp.mongodb.net/test?retryWrites=true", nil)
