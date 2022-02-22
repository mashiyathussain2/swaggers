package app

import (
	"context"
	"encoding/json"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"go-app/server/config"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Elasticsearch interface {
	GetPebble(opts *schema.GetPebbleFilter) ([]schema.GetPebbleESResp, error)
	GetPebbleV2(opts *schema.GetPebbleFilter) ([]*schema.GetPebbleESResp, error)
	GetPebbleByID(opts *schema.GetPebbleByIDFilter) (*schema.GetPebbleESResp, error)
	GetPebblesByInfluencerID(opts *schema.GetPebbleByInfluencerID) ([]*schema.GetPebbleESResp, error)
	GetPebblesByBrandID(opts *schema.GetPebbleByBrandID) ([]schema.GetPebbleESResp, error)
	GetCatalogsByInfluencerID(opts *schema.GetCatalogsByInfluencerID) ([]primitive.ObjectID, error)
	// GetPebbleSeries(opts *schema.GetPebbleSeriesFilter) ([]schema.GetPebbleSeriesESResp, error)
	GetPebbleAndSeries(opts *schema.GetPebbleFilter) ([]schema.GetPebbleSeriesESResp, error)
	GetPebbleCollections(opts *schema.GetCollectionFilter) ([]schema.GetPebbleCollectionESResp, error)
	GetPebblesByHashtag(opts *schema.GetPebbleByHashtag) ([]schema.GetPebbleESResp, error)
	GetSeriesByIDs(opts *schema.GetSeriesByIDs) ([]schema.GetPebbleSeriesESResp, error)
	GetPebblesInfoByCategoryID(opts *schema.GetPebbleByCategoryIDOpts) ([]schema.GetPebbleESResp, error)
	GetPebblesForCreator(opts *schema.GetPebbleByInfluencerID) ([]schema.GetPebbleESResp, error)
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

// GetPebble returns pebble with matching filter
func (ei *ElasticsearchImpl) GetPebble(opts *schema.GetPebbleFilter) ([]schema.GetPebbleESResp, error) {
	var queries []elastic.Query
	if len(opts.Genders) > 0 {
		queries = append(queries, elastic.NewTermsQueryFromStrings("label.genders", opts.Genders...))
	}
	if len(opts.Interests) > 0 {
		queries = append(queries, elastic.NewTermsQueryFromStrings("label.interests", opts.Interests...))
	}
	queries = append(queries, elastic.NewTermQuery("type", model.PebbleType))
	queries = append(queries, elastic.NewTermQuery("media_type", model.VideoType))
	queries = append(queries, elastic.NewTermQuery("is_active", true))
	boolQuery := elastic.NewBoolQuery().Must(queries...)

	sf := elastic.NewScriptField("is_liked_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['liked_by'].contains('%s')) {return true} return false`, opts.UserID)))
	builder := elastic.NewSearchSource().Query(boolQuery).FetchSource(true).ScriptFields(sf)
	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	res, err := ei.Client.Search().Index(ei.Config.ContentFullIndex).SearchSource(builder).Size(10).From(from).Sort("id", false).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble")
		return nil, errors.Wrap(err, "failed to get pebbles")
	}

	var resp []schema.GetPebbleESResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetPebbleESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		if ilbu, ok := hit.Fields["is_liked_by_user"]; ok {
			if len(ilbu.([]interface{})) != 0 {
				if isLikedByUser, ok := ilbu.([]interface{})[0].(bool); ok {
					s.IsLikedByUser = isLikedByUser
				}

			}
		}
		resp = append(resp, s)
	}

	return resp, nil
}

// GetPebbleV2 returns pebble with matching filter
func (ei *ElasticsearchImpl) GetPebbleV2(opts *schema.GetPebbleFilter) ([]*schema.GetPebbleESResp, error) {
	var wg sync.WaitGroup
	res, err := ei.getActivePebbles(opts)
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble")
		return nil, errors.Wrap(err, "failed to get pebbles")
	}

	var resp []*schema.GetPebbleESResp
	var pebbleIDs []interface{}
	pebbleIDResp := make(map[string]*schema.GetPebbleESResp)
	for _, hit := range res.Hits.Hits {
		var s schema.GetPebbleESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		pebbleIDs = append(pebbleIDs, s.ID.Hex())
		pebbleIDResp[s.ID.Hex()] = &s
		resp = append(resp, &s)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		res1, err1 := ei.getPebblesLikeCount(pebbleIDs, opts.UserID)
		if err1 != nil {
			ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble likes")
			return
		}
		if likeCount, ok := res1.Aggregations.Terms("like_count"); ok {
			for _, hit := range likeCount.Buckets {
				pebbleIDResp[hit.Key.(string)].LikeCount = int(hit.DocCount)
				if isLikedByUser, ok := hit.Aggregations.Filter("is_liked_by_user"); ok {
					if isLikedByUser.DocCount != 0 {
						pebbleIDResp[hit.Key.(string)].IsLikedByUser = true
					}
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res1, err1 := ei.getPebblesViewCount(pebbleIDs)
		if err1 != nil {
			ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble views")
			return
		}
		if viewCount, ok := res1.Aggregations.Terms("view_count"); ok {
			for _, hit := range viewCount.Buckets {
				pebbleIDResp[hit.Key.(string)].ViewCount = int(hit.DocCount)
			}
		}
	}()
	wg.Wait()
	return resp, nil
}

func (ei *ElasticsearchImpl) getActivePebbles(opts *schema.GetPebbleFilter) (*elastic.SearchResult, error) {
	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("type", model.PebbleType))
	queries = append(queries, elastic.NewTermQuery("media_type", model.VideoType))
	queries = append(queries, elastic.NewTermQuery("is_active", true))
	boolQuery := elastic.NewBoolQuery().Must(queries...)
	if !opts.IsSeries {
		boolQuery.MustNot(elastic.NewExistsQuery("series_ids"))
	}
	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	return ei.Client.Search().Index(ei.Config.ContentFullIndex).Query(boolQuery).Size(10).From(from).Sort("id", false).Do(context.Background())
}

func (ei *ElasticsearchImpl) getPebblesInfoByInfluencerID(opts *schema.GetPebbleByInfluencerID) (*elastic.SearchResult, error) {
	var queries []elastic.Query

	queries = append(queries, elastic.NewTermQuery("type", model.PebbleType))
	queries = append(queries, elastic.NewTermQuery("media_type", model.VideoType))
	if opts.IsActive {
		queries = append(queries, elastic.NewTermQuery("is_active", true))
	}
	queries = append(queries, elastic.NewTermQuery("influencer_ids", opts.InfluencerID))

	boolQuery := elastic.NewBoolQuery().Must(queries...)

	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	return ei.Client.Search().Index(ei.Config.ContentFullIndex).Query(boolQuery).Size(10).From(from).Sort("id", false).Do(context.Background())
}

func (ei *ElasticsearchImpl) getPebblesInfoByBrandID(opts *schema.GetPebbleByBrandID) (*elastic.SearchResult, error) {
	var queries []elastic.Query

	queries = append(queries, elastic.NewTermQuery("type", model.PebbleType))
	queries = append(queries, elastic.NewTermQuery("media_type", model.VideoType))
	queries = append(queries, elastic.NewTermQuery("is_active", true))

	queries = append(queries, elastic.NewTermQuery("brand_ids", opts.BrandID))

	boolQuery := elastic.NewBoolQuery().Must(queries...)

	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	return ei.Client.Search().Index(ei.Config.ContentFullIndex).Query(boolQuery).Size(10).From(from).Sort("id", false).Do(context.Background())
}

func (ei *ElasticsearchImpl) getPebblesByInfluencerID(opts *schema.GetPebbleByInfluencerID) ([]schema.GetPebbleESResp, error) {
	var queries []elastic.Query

	queries = append(queries, elastic.NewTermQuery("type", model.PebbleType))
	queries = append(queries, elastic.NewTermQuery("media_type", model.VideoType))
	if opts.IsActive {
		queries = append(queries, elastic.NewTermQuery("is_active", true))
	}
	queries = append(queries, elastic.NewTermQuery("influencer_ids", opts.InfluencerID))

	boolQuery := elastic.NewBoolQuery().Must(queries...)

	var from int
	fmt.Println(" getPebblesByInfluencerID Page ", opts.Page)
	if opts.Page > 0 {
		from = int(opts.Page) * 10
	}
	fmt.Println(" getPebblesByInfluencerID From ", from)

	resp, err := ei.Client.Search().Index(ei.Config.ContentFullIndex).Query(boolQuery).From(from).Size(10).Sort("id", false).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble by influencer id")
		return nil, errors.Wrap(err, "failed to get pebbles by influencer id")
	}
	fmt.Println(" getPebblesByInfluencerID From ", resp)

	var res []schema.GetPebbleESResp
	for _, hit := range resp.Hits.Hits {
		var s schema.GetPebbleESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		res = append(res, s)
	}

	return res, nil
}

func (ei *ElasticsearchImpl) getPebblesLikeCount(ids []interface{}, userID string) (*elastic.SearchResult, error) {
	termQuery := elastic.NewTermsQuery("resource_id", ids...)
	// Query like counts
	isLikedByUserAggsQuery := elastic.NewFilterAggregation().Filter(elastic.NewTermQuery("user_id", userID))
	likeCountAggsQuery := elastic.NewTermsAggregation().Field("resource_id").IncludeValues(ids...).SubAggregation("is_liked_by_user", isLikedByUserAggsQuery).Size(50)
	return ei.Client.Search().Index(config.GetConfig().ElasticsearchConfig.LikeIndex).Size(0).Query(termQuery).Aggregation("like_count", likeCountAggsQuery).Do(context.TODO())
}

func (ei *ElasticsearchImpl) getPebblesViewCount(ids []interface{}) (*elastic.SearchResult, error) {
	termQuery := elastic.NewTermsQuery("resource_id", ids...)
	// Query like counts
	viewCountAggsQuery := elastic.NewTermsAggregation().Field("resource_id").IncludeValues(ids...).Size(50)
	return ei.Client.Search().Index(config.GetConfig().ElasticsearchConfig.ViewIndex).Size(0).Query(termQuery).Aggregation("view_count", viewCountAggsQuery).Do(context.TODO())
}

// GetPebbleByID returns pebble with matching id
func (ei *ElasticsearchImpl) GetPebbleByID(opts *schema.GetPebbleByIDFilter) (*schema.GetPebbleESResp, error) {
	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("type", model.PebbleType))
	queries = append(queries, elastic.NewTermQuery("media_type", model.VideoType))
	queries = append(queries, elastic.NewTermQuery("is_active", true))
	queries = append(queries, elastic.NewTermQuery("id", opts.ID))
	boolQuery := elastic.NewBoolQuery().Must(queries...)

	resp, err := ei.Client.Search().Index(ei.Config.ContentFullIndex).Query(boolQuery).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble by id")
		return nil, errors.Wrap(err, "failed to get pebbles by id")
	}

	var finalResp []*schema.GetPebbleESResp
	var pebbleIDs []interface{}
	pebbleIDResp := make(map[string]*schema.GetPebbleESResp)

	for _, hit := range resp.Hits.Hits {
		var s schema.GetPebbleESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, err
		}
		pebbleIDs = append(pebbleIDs, s.ID.Hex())
		pebbleIDResp[s.ID.Hex()] = &s
		finalResp = append(finalResp, &s)
	}

	// Fetching Likes and Views for pebble
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		res, err := ei.getPebblesLikeCount(pebbleIDs, opts.UserID)
		if err != nil {
			ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble likes")
			return
		}
		if likeCount, ok := res.Aggregations.Terms("like_count"); ok {
			for _, hit := range likeCount.Buckets {
				if _, ok := pebbleIDResp[hit.Key.(string)]; ok {
					pebbleIDResp[hit.Key.(string)].LikeCount = int(hit.DocCount)
					if isLikedByUser, ok := hit.Aggregations.Filter("is_liked_by_user"); ok {
						if isLikedByUser.DocCount != 0 {
							pebbleIDResp[hit.Key.(string)].IsLikedByUser = true
						}
					}
				}
			}
		}
	}()

	wg1.Wait()
	if len(finalResp) == 0 {
		return nil, nil
	}
	return finalResp[0], nil
}

