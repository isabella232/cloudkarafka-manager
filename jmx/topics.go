package jmx

import (
	"cloudkarafka-mgmt/zookeeper"

	"fmt"
)

func TopicMetrics(t string) (TopicMetric, error) {
	var tm TopicMetric
	bi, _ := BrokerTopicMetric("BytesInPerSec", t)
	bo, _ := BrokerTopicMetric("BytesOutPerSec", t)
	mi, _ := BrokerTopicMetric("MessagesInPerSec", t)
	topic, _ := zookeeper.Topic(t)
	var partitions []string
	for p, _ := range topic.Partitions {
		partitions = append(partitions, p)
	}
	tm = TopicMetric{
		TransferMetric: TransferMetric{
			BytesInPerSec:    bi,
			BytesOutPerSec:   bo,
			MessagesInPerSec: mi,
		},
		MessageCount: TopicMessageCount(t, partitions),
	}
	return tm, nil
}

func LogOffsetMetric(t, p string) (OffsetMetric, error) {
	var om OffsetMetric
	so, err := LogOffset("LogStartOffset", t, p)
	if err != nil {
		return om, nil
	}
	eo, err := LogOffset("LogEndOffset", t, p)
	if err != nil {
		return om, nil
	}
	om = OffsetMetric{
		LogStartOffset: so,
		LogEndOffset:   eo,
	}
	return om, nil
}

func TopicMessageCount(topic string, partitions []string) int {
	msgs := 0
	for _, p := range partitions {
		s, err := LogOffset("LogStartOffset", topic, p)
		if err != nil {
			fmt.Println(err)
		}
		e, err := LogOffset("LogEndOffset", topic, p)
		if err != nil {
			fmt.Println(err)
		}
		msgs += e - s
	}
	return msgs
}
