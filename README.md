# Go/Golang package for Crontab tickers [![GoDoc](https://godoc.org/github.com/mileusna/crontab?status.svg)](https://godoc.org/github.com/mileusna/crontab)

This package provides crontab tickers to golang apps, supporting crontab-like syntax like `* * * * *` or `*/2 * * * *` etc.

## Installation <a id="installation"></a>
```
go get github.com/mileusna/crontab
```

## Example<a id="example"></a>

```go
package main

import (
	"fmt"
	"log"

	"github.com/mileusna/crontab"
)

func main() {

    ctab := crontab.New() // create cron table

    // MustAddJob panics on wrong syntax or problem with func and args
    // use for easier initialization 
    ctab.MustAddJob("* * * * *", myFunc) // every minute
    ctab.MustAddJob("0 12 * * *", myFunc3) // noon lauch

    // fn with args
    ctab.MustAddJob("0 0 * * 1,2", myFunc2, "Monday and Tuesday midnight", 123) 
    ctab.MustAddJob("*/5 * * * *", myFunc2, "every five min", 0)

    // or use AddJob if you want to test the error
    err := ctab.AddJob("0 12 1 * *", myFunc) // on 1st day of month
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

```

## Crontab syntax <a id="syntax"></a>

If you are not faimiliar with crontab syntax you might be better with other similar packages. But here are few references about crontab syntax.

```
*     *     *     *     *        
-     -     -     -     -
|     |     |     |     |
|     |     |     |     +----- day of week (0-6) (Sunday=0)
|     |     |     +------- month (1-12)
|     |     +--------- day of month (1-31)
|     +----------- hour (0-23)
+------------- min (0-59)
```

### Examples

+ `* * * * *` run on every minuta
+ `10 * * * *` run on 0:10, 1:10 etc
+ `10 15 * * *` run on 10:15 every day
+ `* * 1 * *` run on every minute in 1st day of month
+ `0 0 1 1 *` Happy new year schedule
+ `0 0 * * 1` Run on every Monday in midnight

### Lists

+ `* 10,15,19 * * *` run in 10:00, 15:00 and 19:00
+ `1-15 * * * *` run on 1, 2, 3...15 minute of each hour
+ `0 0-5,10 * * *` run on every hour from 0-5 and in 10 oclock

### Steps
+ `*/2 * * * *` run every two minutes
+ `10 */3 * * *` run every 3 hours on 10th min
+ `0 12 */2 * *` run in noon every two days
+ `1-59/2 * * *` run every two minutes, but on odd minutes

## Notice

There is no way to reschedule or to remove single job from crontab during runtime. (Re)create new instance of crontab to reschedule during runtime or use `crontab.Clear()` function and add jobs if you prefere.



