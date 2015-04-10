package main

import (
	"fmt"
	"time"

	"github.com/gamedevja/backend/message"
)

func main() {
	for {
		fmt.Println(message.Hello)
		time.Sleep(10 * time.Second)
	}
}
