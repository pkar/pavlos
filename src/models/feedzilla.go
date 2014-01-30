package models

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

var client *http.Client = &http.Client{}

type Category struct {
	ID          bson.ObjectId "_id,omitempty"
	CategoryID  uint16        `json:"category_id"`
	DisplayName string        `json:"display_category_name"`
	EnglishName string        `json:"english_category_name"`
	UrlName     string        `json:"url_category_name"`
}

type SubCategory struct {
	ID            bson.ObjectId "_id,omitempty"
	CategoryID    uint16        `json:"category_id"`
	SubCategoryID uint16        `json:"subcategory_id"`
	DisplayName   string        `json:"display_subcategory_name"`
	EnglishName   string        `json:"english_subcategory_name"`
	UrlName       string        `json:"url_subcategory_name"`
}

type Article struct {
	ID          bson.ObjectId "_id,omitempty"
	CategoryID  uint16
	PublishDate string `json:"publish_date"`
	Source      string `json:"source"`
	SourceURL   string `json:"source_url"`
	Summary     string `json:"summary"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

type Articles struct {
	List []*Article `json:"articles"`
}

// doRequest [...]
var doRequest = func(url string) ([]byte, error) {
	//request, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return body, nil
}

// FeedzillaCategories [...]
func FeedzillaCategories() ([]*Category, error) {
	url := "http://api.feedzilla.com/v1/categories.json"
	body, err := doRequest(url)
	if err != nil {
		glog.Error(err)
		return nil, nil
	}
	var ret []*Category
	err = json.Unmarshal(body, &ret)
	if err != nil {
		glog.Error(err, string(body))
		return nil, err
	}

	return ret, nil
}

// FeedzillaArticles [...]
func FeedzillaArticles(categoryID uint16, count uint32) (*Articles, error) {
	url := fmt.Sprintf("http://api.feedzilla.com/v1/categories/%d/articles.json?count=%d", categoryID, count)
	body, err := doRequest(url)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	var ret Articles
	err = json.Unmarshal(body, &ret)
	if err != nil {
		glog.Error(err, string(body))
		return nil, err
	}

	return &ret, nil
}

// FeedzillaSearchArticles [...]
func FeedzillaSearchArticles(categoryID uint16, query string, count uint32) (*Articles, error) {
	url := fmt.Sprintf("http://api.feedzilla.com/v1/categories/%d/articles/search.json?q=%s&count=%d", categoryID, query, count)
	body, err := doRequest(url)
	if err != nil {
		glog.Error(err, string(body))
		return nil, err
	}
	var ret Articles
	err = json.Unmarshal(body, &ret)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return &ret, nil
}
