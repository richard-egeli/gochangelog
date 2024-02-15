package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
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
	OTHER    CommitType = "other"
)

type Tag struct {
	Raw  string
	Tag  string
	Date string
}

func GetCommitTypeName(c CommitType) string {
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

func GetCommits(t1 string, t2 string, filters []string) ([]string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("git", "log", "--pretty=format:[%cd] %s [%h]", "--date=format:%d-%m-%Y", t1+".."+t2)
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

	return filteredLines, nil
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

func SortCommits(commits []string) map[CommitType][]string {
	result := make(map[CommitType][]string)

	commitTypes := []CommitType{
		FIX, FEAT, BREAKING, BUILD, CHORE, CI, DOCS, STYLE, REFACTOR, PERF, TEST,
	}

	for _, line := range commits {
		isValid := false
		for _, commitType := range commitTypes {
			typeWithBreak := commitType + ":"
			line = strings.ToLower(line)

			if strings.Contains(line, string(typeWithBreak)) {
				line = strings.Replace(line, " "+string(typeWithBreak), "", -1)
				result[commitType] = append(result[commitType], line)
				isValid = true
				break
			}
		}

		if !isValid {
			result[OTHER] = append(result[OTHER], line)
		}
	}

	return result
}
