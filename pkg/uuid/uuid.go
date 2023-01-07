package uuid

import (
	uu "github.com/V-H-R-Oliveira/simple-uuid/uuid"
)

type UUID struct {
	options
}

type options struct {
	Version   int
	Variant   string
	Namespace string
	Name      string
}

const (
	defaultVersion = 1
	defaultVariant = "dce"
)

var supportedVariants = map[string]bool{
	"dce":       true,
	"microsoft": true,
	"future":    true,
}

type Option func(UUID) UUID

func WithVersion(v int) Option {
	return func(o UUID) UUID {
		o.Version = v
		return o
	}
}

func WithVariant(v string) Option {
	return func(o UUID) UUID {
		o.Variant = v
		return o
	}
}

func WithNamespace(v string) Option {
	return func(o UUID) UUID {
		o.Namespace = v
		return o
	}
}

func WithName(v string) Option {
	return func(o UUID) UUID {
		o.Name = v
		return o
	}
}

func New(opts ...Option) (*UUID, error) {
	u := UUID{}
	u.withOpts(opts)

	// Generate a UUID to confirm options are valid
	_, err := uu.NewUUID(u.Version, u.toArgs())
	if err != nil {
		return nil, err
	}
	return &u, err
}

func (u *UUID) withOpts(opts []Option) {
	for _, opt := range opts {
		*u = opt(*u)
	}
	if u.Version < 1 || u.Version > 4 {
		u.Version = defaultVersion
	}
	if !supportedVariants[u.Variant] {
		u.Variant = defaultVariant
	}
}

func (u UUID) Generate() string {
	v, err := uu.NewUUID(4, u.toArgs())
	if err != nil {
		return ""
	}
	return v.Stringify()
}

func (u UUID) toArgs() map[string]string {
	args := make(map[string]string, 0)
	if u.Name != "" {
		args["name"] = u.Name
	}
	if u.Namespace != "" {
		args["namespace"] = u.Namespace
	}
	args["variant"] = u.Variant

	return args
}
