package main

import (
	controller "getir1/controllers"
	"net/http"
)

func main() {
	http.HandleFunc("/records", controller.NewMongoRecordHandler().Records)
	http.HandleFunc("/in-memory", controller.NewRedisHandler().Cache)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
