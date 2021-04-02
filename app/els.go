package app

import (
	"context"
	"encoding/json"
	"go-app/model"
	"go-app/schema"
	"go-app/server/config"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Elasticsearch interface {
	GetActiveCollections() ([]schema.GetCollectionESResp, error)
	GetCatalogByIDs([]string) ([]schema.GetCatalogBasicResp, error)
	GetCatalogInfoByID(string) (*schema.GetCatalogInfoResp, error)
	GetCatalogInfoByCategoryID(string) ([]schema.GetCatalogBasicResp, error)
}

type ElasticsearchImpl struct {
	Client *elastic.Client
	Config *config.ElasticsearchConfig
	Logger *zerolog.Logger
}

type ElasticsearchOpts struct {
	Config *config.ElasticsearchConfig
	Logger *zerolog.Logger
}

func InitElasticsearch(opts *ElasticsearchOpts) Elasticsearch {
	c, err := elastic.NewClient(
		elastic.SetURL(opts.Config.Endpoint),
		elastic.SetBasicAuth(opts.Config.Username, opts.Config.Password),
		elastic.SetHealthcheck(false),
		elastic.SetSniff(false),
		elastic.SetTraceLog(opts.Logger),
	)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	ei := ElasticsearchImpl{
		Client: c,
		Config: opts.Config,
		Logger: opts.Logger,
	}
	return &ei
}

func (ei *ElasticsearchImpl) GetActiveCollections() ([]schema.GetCollectionESResp, error) {
	query := elastic.NewTermQuery("status", model.Publish)
	res, err := ei.Client.Search().Index(ei.Config.CollectionFullIndex).Query(query).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get active collections")
		return nil, errors.Wrap(err, "failed to get active collections")
	}
	var resp []schema.GetCollectionESResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetCollectionESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		resp = append(resp, s)
	}

	return resp, nil
}

func (ei *ElasticsearchImpl) GetCatalogByIDs(ids []string) ([]schema.GetCatalogBasicResp, error) {
	query := elastic.NewTermsQueryFromStrings("id", ids...)
	res, err := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get active collections")
		return nil, errors.Wrap(err, "failed to get active collections")
	}

	var resp []schema.GetCatalogBasicResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetCatalogBasicResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		resp = append(resp, s)
	}

	return resp, nil
}

func (ei *ElasticsearchImpl) GetCatalogInfoByID(id string) (*schema.GetCatalogInfoResp, error) {
	query := elastic.NewTermQuery("id", id)
	res, err := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get active collections")
		return nil, errors.Wrap(err, "failed to get active collections")
	}

	var resp []schema.GetCatalogInfoResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetCatalogInfoResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		resp = append(resp, s)
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return &resp[0], nil
}

func (ei *ElasticsearchImpl) GetCatalogInfoByCategoryID(id string) ([]schema.GetCatalogBasicResp, error) {
	query := elastic.NewTermQuery("category_path", id)
	res, err := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get catalogs")
		return nil, errors.Wrap(err, "failed to get catalogs")
	}

	var resp []schema.GetCatalogBasicResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetCatalogBasicResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode catalog basic json")
		}
		resp = append(resp, s)
	}

	return resp, nil
}
