package models

import (
	"log"
	"testing"
)

var d *DB

func init() {
	var err error
	d, err = New("localhost")
	if err != nil {
		log.Fatal(err)
	}
	d.DB = d.Session.DB("pavlos_test")
	doRequest = doRequestCategories
	go d.InitCategories()
}

func TestItemUpdaterHelper(t *testing.T) {
	doRequest = doRequestCategories
	currentCategories, err := d.GetCategories()
	if err != nil {
		t.Fatal(err)
	}

	first := currentCategories[0]
	t.Logf("%+v", first.ID)

	doRequest = doRequestArticles
	d.itemUpdater(first.CategoryID)
}

func TestUser(t *testing.T) {
	// Test new user generation
	for _, n := range []string{"a", "b", "c", "d"} {
		user, err := d.User(n)
		if err != nil {
			t.Fatal(err)
		}
		if user.Name != n {
			t.Fatalf("user name not set right %s %+v", n, user)
		}
	}

	// Test get previous users
	for _, n := range []string{"a", "b", "c", "d"} {
		user, err := d.User(n)
		if err != nil {
			t.Fatal(err)
		}
		if user.Name != n {
			t.Fatalf("user name not set right %s %+v", n, user)
		}
	}
}
