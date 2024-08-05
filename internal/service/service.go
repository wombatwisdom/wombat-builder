package service

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
	"github.com/rs/zerolog/log"
	"github.com/wombatwisdom/wombat-builder/internal/shared"
	"github.com/wombatwisdom/wombat-builder/internal/store"
)

func NewService(nc *nats.Conn, s *store.Store) (*Service, error) {
	return &Service{
		nc: nc,
		s:  s,
	}, nil
}

type Service struct {
	nc *nats.Conn
	s  *store.Store
}

func (s *Service) Run(ctx context.Context) error {
	cfg := micro.Config{
		Name:        "build",
		Version:     "0.1.0",
		Description: "The build service for the wombat project",
		DoneHandler: func(service micro.Service) {
			log.Info().Msg("service stopped")
		},
		ErrorHandler: func(service micro.Service, natsError *micro.NATSError) {
			log.Error().Err(natsError).Msg("service error")
		},
	}

	svc, err := micro.AddService(s.nc, cfg)
	if err != nil {
		return err
	}

	buildGrp := svc.AddGroup("build")
	registerEndpointOrDie(buildGrp, "request", getBuildRequestHandler(s.s), micro.WithEndpointMetadata(map[string]string{
		"description":     "Request a build",
		"request-schema":  shared.SchemaForOrDie(&BuildRequestRequest{}),
		"response-schema": shared.SchemaForOrDie(&BuildRequestResponse{}),
	}))

	registerEndpointOrDie(buildGrp, "list", getBuildListHandler(s.s), micro.WithEndpointMetadata(map[string]string{
		"description":     "List builds",
		"request-schema":  shared.SchemaForOrDie(&BuildListRequest{}),
		"response-schema": shared.SchemaForOrDie(&BuildListResponse{}),
	}))

	log.Info().Msgf("service started: %v", svc.Info().ID)

	// -- wait for the context to complete
	select {
	case <-ctx.Done():
		return nil
	}
}

func registerEndpointOrDie(gr micro.Group, name string, handler micro.Handler, opts ...micro.EndpointOpt) {
	err := gr.AddEndpoint(name, handler, opts...)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to add %s endpoint", name)
	}
}
