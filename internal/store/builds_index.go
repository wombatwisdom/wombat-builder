package store

import (
  "context"
  "encoding/json"
  "github.com/benthosdev/benthos-builder/public/model"
  "github.com/blevesearch/bleve/v2"
  "github.com/nats-io/nats.go/jetstream"
  "github.com/rs/zerolog/log"
)

func NewBuildIndex(ctx context.Context, js jetstream.JetStream) (*BuildIndex, error) {
  index, err := bleve.Open("builds.bleve")
  if err != nil {
    index, err = bleve.New("builds.bleve", bleve.NewIndexMapping())
    if err != nil {
      return nil, err
    }
  }

  kv, err := js.KeyValue(ctx, JetstreamKVBuilds)
  if err != nil {
    return nil, err
  }

  watcher, err := kv.WatchAll(ctx)
  if err != nil {
    return nil, err
  }

  go func() {
    for {
      select {
      case <-ctx.Done():
        return
      case upd := <-watcher.Updates():
        if upd == nil {
          continue
        }

        id := string(upd.Key())

        if upd.Operation() == jetstream.KeyValueDelete {
          if err := index.Delete(id); err != nil {
            // -- log error
            log.Warn().Err(err).Msgf("failed to delete %s from index", id)
          }
        } else if upd.Operation() == jetstream.KeyValuePut {
          var build model.Build
          if err := json.Unmarshal(upd.Value(), &build); err != nil {
            log.Warn().Err(err).Msgf("failed to unmarshal %s", id)
            continue
          }

          // -- update the index
          if err := index.Index(id, build); err != nil {
            // -- log error
            log.Warn().Err(err).Msgf("failed to index %s", id)
          }
        }
      }
    }
  }()

  return &BuildIndex{index: index}, nil
}

type BuildIndex struct {
  index bleve.Index
}

func (b *BuildIndex) Search(query string) ([]string, error) {
  q := bleve.NewQueryStringQuery(query)
  srch := bleve.NewSearchRequest(q)
  srch.Fields = []string{"id"}

  res, err := b.index.Search(srch)
  if err != nil {
    return nil, err
  }

  var ids []string
  for _, hit := range res.Hits {
    ids = append(ids, hit.ID)
  }

  return ids, nil
}
