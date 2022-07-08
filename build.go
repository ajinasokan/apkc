package main

import (
	"os"
	"os/exec"
	"strings"
)

// build compiles all the source code and bundles into apk file with dependencies
func build() {
	prepare()
	compileRes()
	bundleRes()
	compileJava()
	bundleJava()
	bundleLibs()
	signAPK()
	alignAPK()
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
		if err := os.Mkdir(path, 0755); err != nil && !os.IsExist(err) {
			LogF("build", err)
		}
	}

	mustMkdir("build")
	mustMkdir("build/res")
}

// compileRes compiles the xml files in res dir
func compileRes() {
	res := getFiles("res", "")
	for _, r := range res {
		LogI("build", "compiling", r)
		cmd := exec.Command(aapt2Path, "compile", r, "-o", "build/res/")
		out, err := cmd.CombinedOutput()
		if err != nil {
			LogF("build", string(out))
		}
	}
	err := copyFiles("AndroidManifest.xml", "build/AndroidManifest.xml")
	if err != nil {
		LogF("build", err)
	}
}

// bundleRes bundles all the flat files into apk and generates R.* id file for java
func bundleRes() {
	LogI("build", "bundling resources")

	flats := getFiles("build/res", "")
	args := []string{"link", "-I", androidJar, "--manifest", "build/AndroidManifest.xml", "-o", "build/bundle.apk", "--java", "src"}
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

	args := []string{"-d", "build", "-classpath", androidJar + ":" + jars}
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

	classes := getFiles("build", ".class")
	jars := getFiles("jar", ".jar")

	args := []string{"--lib", androidJar, "--release", "--output", "build"}
	args = append(args, classes...)
	args = append(args, jars...)
	cmd := exec.Command(d8Path, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out), d8Path, args)
	}

	cmd = exec.Command(aaptPath, "add", "bundle.apk", "classes.dex")
	cmd.Dir = "build"
	out, err = cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}

// bundleLibs bundles all the native libs in lib/ dir to apk
func bundleLibs() {
	LogI("build", "bundling native libs")

	copyFiles("lib", "build/lib")

	files := getFiles("build/lib", "")
	if len(files) == 0 {
		LogI("build", "no native libs")
		return
	}
	for i := 0; i < len(files); i++ {
		files[i] = strings.TrimPrefix(files[i], "build/")
	}

	args := []string{"add", "bundle.apk"}
	args = append(args, files...)

	cmd := exec.Command(aaptPath, args...)
	cmd.Dir = "build"
	out, err := cmd.CombinedOutput()
	if err != nil {
		LogF("build", string(out))
	}
}