// GetPebblesByBrandID returns pebble with matching brand id
func (ei *ElasticsearchImpl) GetPebblesByBrandID(opts *schema.GetPebbleByBrandID) ([]schema.GetPebbleESResp, error) {
	var wg sync.WaitGroup

	var finalResp []schema.GetPebbleESResp
	var pebbleIDs []interface{}

	pebbleIDResp := make(map[string]*schema.GetPebbleESResp)
	// Getting Active Pebble which are not part of the series
	wg.Add(1)

	go func() {
		defer wg.Done()
		if pebbleResp, err := ei.getPebblesInfoByBrandID(opts); err == nil {
			for _, hit := range pebbleResp.Hits.Hits {
				var s schema.GetPebbleESResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					return
				}
				pebbleIDs = append(pebbleIDs, s.ID.Hex())
				pebbleIDResp[s.ID.Hex()] = &s
				finalResp = append(finalResp, s)
			}
		}
	}()
	wg.Wait()

	// Fetching Likes and Views for pebble
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		res, err := ei.getPebblesLikeCount(pebbleIDs, opts.UserID)
		if err != nil {
			ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble likes")
			return
		}
		if likeCount, ok := res.Aggregations.Terms("like_count"); ok {
			for _, hit := range likeCount.Buckets {
				if _, ok := pebbleIDResp[hit.Key.(string)]; ok {
					pebbleIDResp[hit.Key.(string)].LikeCount = int(hit.DocCount)
					if isLikedByUser, ok := hit.Aggregations.Filter("is_liked_by_user"); ok {
						if isLikedByUser.DocCount != 0 {
							pebbleIDResp[hit.Key.(string)].IsLikedByUser = true
						}
					}
				}
			}
		}
	}()

	wg1.Wait()
	return finalResp, nil
}

