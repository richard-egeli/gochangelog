package main

import (
	"os"
	"os/exec"

	"gochangelog/pkg/config"
	"gochangelog/pkg/git"
	"gochangelog/pkg/provider"
	"gochangelog/pkg/readme"
)

func Reverse[T any](data []T) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func main() {
	yaml, err := config.Read()

	if err != nil {
		panic(err)
	}

	tags, err := git.GetTags()

	if err != nil {
		panic(err)
	}

	readme := readme.Create()
	provider := provider.Get(provider.Type(yaml.Provider))

	for index, tag := range tags {
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

		url := provider.Diff(t1, t2, yaml)
		readme.WriteTag(tag.Tag, url, tag.Date)
		commits, err := git.GetCommits(t1, t2, yaml.Filter)
		if err != nil {
			panic(err)
		}

		sortedCommits := git.SortCommits(commits)
		for commitType, commits := range sortedCommits {
			readme.WriteType(git.GetCommitTypeName(commitType))

			for _, commit := range commits {
				readme.WriteCommit(commit.Message, commit.Hash, provider.Commit(yaml.RepoURL))
			}
		}
	}

	cmd := exec.Command("echo", readme.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}
