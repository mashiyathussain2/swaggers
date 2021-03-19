package app

import (
	"context"
	"encoding/json"
	"go-app/schema"
	"go-app/server/config"
	"log"

	"os"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Elasticsearch contains methods to search for external collections such as
// catalog, brand, user, influencer
type Elasticsearch interface {
	// GetBrand([]string) ([]schema.GetBrandSchema, error)
}

type ElasticsearchImpl struct {
	Config *config.ElasticsearchConfig
	Es     *elastic.Client
	Logger *zerolog.Logger
}

type ElasticsearchImplOpts struct {
	Logger *zerolog.Logger
	Config *config.ElasticsearchConfig
}

func InitElasticsearch(opts *ElasticsearchImplOpts) Elasticsearch {
	e := ElasticsearchImpl{
		Logger: opts.Logger,
		Config: opts.Config,
	}
	client, err := elastic.NewClient(
		elastic.SetURL(opts.Config.Endpoint),
		// elastic.SetBasicAuth("admin", "Admin@1234"),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal(err, " failed to load elastic client")
		os.Exit(1)
	}
	e.Es = client

	return e
}

// GetBrand returns brands matching with id
func (ei *ElasticsearchImpl) GetBrand(ids []string) ([]schema.GetBrandSchema, error) {
	var brands []schema.GetBrandSchema
	q := elastic.NewIdsQuery(ids...)
	res, err := ei.Es.Search().
		Index(ei.Config.BrandIndex).
		Query(q).
		Do(context.TODO())
	if err != nil {
		return nil, err
	}

	if res.Hits == nil {
		return nil, errors.Errorf("brand with id:%s not found", ids)
	}

	for i, hit := range res.Hits.Hits {
		var s schema.GetBrandSchema
		if err := json.Unmarshal(*&hit.Source, &s); err != nil {
			ei.Logger.Err(err).Interface("hit", hit.Source).Msgf("failed to decode brand with id:%s", ids[i])
			return nil, errors.Errorf("failed to decode brand with id:%s", ids[i])
		}
		brands = append(brands, s)
	}
	return brands, nil
}
