package graylog

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/Sirupsen/logrus"
)

func TestHook(t *testing.T) {
	reader, _ := gelf.NewReader("127.0.0.1:9999")
	server := httptest.NewServer(handlerWithStatusCode(http.StatusServiceUnavailable))

	config := NewConfig()
	config.HealthCheckInterval = time.Second
	config.StaticMeta = map[string]interface{}{"go_version": "1.7.1"}
	config.Facility = "graylog_hook_test"
	hook := New(config)
	hook.SetNodeConfigs(
		NodeConfig{
			UDPAddress:     reader.Addr(),
			HealthCheckURL: server.URL,
			Weight:         1,
		},
		NodeConfig{
			UDPAddress:     reader.Addr(),
			HealthCheckURL: server.URL,
			Weight:         0,
		},
	)

	log := logrus.New()
	log.Out = ioutil.Discard
	log.Hooks.Add(hook)
	log.Level = logrus.InfoLevel

	// all graylog nodes are alive by default
	//
	msg := "this is a message\nsomething after new line"
	short := "this is a message"
	log.WithFields(logrus.Fields{
		"string":  "value",
		"number":  10.1,
		"boolean": true, // NOTE: no boolean in GELF, only string and number
	}).Info(msg)

	gelfMessage, _ := reader.ReadMessage()
	expectedExtra := map[string]interface{}{
		"_go_version": "1.7.1",
		"_facility":   "graylog_hook_test",
		"_string":     "value",
		"_number":     10.1,
		"_boolean":    "true", // NOTE: no boolean in GELF, only string and number
	}
	if reflect.DeepEqual(gelfMessage.Extra, expectedExtra) {
		t.Errorf("expected %v but get %v", expectedExtra, gelfMessage.Extra)
	}
	if gelfMessage.Full != msg {
		t.Errorf("expected %v but get %v", msg, gelfMessage.Full)
	}
	if gelfMessage.Short != short {
		t.Errorf("expected %v but get %v", short, gelfMessage.Short)
	}

	// after 1 second all graylog nodes are dead
	hook.StartHealthCheck()
	time.Sleep(2 * time.Second)

	if node := hook.nodes.selectNode(); node != nil {
		t.Errorf("expected not able to select a node with all nodes dead but still get %v", node)
	}

	log.WithField("key", "value").Error("cannot send log")
}
