package main

import (
	"os"
	"os/exec"

	"gochangelog/pkg/readme"
)

func main() {
	readme, err := readme.Generate()
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("echo", readme)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
}
