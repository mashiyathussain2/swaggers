package app

import (
	"context"
	"encoding/json"
	"fmt"
	"go-app/schema"
	"go-app/server/config"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Elasticsearch interface {
	GetBrandsByIDBasic(*schema.GetBrandsByIDBasicOpts) ([]schema.GetBrandBasicESEesp, error)
	GetBrandInfoByID(*schema.GetBrandsInfoByIDOpts) (*schema.GetBrandInfoEsResp, error)

	GetInfluencerInfoByID(*schema.GetInfluencerInfoByIDOpts) (*schema.GetInfluencerInfoEsResp, error)
	GetInfluencersByIDBasic(*schema.GetInfluencersByIDBasicOpts) ([]schema.GetInfluencerBasicESEesp, error)
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

func (ei *ElasticsearchImpl) GetBrandsByIDBasic(opts *schema.GetBrandsByIDBasicOpts) ([]schema.GetBrandBasicESEesp, error) {
	query := elastic.NewTermsQueryFromStrings("id", opts.IDs...)
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.UserID.Hex())))
	builder := elastic.NewSearchSource().Query(query).FetchSource(true).ScriptFields(sf)
	res, err := ei.Client.Search().Index(ei.Config.BrandFullIndex).SearchSource(builder).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get brands")
		return nil, errors.Wrap(err, "failed to get brands")
	}
	var resp []schema.GetBrandBasicESEesp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetBrandBasicESEesp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		if ifbu, ok := hit.Fields["is_followed_by_user"]; ok {
			if len(ifbu.([]interface{})) != 0 {
				if isFollowedByUser, ok := ifbu.([]interface{})[0].(bool); ok {
					s.IsFollowedByUser = isFollowedByUser
				}
			}
		}
		resp = append(resp, s)
	}
	return resp, nil
}

func (ei *ElasticsearchImpl) GetBrandInfoByID(opts *schema.GetBrandsInfoByIDOpts) (*schema.GetBrandInfoEsResp, error) {
	query := elastic.NewTermQuery("id", opts.ID.Hex())
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.UserID.Hex())))
	builder := elastic.NewSearchSource().Query(query).FetchSource(true).ScriptFields(sf)
	res, err := ei.Client.Search().Index(ei.Config.BrandFullIndex).SearchSource(builder).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get brands")
		return nil, errors.Wrap(err, "failed to get brands")
	}
	var resp []schema.GetBrandInfoEsResp
	if len(res.Hits.Hits) == 0 {
		return nil, nil
	}
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetBrandInfoEsResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		if ifbu, ok := hit.Fields["is_followed_by_user"]; ok {
			if len(ifbu.([]interface{})) != 0 {
				if isFollowedByUser, ok := ifbu.([]interface{})[0].(bool); ok {
					s.IsFollowedByUser = isFollowedByUser
				}
			}
		}
		resp = append(resp, s)
	}
	return &resp[0], nil
}

func (ei *ElasticsearchImpl) GetInfluencersByIDBasic(opts *schema.GetInfluencersByIDBasicOpts) ([]schema.GetInfluencerBasicESEesp, error) {
	query := elastic.NewTermsQueryFromStrings("id", opts.IDs...)
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.UserID.Hex())))
	builder := elastic.NewSearchSource().Query(query).FetchSource(true).ScriptFields(sf)
	res, err := ei.Client.Search().Index(ei.Config.InfluencerFullIndex).SearchSource(builder).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get influencers")
		return nil, errors.Wrap(err, "failed to get influencers")
	}
	var resp []schema.GetInfluencerBasicESEesp
	for _, hit := range res.Hits.Hits {
		var s schema.GetInfluencerBasicESEesp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		if ifbu, ok := hit.Fields["is_followed_by_user"]; ok {
			if len(ifbu.([]interface{})) != 0 {
				if isFollowedByUser, ok := ifbu.([]interface{})[0].(bool); ok {
					s.IsFollowedByUser = isFollowedByUser
				}
			}
		}
		resp = append(resp, s)
	}
	return resp, nil
}

func (ei *ElasticsearchImpl) GetInfluencerInfoByID(opts *schema.GetInfluencerInfoByIDOpts) (*schema.GetInfluencerInfoEsResp, error) {
	query := elastic.NewTermQuery("id", opts.ID.Hex())
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.UserID.Hex())))
	builder := elastic.NewSearchSource().Query(query).FetchSource(true).ScriptFields(sf)
	res, err := ei.Client.Search().Index(ei.Config.InfluencerFullIndex).SearchSource(builder).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get influencer")
		return nil, errors.Wrap(err, "failed to get influencer")
	}
	var resp []schema.GetInfluencerInfoEsResp

	if len(res.Hits.Hits) == 0 {
		return nil, nil
	}

	for _, hit := range res.Hits.Hits {
		var s schema.GetInfluencerInfoEsResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		if ifbu, ok := hit.Fields["is_followed_by_user"]; ok {
			if len(ifbu.([]interface{})) != 0 {
				if isFollowedByUser, ok := ifbu.([]interface{})[0].(bool); ok {
					s.IsFollowedByUser = isFollowedByUser
				}
			}
		}
		resp = append(resp, s)
	}
	return &resp[0], nil
}