// GetPebblesByInfluencerID returns pebble with matching influencer_id
func (ei *ElasticsearchImpl) GetPebblesByInfluencerID(opts *schema.GetPebbleByInfluencerID) ([]*schema.GetPebbleESResp, error) {
	var wg sync.WaitGroup
	opts.IsActive = true

	var finalResp []*schema.GetPebbleESResp
	var pebbleIDs []interface{}

	pebbleIDResp := make(map[string]*schema.GetPebbleESResp)
	// Getting Active Pebble which are not part of the series
	wg.Add(1)

	go func() {
		defer wg.Done()
		if pebbleResp, err := ei.getPebblesInfoByInfluencerID(opts); err == nil {
			for _, hit := range pebbleResp.Hits.Hits {
				var s schema.GetPebbleESResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					return
				}
				pebbleIDs = append(pebbleIDs, s.ID.Hex())
				pebbleIDResp[s.ID.Hex()] = &s
				finalResp = append(finalResp, &s)
			}
		}
	}()
	wg.Wait()

	// Fetching Likes and Views for pebble
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		res, err := ei.getPebblesLikeCount(pebbleIDs, opts.UserID)
		if err != nil {
			ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble likes")
			return
		}
		if likeCount, ok := res.Aggregations.Terms("like_count"); ok {
			for _, hit := range likeCount.Buckets {
				if _, ok := pebbleIDResp[hit.Key.(string)]; ok {
					pebbleIDResp[hit.Key.(string)].LikeCount = int(hit.DocCount)
					if isLikedByUser, ok := hit.Aggregations.Filter("is_liked_by_user"); ok {
						if isLikedByUser.DocCount != 0 {
							pebbleIDResp[hit.Key.(string)].IsLikedByUser = true
						}
					}
				}
			}
		}
	}()

	wg1.Wait()
	return finalResp, nil
}

