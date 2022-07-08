package main

import (
	_ "embed"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//go:embed template/AndroidManifest.xml
var manifestTmpl string

//go:embed template/MainActivity.java
var activityTmpl string

//go:embed template/styles.xml
var stylesTmpl string

//go:embed template/main.xml
var layoutTmpl string

// create command sets up the project with the template files in template/ dir
func create() {
	rootDir := "."
	if len(os.Args) > 2 && os.Args[2] != "." {
		rootDir = os.Args[2]
	}

	pkg := prompt("Package name (com.myapp): ", "com.myapp")
	name := prompt("App name (My App): ", "My App")

	mustMkdir := func(path string) {
		if err := os.Mkdir(path, 0755); err != nil {
			LogF("create", err)
		}
	}

	mustMkdir(rootDir)
	src := filepath.Join(rootDir, "src")
	mustMkdir(src)
	for _, part := range strings.Split(pkg, ".") {
		src = filepath.Join(src, part)
		mustMkdir(src)
	}
	mustMkdir(filepath.Join(rootDir, "lib"))
	mustMkdir(filepath.Join(rootDir, "lib", "arm64-v8a"))
	mustMkdir(filepath.Join(rootDir, "lib", "armeabi-v7a"))
	mustMkdir(filepath.Join(rootDir, "lib", "x86"))
	mustMkdir(filepath.Join(rootDir, "lib", "x86_64"))
	mustMkdir(filepath.Join(rootDir, "res"))
	mustMkdir(filepath.Join(rootDir, "res", "layout"))
	mustMkdir(filepath.Join(rootDir, "res", "values"))

	mustWriteFile := func(path string, content string) {
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			LogF("create", err)
		}
	}

	activity := strings.ReplaceAll(activityTmpl, "{{pkgname}}", pkg)
	mustWriteFile(filepath.Join(src, "MainActivity.java"), activity)

	manifest := strings.ReplaceAll(manifestTmpl, "{{appname}}", name)
	manifest = strings.ReplaceAll(manifest, "{{pkgname}}", pkg)
	mustWriteFile(filepath.Join(rootDir, "AndroidManifest.xml"), manifest)

	mustWriteFile(filepath.Join(rootDir, "res", "layout", "main.xml"), layoutTmpl)
	mustWriteFile(filepath.Join(rootDir, "res", "values", "styles.xml"), stylesTmpl)
	mustWriteFile(filepath.Join(rootDir, ".gitignore"), "build/\n")
}
