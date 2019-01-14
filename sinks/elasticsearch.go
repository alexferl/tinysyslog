package sinks

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
)

// ElasticsearchSink represents an Elasticsearch sink
type ElasticsearchSink struct {
	client    *elastic.Client
	ctx       context.Context
	Address   string
	IndexName string
}

const mapping = `
{
	"mappings":{
		"syslog":{
			"properties":{
				"timestamp":{
					"type":"date"
				},
				"hostname":{
					"type":"keyword"
				},
				"app_name":{
					"type":"keyword"
				},
				"proc_id":{
					"type":"keyword"
				},
				"severity":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				}
			}
		}
	}
}`

// NewElasticsearchSink creates a ElasticsearchSink instance
func NewElasticsearchSink(address, indexName string) Sink {
	ctx := context.Background()

	es := ElasticsearchSink{
		Address:   address,
		ctx:       ctx,
		IndexName: indexName,
	}

	client, err := elastic.NewClient(elastic.SetURL(es.Address))
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

	exists, err := client.IndexExists(es.IndexName).Do(ctx)
	if err != nil {
		logrus.Panicf("Error checking if index exists: %v", err)
		panic(err)
	}
	if !exists {
		createIndex, err := client.CreateIndex(es.IndexName).BodyString(mapping).Do(ctx)
		if err != nil {
			logrus.Panicf("Error creating index %s: %v", es.IndexName, err)
			panic(err)
		}
		if !createIndex.Acknowledged {
			logrus.Panicf("Error creating index %s: %v", es.IndexName, err)
			panic(err)
		}
	}

	return Sink(&es)
}

// Write writes to an Elasticsearch server
func (es *ElasticsearchSink) Write(output []byte) error {
	log, err := es.client.Index().
		Index(es.IndexName).
		Type("syslog").
		BodyJson(string(output)).
		Do(es.ctx)
	if err != nil {
		return err
	}

	logrus.Debugf("Indexed log %s to index %s, type %s\n", log.Id, log.Index, log.Type)
	return nil
}
