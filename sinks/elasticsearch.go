package sinks

import (
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

type Elasticsearch struct {
	client *elasticsearch.TypedClient
	config ElasticsearchConfig
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
		log.Panic().Err(err).Msgf("error creating elasticsearch client")
	}

	es := &Elasticsearch{
		client: client,
		config: conf,
	}
	idx := es.config.IndexName

	ctx := context.Background()

	cctx, cancel := context.WithTimeout(ctx, es.config.Timeout)
	defer cancel()

	exists, err := client.Indices.ExistsIndexTemplate(idx).Do(cctx)
	if err != nil {
		log.Panic().Err(err).Msgf("error checking if index template exists")
	}

	if !exists {
		cctx, cancel = context.WithTimeout(ctx, es.config.Timeout)
		defer cancel()
		pattern := fmt.Sprintf("%s-*", idx)

		resp, err := client.Indices.PutTemplate(idx).IndexPatterns(pattern).Do(cctx)
		if err != nil {
			log.Panic().Err(err).Msgf("error creating index template '%s'", idx)
		}
		if !resp.Acknowledged {
			log.Panic().Err(err).Msgf("error acknowledging index template '%s'", idx)
		}
	}

	return Sink(es)
}

func (es *Elasticsearch) Write(output []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), es.config.Timeout)
	defer cancel()

	resp, err := es.client.Index(es.getCurrentDayIndex()).Request(string(output)).Do(ctx)
	if err != nil {
		return err
	}

	log.Debug().Msgf("elasticsearch indexed log '%s' to index '%s'", resp.Id_, resp.Index_)

	return nil
}

func (es *Elasticsearch) getCurrentDayIndex() string {
	t := time.Now()
	return fmt.Sprintf("%s-%s", es.config.IndexName, t.Format("2006-01-02"))
}
