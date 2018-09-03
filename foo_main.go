package main

import (
	"fmt"

	"crontab"
)

func main() {
	fmt.Println(crontab.Fake_Msg("hello go crontab .."))
}
