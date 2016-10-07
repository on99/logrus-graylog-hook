package graylog

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewNode(t *testing.T) {
	Convey("Given node config with invalid udp address", t, func() {
		config := NodeConfig{UDPAddress: "invalid-udp-address"}

		Convey("When create new node", func() {
			_, err := newNode(config)

			Convey("Error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestSetStatus(t *testing.T) {
	Convey("Given a graylog node", t, func() {
		node := &node{status: new(uint32)}

		Convey("When set invalid status", func() {
			f := func() {
				var status uint32 = 1 << 31
				node.setStatus(status)
			}

			Convey("It should panic", func() {
				So(f, ShouldPanic)
			})
		})

		Convey("When set alive status", func() {
			node.setStatus(nodeStatusAlive)
			status := *node.status

			Convey("status should be alive", func() {
				So(status, ShouldEqual, nodeStatusAlive)
			})
		})

		Convey("When set dead status", func() {
			node.setStatus(nodeStatusDead)
			status := *node.status

			Convey("status should be alive", func() {
				So(status, ShouldEqual, nodeStatusDead)
			})
		})
	})
}

func TestGetWeight(t *testing.T) {
	Convey("Given a dead graylog node", t, func() {
		status := new(uint32)
		*status = nodeStatusDead
		node := &node{status: status, weight: 999}

		Convey("When get weight", func() {
			weight := node.getWeight()

			Convey("Weight should be equal to 0", func() {
				So(weight, ShouldEqual, 0)
			})
		})
	})

	Convey("Given an alive graylog node", t, func() {
		status := new(uint32)
		*status = nodeStatusAlive
		node := &node{status: status, weight: 999}

		Convey("When get weight", func() {
			weight := node.getWeight()

			Convey("Weight should be equal to 999", func() {
				So(weight, ShouldEqual, 999)
			})
		})
	})
}

func TestNodesString(t *testing.T) {
	ns := nodes{}
	statusAlive := new(uint32)
	statusDead := new(uint32)
	*statusDead = nodeStatusDead
	ns = append(ns,
		&node{
			udpAddress:     "127.0.0.1:9998",
			healthCheckURL: "url1",
			status:         statusAlive,
			weight:         1,
		},
		&node{
			udpAddress:     "127.0.0.1:9999",
			healthCheckURL: "url2",
			status:         statusDead,
			weight:         2,
		},
	)

	actual := fmt.Sprint(ns)
	expected := fmt.Sprintf(
		"%v\n%v\n",
		"udp_address=127.0.0.1:9998 weight=1 alive=true health_check_url=url1",
		"udp_address=127.0.0.1:9999 weight=2 alive=false health_check_url=url2",
	)
	if actual != expected {
		t.Errorf("expected %v but get %v", expected, actual)
	}
}
