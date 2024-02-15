package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	FIX      string = "fix"
	FEAT     string = "feat"
	BREAKING string = "BREAKING CHANGE"
	BUILD    string = "build"
	CHORE    string = "chore"
	CI       string = "ci"
	DOCS     string = "docs"
	STYLE    string = "style"
	REFACTOR string = "refactor"
	PERF     string = "perf"
	TEST     string = "test"
	OTHER    string = "other"
)

type Commit struct {
	Date    string
	Type    string
	Message string
	Hash    string
}

type Tag struct {
	Raw  string
	Tag  string
	Date string
}

var commitTypes = []string{
	FIX, FEAT, BREAKING, BUILD, CHORE, CI, DOCS, STYLE, REFACTOR, PERF, TEST,
}

func GetCommitTypeName(c string) string {
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
		return "Chores"
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

func GetCommits(t1 string, t2 string, filters []string) ([]Commit, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("git", "log", "--pretty=format:%cd <:::> %s <:::> %H", "--date=format:%d-%m-%Y", t1+".."+t2)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprint(err) + ": " + stderr.String())
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

	var commits []Commit
	for _, line := range filteredLines {
		split := strings.Split(line, "<:::>")
		if len(split) < 3 {
			continue
		}

		commit := Commit{
			Date:    strings.Trim(split[0], " "),
			Message: strings.Trim(split[1], " "),
			Hash:    strings.Trim(split[2], " "),
		}

		for _, commitType := range commitTypes {
			m := strings.ToLower(commit.Message)
			t := strings.ToLower(commitType) + ":"

			if strings.Contains(m, t) {
				commit.Message = strings.Trim(strings.ReplaceAll(m, t, ""), " ")
				commit.Type = commitType
				break
			}
		}

		if len(commit.Type) <= 0 {
			commit.Type = OTHER
		}

		commits = append(commits, commit)
	}

	return commits, nil
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

func SortCommits(commits []Commit) map[string][]Commit {
	result := make(map[string][]Commit)

	for _, commit := range commits {
		result[commit.Type] = append(result[commit.Type], commit)
	}

	return result
}
