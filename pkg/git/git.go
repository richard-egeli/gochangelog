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
	BREAKING string = "breaking change"
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
	Tag     string
	Hash    string
	Date    string
	Type    string
	Message string
}

func (c *Commit) IsTag() bool {
	return len(c.Tag) > 0
}

type Tag struct {
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

func ParseCommitTag(line string) string {
	return strings.Trim(strings.ReplaceAll(line, "tag: ", ""), " ")
}

func ParseCommitMessage(line string) (message string, commitType string) {
	lowerCaseLine := strings.ToLower(line)
	for _, commitType := range commitTypes {
		if strings.Contains(lowerCaseLine, commitType+":") {
			return strings.Trim(strings.ReplaceAll(line, commitType+":", ""), ""), commitType
		}
	}

	return line, OTHER
}

func ParseCommit(line string) (*Commit, error) {
	commit := &Commit{}
	parts := strings.Split(line, "<:::>")
	if len(parts) != 4 {
		return nil, errors.New("Commit line is the wrong length")
	}

	message, commitType := ParseCommitMessage(parts[3])
	commit.Tag = ParseCommitTag(parts[0])
	commit.Hash = parts[1]
	commit.Date = parts[2]
	commit.Message = message
	commit.Type = commitType

	return commit, nil
}

func GetCommits() ([]string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var result []string

	cmd := exec.Command("git", "log", "--oneline", "--tags", "--pretty=format:%D<:::>%H<:::>%cd<:::>%s", "--date=format:%d-%m-%Y")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}

	result = strings.Split(stdout.String(), "\n")
	return result, nil
}

// func GetTags() ([]Tag, error) {
// 	var stdout bytes.Buffer
// 	var stderr bytes.Buffer
// 	var tags []Tag
//
// 	cmd := exec.Command("git", "for-each-ref", "--sort=-committerdate", "refs/tags", "--format=%(committerdate:short) | %(refname)")
// 	cmd.Stdout = &stdout
// 	cmd.Stderr = &stderr
// 	err := cmd.Run()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	lines := strings.Split(stdout.String(), "\n")
// 	for _, line := range lines {
// 		var tag Tag
//
// 		tag.Raw = line
// 		line = strings.Replace(line, "refs/tags/", "", -1)
// 		data := strings.Split(line, "|")
// 		if len(data) < 2 {
// 			continue
// 		}
//
// 		tag.Date = strings.Trim(data[0], " ")
// 		tag.Tag = strings.Trim(data[1], " ")
// 		tags = append(tags, tag)
// 	}
//
// 	return tags, nil
// }

func SortCommits(commits []*Commit) map[string][]*Commit {
	result := make(map[string][]*Commit)

	for _, commit := range commits {
		result[commit.Type] = append(result[commit.Type], commit)
	}

	return result
}
