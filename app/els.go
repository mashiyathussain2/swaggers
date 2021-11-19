package app

import (
	"context"
	"encoding/json"
	"go-app/model"
	"go-app/schema"
	"go-app/server/config"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Elasticsearch interface {
	GetActiveCollections(*schema.GetActiveCollectionsOpts) ([]schema.GetCollectionESResp, error)
	GetCatalogByIDs([]string) ([]schema.GetCatalogBasicResp, error)
	GetSimilarProducts(string) ([]schema.GetCatalogBasicResp, error)
	GetCatalogInfoByID(string) (*schema.GetCatalogInfoResp, error)
	GetCatalogInfoByCategoryID(*schema.GetCatalogByCategoryIDOpts) (*schema.GetCatalogByCategoryIDResp, error)

	GetCatalogBySaleID(*schema.GetCatalogBySaleIDOpts) ([]schema.GetCatalogBasicResp, error)
	SearchBrandCatalogInfluencerContent(opts *schema.SearchOpts) (*schema.SearchResp, error)
	SearchCatalog(opts *schema.SearchOpts) (*schema.SearchResp, error)
	SearchDiscover(opts *schema.SearchOpts) (*schema.DiscoverSearchResp, error)
	SearchBrand(opts *schema.SearchOpts) ([]schema.BrandSearchResp, error)
	SearchInfluencer(opts *schema.SearchOpts) ([]schema.InfluencerSearchResp, error)
	SearchHashtag(opts *schema.SearchOpts) ([]schema.HashtagSearchResp, error)
	SearchSeries(opts *schema.SearchOpts) ([]schema.SeriesSearchResp, error)
	GetReviewsByCatalogID(*schema.GetReviewsByCatalogIDFilter) ([]schema.GetReviewsByCatalogIDResp, error)
	GetCatalogByBrandID(*schema.GetCatalogByBrandIDOpts) ([]schema.GetCatalogBasicResp, error)
	GetCollectionCatalogByIDs(*schema.GetCollectionCatalogByIDs) ([]schema.GetCatalogBasicResp, error)

	// Shop Search API
	ShopSearch(opts *schema.SearchOpts) (*schema.ShopSearchResp, error)
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

func (ei *ElasticsearchImpl) GetActiveCollections(opts *schema.GetActiveCollectionsOpts) ([]schema.GetCollectionESResp, error) {
	var mustQueries []elastic.Query
	var shouldQueries []elastic.Query
	mustQueries = append(mustQueries, elastic.NewTermQuery("status", model.Publish))
	if opts.Gender != "" {
		mustQueries = append(mustQueries, elastic.NewTermsQuery("genders", opts.Gender))
	}
	shouldQueries = append(shouldQueries, elastic.NewNestedQuery(
		"sub_collections",
		elastic.NewNestedQuery(
			"sub_collections.catalog_info",
			elastic.NewMatchQuery("sub_collections.catalog_info.status.value", model.Publish),
		).InnerHit(elastic.NewInnerHit().Size(4)),
	).InnerHit(elastic.NewInnerHit().FetchSource(false)))

	boolQuery := elastic.NewBoolQuery().Must(mustQueries...).Should(shouldQueries...)
	fsctx := elastic.NewFetchSourceContext(true).Include([]string{"id", "name", "title", "type", "sub_collections.featured_catalog_ids", "sub_collections.catalog_ids", "sub_collections.id", "sub_collections.name", "sub_collections.image", "sub_collections.title", "inner_hits.sub_collections"}...)
	var pageFrom int
	var pageSize int = 20
	if opts.Size > 0 {
		pageSize = opts.Size
	}
	if opts.Page > 0 {
		pageFrom = (opts.Page * pageSize) + 1
	}
	res, err := ei.Client.Search().Index(ei.Config.CollectionFullIndex).Query(boolQuery).FetchSourceContext(fsctx).Sort("order", true).From(pageFrom).Size(pageSize).Do(context.Background())
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

		if s.Type == model.ProductCollection {
			// Fetching Inner Hits
			var catalogInfo []schema.GetCollectionCatalogInfoResp
			for _, innerHit := range hit.InnerHits["sub_collections"].Hits.Hits[0].InnerHits["sub_collections.catalog_info"].Hits.Hits {
				var ci schema.GetCollectionCatalogInfoResp
				if err := json.Unmarshal(innerHit.Source, &ci); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					return nil, errors.Wrap(err, "failed to decode content json")
				}
				catalogInfo = append(catalogInfo, ci)
			}
			if len(s.SubCollections) == 1 {
				s.SubCollections[0].CatalogInfo = catalogInfo
			}
		}
		resp = append(resp, s)
	}

	return resp, nil
}

