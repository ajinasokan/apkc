package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
)

var (
	// java sdk paths
	javaPath      string
	javaBinPath   string
	javacPath     string
	jarsignerPath string

	// android sdk paths
	sdkPath      string
	toolsPath    string
	zipAlignPath string
	aapt2Path    string
	d8Path       string
	adbPath      string

	// android java apis
	androidJar string

	// debug key store file for signing
	keyStorePath string
)

// doctor command checks all the sdk and binary paths are valid
func doctor() {
	findSDKs()

	LogI("doctor", "java", javaPath)
	LogI("doctor", "javac", javacPath)
	LogI("doctor", "jarsigner", jarsignerPath)

	LogI("doctor", "sdk", sdkPath)
	LogI("doctor", "aapt2", aapt2Path)
	LogI("doctor", "d8", d8Path)
	LogI("doctor", "zipalign", zipAlignPath)

	LogI("doctor", "android jar", androidJar)
}

// findSDKs finds Android SDK, Java SDK and paths of other required binaries
func findSDKs() {
	home, err := os.UserHomeDir()
	if err != nil {
		LogF("doctor", err)
	}
	keyStorePath = filepath.Join(home, ".android", "debug.keystore")

	sdkPath = os.Getenv("ANDROID_HOME")

	if sdkPath == "" {
		LogW("doctor", "ANDROID_HOME was not found in environment")
		// try default paths
		switch runtime.GOOS {
		case "darwin":
			sdkPath = filepath.Join(home, "Android", "Sdk")
		case "linux":
			sdkPath = filepath.Join(home, "Library", "Android", "sdk")
		case "windows":
			sdkPath = filepath.Join(home, "AppData", "Local", "Android", "sdk")
		default:
			LogW("doctor", "SDK path unknown")
		}
	}

	if sdkPath != "" {
		btVersion, err := latestBuildTools(filepath.Join(sdkPath, "build-tools"))
		if err != nil {
			LogE("doctor", err)
		}
		toolsPath = filepath.Join(sdkPath, "build-tools", btVersion)
		aapt2Path = filepath.Join(toolsPath, "aapt2")
		d8Path = filepath.Join(toolsPath, "d8")
		zipAlignPath = filepath.Join(toolsPath, "zipalign")
		adbPath = filepath.Join(sdkPath, "platform-tools", "adb")
		api := strings.Split(btVersion, ".")[0]
		androidJar = filepath.Join(sdkPath, "platforms", "android-"+api, "android.jar")
	}

	javaPath = os.Getenv("JAVA_HOME")
	if javaPath != "" {
		javaBinPath = filepath.Join(javaPath, "bin")
		javacPath = filepath.Join(javaBinPath, "javac")
		jarsignerPath = filepath.Join(javaBinPath, "jarsigner")
	} else {
		LogW("doctor", "JAVA_HOME was not found in environment")
		// hope the bin is in path
		javacPath = "javac"
		jarsignerPath = "jarsigner"
	}
}

// latestBuildTools checks the available build tools and picks the most recent one
func latestBuildTools(path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	versions := []*version.Version{}
	for _, f := range files {
		n := f.Name()
		// exclude release candidates
		if f.IsDir() && !strings.Contains(n, "rc") {
			v, err := version.NewVersion(n)
			if err != nil {
				LogW("doctor", "error parsing build tools version '"+n+"'")
				continue
			}
			versions = append(versions, v)
		}
	}

	if len(versions) == 0 {
		return "", errors.New("no usable build tools versions")
	}

	sort.Sort(version.Collection(versions))

	return versions[len(versions)-1].String(), nil
}
