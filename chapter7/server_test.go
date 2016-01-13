package main

import (
	"testing"

	"labix.org/v2/mgo"
)

func initMgo(t *testing.T, dbname string) *mgo.Database {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	db := session.DB(dbname)
	db.DropDatabase()

	return db
}

func TestServer(t *testing.T) {
	db := initMgo(t, "short_test")

	coll := db.C("short")
	ist, err := NewInstance(coll)
	if err != nil {
		t.Fatal("new instance error -", err)
	}

	short, err := ist.Short("www.baidu.com")
	if err != nil {
		t.Fatal("short error -", err)
	}

	if short != "o0ve3u" {
		t.Fatal("short unpexted:", short)
	}

	short, err = ist.Short("www.google.com")
	if err != nil {
		t.Fatal("short error -", err)
	}

	if short != "o0ve3v" {
		t.Fatal("short unpexted:", short)
	}

	// test pickIdx
	ist2, err := NewInstance(coll)
	if err != nil {
		t.Fatal("new instance error -", err)
	}

	short, err = ist2.Short("www.qiniu.com")
	if err != nil {
		t.Fatal("short error -", err)
	}
	if short != "o0ve3w" {
		t.Fatal("short unpexted:", short)
	}

	// test redirect
	URL, err := ist.Redirect("o0ve3u")
	if err != nil {
		t.Fatal("redirect err -", err)
	}

	if URL != "www.baidu.com" {
		t.Fatal("redirect unexpected:", URL)
	}

	URL, err = ist.Redirect("o0ve3w")
	if err != nil {
		t.Fatal("redirect err -", err)
	}

	if URL != "www.qiniu.com" {
		t.Fatal("redirect unexpected:", URL)
	}
}
