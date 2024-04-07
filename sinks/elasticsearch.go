package sinks

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

type Elasticsearch struct {
	client *elasticsearch.TypedClient
	config ElasticsearchConfig
	kind   Kind
}

type ElasticsearchConfig struct {
	IndexName    string
	Timeout      time.Duration
	Addresses    []string
	Username     string
	Password     string
	CloudID      string
	APIKey       string
	ServiceToken string
}

func NewElasticsearch(conf ElasticsearchConfig) Sink {
	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses:    conf.Addresses,
		Username:     conf.Username,
		Password:     conf.Password,
		CloudID:      conf.CloudID,
		APIKey:       conf.APIKey,
		ServiceToken: conf.ServiceToken,
	})
	if err != nil {
		log.Panic().Err(err).Msgf("failed creating elasticsearch client")
	}

	es := &Elasticsearch{
		client: client,
		config: conf,
		kind:   ElasticsearchKind,
	}
	idx := es.config.IndexName

	ctx := context.Background()
	cctx, cancel := context.WithTimeout(ctx, es.config.Timeout)
	defer cancel()

	exists, err := client.Indices.ExistsIndexTemplate(idx).Do(cctx)
	if err != nil {
		log.Panic().Err(err).Msgf("failed checking if index template exists")
	}

	if !exists {
		pattern := fmt.Sprintf("%s-*", idx)
		resp, err := client.Indices.PutTemplate(idx).IndexPatterns(pattern).Do(cctx)
		if err != nil {
			log.Panic().Err(err).Msgf("failed creating index template '%s'", idx)
		}
		if !resp.Acknowledged {
			log.Panic().Err(err).Msgf("failed acknowledging index template '%s'", idx)
		}
	}

	return Sink(es)
}

func (es *Elasticsearch) Write(output []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), es.config.Timeout)
	defer cancel()

	resp, err := es.client.Index(es.getCurrentDayIndex()).Raw(bytes.NewReader(output)).Do(ctx)
	if err != nil {
		return err
	}

	log.Debug().Msgf("elasticsearch indexed log '%s' to index '%s'", resp.Id_, resp.Index_)

	return nil
}

func (es *Elasticsearch) GetKind() Kind {
	return es.kind
}

func (es *Elasticsearch) getCurrentDayIndex() string {
	t := time.Now()
	return fmt.Sprintf("%s-%s", es.config.IndexName, t.Format("2006-01-02"))
}