// GetCatalogsByInfluencerID returns catalogs with matching influencer_id
func (ei *ElasticsearchImpl) GetCatalogsByInfluencerID(opts *schema.GetCatalogsByInfluencerID) ([]primitive.ObjectID, error) {

	pebblesOpts := schema.GetPebbleByInfluencerID{
		UserID:       opts.UserID,
		InfluencerID: opts.InfluencerID,
		Page:         opts.Page,
		IsActive:     true,
	}

	resp, err := ei.getPebblesByInfluencerID(&pebblesOpts)
	fmt.Println("here")
	fmt.Printf("%+v\n", resp)
	if err != nil {
		return nil, err
	}
	var catIDs []primitive.ObjectID
	for _, r := range resp {
		catIDs = append(catIDs, r.CatalogIDs...)
	}
	return catIDs, nil
}

func (ei *ElasticsearchImpl) getSeriesPebbles(opts *schema.GetPebbleFilter) (*elastic.SearchResult, error) {
	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("is_active", true))
	boolQuery := elastic.NewBoolQuery().Must(queries...)
	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	return ei.Client.Search().Index(ei.Config.PebbleSeriesFullIndex).Query(boolQuery).Size(10).From(from).Sort("id", false).Do(context.Background())
}

func (ei *ElasticsearchImpl) GetPebbleAndSeries(opts *schema.GetPebbleFilter) ([]schema.GetPebbleSeriesESResp, error) {
	var wg sync.WaitGroup

	var finalResp []schema.GetPebbleSeriesESResp
	var pebbleIDs []interface{}

	pebbleIDResp := make(map[string]*schema.GetPebbleESResp)
	// Getting Active Pebble which are not part of the series
	wg.Add(1)
	go func() {
		defer wg.Done()
		if pebbleResp, err := ei.getActivePebbles(opts); err == nil {
			for _, hit := range pebbleResp.Hits.Hits {
				var s schema.GetPebbleESResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					return
				}
				pebbleIDs = append(pebbleIDs, s.ID.Hex())
				pebbleIDResp[s.ID.Hex()] = &s
				finalResp = append(finalResp, schema.GetPebbleSeriesESResp{
					PebbleIds:  []interface{}{s.ID.Hex()},
					PebbleInfo: []*schema.GetPebbleESResp{&s},
				})
			}
		}
	}()

	// Getting Series Pebbles
	wg.Add(1)
	go func() {
		defer wg.Done()
		res, err := ei.getSeriesPebbles(opts)
		if err == nil {
			var resp []schema.GetPebbleSeriesESResp
			for _, hit := range res.Hits.Hits {
				var s schema.GetPebbleSeriesESResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					return
				}
				resp = append(resp, s)
				pebbleIDs = append(pebbleIDs, s.PebbleIds...)
				for _, pebble := range s.PebbleInfo {
					pebbleIDResp[pebble.ID.Hex()] = pebble
				}
			}
			finalResp = append(finalResp, resp...)
		}
	}()
	wg.Wait()

	// Fetching Likes and Views for pebble
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		res, err := ei.getPebblesLikeCount(pebbleIDs, opts.UserID)
		if err != nil {
			ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble likes")
			return
		}
		if likeCount, ok := res.Aggregations.Terms("like_count"); ok {
			for _, hit := range likeCount.Buckets {
				if _, ok := pebbleIDResp[hit.Key.(string)]; ok {
					pebbleIDResp[hit.Key.(string)].LikeCount = int(hit.DocCount)
					if isLikedByUser, ok := hit.Aggregations.Filter("is_liked_by_user"); ok {
						if isLikedByUser.DocCount != 0 {
							pebbleIDResp[hit.Key.(string)].IsLikedByUser = true
						}
					}
				}
			}
		}
	}()

	// wg1.Add(1)
	// go func() {
	// 	defer wg1.Done()
	// 	res, err := ei.getPebblesViewCount(pebbleIDs)
	// 	if err != nil {
	// 		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble views")
	// 		return
	// 	}
	// 	if viewCount, ok := res.Aggregations.Terms("view_count"); ok {
	// 		for _, hit := range viewCount.Buckets {
	// 			if _, ok := pebbleIDResp[hit.Key.(string)]; ok {
	// 				pebbleIDResp[hit.Key.(string)].ViewCount = int(hit.DocCount)
	// 			}
	// 		}
	// 	}
	// }()

	wg1.Wait()

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(finalResp), func(i, j int) { finalResp[i], finalResp[j] = finalResp[j], finalResp[i] })
	return finalResp, nil
}

