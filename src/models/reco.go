package models

import (
	"math/rand"
	"time"

	"github.com/golang/glog"
	"labix.org/v2/mgo/bson"
)

const (
	MaxInt = int(^uint(0) >> 1)
	MinInt = -(MaxInt - 1)
)

type Meta struct {
	Weight float64
	Count  int // number of seen in category
	Score  int // up down voted
}

type MetaUser struct {
	Euclidean float64
	Pearson   float64
}

type User struct {
	ID         bson.ObjectId "_id,omitempty"
	Name       string
	Categories map[string]*Meta
	Items      map[string]int       // item and score
	Sims       map[string]*MetaUser // similarity scores to other users
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// NewUser [...]
func NewUser() *User {
	return &User{
		Name:       "",
		Categories: map[string]*Meta{},
		Items:      map[string]int{},
		Sims:       map[string]*MetaUser{},
	}
}

// WeightedChoice selects a random weighted
// index from a list of weights.
func WeightedChoice(weights []float64) int {
	totals := []float64{}
	runningTotal := 0.0

	for _, w := range weights {
		runningTotal += w
		totals = append(totals, runningTotal)
	}

	rnd := rand.Float64() * runningTotal
	for i, total := range totals {
		if rnd < total {
			return i
		}
	}

	return 0
}

// Recommendation [...]
func (d *DB) Recommendation(user *User) (*Article, error) {
	currentCategories, err := d.GetCategories()
	if err != nil {
		return nil, err
	}

	for _, c := range currentCategories {
		_, ok := user.Categories[c.ID.Hex()]
		// Found a category not explored yet
		if !ok {
			article, err := d.NextArticle(c.CategoryID, 0, user)
			if err != nil {
				glog.Error(err)
				continue
			}
			return article, nil
		}
	}

	// Try 20 times for an article and give up, naive
	for i := 0; i < 20; i++ {
		i++
		// All categories used, use weighted random choice to select next.
		categoryIDList := []string{}
		weightList := []float64{}
		for catName, meta := range user.Categories {
			categoryIDList = append(categoryIDList, catName)
			weightList = append(weightList, meta.Weight)
		}
		catIndex := WeightedChoice(weightList)
		var catID uint16
		for _, c := range currentCategories {
			if len(categoryIDList) > 0 && categoryIDList[catIndex] == c.ID.Hex() {
				catID = c.CategoryID
				break
			}
		}
		article, err := d.NextArticle(catID, 0, user)
		if err != nil {
			glog.Error(err)
			continue
		}
		return article, nil
	}
	article, err := d.NextUnseenArticle(user)
	return article, err
	//return nil, fmt.Errorf("ran out of articles")
}

// Collect [...]
func (d *DB) Collect(user *User, params CollectParams) error {
	var article Article
	err := d.Items.Find(bson.M{"_id": bson.ObjectIdHex(params.ID)}).One(&article)
	if err != nil {
		glog.Errorf("%v %+v %+v", err, user, params)
		return err
	}

	var category Category
	err = d.Categories.Find(bson.M{"categoryid": article.CategoryID}).One(&category)
	if err != nil {
		glog.Errorf("%v %+v %+v", err, user, params)
		return err
	}

	user.Items[params.ID] = params.Like
	hex := category.ID.Hex()
	_, ok := user.Categories[hex]
	switch ok {
	case false:
		user.Categories[hex] = &Meta{}
		user.Categories[hex].Count = 1
		user.Categories[category.ID.Hex()].Score = params.Like
	case true:
		user.Categories[hex].Count += 1
		user.Categories[category.ID.Hex()].Score += params.Like
	}

	// Alg
	// - find min and max
	// - each (elem - min) / (max - min)
	//totalCategories := len(user.Categories)
	max := MinInt
	min := MaxInt
	for _, m := range user.Categories {
		if m.Score > max {
			max = m.Score
		}
		if m.Score < min {
			min = m.Score
		}
	}

	// Get normalized weight between 0 and 1
	// and set categories weight.
	for _, m := range user.Categories {
		norm := float64(m.Score-min) / (1 + float64(max-min))
		if norm == 0 {
			// give it some small chance
			norm = 0.01
		}
		m.Weight = norm
	}

	// Get users similarity scores
	// as an example implementation
	otherUsers, err := d.OtherUsers(user)
	if err == nil {
		p1 := map[string]float64{}
		for c, v := range user.Categories {
			p1[c] = v.Weight
		}

		for _, other := range otherUsers {
			p2 := map[string]float64{}
			for c, v := range other.Categories {
				p2[c] = v.Weight
			}

			wE := Euclidean(p1, p2)
			wP := Pearson(p1, p2)
			user.Sims[other.ID.Hex()] = &MetaUser{Euclidean: wE, Pearson: wP}
		}
	}

	//change := bson.M{"items": user.Items}
	err = d.Users.UpdateId(user.ID, user)
	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}
