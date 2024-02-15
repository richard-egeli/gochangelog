package readme

import "strings"

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
	return r.builder.String()
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

func (r *Readme) WriteCommit(c string) {
	r.Write("- ")
	r.Write(c)
	r.Write("\n")
}

func Create() *Readme {
	readme := &Readme{}
	readme.Write("# Changelog\n")
	return readme
}
