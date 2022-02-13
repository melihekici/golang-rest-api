package controller

import (
	"bytes"
	"encoding/json"
	"getir1/models"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestWriteRedisService(t *testing.T) {

	jsonStr := []byte(`{"key":"testKey", "value":"testValue"}`)

	req, err := http.NewRequest(
		"POST",
		"http://localhost/in-memory",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Error(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()

	var result models.RedisRequest
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err.Error())
	}
	json.Unmarshal(b, &result)

	if result.Key != "testKey" || result.Value != "testValue" {
		t.Error("Key value error")
	}
}

func TestRWRedisClient(t *testing.T) {
	redisClient.Set("testKey2", "testValue2", 0)

	s, err := redisClient.Get("testKey2").Result()
	if err != nil {
		t.Error(err.Error())
	}

	if s != "testValue2" {
		t.Error("expected", "testValue2", "got", s)
	}
}

func TestNoMinCount(t *testing.T) {
	jsonStr := []byte(`{
		"startDate": "2016-01-01",
		"endDate": "2016-02-01",
		"maxCount": 3000
	}`)

	req, err := http.NewRequest(
		"POST",
		"http://localhost/records",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Error(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Error("Wrong Response Code")
	}
}

func TestNoMaxCount(t *testing.T) {
	jsonStr := []byte(`{
		"startDate": "2016-01-01",
		"endDate": "2016-02-01",
		"minCount": 1000
	}`)

	req, err := http.NewRequest(
		"POST",
		"http://localhost/records",
		bytes.NewBuffer(jsonStr),
	)
	if err != nil {
		t.Error(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Error("Wrong Response Code")
	}

}

func TestProperMongoRequest(t *testing.T) {
	jsonStrs := [][]byte{
		[]byte(`{"startDate": "2016-01-01","endDate": "2016-02-01","minCount": 2500,"maxCount": 2700}`),
		[]byte(`{"startDate": "2017-02-03","endDate": "2017-06-21","minCount": 2000,"maxCount": 3100}`),
		[]byte(`{"startDate": "2018-05-01","endDate": "2018-06-14","minCount": 1500,"maxCount": 1700}`),
	}

	for _, jsonStr := range jsonStrs {
		req, err := http.NewRequest(
			"POST",
			"http://localhost/records",
			bytes.NewBuffer(jsonStr),
		)
		if err != nil {
			t.Error(err.Error())
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Error("Wrong Response Code")
		}
	}
}
