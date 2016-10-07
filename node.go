package graylog

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync/atomic"

	"github.com/Graylog2/go-gelf/gelf"
)

var randIntn = rand.Intn

// NOTE: do not change the value of node status
// if you do want the modification happen
// please also update default status in NewNode to Alive
const (
	nodeStatusAlive = 0
	nodeStatusDead  = 1
)

type node struct {
	// health check
	healthCheckURL string
	status         *uint32

	// log write
	udpAddress string
	logWriter  *gelf.Writer

	// weight of node used for node selection
	weight int
}

func newNode(config NodeConfig) (*node, error) {
	writer, err := gelf.NewWriter(config.UDPAddress)
	if err != nil {
		return nil, err
	}

	return &node{
		healthCheckURL: config.HealthCheckURL,
		status:         new(uint32), // default status to Alive
		udpAddress:     config.UDPAddress,
		weight:         config.Weight,
		logWriter:      writer,
	}, nil
}

func (n node) alive() bool {
	return atomic.LoadUint32(n.status) == nodeStatusAlive
}

func (n *node) setStatus(status uint32) {
	switch status {
	case nodeStatusAlive, nodeStatusDead:
		atomic.StoreUint32(n.status, status)
	default:
		errorLogger.Panicf("unknown node status %v", status)
	}
}

func (n node) getWeight() int {
	if !n.alive() {
		return 0
	}
	return n.weight
}

func (n *node) String() string {
	return fmt.Sprintf("udp_address=%v weight=%v alive=%v health_check_url=%v\n", n.udpAddress, n.weight, n.alive(), n.healthCheckURL)
}

type nodes []*node

// selectNode using weighted-random algorithm
// http://eli.thegreenplace.net/2010/01/22/weighted-random-generation-in-python/
func (ns nodes) selectNode() *node {
	total := ns.totalWeight()
	if total <= 0 {
		errorLogger.Printf("total weight of nodes is %v\n", total)
		return nil
	}

	var selectedNode *node
	seed := randIntn(total)
	for _, node := range ns {
		seed -= node.getWeight()
		if seed < 0 {
			selectedNode = node
			break
		}
	}

	return selectedNode
}

func (ns nodes) totalWeight() int {
	sum := 0
	for _, node := range ns {
		sum += node.getWeight()
	}
	return sum
}

func (ns nodes) String() string {
	buf := bytes.NewBuffer(nil)
	for _, node := range ns {
		buf.WriteString(node.String())
	}
	return buf.String()
}
