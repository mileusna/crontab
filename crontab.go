package crontab

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Crontab struct representing cron table
type Crontab struct {
	ticker *time.Ticker
	jobs   []job
}

// job in cron table
type job struct {
	min       map[int]struct{}
	hour      map[int]struct{}
	day       map[int]struct{}
	month     map[int]struct{}
	dayOfWeek map[int]struct{}

	fn   interface{}
	args []interface{}
}

// tick is individual tick that occures each minute
type tick struct {
	min       int
	hour      int
	day       int
	month     int
	dayOfWeek int
}

// New initializes and returns new cron table
func New() *Crontab {
	return new(time.Minute)
}

// new creates new crontab, arg provided for testing purpose
func new(t time.Duration) *Crontab {
	c := &Crontab{
		ticker: time.NewTicker(t),
	}

	go func() {
		for t := range c.ticker.C {
			c.runScheduled(t)
		}
	}()

	return c
}

// AddJob to cron table
//
// Returns error if:
//
// * Cron syntax can't be parsed or out of bounds
//
// * fn is not function
//
// * Provided args don't match the number and/or the type of fn args
func (c *Crontab) AddJob(schedule string, fn interface{}, args ...interface{}) error {
	j, err := parseSchedule(schedule)
	if err != nil {
		return err
	}

	if fn == nil || reflect.ValueOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("Cron job must be func()")
	}

	fnType := reflect.TypeOf(fn)
	if len(args) != fnType.NumIn() {
		return fmt.Errorf("Number of func() params and number of provided params doesn't match")
	}

	for i := 0; i < fnType.NumIn(); i++ {
		a := args[i]
		t1 := fnType.In(i)
		t2 := reflect.TypeOf(a)

		if t1 != t2 {
			if t1.Kind() != reflect.Interface {
				return fmt.Errorf("Param with index %d shold be `%s` not `%s`", i, t1, t2)
			}
			if !t2.Implements(t1) {
				return fmt.Errorf("Param with index %d of type `%s` doesn't implement interface `%s`", i, t2, t1)
			}
		}
	}

	// all checked, add job to cron tab
	j.fn = fn
	j.args = args
	c.jobs = append(c.jobs, j)
	return nil
}

// MustAddJob is like AddJob but panics if there is an problem with job
//
// It simplifies initialization, since we usually add jobs at the beggining so you won't have to check for errors (it will panic when program starts).
// It is a similar aproach as go's std lib package `regexp` and `regexp.Compile()` `regexp.MustCompile()`
// MustAddJob will panic if:
//
// * Cron syntax can't be parsed or out of bounds
//
// * fn is not function
//
// * Provided args don't match the number and/or the type of fn args
func (c *Crontab) MustAddJob(schedule string, fn interface{}, args ...interface{}) {
	if err := c.AddJob(schedule, fn, args...); err != nil {
		panic(err)
	}
}

// Shutdown the cron table schedule
//
// Once stopped, it can't be restarted.
// This function is pre-shuttdown helper for your app, there is no Start/Stop functionallity with crontab package.
func (c *Crontab) Shutdown() {
	c.ticker.Stop()
}

// Clear all jobs from cron table
func (c *Crontab) Clear() {
	c.jobs = []job{}
}

// RunAll jobs in cron table, shcheduled or not
func (c *Crontab) RunAll() {
	for _, j := range c.jobs {
		go j.run()
	}
}

// RunScheduled jobs
func (c *Crontab) runScheduled(t time.Time) {
	tick := getTick(t)
	for _, j := range c.jobs {
		if j.tick(tick) {
			go j.run()
		}
	}
}

// run the job using reflection
// Recover from panic although all functions and params are checked by AddJob, but you never know.
func (j job) run() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Crontab error", r)
		}
	}()
	v := reflect.ValueOf(j.fn)
	rargs := make([]reflect.Value, len(j.args))
	for i, a := range j.args {
		rargs[i] = reflect.ValueOf(a)
	}
	v.Call(rargs)
}

// tick decides should the job be lauhcned at the tick
func (j job) tick(t tick) bool {
	if _, ok := j.min[t.min]; !ok {
		return false
	}

	if _, ok := j.hour[t.hour]; !ok {
		return false
	}

	// cummulative day and dayOfWeek, as it should be
	_, day := j.day[t.day]
	_, dayOfWeek := j.dayOfWeek[t.dayOfWeek]
	if !day && !dayOfWeek {
		return false
	}

	if _, ok := j.month[t.month]; !ok {
		return false
	}

	return true
}

