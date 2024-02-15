package main

import (
	"os"
	"os/exec"
	"strings"

	"gochangelog/pkg/config"
	"gochangelog/pkg/git"
)

type Config struct {
	Yaml *config.YAML
	Tags []string
	Logs []string
}

func Reverse[T any](data []T) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func main() {
	yaml, err := config.Read()
	tags, err := git.GetTags()

	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	builder.WriteString("# Changelog\n")
	for index, tag := range tags {
		builder.WriteString("## " + tag.Tag + " " + tag.Date)
		builder.WriteString("\n\n")

		var t1 string
		var t2 string

		if index == len(tags)-1 {
			t1 = tags[index].Tag
		} else {
			t1 = tags[index].Tag + "^"
		}

		if index == 0 {
			t2 = "HEAD"
		} else {
			t2 = tags[index-1].Tag + "^"
		}

		commits, err := git.GetCommits(t1, t2, yaml.Filter)
		if err != nil {
			panic(err)
		}

		sortedCommits := git.SortCommits(commits)
		for commitType, logs := range sortedCommits {
			builder.WriteString("\n")
			commitTypeName := git.GetCommitTypeName(commitType)
			builder.WriteString("### ")
			builder.WriteString(commitTypeName)
			builder.WriteString("\n\n")

			for _, log := range logs {
				builder.WriteString("- ")
				builder.WriteString(log)
				builder.WriteByte('\n')
			}
		}

		builder.WriteString("\n")
	}

	cmd := exec.Command("echo", builder.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}
