package bencode

import (
	"go_utils/utils"
	"os"
	"testing"
)

func TestBencode(t *testing.T) {
	fileData, err := os.ReadFile("debian-12.0.0-amd64-DVD-1.iso.torrent")
	if err != nil {
		utils.LogPrintError(err)
		t.FailNow()
	}

	d := &BencodeDecoder{}

	ret, err := d.Decode(fileData)
	rr, ok := ret.(map[string]interface{})
	if !ok {
		t.FailNow()
	}
	for k, v := range rr {
		t.Log(k)
		// t.Log(v)
		switch k {
		case "url-list":
			// t.Log(v)
			vv, ok := v.([]interface{})
			if !ok {
				t.Log("Type wrong")
				t.FailNow()
			}
			for _, url := range vv {
				urll, ok := url.([]byte)
				if !ok {
					t.Log("Type wrong")
					t.FailNow()
				}
				t.Log(string(urll))
			}
		case "comment":
			// t.Log(v)
			vv, ok := v.([]byte)
			if !ok {
				t.Log("Type wrong")
				t.FailNow()
			}
			t.Log(string(vv))
		case "announce":
			// t.Log(v)
			vv, ok := v.([]byte)
			if !ok {
				t.Log("Type wrong")
				t.FailNow()
			}
			t.Log(string(vv))
		case "announce-list":
			// t.Log(v)
			vv, ok := v.([]interface{})
			if !ok {
				t.Log("Type wrong")
				t.FailNow()
			}
			for _, v := range vv {
				// t.Log(v)
				vv, ok := v.([]interface{})
				if !ok {
					t.Log("Type wrong")
					t.FailNow()
				}
				for _, v := range vv {
					vv, ok := v.([]byte)
					if !ok {
						t.Log("Type wrong")
						t.FailNow()
					}
					t.Log(string(vv))
				}

			}
		case "created by":
			// t.Log(v)
			vv, ok := v.([]byte)
			if !ok {
				t.Log("Type wrong")
				t.FailNow()
			}
			t.Log(string(vv))

		case "creation date":
			// t.Log(v)
			vv, ok := v.(int64)
			if !ok {
				t.Log("Type wrong")
				t.FailNow()
			}
			t.Log(vv)

		case "info":
			// t.Log(v)
			vv, ok := v.(map[string]interface{})
			if !ok {
				t.Log("Type wrong")
				t.FailNow()
			}
			for k, v := range vv {
				switch k {
				case "length":
					t.Log(v)
				case "name":
					// t.Log(v)
					vv, _ := v.([]byte)
					t.Log(string(vv))
				case "piece length":
					t.Log(v)
				case "pieces":
					// t.Log(v)
					vv, ok := v.([]byte)
					if ok {
						t.Log("is is")
						t.Log(len(vv))
					}
				case "files":
					// t.Log(v)
					vv, ok := v.([]interface{})
					if !ok {
						t.Log("nononon")
						// t.Log(len(vv))
					}
					for _, v := range vv {
						// t.Log(v)
						vv, ok := v.(map[string]interface{})
						if !ok {
							t.Log("nononon")
							// t.Log(len(vv))
						}
						for k, v := range vv {
							t.Log(k)
							switch vv := v.(type) {
							case int64:
								t.Log(vv)
							case []interface{}:
								for _, v := range vv {
									vv, ok := v.([]byte)
									if !ok {
										t.Log("nononon")
										// t.Log(len(vv))
									}
									t.Log(string(vv))
								}
							default:
								t.Log("Unknown type")
							}
						}
					}

				default:
					t.Log("Ignoring: ", k)
				}
			}

		default:
			t.Log("Ignoring: ", k)
			// t.Log("Value: ", v)
		}
	}
	// t.Log(ret, err)

}
