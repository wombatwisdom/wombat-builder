package docgen

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/wombatwisdom/wombat-builder/library"
	"io"
	"os"
	"plugin"
	"runtime"
	"strings"
)

type DocsGenerator interface {
	GenerateForPackage(ctx context.Context, libId string, verId string, pkgId string) error
}

func NewDocsGenerator(lc library.Client) DocsGenerator {
	return &baseDocsGenerator{lc: lc}
}

type baseDocsGenerator struct {
	lc library.Client
}

func (b *baseDocsGenerator) GenerateForPackage(ctx context.Context, libId string, verId string, pkgId string) error {
	// -- download the plugin artifact from the library
	r, err := b.lc.Download(libId, verId, pkgId, runtime.GOOS, runtime.GOARCH, strings.TrimPrefix(runtime.Version(), "go"))
	if err != nil {
		return fmt.Errorf("failed to download artifact for %s/%s/%s (%s/%s): %w", libId, verId, pkgId, runtime.GOOS, runtime.GOARCH, err)
	}
	defer r.Close()

	// -- write the artifact to a tmp file
	f, err := os.CreateTemp(os.TempDir(), "artifact-")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()
	fn := f.Name()

	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("failed to write artifact to temp file: %w", err)
	}

	// -- create a snapshot before the plugin is added
	snapshot := NewSnapshot()

	// load the plugin
	if _, err := plugin.Open(fn); err != nil {
		return fmt.Errorf("failed to load plugin %s: %w", fn, err)
	}

	// crawl the environment
	snapshot.WalkDelta(func(doc *library.DocView) {
		if err := b.lc.UploadDocs(libId, verId, pkgId, *doc); err != nil {
			fmt.Printf("failed to upload doc: %v\n", err)
		}
	})

	return nil
}
