package rollup

import (
	"errors"

	"github.com/ethereum/go-ethereum/eth/tracers"
)

func init() {
	tracers.RegisterLookup(false, lookup)
}

var ctors map[string]func() tracers.Tracer

// register is used by native tracers to register their presence.
func register(name string, ctor func() tracers.Tracer) {
	if ctors == nil {
		ctors = make(map[string]func() tracers.Tracer)
	}
	ctors[name] = ctor
}

// lookup returns a tracer, if one can be matched to the given name.
func lookup(name string, ctx *tracers.Context) (tracers.Tracer, error) {
	if ctors == nil {
		ctors = make(map[string]func() tracers.Tracer)
	}
	if ctor, ok := ctors[name]; ok {
		return ctor(), nil
	}
	return nil, errors.New("no tracer found")
}
