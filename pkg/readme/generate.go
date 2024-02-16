package readme

import (
	"gochangelog/pkg/config"
	"gochangelog/pkg/git"
	"strings"
)

func Generate() (string, error) {
	tags := []string{}
	commits := map[string][]*git.Commit{}
	readme := Create()

	yaml, err := config.Read()
	if err != nil {
		return "", nil
	}

	commitLines, err := git.GetCommits()
	if err != nil {
		return "", err
	}

	activeTag := "HEAD"
	for _, line := range commitLines {
		commit, err := git.ParseCommit(strings.Clone(line))
		if err != nil {
			return "", err
		}

		if commit.IsTag() {
			activeTag = commit.Tag
			tags = append(tags, commit.Tag)
		}

		commits[activeTag] = append(commits[activeTag], commit)
	}

	prevTag := "HEAD"
	for _, tag := range tags {
		commitList := commits[tag]
		if len(commitList) <= 0 {
			continue
		}

		commit := commitList[0]
		if prevTag != commit.Tag {
			readme.WriteTag(commit.Tag, prevTag, commit.Date, yaml)
			prevTag = commit.Tag
		}

		sorted := git.SortCommits(commitList)
		for commitType, commits := range sorted {
			readme.WriteType(git.GetCommitTypeName(commitType))

			for _, commit := range commits {
				readme.WriteCommit(commit, yaml)
			}
		}
	}

	return readme.String(), nil
}
