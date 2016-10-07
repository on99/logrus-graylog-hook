package graylog

import "testing"

func TestSetNodeConfigs(t *testing.T) {
	hook := &Hook{}
	if err := hook.SetNodeConfigs(NodeConfig{UDPAddress: "invalid-udp-address"}); err == nil {
		t.Error("expected non-nil error")
	}
}
