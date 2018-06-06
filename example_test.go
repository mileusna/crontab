package crontab_test

import (
	"fmt"
	"log"
	"testing"
	"time"

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

// go test -v ./... -run TestExampleExecStats
func TestExampleExecStats(t *testing.T) {
	ctab := crontab.New()

	ctab.MustAddJob("* * * * *", myFuncWithStats, ctab.StatsChan())
	log.Println("Waiting a bit for the test to complete...")

	for i := 1; i <= 1; i++ {
		myExecStats := <-ctab.StatsChan()
		if myExecStats.JobType != "myFuncWithStats" {
			t.Errorf("Found an unexpected Job type")
		}
		customStuff := myExecStats.Stats().(*myCustomStats)
		if customStuff.strParam != "foo" {
			t.Errorf("Found an unexpected string parameter in the stats")
		}
		if customStuff.intParam != 42 {
			t.Errorf("Found an unexpected integer parameter in the stats")
		}
		ctab.Shutdown()
		log.Println("Done with the test, the received stats:", customStuff)
	}
}

// custom execution stats depending on the scheduled function
type myCustomStats struct {
	strParam string
	intParam int
}

func myFuncWithStats(statsChan chan crontab.ExecStats) {
	// work a bit...
	time.Sleep(1 * time.Second)
	// publish the execution stats...
	statsChan <- crontab.ExecStats{
		// ID to identify the job
		JobType: "myFuncWithStats",
		// custom execution stats
		Stats: func() interface{} {
			return &myCustomStats{
				strParam: "foo",
				intParam: 42,
			}
		},
	}
}
