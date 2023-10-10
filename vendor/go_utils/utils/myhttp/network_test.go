package myhttp

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
	t.Log("Testing GET")
	if ret, err := c.SendReq("GET", "http://127.0.0.1:5000/get", nil); err == nil {
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}

	// Test POST form
	t.Log("Testing POST Form")
	body := url.Values{}
	body.Set("test", "test")
	c.SetSendForm(true)
	if ret, err := c.SendReq("POST", "http://127.0.0.1:5000/post", body); err == nil {
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}

	// Test POST JSON
	t.Log("Testing POST JSON")
	jsonBody := map[string]string{
		"test": "test",
	}
	recv := map[string]interface{}{}
	c.SetSendForm(false)
	c.SetSendJSON(true)
	if ret, err := c.SendReq("POST", "http://127.0.0.1:5000/post", jsonBody); err == nil {
		t.Log(ret.JSON(&recv))
		t.Log(recv)
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}

	// Test Raw Body
	t.Log("Testing POST RAW")
	c.SetSendJSON(false)
	c.SetSendRawBody(true)
	bodyBytes := []byte("test")
	if ret, err := c.SendReq("POST", "http://127.0.0.1:5000/post", bodyBytes); err == nil {
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}

	// Test basic auth
	t.Log("Testing POST JSON with Auth")
	c.SetBasicAuth("user", "pass")
	c.SetSendRawBody(false)
	c.SetSendJSON(true)
	if ret, err := c.SendReq("POST", "http://127.0.0.1:5000/post", jsonBody); err == nil {
		t.Log(ret.Text())
	} else {
		t.Fatal(err)
	}
}