func (ei *ElasticsearchImpl) GetPebbleCollections(opts *schema.GetCollectionFilter) ([]schema.GetPebbleCollectionESResp, error) {
	var queries []elastic.Query
	if len(opts.Genders) > 0 {
		queries = append(queries, elastic.NewTermsQueryFromStrings("genders", opts.Genders...))
	}
	queries = append(queries, elastic.NewTermQuery("status", model.Publish))
	boolQuery := elastic.NewBoolQuery().Must(queries...)

	builder := elastic.NewSearchSource().Query(boolQuery).FetchSource(true)
	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	res, err := ei.Client.Search().Index(ei.Config.PebbleCollectionFullIndex).SearchSource(builder).Size(10).From(from).Sort("id", false).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble collection")
		return nil, errors.Wrap(err, "failed to get pebbles")
	}

	var resp []schema.GetPebbleCollectionESResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleCollectionESResp
		var s schema.GetPebbleCollectionESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		resp = append(resp, s)
	}
	return resp, nil
}

// GetPebblesByHashtag returns pebble with matching Hashtag
func (ei *ElasticsearchImpl) GetPebblesByHashtag(opts *schema.GetPebbleByHashtag) ([]schema.GetPebbleESResp, error) {

	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("type", model.PebbleType))
	queries = append(queries, elastic.NewTermQuery("media_type", model.VideoType))
	queries = append(queries, elastic.NewTermQuery("is_active", true))
	// queries = append(queries, elastic.NewMatchQuery("label.interests", opts.Hashtag))
	queries = append(queries, elastic.NewTermQuery("hashtags.hashtags", opts.Hashtag))

	boolQuery := elastic.NewBoolQuery().Must(queries...)

	sf := elastic.NewScriptField("is_liked_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['liked_by'].contains('%s')) {return true} return false`, opts.UserID)))

	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	builder := elastic.NewSearchSource().Query(boolQuery).FetchSource(true).ScriptFields(sf)
	res, err := ei.Client.Search().Index(ei.Config.ContentFullIndex).SearchSource(builder).Size(10).From(from).Sort("id", false).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble")
		return nil, errors.Wrap(err, "failed to get pebbles")
	}

	var resp []schema.GetPebbleESResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetPebbleESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		if ilbu, ok := hit.Fields["is_liked_by_user"]; ok {
			if len(ilbu.([]interface{})) != 0 {
				if isLikedByUser, ok := ilbu.([]interface{})[0].(bool); ok {
					s.IsLikedByUser = isLikedByUser
				}

			}
		}
		resp = append(resp, s)
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return resp, nil
}

