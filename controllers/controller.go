package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var redisClient *redis.Client

const MONGO_URI = "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir-case-study?retryWrites=true"
const dbName = "getir-case-study"
const colName = "records"

const REDIS_URI = "localhost:6379"
const REDIS_PW = ""
const REDIS_DB = 0

func init() {
	clientOption := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB")
	collection = client.Database(dbName).Collection(colName)

	redisOption := &redis.Options{
		Addr: func() string {
			uri := os.Getenv("REDIS_URL")
			if uri == "" {
				return REDIS_URI
			}
			return uri
		}(),
		Password: REDIS_PW,
		DB:       REDIS_DB,
	}
	redisClient = redis.NewClient(redisOption)
}
