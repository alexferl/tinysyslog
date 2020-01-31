package sinks

import (
	"context"
	"crypto/tls"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

// ElasticsearchSink represents an Elasticsearch sink
type ElasticsearchSink struct {
	client             *elastic.Client
	ctx                context.Context
	Address            string
	IndexName          string
	Username           string
	Password           string
	InsecureSkipVerify bool
	Sniff              bool
}

// Backoff code taken from https://github.com/olivere/elastic/blob/release-branch.v6/backoff.go
// and modified to retry forever.

// ExponentialBackoff implements the simple exponential backoff described by
// Douglas Thain at http://dthain.blogspot.de/2009/02/exponential-backoff-in-distributed.html.
type ExponentialBackoff struct {
	t float64 // initial timeout (in msec)
	f float64 // exponential factor (e.g. 2)
	m float64 // maximum timeout (in msec)
}

// NewExponentialBackoff returns a ExponentialBackoff backoff policy.
// Use initialTimeout to set the first/minimal interval
// and maxTimeout to set the maximum wait interval.
func NewExponentialBackoff(initialTimeout, maxTimeout time.Duration) *ExponentialBackoff {
	return &ExponentialBackoff{
		t: float64(int64(initialTimeout / time.Millisecond)),
		f: 2.0,
		m: float64(int64(maxTimeout / time.Millisecond)),
	}
}

// Next implements BackoffFunc for ExponentialBackoff.
func (b *ExponentialBackoff) Next(retry int) (time.Duration, bool) {
	r := 1.0 + rand.Float64() // random number in [1..2]
	m := math.Min(r*b.t*math.Pow(b.f, float64(retry)), b.m)
	d := time.Duration(int64(m)) * time.Millisecond
	return d, true
}

type Retrier struct {
	backoff elastic.Backoff
}

func NewRetrier() *Retrier {
	return &Retrier{
		backoff: NewExponentialBackoff(500*time.Millisecond, 10*time.Second),
	}
}

func (r *Retrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	// Let the backoff strategy decide how long to wait and whether to stop
	wait, _ := r.backoff.Next(retry)

	logrus.Errorf("ElasticSearch error: %v, retrying in %s", err, wait)

	return wait, true, nil
}

// NewElasticsearchSink creates a ElasticsearchSink instance
func NewElasticsearchSink(address, indexName, username, password string, insecureSkipVerify, sniff bool) Sink {
	ctx := context.Background()

	es := ElasticsearchSink{
		Address:            address,
		ctx:                ctx,
		IndexName:          indexName,
		Username:           username,
		Password:           password,
		InsecureSkipVerify: insecureSkipVerify,
		Sniff:              sniff,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: es.InsecureSkipVerify},
	}
	httpClient := &http.Client{Transport: tr}

	client, err := elastic.NewClient(
		elastic.SetSniff(es.Sniff),
		elastic.SetURL(es.Address),
		elastic.SetHttpClient(httpClient),
		elastic.SetBasicAuth(es.Username, es.Password),
		elastic.SetRetrier(NewRetrier()),
	)

	if err != nil {
		logrus.Panicf("Error connecting to Elasticsearch (%s): %v", es.Address, err)
		panic(err)
	}

	es.client = client

	info, code, err := es.client.Ping(es.Address).Do(ctx)
	if err != nil {
		logrus.Panicf("Error pinging Elasticsearch: %v", err)
		panic(err)
	}
	logrus.Debugf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	mapping := fmt.Sprintf(`{"template":"%s-*"}`, es.IndexName)

	exists, err := client.IndexTemplateExists(es.IndexName).Do(ctx)
	if err != nil {
		logrus.Panicf("Error checking if index template exists: %v", err)
		panic(err)
	}
	if !exists {
		createIndex, err := client.IndexPutTemplate(es.IndexName).BodyString(mapping).Do(ctx)
		if err != nil {
			logrus.Panicf("Error creating index template %s: %v", es.IndexName, err)
			panic(err)
		}
		if !createIndex.Acknowledged {
			logrus.Panicf("Error acknowledging index template %s: %v", es.IndexName, err)
			panic(err)
		}
	}

	return Sink(&es)
}

// Write writes to an Elasticsearch server
func (es *ElasticsearchSink) Write(output []byte) error {
	log, err := es.client.Index().
		Index(es.getCurrentDayIndex()).
		Type("log").
		BodyJson(string(output)).
		Do(es.ctx)
	if err != nil {
		return err
	}

	logrus.Debugf("Elasticsearch indexed log %s to index %s", log.Id, log.Index)
	return nil
}

func (es *ElasticsearchSink) getCurrentDayIndex() string {
	t := time.Now()
	return fmt.Sprintf("%s-%s", es.IndexName, t.Format("2006-01-02"))
}
