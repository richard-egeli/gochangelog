package readme

import "strings"

type Readme struct {
	builder strings.Builder
}

func (r *Readme) Write(s string) {

}

func (r *Readme) WriteTag(current, next, date string) {
	r.builder.WriteString("## ")
	r.builder.WriteString(current)
	r.builder.WriteString(" (")
	r.builder.WriteString(date)
	r.builder.WriteString(")\n\n")
}

func Create() *Readme {
	return &Readme{}
}
