package shared

import (
  "github.com/invopop/jsonschema"
  "github.com/rs/zerolog/log"
)

func SchemaForOrDie(a any) string {
  b, err := jsonschema.Reflect(a).MarshalJSON()
  if err != nil {
    log.Panic().Err(err).Msg("failed to marshal schema")
  }
  return string(b)
}