func (ei *ElasticsearchImpl) GetCatalogByIDs(ids []string) ([]schema.GetCatalogBasicResp, error) {
	mustQuery := elastic.NewTermsQueryFromStrings("id", ids...)
	filterQuery := elastic.NewTermQuery("status.value", model.Publish)
	query := elastic.NewBoolQuery().Must(mustQuery).Filter(filterQuery)
	res, err := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).From(0).Size(100).Do(context.Background())
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

func (ei *ElasticsearchImpl) GetSimilarProducts(q string) ([]schema.GetCatalogBasicResp, error) {
	similarQuery := elastic.NewMoreLikeThisQuery().Field([]string{"name.name", "description"}...).LikeItems(elastic.NewMoreLikeThisQueryItem().Id(q).Index(ei.Config.CatalogFullIndex)).MinTermFreq(1).MaxQueryTerms(12)
	filterQuery := elastic.NewTermQuery("status.value", model.Publish)
	mustNot := elastic.NewMatchQuery("id", q)
	query := elastic.NewBoolQuery().Must(similarQuery).Filter(filterQuery).MustNot(mustNot)
	res, err := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).Size(7).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get similar products")
		return nil, errors.Wrap(err, "failed to get similar products")
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
	mustQuery := elastic.NewTermQuery("id", id)
	filterQuery := elastic.NewTermQuery("status.value", model.Publish)
	query := elastic.NewBoolQuery().Must(mustQuery).Filter(filterQuery)
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

func (ei *ElasticsearchImpl) GetCatalogInfoByCategoryID(opts *schema.GetCatalogByCategoryIDOpts) (*schema.GetCatalogByCategoryIDResp, error) {
	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("status.value", model.Publish))
	queries = append(queries, elastic.NewTermQuery("category_path", opts.CategoryID))
	query := elastic.NewBoolQuery().Must(queries...)
	if len(opts.BrandName) > 0 {
		query = query.Filter(elastic.NewTermsQueryFromStrings("brand_info.name.name", opts.BrandName...))
	}

	aggs := elastic.NewTermsAggregation().Field("brand_info.name.name").Size(99)
	var fromPage int
	if opts.Page != 0 {
		fromPage = (int(opts.Page) * 20) + 1
	}
	q := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).Aggregation("brands", aggs).From(fromPage).Size(20)
	switch opts.Sort {
	case -1:
		q = q.Sort("retail_price.value", false)
	case 1:
		q = q.Sort("retail_price.value", true)
	}
	res, err := q.Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get catalogs")
		return nil, errors.Wrap(err, "failed to get catalogs")
	}

	result := schema.GetCatalogByCategoryIDResp{}
	var resp []schema.GetCatalogBasicResp
	var filter []schema.GetCatalogByCategoryIDFilterResp
	for _, hit := range res.Hits.Hits {
		// Deserialize hit.Source into a GetPebbleESResp
		var s schema.GetCatalogBasicResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode catalog basic json")
		}
		resp = append(resp, s)
	}
	if aggs, ok := res.Aggregations.Terms("brands"); ok {
		for _, bucket := range aggs.Buckets {
			filter = append(filter, schema.GetCatalogByCategoryIDFilterResp{Key: bucket.Key.(string), Count: int(bucket.DocCount)})
		}
	}
	result.Data = resp
	result.BrandFilter = filter
	return &result, nil
}

func (ei *ElasticsearchImpl) GetCatalogBySaleID(opts *schema.GetCatalogBySaleIDOpts) ([]schema.GetCatalogBasicResp, error) {
	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("status.value", model.Publish))
	queries = append(queries, elastic.NewNestedQuery("discount_info", elastic.NewTermQuery("discount_info.sale_id", opts.SaleID)))
	query := elastic.NewBoolQuery().Must(queries...)
	res, err := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).From(int(opts.Page * 20)).Size(20).Do(context.Background())
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

