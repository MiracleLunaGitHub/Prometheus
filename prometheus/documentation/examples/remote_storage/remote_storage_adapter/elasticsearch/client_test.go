package elasticsearch

import (
	"testing"
	"github.com/prometheus/common/model"
	"reflect"
	"fmt"
)

var (
	metric = model.Metric{
		model.MetricNameLabel: "machine_cpu_cores",
		"beta_kubernetes_io_arch": "amd64",
		"beta_kubernetes_io_os": "linux",
		"instance": "minion2",
		"job": "kubernetes-nodes",
		"kubernetes_io_hostname": "master",
	}
)

func TestTagsFromMetric(t *testing.T) {
	expected := map[string]string{
		"beta_kubernetes_io_arch":  "amd64",
		"beta_kubernetes_io_os": "linux",
		"instance": "minion2",
		"job": "kubernetes-nodes",
		"kubernetes_io_hostname": "master",
	}
	actual := tagsFromMetric(metric)
	fmt.Printf("expected: %s\n", expected)
	fmt.Printf("actual: %s\n", actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func TestWrite(t *testing.T) {
	// node cpu usage ratio
	/*samples := model.Samples{
		{
			Metric: model.Metric{
				model.MetricNameLabel: "node_cpu_usage_ratio",
				"kubernetes_io_hostname":   "master",
			},
			Timestamp: model.Time(1503368411461),
			Value:     1.9170004099999953,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "node_cpu_usage_ratio",
				"kubernetes_io_hostname":   "minion1",
			},
			Timestamp: model.Time(1503368411461),
			Value:     1.1616371858333558,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "node_cpu_usage_ratio",
				"kubernetes_io_hostname":   "minion2",
			},
			Timestamp: model.Time(1503368411461),
			Value:     0.8047755066666825,
		},
	}
	*/

	// machine cpu cores
	/*samples := model.Samples{
		{
			Metric: model.Metric{
				model.MetricNameLabel: "machine_cpu_cores",
				"beta_kubernetes_io_arch": "amd64",
				"beta_kubernetes_io_os": "linux",
				"instance": "minion2",
				"job": "kubernetes-nodes",
				"kubernetes_io_hostname": "master",
			},
			Timestamp: model.Time(1503368792),
			Value:     16,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "machine_cpu_cores",
				"beta_kubernetes_io_arch": "amd64",
				"beta_kubernetes_io_os": "linux",
				"instance": "master",
				"job": "kubernetes-nodes",
				"kubernetes_io_hostname": "minion2",
			},
			Timestamp: model.Time(1503368792),
			Value:     8,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "machine_cpu_cores",
				"beta_kubernetes_io_arch": "amd64",
				"beta_kubernetes_io_os": "linux",
				"instance": "minion1",
				"job": "kubernetes-nodes",
				"kubernetes_io_hostname": "minion2",
			},
			Timestamp: model.Time(1503368792),
			Value:     16,
		},
	}*/

	// node memory MemTotal
	samples := model.Samples{
		{
			Metric: model.Metric{
				model.MetricNameLabel: "node_memory_MemTotal",
				"app": "prometheus",
				"component": "node-exporter",
				"instance": "192.168.254.44:9100",
				"job": "kubernetes-service-endpoints",
				"kubernetes_io_hostname": "minion1",
				"kubernetes_name": "prometheus-node-exporter",
				"kubernetes_namespace": "sky-firmament",
			},
			Timestamp: model.Time(1503366851),
			Value:     50467577856,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "node_memory_MemTotal",
				"app": "prometheus",
				"component": "node-exporter",
				"instance": "192.168.254.44:9100",
				"job": "kubernetes-service-endpoints",
				"kubernetes_io_hostname": "master",
				"kubernetes_name": "prometheus-node-exporter",
				"kubernetes_namespace": "sky-firmament",
			},
			Timestamp: model.Time(1503366851),
			Value:     8202465280,
		},
		{
			Metric: model.Metric{
				model.MetricNameLabel: "node_memory_MemTotal",
				"app": "prometheus",
				"component": "node-exporter",
				"instance": "192.168.254.44:9100",
				"job": "kubernetes-service-endpoints",
				"kubernetes_io_hostname": "minion2",
				"kubernetes_name": "prometheus-node-exporter",
				"kubernetes_namespace": "sky-firmament",
			},
			Timestamp: model.Time(1503366851),
			Value:     50467577856,
		},
	}

	c := NewClient("http://192.168.253.161:9200","prometheus","prometheus")

	if err := c.Write(samples); err != nil {
		t.Fatalf("Error sending samples: %s", err)
	}
}
