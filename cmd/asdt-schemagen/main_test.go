package main

import (
	"os/exec"
	"testing"
)

func TestBuild_ExitsZero(t *testing.T) {
	cmd := exec.Command("go", "build", "./cmd/asdt-schemagen/")
	cmd.Dir = "../.." // repo root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./cmd/asdt-schemagen/ failed:\n%s\n%v", out, err)
	}
}
