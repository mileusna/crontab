package crontab_test

import (
	"fmt"
	"log"

	"github.com/mileusna/crontab"
)

func ExampleCrontab() {

	ctab := crontab.New() // create cron table

	// MustAddJob panics on wrong syntax or problem with func and args for easier initialization
	ctab.MustAddJob("0 12 * * *", myFunc3)
	ctab.MustAddJob("* * * * *", myFunc2, "on every minute", 123) // fn with args
	ctab.MustAddJob("*/2 * * * *", myFunc2, "every two min", 18)

	// or use AddJob if you want to test the error
	err := ctab.AddJob("* * * * *", myFunc)
	if err != nil {
		log.Println(err)
		return
	}

	// all your other app code as usual, or put sleep timer for example
	// time.Sleep(5 * time.Minute)
}

func myFunc() {
	fmt.Println("Helo, world")
}

func myFunc3() {
	fmt.Println("Noon!")
}

func myFunc2(s string, n int) {
	fmt.Printf("We have params here, string `%s` and nymber %d\n", s, n)
}

type MyTypeInterface struct {
	ID   int
	Name string
}

func (m MyTypeInterface) Bar() string {
	return "OK"
}

type MyTypeNoInterface struct {
	ID   int
	Name string
}

func myFuncStruct(m MyTypeInterface) {
	fmt.Println("Custom type as param")
}

func myFuncInterface(i Foo) {
	i.Bar()
}

type Foo interface {
	Bar() string
}
