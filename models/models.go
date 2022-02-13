package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoFilter struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type MongoRecord struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Counts    []int              `json:"counts" bson:"counts"`
	CreatedAt primitive.DateTime `json:"createdAt" bson:"createdAt"`
	Key       string             `json:"key" bson:"key"`
	Value     string             `json:"value" bson:"value"`
}

type MongoResponse struct {
	Code    int                   `json:"code"`
	Msg     string                `json:"msg"`
	Records []MongoResponseRecord `json:"records"`
}

type MongoResponseRecord struct {
	Key        string `json:"key"`
	CreatedAt  string `json:"createdAt"`
	TotalCount int    `json:"totalCount"`
}

type RedisRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
