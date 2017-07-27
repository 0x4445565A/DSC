package main

import (
	"fmt"
	"github.com/hpcloud/tail"
)

func main() {
	logFile := "./access.log"

	t, err := tail.TailFile(logFile, tail.Config{Follow: true})
	if err != nil {
		panic(err)
	}

	// This loops forever
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}
