package main

import (
	"sort"
	"encoding/csv"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"time"
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

	localmongo.SetMode(mgo.Monotonic, true)

	//Now we have the collection names in collections
	for _, collectionName := range collections {

		if collectionName == "RBShop" {
			var line bson.M
			currentCollection := ramenbeastdb.C(collectionName)
			err = currentCollection.Find(nil).One(&line)
			if err != nil {
				log.Fatal(err)
			}
			
			var headers []string
			for column := range line {
				headers = append(headers, fmt.Sprint(column))
			}

			//Sort headers otherwise they'll be random?
			sort.Strings(headers)

			//Create a file
			file, err := os.Create(collectionName+".csv")
			if err != nil {
				fmt.Println("ERROR: Cannot create file")
			}
			defer file.Close()

			writer := csv.NewWriter(file)
			defer writer.Flush()

			writer.Write(headers)
			fmt.Printf("Printing headers for collection %v\n", collectionName)
			fmt.Println(headers)

			iteratingCursor := currentCollection.Find(nil).Iter()
			for iteratingCursor.Next(&line) {
				fmt.Println(line)

				valueArray := make([]string, len(line))
				for _, value := range line {
					valueArray = append(valueArray, value.(string))
				}

				writer.Write(valueArray)
			}

			if err := iteratingCursor.Close(); err != nil {
				panic(err)
			}
		}
	}
	
}

//Function below inspired by https://github.com/Zenithar/mgoexport/blob/master/main.go
func flatten(input bson.M, lkey string, flattened *map[string]interface{}) {
	for rkey, value := range input {
		key := lkey + rkey

		if _, ok := value.(string); ok {
			(*flattened)[key] = value.(string)
		} else if _, ok := value.(float64); ok {
			(*flattened)[key] = value.(float64)
		} else if _, ok = value.(int); ok {
			(*flattened)[key] = value.(int)
		} else if _, ok = value.(int64); ok {
			(*flattened)[key] = value.(int64)
		} else if _, ok = value.(bool); ok {
			(*flattened)[key] = value.(bool)
		} else if _, ok = value.(time.Time); ok {
			(*flattened)[key] = value.(time.Time).Format("2006-01-02T15:04:05Z07:00")
		} else if _, ok := value.(bson.ObjectId); ok {
			(*flattened)[key] = value.(bson.ObjectId).Hex
		} else if _, ok := value.([]interface{}); ok {

		}
	}
}