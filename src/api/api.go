package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/drone/routes"
	"github.com/golang/glog"
	"models"
)

type APIServer struct {
	Port string
	Env  string
	DB   *models.DB
}

// PingHandler responds with PONG
func (a *APIServer) IndexHandler(w http.ResponseWriter, req *http.Request) {
	//fmt.Fprint(w, a.Index)
	http.Redirect(w, req, "/static/", http.StatusTemporaryRedirect)
}

// IndexHandler responds with index.html
func (a *APIServer) PingHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "PONG")
}

// CollectHandler collects user clicks
func (a *APIServer) CollectHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	userID := params.Get(":user")
	if userID == "" {
		http.Error(w, "Missing user", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "POST":
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			glog.Error(err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		user, err := a.DB.User(userID)
		if err != nil {
			glog.Error(err)
			http.Error(w, "Crap", http.StatusInternalServerError)
			return
		}

		var params models.CollectParams
		err = json.Unmarshal(body, &params)
		if err != nil {
			glog.Error(err)
			http.Error(w, "Crap", http.StatusInternalServerError)
			return
		}

		err = a.DB.Collect(user, params)
		if err != nil {
			glog.Error(err)
			http.Error(w, "Crap", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "OK")
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

// RecommendHandler responds an item for the given user.
func (a *APIServer) RecommendHandler(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	userID := params.Get(":user")
	if userID == "" {
		http.Error(w, "Missing user", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "GET":
		user, err := a.DB.User(userID)
		if err != nil {
			glog.Error(err)
			http.Error(w, "Crap", http.StatusInternalServerError)
			return
		}

		item, err := a.DB.Recommendation(user)
		if err != nil {
			glog.Error(err)
			http.Error(w, "Crap you've reached the end", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(item)
		if err != nil {
			glog.Error(err)
			http.Error(w, "Crap", http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, string(resp))
	default:
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

// Run [...]
func (a *APIServer) Run() {
	err := http.ListenAndServe(":"+a.Port, nil)
	if err != nil {
		glog.Fatal(err)
	}
}

// Close [...]
func (a *APIServer) Close() {
}

// New [...]
func New(port, env string, db *models.DB) (*APIServer, error) {
	a := &APIServer{Port: port, Env: env, DB: db}

	mux := routes.New()
	mux.Get("/ping", a.PingHandler)
	mux.Get("/relevant/:user", a.RecommendHandler)
	mux.Post("/collect/:user", a.CollectHandler)

	http.Handle("/", mux)
	go a.Run()

	return a, nil
}
