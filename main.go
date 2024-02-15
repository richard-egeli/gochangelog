package main

import (
	"bytes"
	"fmt"
	"log"
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

func GetLogs(t1 string, t2 string, filters []string) []string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("git", "log", "--pretty=format:%h - %an, %ar | %s", t1+"..."+t2)
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

func main() {
	yaml, err := config.Read()
	tags, err := GetTags()

	if err != nil {
		panic(err)
	}

	for index, tag := range tags {
		log.Printf("Date %s -- Tag %s", tag.Date, tag.Tag)
		if index+1 >= len(tags) {
			continue
		}

		t1 := tags[index].Tag
		t2 := tags[index+1].Tag

		logs := GetLogs(t2, t1, yaml.Filter)
		for _, log := range logs {
			fmt.Println(log)
		}

		fmt.Printf("\n\n")
	}

	// config := Config{
	// 	Yaml: yaml,
	// 	Tags: GetTag(),
	// 	Logs: GetLog(yaml.Filter),
	// }

	// file, err := os.Create(config.Yaml.Output)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// defer file.Close()
	//
	// _, err = file.WriteString("Testing")
	// if err != nil {
	// 	panic(err)
	// }

}
