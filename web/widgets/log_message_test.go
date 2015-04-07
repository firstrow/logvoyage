package widgets

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBuildLogLine(t *testing.T) {
	Convey("It should build simple log line", t, func() {
		data := make(map[string]interface{})
		data["message"] = "test message"
		r := BuildLogLine(data)
		So(r, ShouldEqual, "test message")
	})
	Convey("It should build simple log and json data", t, func() {
		data := make(map[string]interface{})
		data["message"] = "test"
		data["amount"] = 10
		r := BuildLogLine(data)
		So(r, ShouldEqual, `test {"amount":10}`)
	})
	Convey("It should properly display line with no message", t, func() {
		data := make(map[string]interface{})
		r := BuildLogLine(data)
		So(r, ShouldEqual, "")
	})
	Convey("It should properly display line with no message and json", t, func() {
		data := make(map[string]interface{})
		data["amount"] = 10
		r := BuildLogLine(data)
		So(r, ShouldEqual, ` {"amount":10}`)
	})
	Convey("It should properly display message as integer", t, func() {
		data := make(map[string]interface{})
		data["message"] = 10
		r := BuildLogLine(data)
		So(r, ShouldEqual, `10`)
	})
}
