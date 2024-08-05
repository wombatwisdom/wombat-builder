package docgen

import (
	"github.com/redpanda-data/benthos/v4/public/bloblang"
	"github.com/redpanda-data/benthos/v4/public/service"
	"github.com/rs/zerolog/log"
	"github.com/wombatwisdom/wombat-builder/library"
)

type Kind string

const (
	KindInput     Kind = "input"
	KindOutput    Kind = "output"
	KindCache     Kind = "cache"
	KindBuffer    Kind = "buffer"
	KindProcessor Kind = "processor"
	KindRateLimit Kind = "rate_limit"
	KindTracer    Kind = "tracer"
	KindScanner   Kind = "scanner"
	KindMetric    Kind = "metric"

	KindFunction Kind = "function"
	KindMethod   Kind = "method"
)

type Snapshot struct {
	Inputs     map[string]struct{}
	Outputs    map[string]struct{}
	Caches     map[string]struct{}
	Buffers    map[string]struct{}
	Processors map[string]struct{}
	RateLimits map[string]struct{}
	Tracers    map[string]struct{}
	Scanners   map[string]struct{}
	Metrics    map[string]struct{}

	Functions map[string]struct{}
	Methods   map[string]struct{}
}

func NewSnapshot() *Snapshot {
	snapshot := &Snapshot{
		Inputs:     map[string]struct{}{},
		Outputs:    map[string]struct{}{},
		Caches:     map[string]struct{}{},
		Buffers:    map[string]struct{}{},
		Processors: map[string]struct{}{},
		RateLimits: map[string]struct{}{},
		Tracers:    map[string]struct{}{},
		Scanners:   map[string]struct{}{},
		Metrics:    map[string]struct{}{},
		Functions:  map[string]struct{}{},
		Methods:    map[string]struct{}{},
	}

	env := service.GlobalEnvironment()
	env.WalkBuffers(snapshotWalker(snapshot, KindBuffer))
	env.WalkCaches(snapshotWalker(snapshot, KindCache))
	env.WalkInputs(snapshotWalker(snapshot, KindInput))
	env.WalkOutputs(snapshotWalker(snapshot, KindOutput))
	env.WalkProcessors(snapshotWalker(snapshot, KindProcessor))
	env.WalkRateLimits(snapshotWalker(snapshot, KindRateLimit))
	env.WalkTracers(snapshotWalker(snapshot, KindTracer))
	env.WalkScanners(snapshotWalker(snapshot, KindScanner))
	env.WalkMetrics(snapshotWalker(snapshot, KindMetric))

	benv := bloblang.GlobalEnvironment()
	benv.WalkFunctions(snapshotFunctionWalker(snapshot))
	benv.WalkMethods(snapshotMethodWalker(snapshot))

	return snapshot
}

func (s *Snapshot) WalkDelta(fn func(doc *library.DocView)) {
	benv := bloblang.GlobalEnvironment()
	benv.WalkFunctions(deltaFunctionWalker(s.Functions, fn))
	benv.WalkMethods(deltaMethodWalker(s.Methods, fn))

	env := service.GlobalEnvironment()
	env.WalkBuffers(deltaWalker(s.Buffers, fn))
	env.WalkCaches(deltaWalker(s.Caches, fn))
	env.WalkInputs(deltaWalker(s.Inputs, fn))
	env.WalkOutputs(deltaWalker(s.Outputs, fn))
	env.WalkProcessors(deltaWalker(s.Processors, fn))
	env.WalkRateLimits(deltaWalker(s.RateLimits, fn))
	env.WalkTracers(deltaWalker(s.Tracers, fn))
	env.WalkScanners(deltaWalker(s.Scanners, fn))
	env.WalkMetrics(deltaWalker(s.Metrics, fn))
}

func snapshotWalker(snapshot *Snapshot, kind Kind) func(name string, config *service.ConfigView) {
	return func(name string, config *service.ConfigView) {
		switch kind {
		case KindInput:
			snapshot.Inputs[name] = struct{}{}
		case KindOutput:
			snapshot.Outputs[name] = struct{}{}
		case KindCache:
			snapshot.Caches[name] = struct{}{}
		case KindBuffer:
			snapshot.Buffers[name] = struct{}{}
		case KindProcessor:
			snapshot.Processors[name] = struct{}{}
		case KindRateLimit:
			snapshot.RateLimits[name] = struct{}{}
		case KindTracer:
			snapshot.Tracers[name] = struct{}{}
		case KindScanner:
			snapshot.Scanners[name] = struct{}{}
		case KindMetric:
			snapshot.Metrics[name] = struct{}{}
		}
	}
}

func snapshotFunctionWalker(snapshot *Snapshot) func(name string, spec *bloblang.FunctionView) {
	return func(name string, spec *bloblang.FunctionView) {
		snapshot.Functions[name] = struct{}{}
	}
}

func snapshotMethodWalker(snapshot *Snapshot) func(name string, spec *bloblang.MethodView) {
	return func(name string, spec *bloblang.MethodView) {
		snapshot.Methods[name] = struct{}{}
	}
}

func deltaMethodWalker(exclusions map[string]struct{}, fn func(doc *library.DocView)) func(name string, spec *bloblang.MethodView) {
	return func(name string, spec *bloblang.MethodView) {
		if _, fnd := exclusions[name]; fnd {
			return
		} else {
			exclusions[name] = struct{}{}
		}

		dv, err := library.ParseDocViewFromMethodView(spec)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse method %s", name)
			return
		}

		fn(dv)
	}
}

func deltaFunctionWalker(exclusions map[string]struct{}, fn func(doc *library.DocView)) func(name string, spec *bloblang.FunctionView) {
	return func(name string, spec *bloblang.FunctionView) {
		if _, fnd := exclusions[name]; fnd {
			return
		} else {
			exclusions[name] = struct{}{}
		}

		dv, err := library.ParseDocViewFromFunctionView(spec)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse function %s", name)
			return
		}

		fn(dv)
	}
}

func deltaWalker(exclusions map[string]struct{}, fn func(doc *library.DocView)) func(name string, config *service.ConfigView) {
	return func(name string, config *service.ConfigView) {
		if _, fnd := exclusions[name]; fnd {
			return
		} else {
			exclusions[name] = struct{}{}
		}

		dv, err := library.ParseDocViewFromConfigView(config)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to parse component %s", name)
			return
		}

		fn(dv)
	}
}