func (ei *ElasticsearchImpl) SearchBrandCatalogInfluencerContent(opts *schema.SearchOpts) (*schema.SearchResp, error) {
	mSearch := elastic.NewMultiSearchService(ei.Client)
	var mSearchQuery []*elastic.SearchRequest

	filterQuery := elastic.NewTermQuery("status.value", model.Publish)
	mustQuery := elastic.NewMultiMatchQuery(opts.Query, []string{"brand_info.name.autocomplete", "name.autocomplete", "keywords.autocomplete"}...).Operator("or").Type("cross_fields")
	catalogQuery := elastic.NewBoolQuery().Must(mustQuery).Filter(filterQuery)
	mSearchQuery = append(mSearchQuery, elastic.NewSearchRequest().Index(ei.Config.CatalogFullIndex).Query(catalogQuery).Size(20).FetchSourceIncludeExclude([]string{"id", "name", "featured_image", "base_price", "retail_price", "discount_info", "variants.id", "brand_info"}, nil))

	brandQuery := elastic.NewMultiMatchQuery(opts.Query, []string{"lname.autocomplete"}...).Operator("or").Type("cross_fields")
	mSearchQuery = append(mSearchQuery, elastic.NewSearchRequest().Index(ei.Config.BrandFullIndex).Query(brandQuery).Size(5).FetchSourceIncludeExclude([]string{"id", "name", "logo"}, nil))

	influencerQuery := elastic.NewMultiMatchQuery(opts.Query, []string{"name.autocomplete"}...).Operator("or").Type("cross_fields")
	mSearchQuery = append(mSearchQuery, elastic.NewSearchRequest().Index(ei.Config.InfluencerFullIndex).Query(influencerQuery).Size(5).FetchSourceIncludeExclude([]string{"id", "name", "profile_image"}, nil))

	var subQuery []elastic.Query
	subQuery = append(subQuery, elastic.NewTermQuery("type", "pebble"))
	subQuery = append(subQuery, elastic.NewTermQuery("is_active", true))
	subQuery = append(subQuery, elastic.NewMultiMatchQuery(opts.Query, []string{"label.interests", "caption", "influencer_info.name.autocomplete", "brand_info.name.autocomplete", "catalog_info.name.autocomplete"}...).Operator("or").Type("cross_fields"))

	contentQuery := elastic.NewBoolQuery().Must(subQuery...)
	mSearchQuery = append(mSearchQuery, elastic.NewSearchRequest().Index(ei.Config.ContentFullIndex).Query(contentQuery).Size(5).FetchSourceIncludeExclude([]string{"name", "id", "media_info", "caption"}, nil))

	resp, err := mSearch.Add(mSearchQuery...).Do(context.TODO())
	if err != nil {
		ei.Logger.Err(err).Msgf("failed to get search result for query:%s", opts.Query)
		return nil, errors.Wrap(err, "failed to get search results")
	}

	var influencer []schema.InfluencerSearchResp
	var brand []schema.BrandSearchResp
	var content []schema.ContentSearchResp
	var catalog []schema.CatalogSearchResp
	for _, result := range resp.Responses {
		for _, hit := range result.Hits.Hits {
			if strings.Contains(hit.Index, ei.Config.BrandFullIndex) {
				var s schema.BrandSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}
				brand = append(brand, s)
			} else if strings.Contains(hit.Index, ei.Config.InfluencerFullIndex) {
				var s schema.InfluencerSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}
				influencer = append(influencer, s)
			} else if strings.Contains(hit.Index, ei.Config.CatalogFullIndex) {
				var s schema.CatalogSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}

				catalog = append(catalog, s)
			} else if strings.Contains(hit.Index, ei.Config.ContentFullIndex) {
				var s schema.ContentSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}
				content = append(content, s)
			}
		}
	}
	res := schema.SearchResp{
		Brand:      brand,
		Influencer: influencer,
		Content:    content,
		Catalog:    catalog,
	}

	return &res, nil
}

