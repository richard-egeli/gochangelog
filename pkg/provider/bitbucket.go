package provider

import (
	"gochangelog/pkg/config"
	"strings"
)

type Bitbucket struct{}

func (b *Bitbucket) Diff(prev, next string, config *config.YAML) string {
	var builder strings.Builder

	builder.WriteString(config.RepoURL)
	builder.WriteString("/branches/compare/")
	builder.WriteString(prev)
	builder.WriteString("%0D")
	builder.WriteString(next)
	builder.WriteString("#diff")

	return builder.String()
}

func (b *Bitbucket) Commit(url string) string {
	return url + "/commits"
}
