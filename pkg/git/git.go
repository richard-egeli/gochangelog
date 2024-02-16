package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
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
	re := regexp.MustCompile(`tag:\s(\S+)`)
	result := re.FindString(line)
	re = regexp.MustCompile(`[^0-9.\nv]`)
	result = re.ReplaceAllString(result, "")

	return result
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

	cmd := exec.Command("git", "log", "--oneline", "--decorate=no", "--pretty=format:%d<:::>%H<:::>%cd<:::>%s", "--date=format:%d-%m-%Y")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}

	result = strings.Split(stdout.String(), "\n")
	return result, nil
}

func SortCommits(commits []*Commit) map[string][]*Commit {
	result := make(map[string][]*Commit)

	for _, commit := range commits {
		result[commit.Type] = append(result[commit.Type], commit)
	}

	return result
}