func (ei *ElasticsearchImpl) SearchDiscover(opts *schema.SearchOpts) (*schema.DiscoverSearchResp, error) {
	mSearch := elastic.NewMultiSearchService(ei.Client)
	var mSearchQuery []*elastic.SearchRequest

	// Brand search
	brandQuery := elastic.NewMultiMatchQuery(opts.Query, []string{"lname.lname^2", "lname.autocomplete"}...).Fuzziness("2")
	mSearchQuery = append(mSearchQuery, elastic.NewSearchRequest().Index(ei.Config.BrandFullIndex).Query(brandQuery).Size(3).FetchSourceIncludeExclude([]string{"id", "name", "logo"}, nil))

	// Influencer search
	influencerQuery := elastic.NewMultiMatchQuery(opts.Query, []string{"name.autocomplete"}...).Operator("or").Type("best_fields")
	mSearchQuery = append(mSearchQuery, elastic.NewSearchRequest().Index(ei.Config.InfluencerFullIndex).Query(influencerQuery).Size(3).FetchSourceIncludeExclude([]string{"id", "name", "profile_image"}, nil))

	// Pebble search
	var subQuery []elastic.Query
	filterSubQuery := elastic.NewTermQuery("is_active", true)
	filterSubQuery1 := elastic.NewTermQuery("type", "pebble")
	// subQuery = append(subQuery, elastic.NewMultiMatchQuery(opts.Query, []string{"name.autocomplete", "influencer_info.name.autocomplete", "brand_info.name.autocomplete"}...).Operator("or").Type("best_fields"))
	subQuery = append(subQuery, elastic.NewMatchQuery("name.autocomplete", opts.Query))
	subQuery = append(subQuery, elastic.NewNestedQuery("pebble_info", elastic.NewNestedQuery("pebble_info.influencer_info", elastic.NewMatchQuery("name.autocomplete", opts.Query))))
	seriesQuery := elastic.NewBoolQuery().Should(subQuery...).Filter(filterSubQuery, filterSubQuery1)
	mSearchQuery = append(mSearchQuery, elastic.NewSearchRequest().Index(ei.Config.SeriesFullIndex).Query(seriesQuery).Size(3).FetchSourceIncludeExclude([]string{"name", "id", "thumbnail"}, nil))

	// Hashtag search
	suggestQuery := elastic.NewCompletionSuggester("hashtag").SkipDuplicates(true).Field("hashtags.suggest").Prefix(opts.Query).Size(3)
	mSearchQuery = append(mSearchQuery, elastic.NewSearchRequest().Index(ei.Config.ContentFullIndex).Suggester(suggestQuery).FetchSource(false))

	resp, err := mSearch.Add(mSearchQuery...).Do(context.TODO())
	if err != nil {
		ei.Logger.Err(err).Msgf("failed to get search result for query:%s", opts.Query)
		return nil, errors.Wrap(err, "failed to get search results")
	}

	var influencer []schema.InfluencerSearchResp
	var brand []schema.BrandSearchResp
	var series []schema.SeriesSearchResp
	var hashtags []schema.HashtagSearchResp

	for _, result := range resp.Responses {
		for _, hit := range result.Hits.Hits {
			if strings.Contains(hit.Index, ei.Config.BrandFullIndex) {
				var s schema.BrandSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}
				brand = append(brand, s)
			} else if strings.Contains(hit.Index, ei.Config.InfluencerFullIndex) {
				var s schema.InfluencerSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}
				influencer = append(influencer, s)
			} else if strings.Contains(hit.Index, ei.Config.SeriesFullIndex) {
				var s schema.SeriesSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}
				series = append(series, s)
			}
		}

		for _, hit := range result.Suggest["hashtag"] {
			for _, opt := range hit.Options {
				hashtags = append(hashtags, schema.HashtagSearchResp{Text: opt.Text})
			}
		}
	}
	res := schema.DiscoverSearchResp{
		Brand:      brand,
		Influencer: influencer,
		Series:     series,
		Hashtag:    hashtags,
	}

	return &res, nil
}

func (ei *ElasticsearchImpl) SearchCatalog(opts *schema.SearchOpts) (*schema.SearchResp, error) {
	mustQuery := elastic.NewMultiMatchQuery(opts.Query, "keywords.*^3", "name.*^2", "brand_info.name^2").Operator("or").Type("best_fields")
	filterQuery := elastic.NewTermQuery("status.value", model.Publish)
	query := elastic.NewBoolQuery().Must(mustQuery).Filter(filterQuery)
	var fromPage int
	if opts.Page != 0 {
		fromPage = (int(opts.Page) * 20) + 1
	}
	res, err := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).From(fromPage).Size(20).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get products")
		return nil, errors.Wrap(err, "failed to get products")
	}

	// var resp []schema.GetCatalogInfoResp
	var resp []schema.CatalogSearchResp
	for _, hit := range res.Hits.Hits {
		var s schema.CatalogSearchResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode catalog basic json")
		}
		resp = append(resp, s)
	}

	result := schema.SearchResp{
		Catalog: resp,
	}

	return &result, nil
}

