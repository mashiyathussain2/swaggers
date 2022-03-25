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
	GetBrandsByUsernameBasic(opts *schema.GetBrandsByUsernameBasicOpts) ([]schema.GetBrandBasicESEesp, error)
	GetBrandsList(*schema.GetBrandsListOpts) ([]schema.GetActiveBrandsListESEesp, error)
	GetBrandInfoByUsername(opts *schema.GetBrandsInfoByUsernameOpts) (*schema.GetBrandInfoEsResp, error)

	GetInfluencerInfoByID(*schema.GetInfluencerInfoByIDOpts) (*schema.GetInfluencerInfoEsResp, error)
	GetInfluencersByIDBasic(*schema.GetInfluencersByIDBasicOpts) ([]schema.GetInfluencerBasicESEesp, error)
	GetInfluencersByUserameBasic(*schema.GetInfluencersByUsernameBasicOpts) ([]schema.GetInfluencerBasicESEesp, error)
	GetInfluencerInfoByUsername(opts *schema.GetInfluencerInfoByUsernameOpts) (*schema.GetInfluencerInfoEsResp, error)
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
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.CustomerID.Hex())))
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
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.CustomerID.Hex())))
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

func (ei *ElasticsearchImpl) GetBrandsByUsernameBasic(opts *schema.GetBrandsByUsernameBasicOpts) ([]schema.GetBrandBasicESEesp, error) {
	query := elastic.NewTermsQueryFromStrings("username", opts.Usernames...)
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.CustomerID.Hex())))
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

func (ei *ElasticsearchImpl) GetBrandInfoByUsername(opts *schema.GetBrandsInfoByUsernameOpts) (*schema.GetBrandInfoEsResp, error) {
	query := elastic.NewTermQuery("username", opts.Username)
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.CustomerID.Hex())))
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
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.CustomerID.Hex())))
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
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.CustomerID.Hex())))
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

func (ei *ElasticsearchImpl) GetInfluencersByUserameBasic(opts *schema.GetInfluencersByUsernameBasicOpts) ([]schema.GetInfluencerBasicESEesp, error) {
	query := elastic.NewTermsQueryFromStrings("username", opts.Usernames...)
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.CustomerID.Hex())))
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

func (ei *ElasticsearchImpl) GetInfluencerInfoByUsername(opts *schema.GetInfluencerInfoByUsernameOpts) (*schema.GetInfluencerInfoEsResp, error) {
	query := elastic.NewTermQuery("username", opts.Username)
	sf := elastic.NewScriptField("is_followed_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['followers_id'] == null) {return false} if (doc['followers_id'].contains('%s')) {return true} return false`, opts.CustomerID.Hex())))
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
	cres, err := ei.GetInfluencerContentCount(&schema.GetInfluencerContentCount{ID: resp[0].ID.Hex()})
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get content count")
		return nil, errors.Wrap(err, "failed to get influencer content count")
	}
	resp[0].ContentCount = cres
	return &resp[0], nil
}

func (ei *ElasticsearchImpl) GetBrandsList(opts *schema.GetBrandsListOpts) ([]schema.GetActiveBrandsListESEesp, error) {
	var from int
	if opts.Page > 0 {
		from = opts.Page*opts.Size + 1
	}
	query := elastic.NewMatchAllQuery()
	resp, err := ei.Client.Search().Index(ei.Config.BrandFullIndex).Query(query).Size(opts.Size).Sort("name", true).From(from).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get brands list")
		return nil, errors.Wrap(err, "failed to get brands list")
	}
	var res []schema.GetActiveBrandsListESEesp
	for _, hit := range resp.Hits.Hits {
		var s schema.GetActiveBrandsListESEesp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		res = append(res, s)
	}
	return res, nil
}

func (ei *ElasticsearchImpl) GetInfluencerContentCount(opts *schema.GetInfluencerContentCount) (*schema.GetInfluencerContentCountResp, error) {
	query := elastic.NewTermsQueryFromStrings("influencer_id", opts.ID)
	// builder := elastic.NewSearchSource().Query(query).FetchSource(true)
	s := schema.GetInfluencerContentCountResp{}
	query_content := elastic.NewTermsQueryFromStrings("influencer_ids", opts.ID)
	res, err := ei.Client.Count().Index(ei.Config.ContentFullIndex).Query(query_content).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get influencers")
		return nil, errors.Wrap(err, "failed to get influencers")
	}
	s.Pebbles = res

	// query = elastic.NewTermsQueryFromStrings("influencer_id", opts.ID)
	sf := elastic.NewScriptField("count", elastic.NewScript(`return doc['catalog_ids'].size()`))
	builder := elastic.NewSearchSource().Query(query).FetchSource(true).ScriptFields(sf)

	resp, err := ei.Client.Search().Index(ei.Config.InfluencerProductIndex).SearchSource(builder).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get influencer product")
		return nil, errors.Wrap(err, "failed to get influencer product")
	}
	if len(resp.Hits.Hits) > 0 {
		s.Products = resp.Hits.Hits[0].Fields["count"].([]interface{})[0].(float64)
	} else {
		s.Products = 0
	}

	res, err = ei.Client.Count().Index(ei.Config.InfluencerCollectionIndex).Query(query).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get influencers")
		return nil, errors.Wrap(err, "failed to get influencers")
	}
	s.Collections = res
	return &s, nil
}
