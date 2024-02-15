package readme

import (
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

func (r *Readme) WriteTag(tag, url, date string) {
	r.Write("\n## [")
	r.Write(tag)
	r.Write("](")
	r.Write(url)
	r.Write(") (")
	r.Write(date)
	r.Write(")\n\n")
}

func (r *Readme) WriteType(t string) {
	r.Write("\n### ")
	r.Write(t)
	r.Write("\n\n")
}

func (r *Readme) WriteCommit(c string, hash string, url string) {
	hashSlice := hash[:7]
	r.Write("- ")
	r.Write(c)
	r.Write(" ([")
	r.Write(hashSlice)
	r.Write("](")
	r.Write(url)
	r.Write("/")
	r.Write(hash)
	r.Write("))\n")
}

func Create() *Readme {
	readme := &Readme{}
	readme.Write("# Changelog\n")
	return readme
}
