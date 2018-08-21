package api

import (
	//"cloudkarafka-mgmt/kafka"
	"cloudkarafka-mgmt/dm"
	"cloudkarafka-mgmt/store"
	"cloudkarafka-mgmt/zookeeper"

	"github.com/gorilla/mux"

	"net/http"
)

type consumerVM struct {
	Name         string            `json:"name"`
	Topics       []consumedTopicVM `json:"topics"`
	PartitionLag []partitionLag    `json:"partition_lag"`
}

type consumedTopicVM struct {
	Name     string `json:"name"`
	Coverage int    `json:"coverage"`
}

type consumerMetric struct {
	Topics map[string]partitionLag `json:"topics"`
}

type partitionLag struct {
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Lag       int    `json:"lag"`
}

func Consumers(w http.ResponseWriter, r *http.Request, p zookeeper.Permissions) {
	allConsumers := store.Consumers()
	var myConsumers []string
	for _, c := range allConsumers {
		if p.GroupRead(c) {
			myConsumers = append(myConsumers, c)
		}
	}
	writeJson(w, myConsumers)
}

func Consumer(w http.ResponseWriter, r *http.Request, p zookeeper.Permissions) {
	vars := mux.Vars(r)
	var data interface{}
	topicName := vars["name"]
	if p.TopicRead(topicName) {
		data = dm.ConsumerMetrics(topicName)
	}
	writeJson(w, data)
}
