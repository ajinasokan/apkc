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
func signAPK(keyStore, storePass, keyAlias *string) {
	LogI("build", "signing app")

	cmd := exec.Command(apksignerPath, "sign", "--ks-pass", "pass:"+*storePass, "--ks", *keyStore, "--ks-key-alias", *keyAlias, filepath.Join("build", "app.apk"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}
