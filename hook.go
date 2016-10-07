package graylog

import (
	"bytes"
	"errors"
	"sync"
	"time"

	"github.com/Graylog2/go-gelf/gelf"
	"github.com/Sirupsen/logrus"
)

var (
	sleep = time.Sleep
	now   = time.Now
)

// Hook implements logrus.Hook for sending log to graylog
type Hook struct {
	config Config // only assign once, wont change

	nodes      nodes
	nodesMutex sync.RWMutex
}

// New returns a Hook with config
func New(config Config) *Hook {
	return &Hook{config: config}
}

// SetNodeConfigs set graylog nodes with node configs given
func (h *Hook) SetNodeConfigs(configs ...NodeConfig) error {
	ns := nodes{}
	for _, config := range configs {
		node, err := newNode(config)
		if err != nil {
			return err
		}
		ns = append(ns, node)
	}

	// mutex protect nodes from concurrent read-write access
	h.nodesMutex.Lock()
	h.nodes = ns
	h.nodesMutex.Unlock()

	return nil
}

// Levels implements logrus.Hook interface
func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire implements logrus.Hook interface
func (h *Hook) Fire(entry *logrus.Entry) error {
	// make gelf message
	m := h.makeGELFMessage(entry)

	// select a node to send message
	h.nodesMutex.RLock()
	node := h.nodes.selectNode()
	h.nodesMutex.RUnlock()
	if node == nil {
		return errors.New("fail to select a graylog node for sending log message")
	}

	return node.logWriter.WriteMessage(m)
}

// StartHealthCheck checks if graylog nodes are alive periodically
func (h *Hook) StartHealthCheck() {
	go h.loopCheckNodesHealth(checkNodeStatus)
}

func (h *Hook) makeGELFMessage(entry *logrus.Entry) *gelf.Message {
	// short & full message
	p := bytes.TrimSpace([]byte(entry.Message))
	short := bytes.NewBuffer(p)
	full := ""
	if i := bytes.IndexRune(p, '\n'); i > 0 {
		full = short.String()
		short.Truncate(i)
	}

	// merge entry.Data & StaticMeta & Facility into extra
	extra := map[string]interface{}{}
	for k, v := range entry.Data {
		extra["_"+k] = v // prefix with _ will be treated as an additional field
	}
	for k, v := range h.config.StaticMeta {
		extra["_"+k] = v // prefix with _ will be treated as an additional field
	}
	extra["_facility"] = h.config.Facility

	return &gelf.Message{
		Version:  "1.1",
		Host:     h.config.Hostname,
		Short:    short.String(),
		Full:     full,
		TimeUnix: float64(now().UnixNano()) / 1e9,
		Level:    int32(entry.Level),
		Extra:    extra,
	}
}

func (h *Hook) loopCheckNodesHealth(checkNodeAlive func(url string) uint32) {
	for {
		sleep(h.config.HealthCheckInterval)

		// mutex protect nodes from concurrent read-write access
		func() {
			h.nodesMutex.RLock()
			defer h.nodesMutex.RUnlock()
			for _, node := range h.nodes {
				node.setStatus(checkNodeAlive(node.healthCheckURL))
			}
		}()
	}
}
