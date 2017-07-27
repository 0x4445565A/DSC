package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"regexp"
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
	logFile := "./access.log"

	t, err := tail.TailFile(logFile, tail.Config{Follow: true})
	if err != nil {
		panic(err)
	}

	errorPaths := map[string]int{}
	statusCodes := map[string]int{}

	// This loops forever
	for line := range t.Lines {
		fields := logFormatRegex.FindStringSubmatch(line.Text)
		if fields == nil {
			panic(fmt.Errorf("access log line '%v' does not match given format '%v'", line, logFormatRegex))
		}

		logEntry := map[string]string{}

		for i, name := range logFormatRegex.SubexpNames() {
			logEntry[name] = fields[i]
		}

		code, ok := logEntry["status"]
		if !ok {
			return
		}

		code = grabStatusCodeClass(code)

		statusCodes[code]++
		if code == "50x" {
			req, ok := logEntry["request"]
			if !ok {
				return
			}
			errorPaths[grabPathFromRequest(req)]++
		}

		fmt.Println(statusCodes, errorPaths)
	}
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

func requestRegex() *regexp.Regexp {
	re := regexp.MustCompile(`\s(\/.*)\s(?:HTTP)`)
	return re
}

func logRegex() *regexp.Regexp {
	re := regexp.MustCompile(`^(?P<remote_addr>[^ ]*) - (?P<http_x_forwarded_for>[^-]*) - (?P<http_x_realip>[^ ]*) - \[(?P<time_local>[^]]*)\] (?P<scheme>[^ ]*) (?P<http_x_forwarded_proto>[^ ]*) (?P<x_forwarded_proto_or_scheme>[^ ]*) (?P<x_forwarded_proto_or_scheme>[^ ]*) "(?P<request>[^"]*)" (?P<status>[^ ]*) (?P<body_bytes_sent>[^ ]*) "(?P<http_referer>[^"]*)" "(?P<http_user_agent>[^"]*)"`)
	return re
}
