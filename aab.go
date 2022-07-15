package main

import (
	"os/exec"
)

func buildAAB() {
	cmd := exec.Command("bundletool", "build-bundle", "--modules=bundle.zip", "--output=app.aab")
	cmd.Dir = "build"
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}