func (ei *ElasticsearchImpl) SearchBrand(opts *schema.SearchOpts) ([]schema.BrandSearchResp, error) {
	query := elastic.NewMultiMatchQuery(opts.Query, []string{"lname.lname^2", "lname.autocomplete"}...).Fuzziness("2")
	var fromPage int
	if opts.Page != 0 {
		fromPage = (int(opts.Page) * 10) + 1
	}
	res, err := ei.Client.Search().Index(ei.Config.BrandFullIndex).Query(query).From(fromPage).Size(10).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get brands")
		return nil, errors.Wrap(err, "failed to get brands")
	}

	var resp []schema.BrandSearchResp
	for _, hit := range res.Hits.Hits {
		var s schema.BrandSearchResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode brand search json")
		}
		resp = append(resp, s)
	}

	return resp, nil
}

func (ei *ElasticsearchImpl) SearchInfluencer(opts *schema.SearchOpts) ([]schema.InfluencerSearchResp, error) {
	query := elastic.NewMatchQuery("name.autocomplete", opts.Query)
	var fromPage int
	if opts.Page != 0 {
		fromPage = (int(opts.Page) * 20) + 1
	}
	res, err := ei.Client.Search().Index(ei.Config.InfluencerFullIndex).Query(query).From(fromPage).Size(20).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get influencer")
		return nil, errors.Wrap(err, "failed to get influencer")
	}

	var resp []schema.InfluencerSearchResp
	for _, hit := range res.Hits.Hits {
		var s schema.InfluencerSearchResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode influencer search json")
		}
		resp = append(resp, s)
	}

	return resp, nil
}

func (ei *ElasticsearchImpl) SearchSeries(opts *schema.SearchOpts) ([]schema.SeriesSearchResp, error) {
	query := elastic.NewMatchQuery("name.autocomplete", opts.Query)
	var fromPage int
	if opts.Page != 0 {
		fromPage = (int(opts.Page) * 20) + 1
	}
	res, err := ei.Client.Search().Index(ei.Config.SeriesFullIndex).Query(query).From(fromPage).Size(20).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get series")
		return nil, errors.Wrap(err, "failed to get series")
	}

	var resp []schema.SeriesSearchResp
	for _, hit := range res.Hits.Hits {
		var s schema.SeriesSearchResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode series search json")
		}
		resp = append(resp, s)
	}

	return resp, nil
}

func (ei *ElasticsearchImpl) SearchHashtag(opts *schema.SearchOpts) ([]schema.HashtagSearchResp, error) {
	var resp []schema.HashtagSearchResp
	if opts.Page == 1 {
		return resp, nil
	}
	query := elastic.NewCompletionSuggester("hashtag").SkipDuplicates(true).Field("hashtags.suggest").Prefix(opts.Query).Size(20)
	res, err := ei.Client.Search().Index(ei.Config.ContentFullIndex).FetchSource(false).Suggester(query).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msg("failed to get hashtag")
		return nil, errors.Wrap(err, "failed to get hashtag")
	}

	for _, hit := range res.Suggest["hashtag"] {
		for _, opt := range hit.Options {
			resp = append(resp, schema.HashtagSearchResp{Text: opt.Text})
		}
	}

	return resp, nil
}

func (ei *ElasticsearchImpl) GetReviewsByCatalogID(opts *schema.GetReviewsByCatalogIDFilter) ([]schema.GetReviewsByCatalogIDResp, error) {
	var resp []schema.GetReviewsByCatalogIDResp

	mustQuery := elastic.NewTermQuery("catalog_id", opts.CatalogID)
	filterQuery := elastic.NewTermQuery("is_processed", true)
	query := elastic.NewBoolQuery().Must(mustQuery).Filter(filterQuery)
	res, err := ei.Client.Search().Index(ei.Config.ReviewFullIndex).Query(query).Sort("created_at", false).From(int(opts.Page) * 10).Size(10).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msgf("failed to get reviews for catalog id: %s", opts.CatalogID)
		return nil, errors.Wrap(err, "failed to get search results")
	}
	for _, hit := range res.Hits.Hits {
		var s schema.GetReviewsByCatalogIDResp
		if err := json.Unmarshal(hit.Source, &s); err != nil {
			ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
			return nil, errors.Wrap(err, "failed to decode review")
		}
		resp = append(resp, s)
	}
	return resp, nil
}