// GetSeriesByIDs returns series with matching ids
func (ei *ElasticsearchImpl) GetSeriesByIDs(opts *schema.GetSeriesByIDs) ([]schema.GetPebbleSeriesESResp, error) {

	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("is_active", true))
	queries = append(queries, elastic.NewTermsQueryFromStrings("id", opts.ID...))

	boolQuery := elastic.NewBoolQuery().Must(queries...)

	// sf := elastic.NewScriptField("is_liked_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['liked_by'].contains('%s')) {return true} return false`, opts.UserID)))

	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	builder := elastic.NewSearchSource().Query(boolQuery).FetchSource(true)
	res, err := ei.Client.Search().Index(ei.Config.PebbleSeriesFullIndex).SearchSource(builder).Size(10).From(from).Sort("id", false).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble series")
		return nil, errors.Wrap(err, "failed to get pebbles")
	}

	var resp []schema.GetPebbleSeriesESResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetPebbleSeriesESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode series json")
		}

		resp = append(resp, s)
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return resp, nil
}

//GetPebblesInfoByCategoryID return pebbles based on category id
func (ei *ElasticsearchImpl) GetPebblesInfoByCategoryID(opts *schema.GetPebbleByCategoryIDOpts) ([]schema.GetPebbleESResp, error) {
	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("type", model.PebbleType))
	queries = append(queries, elastic.NewTermQuery("media_type", model.VideoType))
	queries = append(queries, elastic.NewTermQuery("is_active", true))
	queries = append(queries, elastic.NewTermQuery("category_path", opts.CategoryID))
	boolQuery := elastic.NewBoolQuery().Must(queries...)

	sf := elastic.NewScriptField("is_liked_by_user", elastic.NewScript(fmt.Sprintf(`if (doc['liked_by'].contains('%s')) {return true} return false`, opts.UserID)))

	var from int
	if opts.Page > 0 {
		from = int(opts.Page)*10 + 1
	}
	builder := elastic.NewSearchSource().Query(boolQuery).FetchSource(true).ScriptFields(sf)
	res, err := ei.Client.Search().Index(ei.Config.ContentFullIndex).SearchSource(builder).Size(10).From(from).Sort("id", false).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble")
		return nil, errors.Wrap(err, "failed to get pebbles")
	}

	var resp []schema.GetPebbleESResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetPebbleESResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode content json")
		}
		if ilbu, ok := hit.Fields["is_liked_by_user"]; ok {
			if len(ilbu.([]interface{})) != 0 {
				if isLikedByUser, ok := ilbu.([]interface{})[0].(bool); ok {
					s.IsLikedByUser = isLikedByUser
				}

			}
		}
		resp = append(resp, s)
	}
	if len(resp) == 0 {
		return nil, nil
	}
	return resp, nil
}

