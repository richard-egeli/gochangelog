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

type Config struct {
	Yaml *config.YAML
	Tags []string
	Logs []string
}

func GetLog(filters []string) []string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("git", "log", "--pretty=format:%h - %an, %ar | %s")
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

func GetTag() []string {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("git", "tag")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	config.Read()

	err := cmd.Run()
	if err != nil {
		log.Fatal(fmt.Println(fmt.Sprint(err) + ": " + stderr.String()))
	}

	result := stdout.String()
	return strings.Split(result, "\n")
}

func main() {
	yaml, err := config.Read()
	if err != nil {
		panic(err)
	}

	config := Config{
		Yaml: yaml,
		Tags: GetTag(),
		Logs: GetLog(yaml.Filter),
	}

	file, err := os.Create(config.Yaml.Output)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString("Testing")
	if err != nil {
		panic(err)
	}

	for index := range config.Tags {
		next := index + 1
		if next >= len(config.Tags) {
			break
		}

		curTag := config.Tags[index]
		nexTag := config.Tags[next]

		var stdout bytes.Buffer
		var stderr bytes.Buffer

		cmd := exec.Command("git", "log", "--pretty=format:%h - %an, %ar | %s", curTag+".."+nexTag)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}

		result := stdout.String()
		log.Println(curTag)
		log.Println(result)
		log.Println("\n\n")
	}
}
