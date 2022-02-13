package controller

import (
	"encoding/json"
	"getir1/models"
	"io/ioutil"
	"net/http"
	"time"
)

type RedisHandler struct {
}

// Handle incoming requests
func (red *RedisHandler) Cache(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		red.post(w, r)
		return
	case "GET":
		red.get(w, r)
		return
	default: // Do not accept request methods except POST and GET
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only Get&Post Methods are Allowed."))
		return
	}
}

// Handle get requests
func (red *RedisHandler) get(w http.ResponseWriter, r *http.Request) {
	// Get key parameter from url
	keys, ok := r.URL.Query()["key"]

	// If key parameter does not exist, return an error
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("key parameter is missing in url"))
		return
	}

	// define key
	key := keys[0]

	// Read value from redis database with key
	val, err := redisClient.Get(key).Result()
	// If key does not exist in redis database, return error
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not found in database"))
		return
	}

	// Prepare response model
	response := models.RedisRequest{
		Key:   key,
		Value: val,
	}
	// Marshal response into json
	responseBody, _ := json.Marshal(response)

	// Return response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// Handle post requests
func (red *RedisHandler) post(w http.ResponseWriter, r *http.Request) {
	// Read request body
	bodyBytes, err := ioutil.ReadAll(r.Body)
	// Send error with message if cant read body
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Can't read request body."))
		return
	}

	// Define redis request struct for redis (will be filled from request body)
	var keyValue models.RedisRequest

	// Unmarshal request body to keyValue
	err = json.Unmarshal(bodyBytes, &keyValue)
	// Send error with message if request body cant be unmarshaled into keyValue struct
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key and value should be type string"))
		return
	}

	// Set key-value pair in redis database with expiration of 30 mins
	err = redisClient.Set(keyValue.Key, keyValue.Value, time.Minute*30).Err()
	// If unable to set key-value, return error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to write key-value to in-memory database. Please check if your database is up and running."))
		return
	}

	// Prepare response struct
	response := models.RedisRequest{
		Key:   keyValue.Key,
		Value: keyValue.Value,
	}
	// Marshal response into json
	responseBody, _ := json.Marshal(response)

	// Return response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

// Constructor for RedisHandler
func NewRedisHandler() *RedisHandler {
	return &RedisHandler{}
}
