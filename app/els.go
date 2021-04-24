package app

import (
	"context"
	"encoding/json"
	"fmt"
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
	GetPebble(opts *schema.GetPebbleFilter) ([]schema.GetPebbleESResp, error)
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
	res, err := ei.Client.Search().Index(ei.Config.ContentFullIndex).From(from).Size(10).SearchSource(builder).Do(context.Background())
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
