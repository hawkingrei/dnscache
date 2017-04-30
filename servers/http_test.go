/* http_test.go
 *
 * Copyright (C) 2016 Alexandre ACEBEDO
 *
 * This software may be modified and distributed under the terms
 * of the MIT license.  See the LICENSE file for details.
 */

package servers

import (
	"encoding/json"
	"github.com/hawkingrei/g53/utils"
	"github.com/hawkingrei/g53/version"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestServiceRequests(t *testing.T) {
	const TestAddr = "127.0.0.1:9981"

	config := utils.NewConfig()
	config.HttpAddr = TestAddr

	server := NewHTTPServer(config, NewDNSServer(config))
	go server.Start()

	// Allow some time for server to start
	time.Sleep(250 * time.Millisecond)

	version, _ := json.Marshal(version.VersionOptions{
		GitCommit: version.GitCommit,
		Version:   version.Version,
		BuildTime: version.BuildTime,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	})
	var tests = []struct {
		method, url, body, expected string
		status                      int
	}{
		{"GET", "/version", "", string(version), 200},
		{"GET", "/services", "", "[]", 200},
		{"GET", "/service", `{"Aliaseo"%！@#！@#！@#！@#！@#！@#}`, "", 500},
		{"PUT", "/service", `{"Aliaseo"%！@#！@#！@#！@#！@#！@#}`, "", 500},
		{"GET", "/service", `{"RecordType":"A","Value":"127.0.0.1","TTL":0,"Aliases":"foo.duitang.com."}`, "", 404},
		{"DELETE", "/service", `{"RecordType":"A","Value":"127.0.0.10","Aliases":"foo.duitang.com."}`, "", 400},
		{"GET", "/services/foo", "", "", 404},
		{"PUT", "/service", `{"RecordType":"A","Value":"127.0.0.1","TTL":3600,"Aliases":"foo.duitang.com."}`, "", 200},
		{"PUT", "/service", `{"RecordType":"A","Value":"127.0.0","TTL":3600,"Aliases":"foo.duitang.com."}`, "", 500},
		{"PUT", "/service", `{"RecordType":"A","Value":"127.0.0.1","TTL":3600,"Aliases":""}`, "", 500},
		{"PUT", "/service", `{"RecordType":"A","Value":"127.0.0.1","TTL":0,"Aliases":"foo.duitang.com."}`, "", 500},
		{"PUT", "/service", `{"RecordType":"A","Value":"127.0.0.1","TTL":,"Aliases":"foo.duitang.com."}`, "", 500},
		{"GET", "/service", `{"RecordType":"A","Value":"127.0.0.1","TTL":0,"Aliases":"foo.duitang.com."}`, `[{"RecordType":"A","Value":"127.0.0.1","TTL":3600,"Aliases":"foo.duitang.com."}]`, 200},
		//{"PATCH", "/service", `{"originalValue":{"RecordType":"A","Value":"127.0.0.1","TTL":3600,"Aliases":"foo.duitang.com."},"modifyValue":abc}`, ``, 500},
		//{"PATCH", "/service", `{"originalValue":abc,"modifyValue":{"RecordType":"A","Value":"127.0.0.10","TTL":3600,"Aliases":"foo.duitang.com."}}`, ``, 500},
		//{"PATCH", "/service", `{"originalValue":{"RecordType":"A","Value":"127.0.0.1","TTL":3600,"Aliases":"foo.duitang.com."},"modifyValue":{"RecordType":"A","Value":"127.0.0.10","TTL":3600,"Aliases":"foo.duitang.com."}}`, ``, 200},
		{"GET", "/service", `{"RecordType":"A","Aliases":"foo.duitang.com."}`, `[{"RecordType":"A","Value":"127.0.0.1","TTL":3600,"Aliases":"foo.duitang.com."}]`, 200},
		{"DELETE", "/service", `{"RecordType":"A","Value":"127.0.0.1","Aliases":"foo.duitang.com."}`, "", 200},
		{"PUT", "/service", `{"RecordType":"MX","Value":"www.google.com.","TTL":3600,"Aliases":"www.aws.com."}`, "", 500},
		//{"PUT", "/service", `{"RecordType":"CNAME","Value":"10.0.0.0","TTL":3600,"Aliases":"www.aws.com"}`, "", 500},
		{"PUT", "/service", `{"RecordType":"CNAME","Value":"www.google.com","TTL":3600,"Aliases":"www.aws.com"}`, "", 200},
		//{"PUT", "/service", `{"RecordType":"CNAME","Value":"www.baidu.com.","TTL":3600,"Aliases":"www.aws.com"}`, "", 500},
		//{"PATCH", "/services/abc.duitang.com.", `{"RecordType":"A","Value":"127.0.0.1","TTL":3600,"Aliases":"foo.duitang.com."}`, ``, 400},
		//{"PUT", "/services", `{"RecordType":"A","Value":"127.0.0.2","TTL":3600,"Aliases":"boo.duitang.com."}`, "", 200},

		//{"PATCH", "/services/boo.duitang.com.", `{"RecordType":"A","Value":"127.0.0.3","TTL":3600,"Aliases":"boo.duitang.com."}`, "", 200},
		//{"PATCH", "/services/boo.duitang.com.", `{"RecordType":"MX","Value":"127.0.0.3","TTL":3600,"Aliases":"boo.duitang.com."}`, "", 500},
		//{"PATCH", "/services/boo.duitang.com.", `{"RecordType":"MX","Value":"127.0.0.3","TTL":"ASDF","Aliases":"boo.duitang.com."}`, "", 500},
		//{"PATCH", "/services/boo.duitang.com.", `{"RecordType":"A","Value":"127.0.0.3","TTL":3600,"Aliases":"boo.duitang.com."}`, "", 200},
		//{"PATCH", "/services/boo.duitang.com.", `{"RecordType":"A","Value":"127.0.0.3","TTL":"3600e123","Aliases":"boo.duitang.com."}`, "", 500},
		//{"PATCH", "/services/boo.duitang.com.", `{"RecordType":"A","Value":"127.0.0.3","TTL":3600e123,"Aliases":"boo.duitang.com."}`, "", 500},
		//{"PATCH", "/services/boo.duitang.com.", `{"RecordType":"A","Value":"127.0.0.3","T":3600e123,"Aliases":"boo.duitang.com."}`, "", 200},
		//{"PATCH", "/services/boo.duitang.com.", `{"RecordType":"A","Value":"","TTL":3600.123,"Aliases":"boo.duitang.com."}`, "", 500},
		//{"GET", "/services", "", `{"boo.duitang.com.":{"RecordType":"A","Value":"127.0.0.3","TTL":3600,"Private":false,"Aliases":"boo.duitang.com."},"foo.duitang.com.":{"RecordType":"A","Value":"127.0.0.1","TTL":3600,"Private":false,"Aliases":"foo.duitang.com."}}`, 200},
		//{"PUT", "/services", `{"RecordType":"CNAME","Value":"www.google.com.","TTL":3600,"Aliases":"boo.duitang.com."}`, "", 500},
		//{"PUT", "/services", `{"RecordType":"CNAME","Value":"www.google.com","TTL":3600,"Aliases":"boo.duitang.com."}`, "", 200},
		//{"DELETE", "/services/foo.duitang.com.", ``, "", 200},
		//{"DELETE", "/services/foo", ``, "", 400},
		{"PUT", "/set/ttl", `AB`, "", 500},
	}

	for _, input := range tests {
		t.Log(input.method, input.url)
		req, err := http.NewRequest(input.method, "http://"+TestAddr+input.url, strings.NewReader(input.body))
		if err != nil {
			t.Error(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
		}
		defer resp.Body.Close()

		if input.status != resp.StatusCode {
			t.Log(input, "Expected status:", input.status, "Got:", resp.StatusCode)
			t.Error(resp.Body)
			return
		}

		if input.status != 200 {
			continue
		}

		actual, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		actualStr := strings.Trim(string(actual), " \n")
		if actualStr != input.expected {
			t.Error(input, "Expected:", input.expected, "Got:", actualStr)
		}
	}

	t.Log("Test TTL setter")
	if config.Ttl != 0 {
		t.Error("Default TTL is not 0")
	}
	req, err := http.NewRequest("PUT", "http://"+TestAddr+"/set/ttl", strings.NewReader("12"))
	if err != nil {
		t.Error(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if config.Ttl != 12 {
		t.Error("TTL not updated. Expected: 12 Got:", config.Ttl)
	}
}
