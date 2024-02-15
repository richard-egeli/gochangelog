package provider

import "gochangelog/pkg/config"

type Provider interface {
	Diff(prev, new string, config *config.YAML) string
}

type Type string

const (
	BITBUCKET Type = "bitbucket"
)

func Get(p Type) Provider {
	switch p {
	case BITBUCKET:
		return &Bitbucket{}
	}

	panic("Provider is not known, you must define a valid 'provider' in the config")
}
