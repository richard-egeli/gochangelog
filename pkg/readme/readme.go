package readme

import (
	"gochangelog/pkg/config"
	"gochangelog/pkg/git"
	"regexp"
	"strings"
)

type Readme struct {
	builder strings.Builder
}

func (r *Readme) Write(s string) {
	_, err := r.builder.WriteString(s)
	if err != nil {
		panic(err)
	}
}

func (r *Readme) String() string {
	input := r.builder.String()
	re := regexp.MustCompile(`(\n){3,}`)
	output := re.ReplaceAllString(input, "\n\n")

	return output
}

func (r *Readme) WriteTag(prev, cur, date string, config *config.YAML) {
	host := strings.Trim(config.RepoURL, "/")
	path := strings.Trim(config.ComparePath, "/")
	url := host + "/" + path + "/"

	r.Write("\n## [")
	r.Write(prev)
	r.Write("](")
	r.Write(url)
	r.Write(prev)
	r.Write("%0D")
	r.Write(cur)
	r.Write(") (")
	r.Write(date)
	r.Write(")\n\n")
}

func (r *Readme) WriteType(t string) {
	r.Write("\n### ")
	r.Write(t)
	r.Write("\n\n")
}

func (r *Readme) WriteCommit(commit *git.Commit, config *config.YAML) {
	hashSlice := commit.Hash[:7]
	host := strings.Trim(config.RepoURL, "/")
	path := strings.Trim(config.CommitPath, "/")
	url := host + "/" + path + "/"

	r.Write("- ")
	r.Write(commit.Message)
	r.Write(" ([")
	r.Write(hashSlice)
	r.Write("](")
	r.Write(url)
	r.Write(commit.Hash)
	r.Write("))\n")
}

func Create() *Readme {
	readme := &Readme{}
	readme.Write("# Changelog\n")
	return readme
}
