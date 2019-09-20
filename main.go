package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name      string
	Phone     string
	Timestamp time.Time
}

var database string = "go-development"

func main() {

	session, err := mgo.Dial("mongodb://localhost:27017/" + database)
	if err != nil {
		fmt.Println(err)
	}

	// Cleanup
	defer session.Close()

	c := bootstrap(session)

	create(c)
	read(c)
	update(c)
	delete(c)
}

func bootstrap(s *mgo.Session) *mgo.Collection {

	s.DB(database).DropDatabase()
	c := s.DB(database).C("people")
	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		Background: true,
	}

	c.EnsureIndex(index)

	return c
}

func create(c *mgo.Collection) {

	err := c.Insert(
		&Person{Name: "Eli", Phone: "111-239-333", Timestamp: time.Now()},
		&Person{Name: "Ben", Phone: "111-239-331", Timestamp: time.Now()},
		&Person{Name: "Jun", Phone: "111-239-222", Timestamp: time.Now()},
		&Person{Name: "Len", Phone: "111-231-333", Timestamp: time.Now()},
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n\n")
}

func read(c *mgo.Collection) {

	// Read Once
	result := Person{}
	query := c.Find(bson.M{"name": "Eli"})
	query = query.Select(bson.M{"phone": 0})
	err := query.One(&result)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)

	// Read alot
	var results []Person
	query = c.Find(bson.M{}).Sort("-timestamp")
	err = query.All(&results)

	for _, value := range results {
		fmt.Println(value)
	}

	fmt.Printf("\n\n")
}

func update(c *mgo.Collection) {

	filter := bson.M{"name": "Eli"}
	change := bson.M{
		"$set": bson.M{
			"phone":     "+86 99 8888 7707",
			"timestamp": time.Now(),
		},
	}

	err := c.Update(filter, change)
	if err != nil {
		fmt.Println(err)
		return
	}

	read(c)
	fmt.Printf("\n\n")
}

func delete(c *mgo.Collection) {

	filter := bson.M{"name": "Ben"}
	err := c.Remove(filter)

	if err != nil {
		fmt.Println(err)
		return
	}

	read(c)
	fmt.Printf("\n\n")
}
