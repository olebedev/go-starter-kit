package http2curl

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ExampleGetCurlCommand() {
	req, _ := http.NewRequest("PUT", "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu", bytes.NewBufferString(`{"hello":"world","answer":42}`))
	req.Header.Set("Content-Type", "application/json")

	command, _ := GetCurlCommand(req)
	fmt.Println(command)

	// Output:
	// curl -X PUT -d "{\"hello\":\"world\",\"answer\":42}" -H "Content-Type: application/json" http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu
}

func TestGetCurlCommand(t *testing.T) {
	Convey("Testing http2curl", t, func() {
		uri := "http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu"
		payload := new(bytes.Buffer)
		payload.Write([]byte(`{"hello":"world","answer":42}`))
		req, err := http.NewRequest("PUT", uri, payload)
		So(err, ShouldBeNil)
		req.Header.Set("X-Auth-Token", "private-token")
		req.Header.Set("Content-Type", "application/json")

		command, err := GetCurlCommand(req)
		So(err, ShouldBeNil)
		expected := `curl -X PUT -d "{\"hello\":\"world\",\"answer\":42}" -H "Content-Type: application/json" -H "X-Auth-Token: private-token" http://www.example.com/abc/def.ghi?jlk=mno&pqr=stu`
		So(command.String(), ShouldEqual, expected)
	})
}
