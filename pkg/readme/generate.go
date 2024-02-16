package readme

import (
	"gochangelog/pkg/config"
	"gochangelog/pkg/git"
)

func Generate() (string, error) {

	yaml, err := config.Read()

	if err != nil {
		return "", nil
	}

	tags, err := git.GetTags()

	if err != nil {
		return "", nil
	}

	readme := Create()
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

		readme.WriteTag(t1, t2, &tag, yaml)
		commits, err := git.GetCommits(t1, t2, yaml.Filter)
		if err != nil {
			return "", err
		}

		sortedCommits := git.SortCommits(commits)
		for commitType, commits := range sortedCommits {
			readme.WriteType(git.GetCommitTypeName(commitType))

			for _, commit := range commits {
				readme.WriteCommit(&commit, yaml)
			}
		}
	}

	return readme.String(), nil
}
