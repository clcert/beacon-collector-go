package collectors

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type EthereumCollector struct{}

func (e EthereumCollector) sourceName() string {
	return "ethereum"
}

func (e EthereumCollector) collectEvent() (string, string) {
	ethAPI := "https://eth.labs.clcert.cl"
	jsonStr := []byte(`{"jsonrpc": "2.0", "method": "eth_getBlockByNumber", "id": "1", "params": ["latest", false]}`)
	resp, err := http.Post(ethAPI, "application/json", bytes.NewReader(jsonStr))

	if err != nil {
		log.Error("Failed to get Ethereum event")
		return "0", "0"
	}

	if resp.StatusCode != 200 {
		return "0", "0"
	}

	body := resp.Body
	defer body.Close()

	response, _ := ioutil.ReadAll(body)
	blockInfo := make(map[string]map[string]string)
	_ = json.Unmarshal(response, &blockInfo)
	if _, ok := blockInfo["error"]; ok {
		log.Error("Ethereum response with error")
		log.Error(blockInfo["error"])
		return "0", "0"
	} else {
		lastBlockHash := blockInfo["result"]["hash"][2:]
		lastBlockNumber := blockInfo["result"]["number"][2:]
		return lastBlockHash, lastBlockNumber
	}
}

func (e EthereumCollector) estimateEntropy() int {
	return 0
}

func (e EthereumCollector) processForDigest(s string) string {
	return s
}
