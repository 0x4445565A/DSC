package main

import (
	//"fmt"
	"testing"
)

func TestLogRegex(t *testing.T) {
	re := logRegex()

	logExamples := []string{
		`10.10.180.161 - 72.34.110.66, 192.33.28.238 - - - [03/Aug/2015:15:50:06 +0000]  https https https "GET / HTTP/1.1" 200 20027 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.107 Safari/537.36"`,
		`10.10.180.161 - - - - - [03/Aug/2015:15:50:06 +0000]  https https https "GET /somePage?foo=bar HTTP/1.1" 400 20027 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.107 Safari/537.36"`,
		`172.17.0.1 - - - - - [30/Jul/2017:22:55:27 +0000]  http - - "GET /500 HTTP/1.1" 500 13 "-" "curl/7.52.1"`,
	}

	logCodeResults := []string{
		`20x`,
		`40x`,
		`50x`,
	}

	logPathResults := []string{
		`/`,
		`/somePage?foo=bar`,
		`/500`,
	}

	for i, logExample := range logExamples {

		fields := re.FindStringSubmatch(logExample)
		if fields == nil {
			t.Error("Regex didn't satisfy log example")
			return
		}

		if ret := grabStatusCodeClass(fields[12]); ret != logCodeResults[i] {
			t.Errorf("returned incorrect code: expect %s got %s", logCodeResults[i], ret)
		}

		if ret := grabPathFromRequest(fields[11]); ret != logPathResults[i] {
			t.Errorf("returned incorrect path: expect %s got %s", logPathResults[i], ret)
		}
	}
}
