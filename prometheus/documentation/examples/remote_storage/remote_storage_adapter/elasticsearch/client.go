package elasticsearch

import (
	"fmt"
	"github.com/olivere/elastic"
	"github.com/prometheus/common/model"
	"math"
	"context"
	"github.com/prometheus/common/log"
)

// Client allows sending batches of Prometheus samples to Elasticsearch.
type Client struct {
	url        string
	index      string
	indexType  string
}

// NewClient creates a new Client.
func NewClient(url string, index string, indexType string) *Client {
	return &Client{
		url:        url,
		index:      index,
		indexType:  indexType,
	}
}

// ElasticSearchSamplesRequest is used for building a JSON request for storing samples
// via the Elasticsearch.
type ElasticSearchSamplesRequest struct {
	Metric    string      		`json:"metric"`
	Timestamp int64            	`json:"timestamp"`
	Value     float64           `json:"value"`
	Tags      map[string]string `json:"tags"`
}

// tagsFromMetric translates Prometheus metric into Elasticsearch tags.
func tagsFromMetric(m model.Metric) map[string]string {
	tags := make(map[string]string, len(m)-1)
	for l, v := range m {
		if l == model.MetricNameLabel {
			continue
		}
		tags[string(l)] = string(v)
	}
	return tags
}

// Write sends a batch of samples to Elasticsearch.
func (c *Client) Write(samples model.Samples) error {
	// Create a client and connect to url
	client, err := elastic.NewClient(elastic.SetURL(c.url))
	if err != nil {
		return err
	}

	// elastic: client started
	client.Start()

	// elastic: client stopped
	defer client.Stop()

	// Create a bulk request
	bulkRequest := client.Bulk()

	for _, s := range samples {
		v := float64(s.Value)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			log.Warnf("cannot send value %f to Elasticsearch, skipping sample %#v", v, s)
			continue
		}
		metric := string(s.Metric[model.MetricNameLabel])
		doc := ElasticSearchSamplesRequest{
			Metric:    metric,
			Timestamp: s.Timestamp.Unix(),
			Value:     v,
			Tags:      tagsFromMetric(s.Metric),
		}

		// get each indexRequest
		indexRequest := elastic.NewBulkIndexRequest().Index(c.index).Type(c.indexType).Doc(doc)
		if err != nil {
			return err
		}

		// add each indexRequest to the bulkRequest
		bulkRequest = bulkRequest.Add(indexRequest)
	}

	// Execute the bulkRequest
	bulkResponse, err := bulkRequest.Do(context.TODO())

	if err != nil {
		fmt.Errorf("failed to write the samples to Elasticsearch,", err)
	}
	if bulkResponse == nil {
		fmt.Errorf("expected bulkResponse to be != nil; got nil")
	}

	if bulkRequest.NumberOfActions() != 0 {
		fmt.Errorf("expected bulkRequest.NumberOfActions %d; got %d", 0, bulkRequest.NumberOfActions())
	}

	return nil
}

// Name identifies the client as a Elasticsearch client.
func (c Client) Name() string {
	return "elasticsearch"
}

