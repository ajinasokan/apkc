package main

import (
	"os/exec"
)

func buildAAB() {
	cmd := exec.Command("bundletool", "build-bundle", "--overwrite", "--modules=bundle.zip", "--output=app.aab")
	if cmd.Err != nil {
		LogF("build", cmd.Err)
	}
	cmd.Dir = "build"
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}
