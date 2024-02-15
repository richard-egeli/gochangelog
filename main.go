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

type Commit struct {
	Type CommitType
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
			if strings.Contains(line, filter) {
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
	filter := []string{"feat"}

	fmt.Println(GetLog(filter))
	fmt.Println(GetTag())
}
