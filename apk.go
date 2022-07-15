package main

import (
	"os/exec"
)

// alignAPK does zipalign
func alignAPK() {
	LogI("build", "running zipalign")

	cmd := exec.Command(zipAlignPath, "-f", "4", "build/bundle.zip", "build/app.apk")
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}

// signAPK signs apk with jarsigner and default debug keys
func signAPK() {
	LogI("build", "signing app")

	cmd := exec.Command(jarsignerPath, "-verbose", "-sigalg", "SHA1withRSA", "-digestalg", "SHA1", "-storepass", "android", "-keystore", keyStorePath, "build/bundle.zip", "androiddebugkey")
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}
