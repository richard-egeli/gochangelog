package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gochangelog/pkg/config"
)

type CommitType string

const (
	FIX      CommitType = "fix"
	FEAT     CommitType = "feat"
	BREAKING CommitType = "BREAKING CHANGE"
	BUILD    CommitType = "build"
	CHORE    CommitType = "chore"
	CI       CommitType = "ci"
	DOCS     CommitType = "docs"
	STYLE    CommitType = "style"
	REFACTOR CommitType = "refactor"
	PERF     CommitType = "perf"
	TEST     CommitType = "test"
)

type Tag struct {
	Raw  string
	Tag  string
	Date string
}

type Config struct {
	Yaml *config.YAML
	Tags []string
	Logs []string
}

func GetFullCommitType(c CommitType) string {
	switch c {
	case FIX:
		return "Bug Fixes"
	case FEAT:
		return "Features"
	case BREAKING:
		return "Breaking Changes"
	case BUILD:
		return "Build"
	case CHORE:
		return "Chore"
	case CI:
		return "Configuration"
	case DOCS:
		return "Documentation"
	case STYLE:
		return "Styling"
	case REFACTOR:
		return "Refactor"
	case PERF:
		return "Performance"
	case TEST:
		return "Tests"
	default:
		return "Other"
	}
}

func GetLogs(t1 string, t2 string, filters []string) []string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("git", "log", "--pretty=format:[%cd] %s ([%h]())", "--date=format:%d-%m-%Y", "--no-walk", t1+"..."+t2)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(fmt.Println(fmt.Sprint(err) + ": " + stderr.String()))
	}

	lines := strings.Split(stdout.String(), "\n")
	filteredLines := []string{}

	for _, line := range lines {
		for _, filter := range filters {
			if strings.Contains(line, filter+":") {
				filteredLines = append(filteredLines, line)
				break
			}
		}
	}

	return filteredLines
}

func GetTags() ([]Tag, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var tags []Tag

	cmd := exec.Command("git", "for-each-ref", "--sort=-committerdate", "refs/tags", "--format=%(committerdate:short) | %(refname)")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		var tag Tag

		tag.Raw = line
		line = strings.Replace(line, "refs/tags/", "", -1)
		data := strings.Split(line, "|")
		if len(data) < 2 {
			continue
		}

		tag.Date = strings.Trim(data[0], " ")
		tag.Tag = strings.Trim(data[1], " ")
		tags = append(tags, tag)
	}

	return tags, nil
}

func Reverse[T any](data []T) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func SortCommits(commits []string) map[CommitType][]string {
	result := make(map[CommitType][]string)

	commitTypes := []CommitType{
		FIX, FEAT, BREAKING, BUILD, CHORE, CI, DOCS, STYLE, REFACTOR, PERF, TEST,
	}

	for _, line := range commits {
		for _, commitType := range commitTypes {
			typeWithBreak := commitType + ":"
			line = strings.ToLower(line)

			if strings.Contains(line, string(typeWithBreak)) {
				line = strings.Replace(line, " "+string(typeWithBreak), "", -1)
				result[commitType] = append(result[commitType], line)
			}
		}
	}

	return result
}

func main() {
	yaml, err := config.Read()
	tags, err := GetTags()

	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	builder.WriteString("Changelog\n")
	for index, tag := range tags {
		builder.WriteString("## " + tag.Tag + " " + tag.Date)
		builder.WriteString("\n\n")
		if index-1 < 0 {
			log.Println(tag)
			continue
		}

		t1 := tags[index-1].Tag
		t2 := tags[index].Tag

		t1 += "^"
		if index != len(tags)-1 {
			t2 += "^"
		} else {
			t2 = "HEAD"
		}

		sortedCommits := SortCommits(GetLogs(t2, t1, yaml.Filter))

		for commitType, logs := range sortedCommits {
			builder.WriteString("\n")
			commitTypeMessage := GetFullCommitType(commitType)
			builder.WriteString("### ")
			builder.WriteString(commitTypeMessage)
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
