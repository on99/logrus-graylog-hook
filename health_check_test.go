package graylog

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckNodeStatus(t *testing.T) {
	if status := checkNodeStatus("invalid-url"); status != nodeStatusDead {
		t.Errorf("status should be %v but get %v", nodeStatusDead, status)
	}

	Convey("Given a server response with 503", t, func() {
		server := httptest.NewServer(handlerWithStatusCode(http.StatusServiceUnavailable))

		Convey("When check node status", func() {
			status := checkNodeStatus(server.URL)

			Convey("Status should be dead", func() {
				So(status, ShouldEqual, nodeStatusDead)
			})
		})

		Reset(func() {
			server.Close()
		})
	})

	Convey("Given a server response with 200", t, func() {
		server := httptest.NewServer(handlerWithStatusCode(http.StatusOK))

		Convey("When check node status", func() {
			status := checkNodeStatus(server.URL)

			Convey("Status should be alive", func() {
				So(status, ShouldEqual, nodeStatusAlive)
			})
		})

		Reset(func() {
			server.Close()
		})
	})
}
