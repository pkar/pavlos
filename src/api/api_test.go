package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	//"net/http/httptest"

	"models"
)

var goPath = fmt.Sprintf("%s", os.Getenv("GOPATH"))
var db *models.DB
var client *http.Client = &http.Client{}
var port = "8001"
var userName = "pavlos"

func init() {
	var err error
	db, err = models.New(goPath + "/db/test")
	if err != nil {
		log.Fatal(err)
	}

	_, err = New(port, "testing", db)
	if err != nil {
		log.Fatal(err)
	}
}

func TestRecommendHandler(t *testing.T) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/relevant/%s", port, userName), nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestCollectHandler(t *testing.T) {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%s/collect/%s", port, userName), strings.NewReader(`{"id":"larry","like":1}`))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