// GetPebble returns pebble with matching filter
func (ei *ElasticsearchImpl) GetPebblesNotInSeries(opts *schema.GetPebbleFilter) ([]*schema.GetPebbleESResp, error) {
	return ei.GetPebbleV2(opts)
}

// GetPebblesForCreator returns pebble with matching influencer_id
func (ei *ElasticsearchImpl) GetPebblesForCreator(opts *schema.GetPebbleByInfluencerID) ([]schema.GetPebbleESResp, error) {
	var wg sync.WaitGroup

	var finalResp []schema.GetPebbleESResp
	var pebbleIDs []interface{}

	pebbleIDResp := make(map[string]*schema.GetPebbleESResp)
	// Getting Active Pebble which are not part of the series
	wg.Add(1)

	go func() {
		defer wg.Done()
		if pebbleResp, err := ei.getPebblesInfoByInfluencerID(opts); err == nil {
			for _, hit := range pebbleResp.Hits.Hits {
				var s schema.GetPebbleESResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					return
				}
				pebbleIDs = append(pebbleIDs, s.ID.Hex())
				pebbleIDResp[s.ID.Hex()] = &s
				finalResp = append(finalResp, s)
			}
		}
	}()
	wg.Wait()

	// Fetching Likes and Views for pebble
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		res, err := ei.getPebblesLikeCount(pebbleIDs, opts.UserID)
		if err != nil {
			ei.Logger.Err(err).Interface("opts", opts).Msg("failed to get pebble likes")
			return
		}
		if likeCount, ok := res.Aggregations.Terms("like_count"); ok {
			for _, hit := range likeCount.Buckets {
				if _, ok := pebbleIDResp[hit.Key.(string)]; ok {
					pebbleIDResp[hit.Key.(string)].LikeCount = int(hit.DocCount)
					if isLikedByUser, ok := hit.Aggregations.Filter("is_liked_by_user"); ok {
						if isLikedByUser.DocCount != 0 {
							pebbleIDResp[hit.Key.(string)].IsLikedByUser = true
						}
					}
				}
			}
		}
	}()

	wg1.Wait()
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(finalResp), func(i, j int) { finalResp[i], finalResp[j] = finalResp[j], finalResp[i] })
	return finalResp, nil

}