// regexps for parsing schedyle string
var (
	matchSpaces = regexp.MustCompile("\\s+")
	matchN      = regexp.MustCompile("(.*)/(\\d+)")
	matchRange  = regexp.MustCompile("^(\\d+)-(\\d+)$")
)

// parseSchedule string and creates job struct with filled times to launch, or error if synthax is wrong
func parseSchedule(s string) (j job, err error) {
	s = matchSpaces.ReplaceAllLiteralString(s, " ")
	parts := strings.Split(s, " ")
	if len(parts) != 5 {
		return job{}, errors.New("Schedule string must have five components like * * * * *")
	}

	j.min, err = parsePart(parts[0], 0, 59)
	if err != nil {
		return j, err
	}

	j.hour, err = parsePart(parts[1], 0, 23)
	if err != nil {
		return j, err
	}

	j.day, err = parsePart(parts[2], 1, 31)
	if err != nil {
		return j, err
	}

	j.month, err = parsePart(parts[3], 1, 12)
	if err != nil {
		return j, err
	}

	j.dayOfWeek, err = parsePart(parts[4], 0, 6)
	if err != nil {
		return j, err
	}

	//  day/dayOfWeek combination
	switch {
	case len(j.day) < 31 && len(j.dayOfWeek) == 7: // day set, but not dayOfWeek, clear dayOfWeek
		j.dayOfWeek = make(map[int]struct{})
	case len(j.dayOfWeek) < 7 && len(j.day) == 31: // dayOfWeek set, but not day, clear day
		j.day = make(map[int]struct{})
	default:
		// both day and dayOfWeek are * or both are set, use combined
		// i.e. don't do anything here
	}

	return j, nil
}

// parsePart parse individual schedule part from schedule string
func parsePart(s string, min, max int) (map[int]struct{}, error) {

	r := make(map[int]struct{}, 0)

	// wildcard pattern
	if s == "*" {
		for i := min; i <= max; i++ {
			r[i] = struct{}{}
		}
		return r, nil
	}

	// */2 1-59/5 pattern
	if matches := matchN.FindStringSubmatch(s); matches != nil {
		localMin := min
		localMax := max
		if matches[1] != "" && matches[1] != "*" {
			if rng := matchRange.FindStringSubmatch(matches[1]); rng != nil {
				localMin, _ = strconv.Atoi(rng[1])
				localMax, _ = strconv.Atoi(rng[2])
				if localMin < min || localMax > max {
					return nil, fmt.Errorf("Out of range for %s in %s. %s must be in range %d-%d", rng[1], s, rng[1], min, max)
				}
			} else {
				return nil, fmt.Errorf("Unable to parse %s part in %s", matches[1], s)
			}
		}
		n, _ := strconv.Atoi(matches[2])
		for i := localMin; i <= localMax; i += n {
			r[i] = struct{}{}
		}
		return r, nil
	}

	// 1,2,4  or 1,2,10-15,20,30-45 pattern
	parts := strings.Split(s, ",")
	for _, x := range parts {
		if rng := matchRange.FindStringSubmatch(x); rng != nil {
			localMin, _ := strconv.Atoi(rng[1])
			localMax, _ := strconv.Atoi(rng[2])
			if localMin < min || localMax > max {
				return nil, fmt.Errorf("Out of range for %s in %s. %s must be in range %d-%d", x, s, x, min, max)
			}
			for i := localMin; i <= localMax; i++ {
				r[i] = struct{}{}
			}
		} else if i, err := strconv.Atoi(x); err == nil {
			if i < min || i > max {
				return nil, fmt.Errorf("Out of range for %d in %s. %d must be in range %d-%d", i, s, i, min, max)
			}
			r[i] = struct{}{}
		} else {
			return nil, fmt.Errorf("Unable to parse %s part in %s", x, s)
		}
	}

	if len(r) == 0 {
		return nil, fmt.Errorf("Unable to parse %s", s)
	}

	return r, nil
}

// getTick returns the tick struct from time
func getTick(t time.Time) tick {
	return tick{
		min:       t.Minute(),
		hour:      t.Hour(),
		day:       t.Day(),
		month:     int(t.Month()),
		dayOfWeek: int(t.Weekday()),
	}
}
