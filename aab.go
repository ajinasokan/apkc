package main

import (
	"os/exec"
	"path/filepath"
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

// signAAB signs aab with jarsigner and provided debug keys
func signAAB(keyStore, storePass, keyAlias, sigAlg *string) {
	LogI("build", "signing aab")

	// cmd := exec.Command("jarsigner", "-verbose", "-sigalg", "SHA256withRSA", "pass:"+*storePass, "--ks", *keyStore, "--ks-key-alias", *keyAlias, filepath.Join("build", "app.apk"))
	cmd := exec.Command("jarsigner", "-verbose", "-sigalg", *sigAlg, "-digestalg", "SHA1", "-storepass", *storePass, "-keystore", *keyStore, filepath.Join("build", "app.aab"), *keyAlias)
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}
