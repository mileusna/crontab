package crontab_test

import (
	"testing"
	"time"

	"github.com/mileusna/crontab"
)

func TestJobError(t *testing.T) {

	ctab := crontab.New()

	if err := ctab.AddJob("* * * * *", myFunc, 10); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	var x int
	if err := ctab.AddJob("* * * * *", x); err == nil {
		t.Error("This AddJob should return Error, fn is not func kind")
	}

	if err := ctab.AddJob("* * * * *", myFunc2, "s", 10, 12); err == nil {
		t.Error("This AddJob should return Error, wrong number of args")
	}

	if err := ctab.AddJob("* * * * *", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, arg are not the correct type")
	}

	if err := ctab.AddJob("* * * * * *", myFunc2, "s", "s2"); err == nil {
		t.Error("This AddJob should return Error, syntax error")
	}

	ctab.Shutdown()
}

var testN int
var testS string

func TestCrontab(t *testing.T) {

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
