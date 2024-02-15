package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func main() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	args := os.Args[1:] // The first arg is rubbish
	log.Println(args[0])

	cmd := exec.Command("git", "log", "--pretty=format:%h - %an, %ar : %s")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}

	fmt.Println(stdout.String())
}
