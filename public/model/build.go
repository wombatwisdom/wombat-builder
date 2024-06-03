package model

import (
  "fmt"
  "github.com/mitchellh/hashstructure/v2"
  "github.com/rs/zerolog/log"
  "sort"
  "strings"
)

type (
  BuildIdentity struct {
    Goos      string    `json:"goos"`
    Goarch    string    `json:"goarch"`
    GoVersion string    `json:"goversion"`
    Packages  []Package `json:"packages"`
  }
  Build struct {
    BuildIdentity
    Artifact ArtifactReference `json:"artifact,omitempty"`
    Builder  string            `json:"builder,omitempty"`
    Status   BuildStatus       `json:"status"`
    Error    string            `json:"error,omitempty"`
  }

  ArtifactReference string
  BuildStatus       string

  BuildOpt func(*Build)
)

const (
  BuildStatusNew      BuildStatus = "new"
  BuildStatusBuilding BuildStatus = "building"
  BuildStatusSuccess  BuildStatus = "success"
  BuildStatusFailed   BuildStatus = "failed"
)

func WithPackage(pkg ...Package) BuildOpt {
  return func(b *Build) {
    b.Packages = append(b.Packages, pkg...)
  }
}

func WithPackageUrls(pkg ...string) BuildOpt {
  return func(b *Build) {
    for _, p := range pkg {
      b.Packages = append(b.Packages, Package{Url: p})
    }
  }
}

func WithGoVersion(version string) BuildOpt {
  return func(b *Build) {
    b.GoVersion = version
  }
}

func WithGoos(goos string) BuildOpt {
  return func(b *Build) {
    b.Goos = goos
  }
}

func WithGoarch(goarch string) BuildOpt {
  return func(b *Build) {
    b.Goarch = goarch
  }
}

func NewBuild(opts ...BuildOpt) (*Build, error) {
  b := &Build{
    Status: BuildStatusNew,
  }

  for _, opt := range opts {
    opt(b)
  }

  return b, nil
}

func (b Build) Id() string {
  sort.Slice(b.Packages, func(i, j int) bool {
    return b.Packages[i].Url < b.Packages[j].Url
  })

  // -- create a hash for the build
  hash, err := hashstructure.Hash(b.Packages, hashstructure.FormatV2, nil)
  if err != nil {
    log.Panic().Err(err).Msg("failed to hash build")
  }

  return fmt.Sprintf("build.%s.%s.%s.%x", b.Goos, b.Goarch, strings.ReplaceAll(b.GoVersion, ".", "_"), hash)
}
