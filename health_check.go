package graylog

import "net/http"

// issue a http get request to graylog-server
// to check if the graylog node is alive or not
// http://docs.graylog.org/en/2.1/pages/configuration/load_balancers.html#load-balancer-state
func checkNodeStatus(url string) uint32 {
	resp, err := http.Get(url)
	if err != nil {
		errorLogger.Printf("check graylog node %v status error: %v\n", url, err)
		return nodeStatusDead
	}

	defer resp.Body.Close()

	if statusCode := resp.StatusCode; statusCode != http.StatusOK {
		errorLogger.Printf("graylog node %v dead with status code %v returned\n", url, statusCode)
		return nodeStatusDead
	}

	debugLogger.Printf("graylog node %v alive\n", url)
	return nodeStatusAlive
}
