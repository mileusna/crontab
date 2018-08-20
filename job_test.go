package crontab_test

import (
	"sync"
	"testing"
	"time"

	"github.com/mileusna/crontab"
)

func TestJobError(t *testing.T) {

	ctab := crontab.New()

	if err := ctab.AddJob("* * * * *", myFunc, 10); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := ctab.AddJob("* * * * *", nil); err == nil {
		t.Error("This AddJob should return Error, fn is nil")
	}

	var x int
	if err := ctab.AddJob("* * * * *", x); err == nil {
		t.Error("This AddJob should return Error, fn is not func kind")
	}

	if err := ctab.AddJob("* * * * *", myFunc2, "s", 10, 12); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := ctab.AddJob("* * * * *", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, args are not the correct type")
	}

	if err := ctab.AddJob("* * * * * *", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, syntax error")
	}

	// custom types and interfaces as function params
	var m MyTypeInterface
	if err := ctab.AddJob("* * * * *", myFuncStruct, m); err != nil {
		t.Error(err)
	}

	if err := ctab.AddJob("* * * * *", myFuncInterface, m); err != nil {
		t.Error(err)
	}

	var mwo MyTypeNoInterface
	if err := ctab.AddJob("* * * * *", myFuncInterface, mwo); err == nil {
		t.Error("This should return error, type that don't implements interface assigned as param")
	}

	ctab.Shutdown()
}

var testN int
var testS string

func TestCrontab(t *testing.T) {
	testN = 0
	testS = ""

	ctab := crontab.Fake(2) // fake crontab wiht 2sec timer to speed up test

	var wg sync.WaitGroup
	wg.Add(2)

	if err := ctab.AddJob("* * * * *", func() { testN++; wg.Done() }); err != nil {
		t.Fatal(err)
	}

	if err := ctab.AddJob("* * * * *", func(s string) { testS = s; wg.Done() }, "param"); err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}

	if testN != 1 {
		t.Error("func 1 not executed as scheduled")
	}

	if testS != "param" {
		t.Error("func 2 not executed as scheduled")
	}
	ctab.Shutdown()
}

func TestRunAll(t *testing.T) {
	testN = 0
	testS = ""

	ctab := crontab.New()

	if err := ctab.AddJob("* * * * *", func() { testN++ }); err != nil {
		t.Fatal(err)
	}

	if err := ctab.AddJob("* * * * *", func(s string) { testS = s }, "param"); err != nil {
		t.Fatal(err)
	}

	ctab.RunAll()
	time.Sleep(time.Second)

	if testN != 1 {
		t.Error("func not executed on RunAll()")
	}

	if testS != "param" {
		t.Error("func not executed on RunAll() or arg not passed")
	}

	ctab.Clear()
	ctab.RunAll()

	if testN != 1 {
		t.Error("Jobs not cleared")
	}

	if testS != "param" {
		t.Error("Jobs not cleared")
	}

	ctab.Shutdown()
}
