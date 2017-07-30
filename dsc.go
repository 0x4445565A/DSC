package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"regexp"
	"time"
	"os"
	"flag"
	"log"
)

// Make this global so we only have to compile it once
var requestFormatRegex *regexp.Regexp
var logFormatRegex *regexp.Regexp

// Compile our regex
func init() {
	requestFormatRegex = requestRegex()
	logFormatRegex = logRegex()
}

func main() {
	var logFile string
	var outFile string

	flag.StringVar(&logFile, "log", "/var/log/nginx/access.log", "Location of nginx log")
	flag.StringVar(&outFile, "output", "/var/log/stats.log", "Location of output log")

    flag.Parse()

	errorPaths := map[string]int{}
	statusCodes := map[string]int{}

	go logCollector(outFile, &statusCodes, &errorPaths)
	processLines(logFile, &statusCodes, &errorPaths)

}

func logCollector(outfile string, statusCodes *map[string]int, errorPaths *map[string]int) {
	tickChan := time.NewTicker(time.Second * 5).C
	for {
		select {
		case <-tickChan:
			errorPathsCopy := *errorPaths
			statusCodesCopy := *statusCodes
			*errorPaths = map[string]int{}
			*statusCodes = map[string]int{}

			//fmt.Println(outfile, outputDataToString(statusCodesCopy, errorPathsCopy))
			err := writeDataToLogfile(outfile, outputDataToString(statusCodesCopy, errorPathsCopy))
			if err != nil {
				log.Println(err)
			}

		default:
			continue
		}
	}
}

func processLines(logFile string, statusCodes *map[string]int, errorPaths *map[string]int) {
	t, err := tail.TailFile(logFile, tail.Config{Follow: true})

	// There is something seriously wrong since tail waits for files, probably should just panic
	if err != nil {
		log.Panic(err)
	}

	// This loops forever
	for line := range t.Lines {
		err := processLine(line.Text, statusCodes, errorPaths)
		if err != nil {
			log.Println(err)
		}
	}
}

func processLine(line string, statusCodes *map[string]int, errorPaths *map[string]int) error {
	var code, request string

	fields := logFormatRegex.FindStringSubmatch(line)
	if fields == nil {
		return fmt.Errorf("access log line '%v' does not match given format '%v'", line, logFormatRegex)
	}

	// This wont change unless the format changes so I'm happy for the efficientcy gain
	request = fields[9]
	code = grabStatusCodeClass(fields[10])

	(*statusCodes)[code]++
	if code == "50x" {
		(*errorPaths)[grabPathFromRequest(request)]++
	}

	return nil
}

// Grab the status code and make it generic
func grabStatusCodeClass(code string) string {
	bytes := []byte(code)
	bytes[1] = byte('0')
	bytes[2] = byte('x')
	return string(bytes)
}

// Given a request use a regex to load the path
func grabPathFromRequest(request string) string {
	return requestFormatRegex.FindStringSubmatch(request)[1]
}

// Loop through our maps for errorPaths and statusCode counts and output in the proper format
func outputDataToString(statusCodes map[string]int, errorPaths map[string]int) (out string) {
	// Loop through our valid codes
	for i := 5; i > 1; i-- {
		codeClass := fmt.Sprintf("%d0x", i)
		if val, ok := statusCodes[codeClass]; ok {
			out += fmt.Sprintf("%s:%d|s\n", codeClass, val)
		}
		// else {
		// 	out += fmt.Sprintf("%s:%d|s\n", codeClass, 0)
		// }
	}

	// All of that paths that are 50x and their counts
	for path, val := range errorPaths {
		out += fmt.Sprintf("%s:%d|s\n", path, val)
	}

	// Output to STDOUT so kubectl log can find it
	fmt.Print(out)
	return
}

func writeDataToLogfile(outfile string, data string) error {
	f, err := os.OpenFile(outfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}


func requestRegex() *regexp.Regexp {
	re := regexp.MustCompile(`\s(\/.*)\s(?:HTTP)`)
	return re
}

func logRegex() *regexp.Regexp {
	re := regexp.MustCompile(`^(?P<remote_addr>[^ ]*) - (?P<http_x_forwarded_for>[^ ]*) - (?P<http_x_realip>[^ ]*) - \[(?P<time_local>[^]]*)\]( ){1,2}(?P<scheme>[^ ]*) (?P<http_x_forwarded_proto>[^ ]*) (?P<x_forwarded_proto_or_scheme>[^ ]*) "(?P<request>[^"]*)" (?P<status>[^ ]*) (?P<body_bytes_sent>[^ ]*) "(?P<http_referer>[^"]*)" "(?P<http_user_agent>[^"]*)"`)
	return re
}
