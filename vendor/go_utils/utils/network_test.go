package utils

import (
	"net/url"
	"testing"
)

/*
   Set up httpbin env
   apt install python3-httpbin
   python3 -m httpbin.core
*/

func TestHTTPClient(t *testing.T) {
	c := NewHTTPClient()
	// Test GET
	if ret, err := c.SendReq("GET", "http://127.0.0.1:5000/get", nil); err == nil {
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}
	// Test POST form
	body := url.Values{}
	body.Set("test", "test")
	c.SetSendForm(true)
	if ret, err := c.SendReq("POST", "http://127.0.0.1:5000/post", body); err == nil {
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}
	// Test POST JSON
	jsonBody := map[string]string{
		"test": "test",
	}
	c.SetSendForm(false)
	c.SetSendJSON(true)
	if ret, err := c.SendReq("POST", "http://127.0.0.1:5000/post", jsonBody); err == nil {
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}
	// Test basic auth
	c.SetBasicAuth("user", "pass")
	if ret, err := c.SendReq("POST", "http://127.0.0.1:5000/post", jsonBody); err == nil {
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}
}
