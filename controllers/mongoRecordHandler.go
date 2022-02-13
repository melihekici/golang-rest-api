package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"getir1/models"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoRecordHandlers struct {
}

func (m *MongoRecordHandlers) post(w http.ResponseWriter, r *http.Request) {
	// Read request body
	bodyBytes, err := ioutil.ReadAll(r.Body)
	// Send error with message if cant read body
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Can't read request body."))
		return
	}

	// Define filter struct for mongoDB (will be filled from request body)
	var filter models.MongoFilter

	// Unmarshal request body to filter
	err = json.Unmarshal(bodyBytes, &filter)
	// Send error with message if request body cant be unmarshaled into filter struct
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

	// Get the records from mongodb which satisfies the filter
	res := GetRecordsByFilter(filter)

	// Marshal response body to json
	// We are throwing error away because it is handled inside above function
	// And error message is already exist in responseBody msg field.
	responseBody, _ := json.Marshal(res)

	// Add headers and send response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// Attach Records method to MongoRecordHandlers struct which will handle requests
func (m *MongoRecordHandlers) Records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST": // In case of POST request, call post method of struct
		m.post(w, r)
		return
	default: // Any requests except POST will be denied
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only Post Method is Allowed."))
		return
	}
}

// Constructor for MongoRecordHandlers struct
func NewMongoRecordHandler() *MongoRecordHandlers {
	return &MongoRecordHandlers{}
}

// function that takes filter and returns the query result from mongoDB
func GetRecordsByFilter(filter models.MongoFilter) models.MongoResponse {
	// Parse startDate string to date
	start, err := time.Parse("2006-01-02", filter.StartDate)
	// If parsing fails, return response with error code and message
	if err != nil {
		return models.MongoResponse{
			Code:    1,
			Msg:     "Unable to parse startDate, please check format(should be YYYY-MM-DD)",
			Records: nil,
		}
	}

	// Parse endDate string to date
	end, err := time.Parse("2006-01-02", filter.EndDate)
	// If parsing fails, return response with error code and message
	if err != nil {
		return models.MongoResponse{
			Code:    1,
			Msg:     "Unable to parse endDate, plase check format(should be YYYY-MM-DD)",
			Records: nil,
		}
	}

	// If startDate is a latter date than endDate, return a response with error and message
	if start.Unix() > end.Unix() {
		return models.MongoResponse{
			Code:    2,
			Msg:     "endDate should be a latter date than startDate",
			Records: nil,
		}
	}

	// If minCount is less than maxCount, return a response with error and message
	if filter.MaxCount < filter.MinCount {
		return models.MongoResponse{
			Code:    3,
			Msg:     "maxCount should be greater than minCount",
			Records: nil,
		}
	}

	// Construct query for mongoDB
	mongoFilter := bson.M{
		"$and": []bson.M{
			{"createdAt": bson.M{"$gte": start}},
			{"createdAt": bson.M{"$lte": end}},
		},
	}

	// Get cursor for the query result
	cur, err := collection.Find(context.Background(), mongoFilter)
	// If query fails, return a response with error and message
	if err != nil {
		return models.MongoResponse{
			Code:    4,
			Msg:     "Failure during fetching data from mongoDB collection.",
			Records: nil,
		}
	}
	// Close the cursor after function execution
	defer cur.Close(context.Background())

	// Create a slice of type MongoResponseRecord to keep query results
	var records []models.MongoResponseRecord

	// While cursor has another record
	for cur.Next(context.Background()) {
		// Create a variable for the current record
		var record models.MongoRecord

		// Decode the current record into struct MongoRecord
		err := cur.Decode(&record)
		// In case of an error during decoding of current record
		// I have decided to move to the next record instead of returning an error
		if err != nil {
			log.Println("Unable to decode record")
			continue
		}

		// Reading createdAt date from mongoDB and changing the string format
		var timeString string
		t := record.CreatedAt
		_, offset := t.Time().UTC().Zone()
		timeString = fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ%d",
			t.Time().Year(), t.Time().Month(), t.Time().Day(), t.Time().Hour(), t.Time().Minute(), t.Time().Second(), offset)

		// Calculating total count for current record
		totalCount := 0
		for _, v := range record.Counts {
			totalCount += v
		}

		// If the total count is inbetween the boundaries, add current record to the slice
		if totalCount >= filter.MinCount && totalCount <= filter.MaxCount {
			records = append(records, models.MongoResponseRecord{
				Key:        record.Key,
				CreatedAt:  timeString,
				TotalCount: totalCount,
			})
		}
	}

	// Send the response
	return models.MongoResponse{
		Code:    0,
		Msg:     "Success",
		Records: records,
	}
}
