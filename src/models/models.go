package models

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/glog"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	NumberToGet uint32 = 100
)

type DB struct {
	Path       string
	Session    *mgo.Session
	DB         *mgo.Database
	Users      *mgo.Collection
	Categories *mgo.Collection
	Items      *mgo.Collection
}

type CollectParams struct {
	ID   string `json:"id"`
	Like int    `json:"like"`
}

// New [...]
func New(dbPath string) (*DB, error) {
	session, err := mgo.Dial(dbPath)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	d := session.DB("pavlos")
	db := &DB{Path: dbPath, Session: session, DB: d}

	db.Users = d.C("Users")
	db.Categories = d.C("Categories")
	db.Items = d.C("Items")

	index := mgo.Index{
		Key: []string{"categoryid"},
		//Unique:     true,
		//DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	db.Categories.EnsureIndex(index)

	index = mgo.Index{
		Key: []string{"categoryid"},
		//Unique:     true,
		//DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	db.Items.EnsureIndex(index)

	return db, nil
}

// Get one user unseen article
func (d *DB) NextArticle(categoryID uint16, count int, user *User) (*Article, error) {
	articles, err := d.Articles(categoryID, count, user)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	for _, article := range articles {
		score, alreadySeen := user.Items[article.ID.Hex()]
		if !alreadySeen {
			return article, nil
		}
		// maybe show again, say after some time....
		if score == 0 {
		}
	}
	return nil, fmt.Errorf("none")
}

// OtherUsers
func (d *DB) OtherUsers(user *User) ([]*User, error) {
	users := []*User{}
	iter := d.Users.Find(bson.M{"_id": bson.M{"$ne": user.ID}}).Iter()
	err := iter.All(&users)
	return users, err
}

// In the event of running out of things to show.
func (d *DB) NextUnseenArticle(user *User) (*Article, error) {
	article := Article{}
	usersSeen := []bson.ObjectId{}
	for id, _ := range user.Items {
		usersSeen = append(usersSeen, bson.ObjectIdHex(id))
	}
	err := d.Items.Find(bson.M{"_id": bson.M{"$nin": usersSeen}}).One(&article)
	return &article, err
}

// NextUnseenArticles [...]
func (d *DB) NextUnseenArticles(user *User) ([]*Article, error) {
	articles := []*Article{}
	usersSeen := []bson.ObjectId{}
	for id, _ := range user.Items {
		usersSeen = append(usersSeen, bson.ObjectIdHex(id))
	}
	iter := d.Items.Find(bson.M{"_id": bson.M{"$nin": usersSeen}}).Iter()
	err := iter.All(&articles)
	return articles, err
}

// Articles given a categoryID and count returns a list of articles.
// Count of 0 means all
func (d *DB) Articles(categoryID uint16, count int, user *User) ([]*Article, error) {
	results := []*Article{}
	query := bson.M{"categoryid": categoryID}

	var iter *mgo.Iter
	switch count {
	case 0:
		iter = d.Items.Find(query).Iter()
	default:
		iter = d.Items.Find(bson.M{"categoryid": categoryID}).Limit(count).Iter()
	}
	err := iter.All(&results)
	return results, err
}

// User creates a user if missing or returns a found user.
func (d *DB) User(name string) (*User, error) {
	user := User{}
	err := d.Users.Find(bson.M{"name": name}).One(&user)
	if err != nil {
		glog.Error(err, name)

		user := NewUser()
		user.Name = name
		err := d.Users.Insert(user)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
	}

	return &user, nil
}

// GetCategories gets a list of all currently loaded
// categories.
func (d *DB) GetCategories() ([]*Category, error) {
	results := []*Category{}
	var iter *mgo.Iter = d.Categories.Find(nil).Iter()
	err := iter.All(&results)
	return results, err
}

// InitCategories creates a collection of categories
// It first removes all previously defined ones.
// TODO This should probably sync instead.
func (d *DB) InitCategories() {
	d.Categories.RemoveAll(nil)

	// Get latest categories.
	categories, err := FeedzillaCategories()
	if err != nil {
		glog.Error(err)
		return
	}

	for _, cat := range categories {
		err := d.Categories.Insert(cat)
		if err != nil {
			glog.Error(err)
			continue
		}
		glog.Infof("Loaded %+v", cat)
	}
}

// itemUpdater is a helper function to go get articles for each category
func (d *DB) itemUpdater(catogoryID uint16) {
	articles, err := FeedzillaArticles(catogoryID, NumberToGet)
	if err != nil {
		glog.Error(err)
		return
	}

	for _, art := range articles.List {
		art.CategoryID = catogoryID
		_, err := d.Items.Upsert(bson.M{"url": art.URL}, art)
		if err != nil {
			glog.Error(err)
			continue
		}
		glog.Infof("%+v", art)
	}
}

// ItemUpdaterCron fetches the latest articles for each collection
// periodically every hour.
func (d *DB) ItemUpdaterCron() {
	for {
		time.Sleep(time.Hour)
		d.ItemUpdater()
	}
}

// ItemUpdater fetches the latest articles for each collection
func (d *DB) ItemUpdater() {
	currentCategories, err := d.GetCategories()
	if err != nil {
		glog.Error(err)
		return
	}
	i := 0
	wg := &sync.WaitGroup{}
	for _, category := range currentCategories {
		i++
		/*
			if i > 10 {
				break
			}
		*/
		wg.Add(1)
		go func(cat *Category) {
			glog.Info("Getting category ", cat)
			d.itemUpdater(cat.CategoryID)
			wg.Done()
		}(category)
	}
	wg.Wait()
}
