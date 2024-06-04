package builder

import (
  "context"
  _ "embed"
  "fmt"
  "github.com/benthosdev/benthos-builder/internal/command"
  "github.com/benthosdev/benthos-builder/public/model"
  "github.com/rs/zerolog"
  "github.com/rs/zerolog/log"
  "os"
  "path"
  "path/filepath"
  "sort"
  "text/template"
)

//go:embed templates/go.mod.template
var goModTemplate string

//go:embed templates/main.go.template
var mainGoTemplate string

var templates = map[string]string{
  "main.go": mainGoTemplate,
  "go.mod":  goModTemplate,
}

type BuildTask struct {
  *model.Build
}

func (t *BuildTask) Run(ctx context.Context) (string, error) {
  logger := log.With().Str("build", t.Id()).Logger()
  defer func() {
    logger.Info().Msg("build finished")
  }()

  logger.Debug().Msg("sorting imports")
  sort.Slice(t.Packages, func(i, j int) bool {
    return t.Packages[i].Url < t.Packages[j].Url
  })

  logger.Debug().Msg("creating temp dir")
  dir := path.Join(os.TempDir(), "wombat-builder", t.Id())
  if err := os.MkdirAll(dir, 0755); err != nil {
    return "", fmt.Errorf("failed to create temp dir: %w", err)
  }

  logger.Info().Msgf("building in %s", dir)
  if err := t.generate(dir, &logger); err != nil {
    return "", fmt.Errorf("failed to generate module files: %w", err)
  }

  c := command.InDir(dir)
  logger.Info().Msg("pulling in module imports")
  if err := c.GoModTidy(ctx); err != nil {
    return "", fmt.Errorf("failed to tidy go modules: %w", err)
  }

  logger.Info().Msg("building wombat")
  if err := c.GoBuild(ctx); err != nil {
    return "", fmt.Errorf("failed to build wombat: %w", err)
  }

  return path.Join(dir, "wombat"), nil
}

func (t *BuildTask) generate(dir string, logger *zerolog.Logger) error {
  // -- clean the directory if it exists
  if err := os.RemoveAll(dir); err != nil {
    return fmt.Errorf("failed to clean temp dir: %w", err)
  }

  logger.Info().Msg("generating module files")
  for f, tmpl := range templates {
    logger.Debug().Msgf("generating %v", f)

    if err := os.MkdirAll(path.Dir(filepath.Join(dir, f)), 0755); err != nil {
      return fmt.Errorf("failed to create dir: %w", err)
    }

    outFile, err := os.Create(filepath.Join(dir, f))
    if err != nil {
      return fmt.Errorf("failed to create %v: %w", f, err)
    }
    outTemplate, err := template.New(f).Parse(tmpl)
    if err != nil {
      return fmt.Errorf("failed to initialise %v template: %w", f, err)
    }
    if err := outTemplate.Execute(outFile, t); err != nil {
      return fmt.Errorf("failed to execute %v template: %w", f, err)
    }
    if err := outFile.Close(); err != nil {
      return fmt.Errorf("failed to close %v file: %w", f, err)
    }
  }

  return nil
}
