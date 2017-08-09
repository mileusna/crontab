package crontab_test

import (
	"fmt"
	"log"

	"github.com/mileusna/crontab"
)

func ExampleCrontab() {

	ctab := crontab.New() // create cron table

	// MustAddJob panics on wrong syntax or problem with func and args for easier initialization
	ctab.MustAddJob("0 0 * * *", myFunc3)
	ctab.MustAddJob("* * * * *", myFunc2, "on every minute", 123) // fn with args
	ctab.MustAddJob("*/5 * * * *", myFunc2, "every five min", 0)

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
	fmt.Println("We have params here, string", s, "and number", n)
}
