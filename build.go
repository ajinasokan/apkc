package main

import (
	"archive/zip"
	"flag"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// build compiles all the source code and bundles into apk file with dependencies
func build() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)

	useAAB := buildCmd.Bool("aab", false, "build aab instead of apk")
	keyStore := buildCmd.String("keystore", keyStorePath, "path to keystore")
	storePass := buildCmd.String("storepass", "android", "keystore password")
	keyAlias := buildCmd.String("keyalias", "androiddebugkey", "key alias to use")

	buildCmd.Parse(os.Args[2:])

	prepare()
	compileRes()
	bundleRes(*useAAB)
	compileJava()
	bundleJava()
	buildBundle(*useAAB)
	if *useAAB {
		buildAAB()
	} else {
		alignAPK()
		signAPK(keyStore, storePass, keyAlias)
	}
}

// clean simply deletes the build dir
func clean() {
	LogI("clean", "removing build/*")
	os.RemoveAll("build")
}

// prepare verifies the project paths and setup the build dir
func prepare() {
	mustExist := func(path string) {
		if _, err := os.Stat(path); err != nil {
			LogF("build", err)
		}
	}

	mustExist("src")
	mustExist("res")
	mustExist("AndroidManifest.xml")

	mustMkdir := func(path string) {
		// only ignore if error is "already exist"
		if err := os.MkdirAll(path, 0755); err != nil && !os.IsExist(err) {
			LogF("build", err)
		}
	}

	mustMkdir(filepath.Join("build", "flats"))
	mustMkdir(filepath.Join("build", "classes"))
}

// compileRes compiles the xml files in res dir
func compileRes() {
	res := getFiles("res", "")
	LogI("build", "compiling resources")
	args := []string{"compile", "-o", filepath.Join("build", "flats")}
	args = append(args, res...)
	cmd := exec.Command(aapt2Path, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}

// bundleRes bundles all the flat files into apk and generates R.* id file for java
func bundleRes(useAAB bool) {
	LogI("build", "bundling resources")

	flats := getFiles("build/flats", ".flat")
	args := []string{"link", "-I", androidJar, "--manifest", "AndroidManifest.xml", "-o", "build", "--java", "src", "--output-to-dir"}
	if useAAB {
		args = append(args, "--proto-format")
	}
	args = append(args, flats...)
	cmd := exec.Command(aapt2Path, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}

// compileJava compiles java files in src dir and uses jar files in the jar dir as classpath
func compileJava() {
	LogI("build", "compiling java files")

	javas := getFiles("src", "")
	jars := strings.Join(getFiles("jar", "jar"), ":")

	args := []string{"-d", filepath.Join("build", "classes"), "-classpath", androidJar + ":" + jars}
	args = append(args, javas...)
	cmd := exec.Command(javacPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out), args)
	}
}

// bundleJava bundles compiled java class files and external jar files into apk
func bundleJava() {
	LogI("build", "bundling classes and jars")

	classes := getFiles(filepath.Join("build", "classes"), ".class")
	jars := getFiles("jar", ".jar")

	args := []string{"--lib", androidJar, "--release", "--output", "build"}
	args = append(args, classes...)
	args = append(args, jars...)
	cmd := exec.Command(d8Path, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out), d8Path, args)
	}
}

func buildBundle(useAAB bool) {
	outFile, err := os.Create(filepath.Join("build", "bundle.zip"))
	if err != nil {
		LogF("build", err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)

	addFileToZip := func(s, d string) {
		dst, err := w.Create(d)
		if err != nil {
			LogF("build", err)
		}

		src, err := os.Open(s)
		if err != nil {
			LogF("build", err)
		}

		_, err = io.Copy(dst, src)
		if err != nil {
			LogF("build", err)
		}
	}

	if useAAB {
		addFileToZip(filepath.Join("build", "AndroidManifest.xml"), filepath.Join("manifest", "AndroidManifest.xml"))
		addFileToZip(filepath.Join("build", "classes.dex"), filepath.Join("dex", "classes.dex"))
		addFileToZip(filepath.Join("build", "resources.pb"), "resources.pb")
	} else {
		addFileToZip(filepath.Join("build", "AndroidManifest.xml"), "AndroidManifest.xml")
		addFileToZip(filepath.Join("build", "classes.dex"), "classes.dex")
		addFileToZip(filepath.Join("build", "resources.arsc"), "resources.arsc")
	}

	files := getFiles(filepath.Join("build", "res"), "")
	for _, f := range files {
		r, err := filepath.Rel("build", f)
		if err != nil {
			LogF("build", err)
		}
		addFileToZip(f, r)
	}

	files = getFiles("assets", "")
	if len(files) > 0 {
		LogI("build", "bundling assets")
	}
	for _, f := range files {
		addFileToZip(f, f)
	}

	files = getFiles("lib", "")
	if len(files) > 0 {
		LogI("build", "bundling native libs")
	}
	for _, f := range files {
		addFileToZip(f, f)
	}

	err = w.Close()
	if err != nil {
		LogF("build", err)
	}
}
