package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nats-io/nats.go/micro"
	"github.com/wombatwisdom/wombat-builder/internal/store"
	"github.com/wombatwisdom/wombat-builder/public/model"
)

type (
	BuildRequestRequest struct {
		Goos      string   `json:"goos" jsonschema_description:"The target operating system"`
		Goarch    string   `json:"goarch" jsonschema_description:"The target architecture"`
		GoVersion string   `json:"goVersion" jsonschema_description:"The Go version to use"`
		Packages  []string `json:"packages" jsonschema_description:"The packages to build"`
		Force     bool     `json:"force" jsonschema_description:"Whether to force a rebuild"`
	}

	BuildRequestResponse struct {
		Id     string            `json:"id" jsonschema_description:"The ID of the build"`
		Status model.BuildStatus `json:"status" jsonschema_description:"The status of the build"`
	}
)

func (r *BuildRequestRequest) Validate() error {
	if r.Goos == "" {
		return ErrMissingField("goos")
	}

	if r.Goarch == "" {
		return ErrMissingField("goarch")
	}

	if r.GoVersion == "" {
		return ErrMissingField("goVersion")
	}

	if len(r.Packages) == 0 {
		return ErrMissingField("packages")
	}

	return nil
}

func ErrMissingField(s string) error {
	return errors.New("missing required field " + s)
}

func getBuildRequestHandler(s *store.Store) micro.HandlerFunc {
	return func(request micro.Request) {
		var req BuildRequestRequest
		if err := json.Unmarshal(request.Data(), &req); err != nil {
			_ = request.Error("BAD_REQUEST", "failed to parse request", []byte(err.Error()))
			return
		}

		if err := req.Validate(); err != nil {
			_ = request.Error("BAD_REQUEST", "invalid request", []byte(err.Error()))
			return
		}

		// -- create a build out of the request
		build, err := model.NewBuild(
			model.WithGoVersion(req.GoVersion),
			model.WithGoos(req.Goos),
			model.WithGoarch(req.Goarch),
			model.WithPackageUrls(req.Packages...),
		)
		if err != nil {
			_ = request.Error("BAD_REQUEST", "invalid request", []byte(err.Error()))
			return
		}

		// -- check if the build already exists
		existing, err := s.Builds.Get(context.Background(), build.Id())
		if err != nil {
			_ = request.Error("BACKBONE_ERROR", "failed to check if build exists", []byte(err.Error()))
			return
		}

		if !req.Force && existing != nil {
			result := BuildRequestResponse{
				Id:     existing.Id(),
				Status: existing.Status,
			}
			if err := request.RespondJSON(result); err != nil {
				_ = request.Error("INTERNAL_ERROR", "failed to respond", []byte(err.Error()))
				return
			}
		} else {
			// -- store the build
			_, err := s.Builds.Set(context.Background(), build)
			if err != nil {
				_ = request.Error("BACKBONE_ERROR", "failed to store build", []byte(err.Error()))
				return
			}

			result := BuildRequestResponse{
				Id:     build.Id(),
				Status: build.Status,
			}
			if err := request.RespondJSON(result); err != nil {
				_ = request.Error("INTERNAL_ERROR", "failed to respond", []byte(err.Error()))
				return
			}
		}
	}
}
