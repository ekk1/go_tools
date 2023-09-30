package utils

import "testing"

func TestHTTPClient(t *testing.T) {
	c := NewHTTPClient()
	ret, err := c.SendReq("GET", "https://httpbin.org/get", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret.Text())
}
