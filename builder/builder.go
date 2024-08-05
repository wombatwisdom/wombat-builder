package builder

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/wombatwisdom/wombat-builder/library"

	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"path/filepath"
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

type templateVars struct {
	GoVersion string
	Path      string
	Module    string
}

type Builder interface {
	BuildPackage(ctx context.Context, spec PackageBuildSpec, goExec string) error
}

func NewBuilder(lc library.Client) Builder {
	return &baseBuilder{lc: lc}
}

type baseBuilder struct {
	lc library.Client
}

func (bb *baseBuilder) BuildPackage(ctx context.Context, spec PackageBuildSpec, goExec string) error {
	logger := log.With().Str("library", spec.Library).Str("version", spec.Version).Str("package", spec.Package).Logger()

	// get the library
	lib, err := bb.lc.Library(spec.Library)
	if err != nil {
		return fmt.Errorf("failed to get library: %w", err)
	}

	dir := path.Join(os.TempDir(), xid.New().String())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			logger.Error().Err(err).Msg("failed to clean temp dir")
		}
	}()

	c := InDir(dir, goExec)
	gover, err := c.GoVersion(ctx)

	// -- build the template vars
	vars := templateVars{
		GoVersion: gover,
		Module:    spec.Module,
		Path:      spec.Path,
	}

	if err := bb.generate(vars, dir); err != nil {
		return fmt.Errorf("failed to generate module files: %w", err)
	}

	// -- add the module
	if err := c.GoGet(ctx, fmt.Sprintf("%s@%s", lib.Module, spec.Version)); err != nil {
		return fmt.Errorf("failed to get module: %w", err)
	}

	logger.Debug().Msg("pulling in module imports")
	if err := c.GoModTidy(ctx); err != nil {
		return fmt.Errorf("failed to tidy go modules: %w", err)
	}

	targetFilename := path.Join(dir, spec.Library, spec.Version, spec.Package, fmt.Sprintf("%s_%s_%s_%s_%s_go%s.so", spec.Library, spec.Package, spec.Version, spec.Os, spec.Arch, gover))

	logger.Info().Msg("building shared object file")
	if err := c.GoBuild(ctx, spec.Os, spec.Arch, targetFilename); err != nil {
		return fmt.Errorf("failed to build plugin: %w", err)
	}

	// -- upload the result
	f, err := os.Open(targetFilename)
	if err != nil {
		return fmt.Errorf("failed to open artifact file: %w", err)
	}
	defer f.Close()
	if err := bb.lc.UploadArtifact(spec.Library, spec.Version, spec.Package, spec.Os, spec.Arch, gover, f); err != nil {
		return fmt.Errorf("failed to upload artifact: %w", err)
	}

	return nil
}

func (bb *baseBuilder) generate(vars templateVars, dir string) error {
	// -- clean the directory if it exists
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to clean temp dir: %w", err)
	}

	for f, tmpl := range templates {
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
		if err := outTemplate.Execute(outFile, vars); err != nil {
			return fmt.Errorf("failed to execute %v template: %w", f, err)
		}
		if err := outFile.Close(); err != nil {
			return fmt.Errorf("failed to close %v file: %w", f, err)
		}
	}

	return nil
}
