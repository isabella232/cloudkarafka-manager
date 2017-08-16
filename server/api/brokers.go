package api

import (
	"cloudkarafka-mgmt/jmx"
	"cloudkarafka-mgmt/zookeeper"
	"github.com/gorilla/mux"

	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type broker struct {
	KafkaVersion     string   `json:"kafka_version"`
	Version          int      `json:"-"`
	JmxPort          int      `json:"jmx_port"`
	Timestamp        string   `json:"timestamp"`
	Uptime           string   `json:"uptime"`
	Endpoints        []string `json:"endpoints"`
	Host             string   `json:"host"`
	Port             int      `json:"port"`
	Id               string   `json:"id"`
	BytesInPerSec    float64  `json:"bytes_in_per_sec"`
	BytesOutPerSec   float64  `json:"bytes_out_per_sec"`
	MessagesInPerSec float64  `json:"messages_in_per_sec"`
}

func Brokers(w http.ResponseWriter, r *http.Request) {
	brokers := zookeeper.Brokers()
	writeJson(w, brokers)
}

func Broker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	var b broker
	json.Unmarshal(zookeeper.Broker(vars["id"]), &b)
	ts, err := strconv.ParseInt(b.Timestamp, 10, 64)
	if err != nil {
		internalError(w, err)
	}
	t := time.Unix(ts/1000, 0)
	b.Uptime = time.Since(t).String()
	b.Id = vars["id"]
	b.KafkaVersion = jmx.KafkaVersion(b.Id)
	b.BytesOutPerSec = jmx.BrokerTopicMetric("BytesOutPerSec", "")
	b.BytesInPerSec = jmx.BrokerTopicMetric("BytesInPerSec", "")
	b.MessagesInPerSec = jmx.BrokerTopicMetric("MessagesInPerSec", "")
	writeJson(w, b)
}
