package main

import (
	"os/exec"
	"path/filepath"
)

// alignAPK does zipalign
func alignAPK() {
	LogI("build", "running zipalign")

	cmd := exec.Command(zipAlignPath, "-f", "4", filepath.Join("build", "bundle.zip"), filepath.Join("build", "app.apk"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}

// signAPK signs apk with jarsigner and default debug keys
func signAPK() {
	LogI("build", "signing app")

	cmd := exec.Command(jarsignerPath, "-verbose", "-sigalg", "SHA1withRSA", "-digestalg", "SHA1", "-storepass", "android", "-keystore", keyStorePath, filepath.Join("build", "bundle.zip"), "androiddebugkey")
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}