func (ei *ElasticsearchImpl) GetCatalogByBrandID(opts *schema.GetCatalogByBrandIDOpts) ([]schema.GetCatalogBasicResp, error) {
	var queries []elastic.Query
	queries = append(queries, elastic.NewTermQuery("status.value", model.Publish))
	queries = append(queries, elastic.NewTermQuery("brand_id", opts.BrandID))
	query := elastic.NewBoolQuery().Must(queries...)
	res, err := ei.Client.Search().Index(ei.Config.CatalogFullIndex).Query(query).From(int(opts.Page * 20)).Size(20).Do(context.Background())
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

func (ei *ElasticsearchImpl) GetCollectionCatalogByIDs(opts *schema.GetCollectionCatalogByIDs) ([]schema.GetCatalogBasicResp, error) {
	var finalResp []schema.GetCatalogBasicResp
	var featureResp []schema.GetCatalogBasicResp

	var wg sync.WaitGroup
	var featErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		if len(opts.FeatIDs) != 0 {
			res, err := ei.GetCatalogByIDs(opts.FeatIDs)
			if err != nil {
				featErr = errors.Wrapf(err, "error getting featured catalog ids")
				return
			}
			featureResp = res
		}
	}()

	catalogResp, err := ei.GetCatalogByIDs(opts.IDs)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting featured catalog ids")
	}

	if featErr != nil {
		return nil, featErr
	}

	wg.Wait()
	finalResp = append(featureResp, catalogResp...)
	return finalResp, nil
}

func (ei *ElasticsearchImpl) ShopSearch(opts *schema.SearchOpts) (*schema.ShopSearchResp, error) {
	mSearch := elastic.NewMultiSearchService(ei.Client)
	var mSearchQuery []*elastic.SearchRequest

	var fromPage int
	if opts.Page != 0 {
		fromPage = (int(opts.Page) * 10) + 1
	}

	// catalog search query
	filterQuery := elastic.NewTermQuery("status.value", model.Publish)
	mustQuery := elastic.NewMultiMatchQuery(
		opts.Query,
		[]string{
			"brand_info.name^2",
			"brand_info.name.autocomplete",
			"name.autocomplete",
			"keywords.keywords",
			"name.name^2"}...,
	).Operator("or").Type("cross_fields")
	boolQuery := elastic.NewBoolQuery().Filter(filterQuery).Must(mustQuery)
	catalogSearchQuery := elastic.NewSearchRequest().Index(ei.Config.CatalogFullIndex).Query(boolQuery).From(fromPage).Size(10).FetchSourceIncludeExclude([]string{"id", "name", "featured_image", "base_price", "retail_price", "discount_info", "variants.id", "brand_info"}, nil)
	mSearchQuery = append(mSearchQuery, catalogSearchQuery)
	// brand search query
	// matchQuery := elastic.NewMultiMatchQuery(opts.Query, []string{"lname.lname^2", "lname.autocomplete"}...).Fuzziness("2")
	// brandSearchQuery := elastic.NewSearchRequest().Index(ei.Config.BrandFullIndex).Query(matchQuery).Size(3).FetchSourceIncludeExclude([]string{"id", "name", "logo", "username"}, nil)

	resp, err := mSearch.Add(mSearchQuery...).Do(context.Background())
	if err != nil {
		ei.Logger.Err(err).Msgf("failed to get shop search result for query:%s", opts.Query)
		return nil, errors.Wrap(err, "failed to get shop search results")
	}

	var brand []schema.BrandSearchResp
	var catalog []schema.CatalogSearchResp
	for _, result := range resp.Responses {
		for _, hit := range result.Hits.Hits {
			if strings.Contains(hit.Index, ei.Config.BrandFullIndex) {
				var s schema.BrandSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}
				brand = append(brand, s)
			} else if strings.Contains(hit.Index, ei.Config.CatalogFullIndex) {
				var s schema.CatalogSearchResp
				if err := json.Unmarshal(hit.Source, &s); err != nil {
					ei.Logger.Err(err).Str("source", string(hit.Source)).Msg("failed to unmarshal struct from json")
					continue
				}
				catalog = append(catalog, s)
			}
		}
	}

	res := schema.ShopSearchResp{
		// Brand:   brand,
		Catalog: catalog,
	}
	return &res, nil
}
