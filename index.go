package main

import (
	"encoding/csv"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

func main() {

	localmongo, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	defer localmongo.Close()

	ramenbeastdb := localmongo.DB("ramen-beast")
	collections, err :=ramenbeastdb.CollectionNames()
	if err != nil {
		panic(err)
	}

	
}
